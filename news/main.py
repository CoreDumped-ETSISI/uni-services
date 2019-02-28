import falcon
import scrapper
import json
import stream

class DirectNews(object):
    def on_get(self, req, resp):
        data = stream.cachedNews

        if data == None:
            resp.body = json.dumps({'message':'Cache sin actualizar. Vuelva a intentar más tarde.'}, ensure_ascii=False)
            resp.status = falcon.HTTP_503
            return

        resp.body = json.dumps(data, ensure_ascii=False)

class DirectEvents(object):
    def on_get(self, req, resp):        
        data = stream.cachedEvents

        if data == None:
            resp.body = json.dumps({'message':'Cache sin actualizar. Vuelva a intentar más tarde.'}, ensure_ascii=False)
            resp.status = falcon.HTTP_503
            return

        resp.body = json.dumps(data, ensure_ascii=False)

class DirectAvisos(object):
    def on_get(self, req, resp):        
        data = stream.cachedAvisos

        if data == None:
            resp.body = json.dumps({'message':'Cache sin actualizar. Vuelva a intentar más tarde.'}, ensure_ascii=False)
            resp.status = falcon.HTTP_503
            return

        resp.body = json.dumps(data, ensure_ascii=False)

print('Starting up!')
stream.start()

app = falcon.API()

app.add_route('/news', DirectNews())
app.add_route('/events', DirectEvents())
app.add_route('/avisos', DirectAvisos())