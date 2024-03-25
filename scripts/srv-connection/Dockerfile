FROM python:3.9.7-slim

ARG MSG
ENV ENV_MSG $MSG

COPY netcat-command.py /app/

WORKDIR /app

RUN echo $MSG

RUN apt-get update && apt-get install -y netcat

CMD ["python", "netcat-command.py"]