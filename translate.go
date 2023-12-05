package translate

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"github.com/smilingpoplar/rss-translate/util"
	"github.com/smilingpoplar/translate/translator/google"
)

func Main() int {
	config, err := util.GetConfig("config.yaml")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fp := filepath.Join(config.Output.Dir, "hash.json")
	hashes, err := util.KVStore(fp)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	if err = process(config, hashes); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err = hashes.Save(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err = writeDesc(config); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func process(config *util.Config, hashes *util.Store) error {
	downloader := util.NewDownloader()
	downloader.SetProxy(config.Proxy)
	parser := gofeed.NewParser()

	for feedName, feedConfig := range config.Feeds {
		data, err := downloader.GetURL(feedConfig.URL)
		if err != nil {
			return err
		}
		hash, err := util.MD5(data)
		if err != nil {
			return err
		}
		if oldHash, ok := hashes.Get(feedName); ok && hash == oldHash {
			fmt.Println("no change, skip feed:", feedName)
			continue
		}
		fmt.Println("processing feed:", feedName)

		from, err := parser.Parse(bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("error parsing feed %s: %w", feedName, err)
		}

		to := transformFeed(from, feedConfig.Max)
		if to, err = translateFeed(to, config.ToLang, config.Proxy); err != nil {
			return err
		}
		if err := writeFeed(to, config.Output.Dir, feedName); err != nil {
			return err
		}
		hashes.Set(feedName, hash)
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

func translateFeed(feed *feeds.Feed, toLang string, proxy string) (*feeds.Feed, error) {
	trans, err := google.New(google.WithProxy(proxy))
	if err != nil {
		return nil, fmt.Errorf("error creating translator: %w", err)
	}

	texts := make([]string, 0, len(feed.Items)*3)
	for _, item := range feed.Items {
		texts = append(texts, item.Title, item.Description, item.Content)
	}

	texts, err = trans.Translate(texts, toLang)
	if err != nil {
		return nil, fmt.Errorf("error translating feed %s: %w", feed.Link.Href, err)
	}

	for i, item := range feed.Items {
		item.Title = texts[i*3]
		item.Description = texts[i*3+1]
		item.Content = texts[i*3+2]
	}
	return feed, nil
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

func writeDesc(config *util.Config) error {
	names := make([]string, 0, len(config.Feeds))
	for name := range config.Feeds {
		names = append(names, name)
	}
	sort.Strings(names)

	fp := filepath.Join(config.Output.Dir, "rss.md")
	f, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", fp, err)
	}
	defer f.Close()

	fmt.Fprintln(f, "# RSS List")
	for _, name := range names {
		p := path.Join(config.Output.Dir, name+".xml")
		fmt.Fprintf(f, "- [%s](%s)\n", name, getPath(p, config.Output.URL))
	}
	fmt.Fprintln(f)
	return nil
}

func getPath(p string, prefix string) string {
	if u, err := url.Parse(prefix); err != nil {
		return p
	} else {
		u.Path = path.Join(u.Path, p)
		return u.String()
	}
}
