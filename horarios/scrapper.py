import re
import requests
import json

def scrap_horarios():
    html = requests.get('https://cic.etsisi.upm.es/horarios/').text
    pattern = r"h_grupo\['([A-Z]+[0-9]*)'\]\['([A-Z]+)\']\s*=\s*'(.+?)';"
    matches = re.findall(pattern, html)

    grupos = {}
    dow_map = {
        'L': 0,
        'M': 1,
        'X': 2,
        'J': 3,
        'V': 4,
        'S': 5,
        'D': 6,
    }

    for m in matches:
        g = m[0]
        asig = m[1]
        horarios = m[2]

        if not g in grupos:
            grupos[g] = [[] for _ in range(5)]
        
        for clase in [{'day': dow_map[a[0]], 'hour': a[1:]} for a in horarios.split(';')]:
            grupos[g][clase['day']].append({
                'name': asig,
                'hour': clase['hour']
            })

    for key in grupos:
        for day in range(len(grupos[key])):
            grupos[key][day].sort(key=lambda x: x['hour'])
    
    return grupos

if __name__ == '__main__':
    data = scrap_horarios()
    print(json.dumps(data))