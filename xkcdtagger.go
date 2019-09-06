package xkcdtagger

// ComicID represents a comic identifier
type ComicID int

// TagID represents tag identifier
type TagID int

// Comic represents a comic
type Comic struct {
	ID   ComicID  `json:"id,omitempty"`
	Tags []string `json:"tags"`
}

// Tag represents a tag
type Tag struct {
	Title   string    `json:"title"`
	ComicID []ComicID `json:"comic_ids"`
}

// StorageService is the interface for stoarge
type StorageService interface {
	// TODO: Split this up into 2 services

	// Comics
	AddComic(Comic) error
	ListComics() ([]Comic, error)
	GetComic(ComicID) (*Comic, error)

	// Tags
	GetTagsForComic(ComicID) ([]string, error)
	GetTagByTitle(string) (*Tag, error)
	ListTags() ([]Tag, error)
	AddTags([]Tag) error
}
