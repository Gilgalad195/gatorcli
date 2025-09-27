package webconn

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/Gilgalad195/gatorcli/internal/rss"
)

var client = &http.Client{
	Timeout: time.Second * 10,
}

func FetchFeed(ctx context.Context, feedURL string) (*rss.RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response rss.RSSFeed
	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	unescapeHelper(&response)

	return &response, nil
}

func unescapeHelper(response *rss.RSSFeed) {
	response.Channel.Title = html.UnescapeString(response.Channel.Title)
	response.Channel.Description = html.UnescapeString(response.Channel.Description)
	for i, item := range response.Channel.Item {
		response.Channel.Item[i].Title = html.UnescapeString(item.Title)
		response.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}
}
