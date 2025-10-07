import os

from utils import randdir

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import yt_dlp
import logging
import uvicorn


# Class for downloading music from Youtube
class YoutubeData(BaseModel):
    URL: str


# Init FastAPI
app = FastAPI()

# Create output directory for Youtube Downloader
OUTPUT_DIR = randdir.create_output_dir()


# Endpoints
@app.get("/")
async def read_root():
    return {"message": "Hello World"}


@app.get("/getMusic")
async def get_music(YTData: YoutubeData):

    # Body
    url = YTData.URL

    if url == "":
        return HTTPException(status_code=400, detail="No URL provided")

    # Download from youtube
    result = download_from_youtube(url)

    if result:
        return {"message": "Music Downloaded"}
    else:
        raise HTTPException(status_code=500, detail="Music Failed Downloaded")


def download_from_youtube(url: str) -> tuple[bool, str]:

    ydl_opts = {
        "format": "bestaudio/best",
        "extractaudio": True,
        "audioformat": "mp3",
        # "outtmpl": "%(title)s.%(ext)s",
        "outtmpl": os.path.join(OUTPUT_DIR, "%(title)s.%(ext)s"),
        "noplaylist": True,
        # Metadata and thumbnail embedding
        "writethumbnail": True,  # Download the thumbnail
        "embedthumbnail": True,  # Embed it in the audio file
        "addmetadata": True,  # Add metadata like title and uploader
        "embedmetadata": True,  # Embed metadata in the file
    }

    if url == "":
        return False, "No URL provided"

    with yt_dlp.YoutubeDL(ydl_opts) as ydl:
        ydl.download([url])
        return True, "Success"
