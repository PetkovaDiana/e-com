FROM python:3.9-slim-buster

RUN apt-get update && apt-get install -y wkhtmltopdf

WORKDIR /code

COPY ./req.txt /code/req.txt

RUN pip install --upgrade pip
RUN pip install --no-cache-dir --upgrade -r /code/req.txt

COPY ./app /code/app

ENV PYTHONUNBUFFERED 1

EXPOSE 79

CMD ["uvicorn", "app.main:app", "--proxy-headers", "--host", "0.0.0.0", "--port", "79"]