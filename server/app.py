import os

from fastapi import FastAPI
from starlette.responses import HTMLResponse, FileResponse


app = FastAPI()


# Get requests
@app.get("/")
async def get():
    return HTMLResponse(open("./public/pages/index.html").read())


@app.get("/{path:path}")
async def static_files(path: str):
    if os.path.isfile(f"./public/{path}"):
        return FileResponse(f"./public/{path}")
    return HTMLResponse("404 Not Found")

