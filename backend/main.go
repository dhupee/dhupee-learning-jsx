package main

// library list
import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/go-fuego/fuego"
	"github.com/joho/godotenv"
	"github.com/lrstanley/go-ytdlp"
)

// Youtube Download Request Body
type youtubeDl struct {
	YoutubeURL string `json:"youtubeUrl"`
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

	// Download the actual video
	result, err := downloadYtVid(yt.YoutubeURL)
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
func downloadYtVid(youtubeLink string) (string, error) {

	// Make sure yt-dlp is installed on the machine
	ytdlp.MustInstall(context.TODO(), nil)

	// initializing ytdlp
	outputFile := "%(extractor)s - %(title)s.%(ext)s"
	outputDir := "./output/"

	// Create the directory and any necessary parent directories
	err := os.MkdirAll(outputFile, os.ModePerm)
	if err != nil {
		log.Printf("Error creating directory: %v\n", err)
	}

	outputPath := strings.Join([]string{outputDir, outputFile}, "")
	dl := ytdlp.New().
		FormatSort("res,ext:mp4:m4a").
		RecodeVideo("mp4").
		NoPlaylist().
		Output(outputPath)

	// Downloading the video
	_, err = dl.Run(context.TODO(), youtubeLink)
	if err != nil {
		// panic(err)
		return "failed", err
	} else {
		return "success", err
	}
}
