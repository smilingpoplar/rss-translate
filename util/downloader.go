package util

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
)

type Downloader struct {
	client *http.Client
}

func NewDownloader() *Downloader {
	return &Downloader{
		client: &http.Client{},
	}
}

func (d *Downloader) SetProxy(addr string) error {
	return SetProxy(addr, d.client)
}

func (d *Downloader) GetURL(url string) ([]byte, error) {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 500 * time.Millisecond
	b.Multiplier = 1.5
	b.RandomizationFactor = 0.5
	b.MaxElapsedTime = 10 * time.Second

	return backoff.RetryWithData(func() ([]byte, error) {
		return d.getURL(url)
	}, b)
}

func (d *Downloader) getURL(url string) ([]byte, error) {
	resp, err := d.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error http status code: %d, url: %s", resp.StatusCode, url)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}
	return data, nil
}
