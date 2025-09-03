package main

// library list
import (
	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer()

	// endpoints
	fuego.Get(s, "/", helloWorld)

	// run server
	s.Run()
}

func helloWorld(c fuego.ContextNoBody) (string, error) {
	return "Hello, World!", nil
}
