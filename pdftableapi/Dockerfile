FROM python:2.7

WORKDIR /usr/src/app

RUN pip install --no-cache-dir requests
RUN pip install --no-cache-dir flask

RUN apt update

RUN yes | apt install python-tk ghostscript

RUN pip install --no-cache-dir camelot-py[all]

COPY . .

CMD ["python", "./main.py"]