import schedule
import threading
import time
import redis
import os
import scrapper
import json

r = None
cachedNews = []
cachedEvents = []
cachedAvisos = []

def do_job(feed, func):
    print('Getting ' + feed + '...')

    news = func()

    if len(news) == 0:
        # somethin' ain't right...
        print('No ' + feed + ' scrapped. Broken?')
        return
    
    lastNews = r.get('LAST_' + feed.upper() + '_FETCH')
    r.set('LAST_' + feed.upper() + '_FETCH', news[0]['link'])

    if lastNews == None:
        # Not saved, maybe first time?
        # Don't send anything
        return news
    
    newsToSend = []

    for n in news:
        if n['link'] == lastNews:
            break
        newsToSend.append(n)
    
    for n in newsToSend:
        r.publish('UNIVERSITY_' + feed.upper() + '_CHANNEL', json.dumps(n, ensure_ascii=False))

    print("Published " + str(len(newsToSend)) + " " + feed.upper() + " stories.")

    return news

def get_news_job():
    global cachedNews

    news = do_job('news', scrapper.news_json_scraper)
    cachedNews = news

def get_events_job():
    global cachedEvents

    news = scrapper.events_json_scraper()
    news = do_job('events', scrapper.events_json_scraper)
    cachedEvents = news

def get_avisos_job():
    global cachedAvisos

    news = do_job('avisos', scrapper.avisos_json_scraper)
    cachedAvisos = news

def schedule_bg():
    print("Starting jobs...")

    global r

    env = os.environ

    # Setup redis connection
    r = redis.Redis(host=env['REDIS_HOST'], port=6379, db=int(env['REDIS_DB']), password=env['REDIS_PASS'], decode_responses=True)

    get_news_job()
    get_events_job()
    get_avisos_job()

    schedule.every(10).minutes.do(get_news_job)
    schedule.every(10).minutes.do(get_events_job)
    schedule.every(10).minutes.do(get_avisos_job)

    while True:
        schedule.run_pending()
        time.sleep(1)

def start():
    print('Starting background worker...')
    thread = threading.Thread(target=schedule_bg, args=())
    thread.daemon = True
    thread.start()
