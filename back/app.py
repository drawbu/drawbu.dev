import os
from typing import Optional
from urllib.request import urlretrieve
from urllib.error import URLError

import markdown
from fastapi import FastAPI, Request
from starlette.responses import HTMLResponse, FileResponse
from fastapi.templating import Jinja2Templates
from dotenv import dotenv_values



ENV = dotenv_values(".env")

app = FastAPI()
templates = Jinja2Templates(directory="./front")


@app.get("/get-profile")
async def get_profile():
    username = ENV.get("GITHUB_PROFILE")
    url = f"https://raw.githubusercontent.com/{username}/{username}/main/README.md"
    try:
        filepath, _ = urlretrieve(url)
    except URLError as e:
        return HTMLResponse(f"An error occured: {e.reason}.")
    filepath = open(filepath).read()
    content = markdown.markdown(filepath)
    return HTMLResponse(content)


def get_file(path: str) -> Optional[str]:
    if os.path.isfile(f"./front/{path}"):
        return f"{path}"
    if os.path.isfile(f"./front/{path}/index.html"):
        return f"{path}/index.html"
    return None


@app.get("/{path:path}")
async def static_files(request: Request, path: str):
    file = get_file(path)
    if not file:
        return HTMLResponse("404 not found")
    if file.endswith(".html"):
        return templates.TemplateResponse(file, {
            "request": request,
        })
    file = "./front/" + file
    return FileResponse(file)
