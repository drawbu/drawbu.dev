FROM ubuntu

# Setup repo
WORKDIR /poll
COPY . .

RUN apt-get update
RUN apt-get install -y \
    python3.12         \
    python3.12-venv

RUN pip install -r requirements.txt

CMD gunicorn                              \
    --workers 3                           \
    --bind unix:/pool/process/server.sock \
    -m 007                                \
    wsgi:app
