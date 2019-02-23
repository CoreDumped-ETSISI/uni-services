import schedule
import scrapper
import time
import threading

cachedHorario = None

def get_horarios_job():
    global cachedHorario

    cachedHorario = scrapper.scrap_horarios()

def schedule_bg():
    print("Starting jobs...")

    get_horarios_job()

    schedule.every(24).hours.do(get_horarios_job)

    while True:
        schedule.run_pending()
        time.sleep(600)

def start():
    print('Starting background worker...')
    thread = threading.Thread(target=schedule_bg, args=())
    thread.daemon = True
    thread.start()