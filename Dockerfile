FROM ubuntu

# Setup repo
WORKDIR /poll
COPY . .

RUN sudo apt update
RUN sudo apt install \
    python3.12       \
    python3.12-venv

RUN pip install -r requirements.txt

CMD gunicorn                              \
    --workers 3                           \
    --bind unix:/pool/process/server.sock \
    -m 007                                \
    wsgi:app
