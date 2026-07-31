package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

type rss struct {
	Version string  `xml:"version,attr"`
	Channel channel `xml:"channel"`
}

type channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Language    string `xml:"language"`
	Items       []item `xml:"item"`
}

type item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

// writeFeed writes an RSS 2.0 feed with absolute links and full HTML content.
func writeFeed(path string, cfg Config, essays []Essay) error {
	items := make([]item, 0, len(essays))
	for _, e := range essays {
		link := cfg.essayURL(e.Slug)
		items = append(items, item{
			Title:       e.Title,
			Link:        link,
			GUID:        link,
			PubDate:     e.Date.Format(time.RFC1123Z),
			Description: string(e.BodyHTML),
		})
	}
	return writeXML(path, rss{
		Version: "2.0",
		Channel: channel{
			Title:       cfg.writingTitle(),
			Link:        cfg.writingURL(),
			Description: cfg.writingDescription(),
			Language:    "en",
			Items:       items,
		},
	})
}

type urlset struct {
	Xmlns string `xml:"xmlns,attr"`
	URLs  []loc  `xml:"url"`
}

type loc struct {
	Loc string `xml:"loc"`
}

func writeSitemap(path string, cfg Config, essays []Essay) error {
	urls := []loc{
		{Loc: cfg.BaseURL + "/"},
		{Loc: cfg.writingURL()},
	}
	for _, e := range essays {
		urls = append(urls, loc{Loc: cfg.essayURL(e.Slug)})
	}
	return writeXML(path, urlset{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	})
}

func writeXML(path string, doc any) error {
	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := enc.Encode(doc); err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}
	buf.WriteByte('\n')
	return os.WriteFile(path, buf.Bytes(), 0o644)
}
