from typing import Optional
from urllib.error import URLError
from urllib.request import urlretrieve
import os

from dotenv import dotenv_values
from fastapi import FastAPI, Request
from fastapi.templating import Jinja2Templates
from starlette.responses import HTMLResponse, FileResponse
import markdown

from back.blog import BlogManager


ENV = dotenv_values(".env")

blog = BlogManager(
    ENV.get("REPO_URL") or "",
    ENV.get("REPO_PATH"),
)
app = FastAPI()
templates = Jinja2Templates(directory="./front/templates")


@app.get("/")
async def homepage(request: Request):
    username = ENV.get("GITHUB_PROFILE")
    url = f"https://raw.githubusercontent.com/{username}/{username}/main/README.md"
    try:
        filepath, _ = urlretrieve(url)
    except URLError as e:
        return HTMLResponse(f"An error occured: {e.reason}.")
    filepath = open(filepath).read()
    return templates.TemplateResponse("index.html", {
        "request": request,
        "content": markdown.markdown(filepath),
    })


@app.get("/blog/{name}")
async def blog_page(request: Request, name: str):
    article = None
    for a in blog.articles:
        if name == a.name[:-3]:
            article = a
            break
    if article is None:
        return HTMLResponse("404 error")
    filepath = article.open().read()
    return templates.TemplateResponse("blog.html", {
        "request": request,
        "name": name,
        "content": markdown.markdown(filepath),
    })


def get_file(path: str) -> Optional[str]:
    if os.path.isfile(f"./front/{path}"):
        return f"{path}"
    if os.path.isfile(f"./front/{path}/index.html"):
        return f"{path}/index.html"
    return None


@app.get("/{path:path}")
async def static_files(request: Request, path: str):
    file = get_file(path)
    if not file or file.startswith("templates/"):
        return HTMLResponse("404 not found")
    if file.endswith(".html"):
        return templates.TemplateResponse(file, {
            "request": request,
        })
    file = "./front/" + file
    return FileResponse(file)
