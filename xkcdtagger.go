package xkcdtagger

// ComicID represents a comic identifier
type ComicID int

// Comic represents a comic
type Comic struct {
	ID   ComicID `json:"id"`
	Tags []Tag   `json:"tags"`
}

// Tag represents a tag
type Tag struct {
	ComicID ComicID `json:"id"`
	Name    string  `json:"name"`
}

// StorageService is the interface for stoarge
type StorageService interface {
	ListComics() ([]Comic, error)
	GetComic(ComicID) (*Comic, error)
	GetTagsForComic(ComicID) ([]Tag, error)
	GetTags() ([]Tag, error)
}
