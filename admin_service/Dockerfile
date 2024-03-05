FROM python:3.9-slim-buster

COPY ./ ./

RUN pip install --upgrade pip
RUN pip install --no-cache-dir --upgrade -r req.txt

ENV PYTHONUNBUFFERED 1

EXPOSE 83

CMD gunicorn main_settings.wsgi:application -c /config/gunicorn.conf.py