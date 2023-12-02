package translate

import (
	"fmt"
	"os"

	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
)

func Main() int {
	fp := gofeed.NewParser()
	from, err := fp.ParseURL("https://www.economist.com/science-and-technology/rss.xml")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing feed:", err)
		return 1
	}

	to := transform(from)
	out, err := to.ToRss()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error transforming feed:", err)
		return 1
	}
	fmt.Println(out)

	return 0
}

func transform(from *gofeed.Feed) *feeds.Feed {
	to := &feeds.Feed{
		Title:       from.Title,
		Link:        &feeds.Link{Href: from.Link},
		Description: from.Description,
	}
	if from.PublishedParsed != nil {
		to.Created = *from.PublishedParsed
	}

	for _, item := range from.Items {
		toItem := &feeds.Item{
			Title:       item.Title,
			Link:        &feeds.Link{Href: item.Link},
			Description: item.Description,
			Content:     item.Content,
			Id:          item.GUID,
		}
		if item.PublishedParsed != nil {
			toItem.Created = *item.PublishedParsed
		}
		to.Items = append(to.Items, toItem)
	}
	return to
}
