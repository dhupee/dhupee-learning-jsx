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

# Configure basic logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
)
# Create a logger instance
logger = logging.getLogger(__name__)

# Create output directory for Youtube Downloader
OUTPUT_DIR = randdir.create_output_dir()
logger.warning("Output Directory: " + OUTPUT_DIR)


# Endpoints
@app.get("/")
async def read_root():
    return {"message": "Hello World"}


@app.get("/getMusic")
async def get_music(YTData: YoutubeData):

    # Body
    url = YTData.URL

    if url == "":
        logger.info("User not adding URL")
        return HTTPException(status_code=400, detail="No URL provided")

    # Download from youtube
    result = download_from_youtube(url)

    # Return result
    if result:
        return {"message": "Music Downloaded"}
    else:
        raise HTTPException(status_code=500, detail="Music Failed Downloaded")


def download_from_youtube(url: str) -> tuple[bool, str]:
    """
    Download MP3 file from Youtube URL

    Args:
        url: string

    Returns:
        success status: bool
        result: string

    """

    # Youtube Downloader Parameters
    ydl_opts = {
        "format": "best",
        # "outtmpl": "%(title)s.%(ext)s",
        "outtmpl": os.path.join(OUTPUT_DIR, "%(title)s.%(ext)s"),
        "noplaylist": True,
        "cookiesfrombrowser": ("firefox",),
        # post process downloaded video, FFmpeg magic
        "postprocessors": [
            {
                "key": "FFmpegExtractAudio",
                "preferredcodec": "mp3",  # Convert to MP3
                "preferredquality": "192",  # Set desired MP3 quality (e.g., 192kbps)
            }
        ],
        # Metadata and thumbnail embedding
        "writethumbnail": True,  # Download the thumbnail
        "embedthumbnail": True,  # Embed it in the audio file
        "addmetadata": True,  # Add metadata like title and uploader
        "embedmetadata": True,  # Embed metadata in the file
    }

    if url == "":
        return False, "No URL provided"

    # Attempt to download
    with yt_dlp.YoutubeDL(ydl_opts) as ydl:
        ydl.download([url])
        return True, "Success"
