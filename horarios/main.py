import falcon
import cache
import json

class Horario(object):
    def on_get(self, req, resp, grupo=None):
        horario = cache.cachedHorario

        if horario == None:
            resp.body = json.dumps({'message':'Cache sin actualizar. Vuelva a intentar más tarde.'}, ensure_ascii=False)
            resp.status = falcon.HTTP_503
            return

        if grupo != None:
            grupo = grupo.upper()
            if not grupo in horario:
                resp.body = json.dumps({'message':'Ese grupo no existe!'}, ensure_ascii=False)
                resp.status = falcon.HTTP_400
                return
            horario = horario[grupo] 

        resp.body = json.dumps(horario, ensure_ascii=False)


print('Starting up!')
cache.start()

app = falcon.API()

app.add_route('/horario', Horario())
app.add_route('/horario/{grupo}', Horario())