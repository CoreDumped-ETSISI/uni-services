from flask import Flask, request, make_response, jsonify

import camelot
import tempfile
import os
import requests

app = Flask(__name__)

def downloadFile(uri):
    f = tempfile.NamedTemporaryFile(delete=False, suffix=".pdf")
    r = requests.get(uri)
    f.write(r.content)    
    f.close()
    return f.name
    


@app.route("/table", methods=['POST'])
def parsetable():
    js = request.get_json(force=True)
    tmpname = downloadFile(js['pdf'])

    settings = {}

    if 'settings' in js:
        settings = js['settings']
    tables = None
    try:        
        tables = camelot.read_pdf(tmpname, **settings)
    except:
        return jsonify({'message': 'internal server error'}), 500
    finally:
        os.remove(tmpname)
    
    if len(tables) == 0:
        return jsonify({
            'message': 'no tables found'
        }), 404

    table_results = []

    for t in tables:
        table_results.append({
            'report': t.parsing_report,
            'data': t.data
        })

    return jsonify(table_results)


def main():
    app.run(host='0.0.0.0', port=8080)

if __name__ == '__main__':
    main()
