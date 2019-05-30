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
cachedCore = []

def do_job(feed, func, send):
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

    if not send:
        return news
    
    newsToSend = []

    for n in news:
        if n['link'] == lastNews:
            break
        newsToSend.append(n)

    if len(newsToSend) == len(news):
        # Assume the format changed, do nothing
        return news
    
    for n in newsToSend:
        r.publish('UNIVERSITY_' + feed.upper() + '_CHANNEL', json.dumps(n, ensure_ascii=False))
        break # Send only one at a time

    print("Published " + str(len(newsToSend)) + " " + feed.upper() + " stories.")

    return news

def get_news_job(send=True):
    global cachedNews

    news = do_job('news', scrapper.news_json_scraper, send)
    cachedNews = news

def get_events_job(send=True):
    global cachedEvents

    news = scrapper.events_json_scraper()
    news = do_job('eventos', scrapper.events_json_scraper, send)
    cachedEvents = news

def get_avisos_job(send=True):
    global cachedAvisos

    news = do_job('avisos', scrapper.avisos_json_scraper, send)
    cachedAvisos = news

def get_core_job(send=True):
    global cachedCore

    news = do_job('coredumped', scrapper.core_dumped_scrapper, send)
    cachedCore = news

def schedule_bg():
    print("Starting jobs...")

    global r

    env = os.environ

    # Setup redis connection
    r = redis.Redis(host=env['REDIS_HOST'], port=6379, db=int(env['REDIS_DB']), password=env['REDIS_PASS'], decode_responses=True)

    get_news_job(False)
    get_events_job(False)
    get_avisos_job(False)
    get_core_job(False)

    schedule.every(10).minutes.do(get_news_job)
    schedule.every(10).minutes.do(get_events_job)
    schedule.every(10).minutes.do(get_avisos_job)
    schedule.every(10).minutes.do(get_core_job)

    while True:
        try:
            schedule.run_pending()
        except Exception as e:
            print('Caught unexpected exception')
            print(e)
        time.sleep(1)

def start():
    print('Starting background worker...')
    thread = threading.Thread(target=schedule_bg, args=())
    thread.daemon = True
    thread.start()
