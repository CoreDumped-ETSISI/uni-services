import schedule
import threading
import time
import redis
import os
import scrapper
import json

r = None

def get_news_job():
    print('Getting news...')

    news = scrapper.news_json_scraper()

    if len(news) == 0:
        # somethin' ain't right...
        print('No news scrapped. Broken?')
        return

    lastNews = r.get('LAST_NEWS_FETCH')
    r.set('LAST_NEWS_FETCH', news[0].link)

    if lastNews == None:
        # Not saved, maybe first time?
        # Don't send anything
        return
    
    newsToSend = []

    for n in news:
        if n.link == lastNews:
            break
        newsToSend.append(n)
    
    for n in newsToSend:
        r.publish('UNIVERISTY_NEWS_CHANNEL', json.dumps(n, ensure_ascii=False))

    print("Published " + str(len(newsToSend)) + " news stories.")

def get_events_job():
    print('Getting events...')

    news = scrapper.events_json_scraper()

    if len(news) == 0:
        # somethin' ain't right...
        print('No events scrapped. Broken?')
        return

    lastNews = r.get('LAST_EVENTS_FETCH')
    r.set('LAST_EVENTS_FETCH', news[0].link)

    if lastNews == None:
        # Not saved, maybe first time?
        # Don't send anything
        return
    
    newsToSend = []

    for n in news:
        if n.link == lastNews:
            break
        newsToSend.append(n)
    
    for n in newsToSend:
        r.publish('UNIVERISTY_EVENTS_CHANNEL', json.dumps(n, ensure_ascii=False))

    print("Published " + str(len(newsToSend)) + " events stories.")

def get_avisos_job():
    print('Getting avisos...')

    news = scrapper.avisos_json_scraper()

    if len(news) == 0:
        # somethin' ain't right...
        print('No avisos scrapped. Broken?')
        return

    lastNews = r.get('LAST_AVISOS_FETCH')
    r.set('LAST_AVISOS_FETCH', news[0].link)

    if lastNews == None:
        # Not saved, maybe first time?
        # Don't send anything
        return
    
    newsToSend = []

    for n in news:
        if n.link == lastNews:
            break
        newsToSend.append(n)
    
    for n in newsToSend:
        r.publish('UNIVERISTY_AVISOS_CHANNEL', json.dumps(n, ensure_ascii=False))

    print("Published " + str(len(newsToSend)) + " avisos stories.")

def schedule_bg():
    while True:
        schedule.run_pending()
        time.sleep(1)

def start():
    env = os.environ

    # Setup redis connection
    r = redis.Redis(host=env['REDIS_HOST'], port=6379, db=int(env['REDIS_DB']), password=env['REDIS_PASS'], decode_responses=True)

    schedule.every(10).minutes.do(get_news_job)
    schedule.every(10).minutes.do(get_events_job)
    schedule.every(10).minutes.do(get_avisos_job)

    thread = threading.Thread(target=schedule_bg, args=())
    thread.daemon = True
    thread.start()