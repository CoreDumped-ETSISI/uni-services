from bs4 import BeautifulSoup
import requests
import feedparser

def news_json_scraper():
    html = requests.get("https://etsisi.upm.es/noticias")
    soup = BeautifulSoup(html.text, 'html.parser')
    allnews = []
    for idx, row in enumerate(soup.find(id='main-content').findAll('a')):
        news = {}
        news["text"] = row.string
        news["a-link"] = '<a href="https://etsisi.upm.es' + row.get("href") + '">' + row.string + '</a>'
        news["link"] = 'https://etsisi.upm.es' + row.get("href")
        allnews.append(news)
    return allnews


# Listado de eventos disponible en la web de la upm
def events_json_scraper():
    """ No funciona por parte de etsisi.upm.es!!! """
    return []

    html = requests.get("https://etsisi.upm.es")
    soup = BeautifulSoup(html.text, 'html.parser')
    allevents = []
    for idx, row in enumerate(soup.find(id='block-views-calendario-block-2--2').findAll('a')):
        event = {}
        event["text"] = row.string
        event["a-link"] = '<a href="https://etsisi.upm.es' + row.get("href") + '">' + row.string + '</a>'
        event["link"] = 'https://etsisi.upm.es' + row.get("href")
        allevents.append(event)
    return allevents


def avisos_json_scraper():
    html = requests.get("https://etsisi.upm.es/alumnos/avisos")
    soup = BeautifulSoup(html.text, 'html.parser')
    avisos = []

    for idx, row in enumerate(soup.find(id='main-content').findAll('a')):
        if idx > 5:
            break
        aviso = {}
        aviso["text"] = row.string
        aviso["a-link"] = '<a href="https://etsisi.upm.es' + row.get("href") + '">' + row.string + '</a>'
        aviso["link"] = 'https://etsisi.upm.es' + row.get("href")
        avisos.append(aviso)

    # more = {}
    # more["text"] = "Más avisos..."
    # more["a-link"] = '<a href="https://etsisi.upm.es/alumnos/avisos">Más avisos...</a>'
    # more["link"] = 'https://etsisi.upm.es/alumnos/avisos'
    # avisos.append(more)

    return avisos

def core_dumped_scrapper():
    feed = feedparser.parse('https://coredumped.es/feed/')

    allentries = []
    for entry in feed.entries:
        event = {}
        event["text"] = entry.title # + '\n\n' + entry.description
        event["a-link"] = '<a href="' + entry.link + '">' + entry.title + '</a>'
        event["link"] = entry.link
        allentries.append(event)
    return allentries

if __name__ == '__main__':
    print('Ultimas noticias:')
    print(news_json_scraper())
    print('')
    print('Eventos:')
    print(events_json_scraper())
    print('')
    print('Avisos:')
    print(avisos_json_scraper())