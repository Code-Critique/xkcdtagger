package main

import (
	"log"
	"sync"

	"os"

	"github.com/Code-Critique/xkcdtagger/http"
	"github.com/Code-Critique/xkcdtagger/redis"
)

var wg = &sync.WaitGroup{}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	s := http.NewServer()
	s.Addr = ":" + port
	s.Handler = &http.Handler{
		ComicHandler: http.NewComicHandler(),
		TagHandler:   http.NewTagHandler(),
	}

	c := redis.NewClient()

	s.Handler.ComicHandler.StorageService = c
	s.Handler.TagHandler.StorageService = c

	err := s.Open()

	wg.Add(1)

	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}

// POST comics/:id {comic-number: int, tag: [string]}
// GET comics/:id/tags
// GET comics
// GET tags
