import os

from urllib.request import urlretrieve
from urllib.error import URLError

import markdown
from fastapi import FastAPI
from starlette.responses import HTMLResponse, FileResponse
from fastapi.templating import Jinja2Templates


app = FastAPI()
templates = Jinja2Templates(directory="./front/templates")

GH_PROFILE = "https://raw.githubusercontent.com/drawbu/drawbu/main/README.md"


@app.get("/get-profile")
async def get_profile():
    try:
        filepath, _ = urlretrieve(GH_PROFILE)
    except URLError as e:
        return HTMLResponse(f"An error occured: {e.reason}.")
    filepath = open(filepath).read()
    content = markdown.markdown(filepath)
    return HTMLResponse(content)


@app.get("/{path:path}")
async def static_files(path: str):
    if os.path.isfile(f"./front/{path}"):
        return FileResponse(f"./front/{path}")
    if os.path.isfile(f"./front/{path}/index.html"):
        return FileResponse(f"./front/{path}/index.html")
    return HTMLResponse("404 not found.")
