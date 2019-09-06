package redis

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/Code-Critique/xkcdtagger"
	"github.com/go-redis/redis"
)

// Client represents the redis client
type Client struct {
	*redis.Client
}

// NewClient returns a new client
func NewClient() *Client {
	return &Client{
		Client: redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0, // use default DB
		}),
	}
}

// ListComics lists all comics in redis
func (c *Client) ListComics() ([]xkcdtagger.Comic, error) {
	sc := c.HVals("comics")
	strings, err := sc.Result()

	if err != nil {
		return nil, err
	}

	var out []xkcdtagger.Comic

	for _, s := range strings {

		log.Println(s)

		var item xkcdtagger.Comic

		err := json.Unmarshal([]byte(s), &item)

		if err != nil {
			return nil, err
		}

		out = append(out, item)
	}

	return out, nil
}

// GetComic gets comic by id
func (c *Client) GetComic(id xkcdtagger.ComicID) (*xkcdtagger.Comic, error) {
	key := strconv.Itoa(int(id))

	sc := c.HGet("comics", key)

	s, err := sc.Result()

	if err != nil {
		return nil, err
	}

	var item xkcdtagger.Comic

	err = json.Unmarshal([]byte(s), &item)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

// GetTagsForComic gets tags for comic
func (c *Client) GetTagsForComic(id xkcdtagger.ComicID) ([]string, error) {
	key := strconv.Itoa(int(id))

	sc := c.HGet("comics", key)

	s, err := sc.Result()

	if err != nil {
		return nil, err
	}

	var item xkcdtagger.Comic

	err = json.Unmarshal([]byte(s), &item)

	if err != nil {
		return nil, err
	}

	return item.Tags, nil
}

// ListTags list all tags
func (c *Client) ListTags() ([]xkcdtagger.Tag, error) {
	strings, err := c.HVals("tags").Result()

	log.Println(strings, err)

	if err != nil {
		return nil, err
	}

	var out []xkcdtagger.Tag

	for _, s := range strings {
		var item xkcdtagger.Tag

		err := json.Unmarshal([]byte(s), &item)

		if err != nil {
			return nil, err
		}

		out = append(out, item)
	}

	return out, nil
}

// GetTagByTitle gets a tag by tag title
func (c *Client) GetTagByTitle(title string) (*xkcdtagger.Tag, error) {
	t, err := c.HGet("tags", title).Result()

	if err != nil {
		// If not found, can use other redis statements
		if err.Error() == "redis: nil" {
			return nil, nil
		}

		return nil, err
	}

	var tag xkcdtagger.Tag

	err = json.Unmarshal([]byte(t), &tag)

	if err != nil {
		return nil, err
	}

	return &tag, nil
}

// AddTags adds tags into redis
func (c *Client) AddTags(tags []xkcdtagger.Tag) error {

	for _, t := range tags {
		existingTag, err := c.GetTagByTitle(t.Title)

		if err != nil {
			return err
		}

		if existingTag != nil {
			temp := *existingTag

			set := make(map[int]bool)

			for _, cid := range temp.ComicID {
				set[int(cid)] = true
			}

			keys := make([]xkcdtagger.ComicID, 0, len(set))
			for k := range set {
				keys = append(keys, xkcdtagger.ComicID(k))
			}

			temp.ComicID = keys

			json, err := json.Marshal(temp)

			if err != nil {
				return err
			}

			_, err = c.HSet("tags", t.Title, json).Result()

			if err != nil {
				return err
			}
		} else {
			json, err := json.Marshal(t)

			if err != nil {
				return err
			}

			_, err = c.HSet("tags", t.Title, json).Result()

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) comicExists(id xkcdtagger.ComicID) (*xkcdtagger.Comic, error) {
	sid := strconv.Itoa(int(id))

	comic, err := c.HGet("comics", sid).Result()

	if err != nil {
		// If not found, can use other redis statements
		if err.Error() == "redis: nil" {
			return nil, nil
		}

		return nil, err
	}

	var outComic xkcdtagger.Comic

	err = json.Unmarshal([]byte(comic), &outComic)

	if err != nil {
		return nil, err
	}

	return &outComic, nil
}

// AddComic adds comic into redis
func (c *Client) AddComic(comic xkcdtagger.Comic) error {
	oComic, err := c.comicExists(comic.ID)

	if err != nil {
		return err
	}

	sid := strconv.Itoa(int(comic.ID))

	if oComic == nil {

		json, err := json.Marshal(comic)

		if err != nil {
			return err
		}

		bc := c.HSet("comics", sid, json)

		_, err = bc.Result()

		if err != nil {
			return err
		}
	} else {
		set := make(map[string]bool)

		for _, t := range oComic.Tags {
			set[t] = true
		}

		for _, t := range comic.Tags {
			set[t] = true
		}

		keys := make([]string, 0, len(set))
		for k := range set {
			keys = append(keys, k)
		}

		outComic := *oComic
		outComic.Tags = keys

		json, err := json.Marshal(outComic)

		if err != nil {
			return err
		}

		_, err = c.HSet("comics", sid, json).Result()

		if err != nil {
			return err
		}
	}

	return nil
}
