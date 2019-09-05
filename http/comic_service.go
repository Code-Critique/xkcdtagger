package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Code-Critique/xkcdtagger"
	"github.com/gorilla/mux"
)

// ComicHandler reprsents the http routes for comics
type ComicHandler struct {
	*mux.Router

	StorageService xkcdtagger.StorageService

	Logger *log.Logger
}

// NewComicHandler returns a new comic handler
func NewComicHandler() *ComicHandler {
	ch := &ComicHandler{
		Router: mux.NewRouter(),
	}

	ch.HandleFunc("/comics", ch.comicListHandler)
	ch.HandleFunc("/comics/{id}", ch.comicHandler)
	ch.HandleFunc("/comics/{id}/tags", ch.comicTagHandler)

	return ch
}

func (ch *ComicHandler) comicListHandler(w http.ResponseWriter, r *http.Request) {
	c, err := ch.StorageService.ListComics()

	if err != nil {
		w.Write([]byte(err.Error()))
	}

	encodeJSON(w, c, ch.Logger)
}

func (ch *ComicHandler) comicHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	i, err := strconv.Atoi(id)

	if err != nil {
		w.Write([]byte(err.Error()))
	}

	c, err := ch.StorageService.GetComic(xkcdtagger.ComicID(i))

	encodeJSON(w, c, ch.Logger)
}

func (ch *ComicHandler) comicTagHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	i, err := strconv.Atoi(id)

	if err != nil {
		w.Write([]byte(err.Error()))
	}

	t, err := ch.StorageService.GetTagsForComic(xkcdtagger.ComicID(i))

	if err != nil {
		w.Write([]byte(err.Error()))
	}

	encodeJSON(w, t, ch.Logger)
}
