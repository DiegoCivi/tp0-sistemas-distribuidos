FROM python:3.9.7-slim

COPY server-connection-test.py /app/

WORKDIR /app

RUN apt-get update && apt-get install -y netcat

CMD ["python", "server-connection-test.py"]
