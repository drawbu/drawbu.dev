FROM ubuntu

# Setup repo
WORKDIR /poll
COPY . .

RUN apt-get update
RUN apt-get install -y \
    python3.11         \
    python3.11-venv

RUN python3.11 -m venv venv
RUN venv/bin/pip install -r requirements.txt

CMD venv/bin/gunicorn                              \
    -k uvicorn.workers.UvicornWorker \
    --workers 3                           \
    --bind unix:/pool/process/server.sock \
    -m 007                                \
    wsgi:app
