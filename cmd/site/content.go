package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"gopkg.in/yaml.v3"
)

// Config is the site-wide content from content/site.yaml.
// Lede, Practice, and WorkItem.Body are Markdown rendered at load time.
type Config struct {
	BaseURL      string        `yaml:"base_url"`
	Name         string        `yaml:"name"`
	Description  string        `yaml:"description"`
	Lede         template.HTML `yaml:"lede"`
	Practice     template.HTML `yaml:"practice"`
	SelectedWork []WorkItem    `yaml:"selected_work"`
	Contact      Contact       `yaml:"contact"`
}

type WorkItem struct {
	Title string        `yaml:"title"`
	Body  template.HTML `yaml:"body"`
	URL   string        `yaml:"url"`
}

type Contact struct {
	GitHub   string `yaml:"github"`
	X        string `yaml:"x"`
	LinkedIn string `yaml:"linkedin"`
	Email    string `yaml:"email"`
}

// Essay is a writing piece loaded from content/writing/<slug>.md.
type Essay struct {
	Slug        string
	Title       string
	Date        time.Time
	Description string
	BodyHTML    template.HTML
}

func (c Config) writingURL() string {
	return c.BaseURL + "/writing/"
}

func (c Config) essayURL(slug string) string {
	return c.BaseURL + "/writing/" + slug + "/"
}

func (c Config) writingTitle() string {
	return "Writing · " + c.Name
}

func (c Config) writingDescription() string {
	return "Writing by " + c.Name + "."
}

func loadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var raw struct {
		BaseURL      string `yaml:"base_url"`
		Name         string `yaml:"name"`
		Description  string `yaml:"description"`
		Lede         string `yaml:"lede"`
		Practice     string `yaml:"practice"`
		SelectedWork []struct {
			Title string `yaml:"title"`
			Body  string `yaml:"body"`
			URL   string `yaml:"url"`
		} `yaml:"selected_work"`
		Contact Contact `yaml:"contact"`
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return Config{}, fmt.Errorf("parse %s: %w", path, err)
	}
	cfg := Config{
		BaseURL:     raw.BaseURL,
		Name:        raw.Name,
		Description: raw.Description,
		Lede:        markdown(raw.Lede),
		Practice:    markdown(raw.Practice),
		Contact:     raw.Contact,
	}
	for _, w := range raw.SelectedWork {
		cfg.SelectedWork = append(cfg.SelectedWork, WorkItem{
			Title: w.Title,
			Body:  markdown(w.Body),
			URL:   w.URL,
		})
	}
	return cfg, nil
}

// loadEssays reads content/writing/*.md, newest first.
func loadEssays(dir string) ([]Essay, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return nil, err
	}
	essays := make([]Essay, 0, len(paths))
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		essay, err := parseEssay(strings.TrimSuffix(filepath.Base(path), ".md"), data)
		if err != nil {
			return nil, err
		}
		essays = append(essays, essay)
	}
	// Glob returns paths in lexical order, so a stable sort breaks ties by slug.
	slices.SortStableFunc(essays, func(a, b Essay) int {
		return b.Date.Compare(a.Date)
	})
	return essays, nil
}

func parseEssay(slug string, data []byte) (Essay, error) {
	front, body, ok := bytes.Cut(data, []byte("\n---\n"))
	if !ok || !bytes.HasPrefix(front, []byte("---\n")) {
		return Essay{}, fmt.Errorf("essay %s: missing --- front matter", slug)
	}
	var fm struct {
		Title       string `yaml:"title"`
		Date        string `yaml:"date"`
		Description string `yaml:"description"`
	}
	if err := yaml.Unmarshal(front, &fm); err != nil {
		return Essay{}, fmt.Errorf("essay %s: parse front matter: %w", slug, err)
	}
	date, err := time.Parse("2006-01-02", fm.Date)
	if err != nil {
		return Essay{}, fmt.Errorf("essay %s: invalid date %q (want YYYY-MM-DD)", slug, fm.Date)
	}
	description := fm.Description
	if description == "" {
		description = fm.Title
	}
	return Essay{
		Slug:        slug,
		Title:       fm.Title,
		Date:        date,
		Description: description,
		BodyHTML:    markdown(string(bytes.TrimSpace(body))),
	}, nil
}

var md = goldmark.New(goldmark.WithExtensions(extension.GFM))

func markdown(src string) template.HTML {
	var buf bytes.Buffer
	// Convert only fails when the writer does, and a bytes.Buffer does not.
	_ = md.Convert([]byte(src), &buf)
	return template.HTML(buf.String())
}
