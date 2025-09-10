package main

// library list
import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/go-fuego/fuego"
	"github.com/joho/godotenv"
	"github.com/lrstanley/go-ytdlp"
)

// Youtube Download Request Body
type youtubeDl struct {
	YoutubeURL     string `json:"youtubeUrl"`
	OutputFormat   string `json:"outputFormat"`
	EmbedThumbnail bool   `json:"embedThumnail"`
}

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load .env")
		log.Println("Using default environment variables")
	}

	// Init Server
	s := fuego.NewServer()

	// Endpoints
	fuego.Get(s, "/", helloWorld)
	fuego.Get(s, "/ping", ping)
	fuego.Get(s, "/getYoutubeVid", getYoutubeVid)

	// Run Server
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// Hello World Endpoint
func helloWorld(c fuego.ContextNoBody) (string, error) {
	return "Hello, World!", nil
}

func ping(c fuego.ContextNoBody) (string, error) {
	return "pong!", nil
}

// Get Youtube Vid Endpoint
func getYoutubeVid(c fuego.ContextWithBody[youtubeDl]) (string, error) {
	// Get the body of request
	yt, err := c.Body()
	if err != nil {
		return "Cannot find body", err
	}

	/*=====================DEFAULT VALUES FOR DOWNLOADER======================================*/

	// Default to mp4 if not declared
	format := yt.OutputFormat
	if format == "" {
		format = "mp4"
	}

	/*=========================DEFAULT VALUES ENDS===========================================*/

	// Download the actual video
	result, err := downloadYtVid(yt.YoutubeURL, format, yt.EmbedThumbnail)
	if err != nil {
		// return "Failed to download video", err
		return result, err
	} else {
		// return strings.Join([]string{"Success Download Video to: ", outputPath}, ""), err
		// return "Video Sucessfully Downloaded", err
		return result, err
	}
}

// Download Youtube Vid
func downloadYtVid(youtubeLink string, outputFormat string, embedThumbnail bool) (string, error) {

	// Make sure yt-dlp, and ffmpeg is installed on the machine
	ytdlp.MustInstall(context.TODO(), nil)
	ytdlp.MustInstallFFmpeg(context.TODO(), nil)

	// Return error message if Youtube link is empty
	if youtubeLink == "" {
		// return "Youtube Link is empty", nil
		err := errors.New("empty youtube link")
		return "", fuego.BadRequestError{Title: "Youtube Link Empty", Err: err}
	}

	// Initializing yt-dlp output
	outputFile := "%(extractor)s - %(title)s.%(ext)s"
	outputDir := "./output/"

	// Create the directory and any necessary parent directories
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Printf("Error creating directory: %v\n", err)
	} else {
		log.Println(strings.Join([]string{"Created directory: ", outputDir}, ""))
	}

	// Initializing yt-dlp parameters
	outputPath := strings.Join([]string{outputDir, outputFile}, "")
	dl := ytdlp.New().
		FormatSort("res,ext:mp4:m4a").
		RecodeVideo(outputFormat).
		NoPlaylist().
		Output(outputPath)

	// Additional parameters
	// TODO: TEST IT
	if embedThumbnail {
		dl = dl.EmbedThumbnail()
	}

	// if embedSubs {
	// 	dl = dl.EmbedSubs()
	// } else if embedThumbnail {
	// 	dl = dl.EmbedThumbnail()
	// }

	// Downloading the video
	_, err = dl.Run(context.TODO(), youtubeLink)
	if err != nil {
		// panic(err)
		return strings.Join([]string{"Failed to download video, ", err.Error()}, ""), err
	} else {
		return "Video sucessfully downloaded", nil
	}
}
