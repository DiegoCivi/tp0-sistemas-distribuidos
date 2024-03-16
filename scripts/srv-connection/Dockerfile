FROM python:3.9.7-slim

COPY netcat-command.py /app/

WORKDIR /app

RUN apt-get update && apt-get install -y netcat

CMD ["python", "netcat-command.py"]
