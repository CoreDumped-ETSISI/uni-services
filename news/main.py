import falcon
import scrapper
import json
import stream

class DirectNews(object):
    def on_get(self, req, resp):        
        resp.body = json.dumps(scrapper.news_json_scraper(), ensure_ascii=False)

class DirectEvents(object):
    def on_get(self, req, resp):        
        resp.body = json.dumps(scrapper.events_json_scraper(), ensure_ascii=False)

class DirectAvisos(object):
    def on_get(self, req, resp):        
        resp.body = json.dumps(scrapper.avisos_json_scraper(), ensure_ascii=False)

stream.start()

app = falcon.API()

app.add_route('/news', DirectNews())
app.add_route('/events', DirectEvents())
app.add_route('/avisos', DirectAvisos())