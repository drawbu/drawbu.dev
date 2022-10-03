from flask import Flask, send_from_directory, request, jsonify


app = Flask(__name__)


# Path for our main Svelte page
@app.route("/")
def home():
    return send_from_directory("client/public", "index.html")


# Path for all the static files (compiled JS/CSS, etc.)
@app.route("/<path:path>")
def static_files(path):
    return send_from_directory("client/public", path)


if __name__ == "__main__":
    app.run()
