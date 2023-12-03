package translate

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"github.com/smilingpoplar/rss-translate/util"
)

func Main() int {
	config, err := util.GetConfig("config.yaml")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err = process(config); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

func process(config *util.Config) error {
	parser := gofeed.NewParser()
	for feedName, feedConfig := range config.Feeds {
		data, err := util.GetURL(feedConfig.URL)
		if err != nil {
			return err
		}
		hash, err := util.MD5(data)
		if err != nil {
			return err
		}
		fmt.Println(hash)

		from, err := parser.Parse(bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("error parsing feed %s: %w", feedName, err)
		}

		to := transformFeed(from, feedConfig.Max)
		if err := writeFeed(to, config.Output.Dir, feedName); err != nil {
			return err
		}
	}
	return nil
}

func transformFeed(from *gofeed.Feed, limit int) *feeds.Feed {
	to := &feeds.Feed{
		Title:       from.Title,
		Link:        &feeds.Link{Href: from.Link},
		Description: from.Description,
	}
	if from.PublishedParsed != nil {
		to.Created = *from.PublishedParsed
	}

	for i, item := range from.Items {
		if i >= limit {
			break
		}

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

func writeFeed(feed *feeds.Feed, rssDir string, feedName string) error {
	if err := os.MkdirAll(rssDir, os.ModePerm); err != nil {
		return fmt.Errorf("error mkdir %s: %w", rssDir, err)
	}

	fp := path.Join(rssDir, feedName+".xml")
	f, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", fp, err)
	}
	defer f.Close()

	if err := feed.WriteAtom(f); err != nil {
		return fmt.Errorf("error writing feed %s: %w", feedName, err)
	}
	return nil
}
