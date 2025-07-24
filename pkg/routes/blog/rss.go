package blog

import (
	"app/pkg/app"
	"fmt"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/feeds"
)

func RssHandler(serv *app.Server, w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	serv.Cache_route(w, r, 3600)
	feed := makeFeed()

	rss, err := feed.ToRss()
	if err != nil {
		return nil, fmt.Errorf("Could not render RSS feed: %s", err)
	}
	w.Header().Set("Content-Type", "application/rss+xml")
	io.WriteString(w, rss)
	return nil, nil
}

func AtomHandler(serv *app.Server, w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	serv.Cache_route(w, r, 3600)
	feed := makeFeed()

	atom, err := feed.ToAtom()
	if err != nil {
		return nil, fmt.Errorf("Could not render Atom feed: %s", err)
	}
	w.Header().Set("Content-Type", "application/atom+xml")
	io.WriteString(w, atom)
	return nil, nil
}

func makeFeed() feeds.Feed {
	feed := feeds.Feed{
		Title:       "drawbu.dev blog",
		Link:        &feeds.Link{Href: "https://drawbu.dev/blog"},
		Description: "software engineering student talking about cool stuff",
		Author:      &feeds.Author{Name: "Clément (drawbu)", Email: "contact@drawbu.dev"},
	}

	var items = make([]*feeds.Item, 0)
	for _, item := range Articles {
		items = append(items, &feeds.Item{
			Title: item.Title,
			Link:  &feeds.Link{Href: fmt.Sprintf("https://drawbu.dev/blog/%s", item.Uri)},
			Description: item.Description,
			Author:  &feeds.Author{Name: item.Author.Name, Email: item.Author.Email},
			Created: item.Date,
		})
	}
	feed.Items = items
	return feed
}
