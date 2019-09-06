package http

import (
	"encoding/json"
	"io/ioutil"
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

	th.Router.HandleFunc("/tags", th.tagHandler).Methods("POST", "GET")

	return th
}

func (th *TagHandler) tagHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		var t xkcdtagger.Tag
		err = json.Unmarshal(b, &t)

		log.Println("post", t)

		temp := []xkcdtagger.Tag{t}

		err = th.StorageService.AddTags(temp)

		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		return
	}

	t, err := th.StorageService.ListTags()

	if err != nil {
		w.Write([]byte(err.Error()))
	}

	encodeJSON(w, t, th.Logger)
}
