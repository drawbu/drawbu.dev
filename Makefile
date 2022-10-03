PM = npm
NPM_FLAGS = --prefix client

V_BIN = venv/bin
J_BIN = client/node_modules/.bin

ENV = .flaskenv

FRONT_BUILD = client/public/build
BUILD_MODE = build

GUNICORN_CMD = $(V_BIN)/gunicorn
GUNICORN_ARGS =


all: start


$(J_BIN)/rollup:
	$(PM) install $(NPM_FLAGS)


$(FRONT_BUILD): $(J_BIN)/rollup
	$(PM) run $(NPM_FLAGS) $(BUILD_MODE)


$(V_BIN):
	python3 -m venv venv
	chmod +x $(V_BIN)/activate
	./$(V_BIN)/activate


$(ENV):
	echo 'DEBUG_MODE=false' > $(ENV)


$(GUNICORN_CMD): $(V_BIN) $(ENV)
	$(V_BIN)/pip install -r requirements.txt


start: $(GUNICORN_CMD) client/public/build
	$(GUNICORN_CMD) wsgi:app $(GUNICORN_ARGS)


dev:
	make -j2 start BUILD_MODE='dev &' GUNICORN_ARGS="--reload"


clean:
	rm -rf */__pycache__
	rm -rf *egg-info
	rm -f .flaskenv
	rm -rf client/node_modules
	rm -rf venv


.PHONY: all clean dev fclean start