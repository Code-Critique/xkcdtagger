package http

import (
	"log"
	"net/http"

	"github.com/Code-Critique/xkcdtagger"
	"github.com/gorilla/mux"
)

// TagHandler reprsents the http routes for tags
type TagHandler struct {
	*mux.Router

	StorageService xkcdtagger.StorageService

	Logger *log.Logger
}

// NewTagHandler returns a new tag handler
func NewTagHandler() *TagHandler {
	th := &TagHandler{
		Router: mux.NewRouter(),
	}

	th.Router.HandleFunc("/tags", th.tagHandler)

	return th
}

func (th *TagHandler) tagHandler(w http.ResponseWriter, r *http.Request) {
	t, err := th.StorageService.GetTags()

	if err != nil {
		w.Write([]byte(err.Error()))
	}

	encodeJSON(w, t, th.Logger)
}
