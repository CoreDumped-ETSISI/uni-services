import falcon
import scrapper
import json
import stream

class DirectNews(object):
    def on_get(self, req, resp):
        stream.cacheLock.acquire()
        data = stream.cachedNews
        stream.cacheLock.release()

        resp.body = json.dumps(data, ensure_ascii=False)

class DirectEvents(object):
    def on_get(self, req, resp):        
        stream.cacheLock.acquire()
        data = stream.cachedEvents
        stream.cacheLock.release()

        resp.body = json.dumps(data, ensure_ascii=False)

class DirectAvisos(object):
    def on_get(self, req, resp):        
        stream.cacheLock.acquire()
        data = stream.cachedAvisos
        stream.cacheLock.release()

        resp.body = json.dumps(data, ensure_ascii=False)

print('Starting up!')
stream.start()

app = falcon.API()

app.add_route('/news', DirectNews())
app.add_route('/events', DirectEvents())
app.add_route('/avisos', DirectAvisos())