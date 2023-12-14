import os

from fastapi import FastAPI
from starlette.responses import HTMLResponse, FileResponse
from fastapi.templating import Jinja2Templates


app = FastAPI()
templates = Jinja2Templates(directory="./front/templates")

@app.get("/{path:path}")
async def static_files(path: str):
    if os.path.isfile(f"./front/{path}"):
        return FileResponse(f"./front/{path}")
    if os.path.isfile(f"./front/{path}/index.html"):
        return FileResponse(f"./front/{path}/index.html")
    return HTMLResponse("404 not found.")
