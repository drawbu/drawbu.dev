V_BIN = venv/bin

SERVER_CMD = $(V_BIN)/uvicorn
SERVER_ARGS = --port 8080
PROD_BIND =

all: run

$(V_BIN):
	@ python3 -m venv venv
	@ chmod +x $(V_BIN)/activate
	@ ./$(V_BIN)/activate

$(SERVER_CMD): $(V_BIN) $(ENV)
	@ $(V_BIN)/pip install -r requirements.txt

run: $(SERVER_CMD)
	@ $(SERVER_CMD) server:app $(SERVER_ARGS)

dev: SERVER_ARGS += --reload
dev: run

prod: SERVER_CMD = $(V_BIN)/gunicorn
prod: SERVER_ARGS = -k uvicorn.workers.UvicornWorker -b $(PROD_BIND)
prod: run

.PHONY: run dev prod

clean:
	@ $(RM) -r */__pycache__
	@ $(RM) .flaskenv

fclean: clean
	@ $(RM) -r venv

.PHONY: clean fclean
