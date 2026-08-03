package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"gopkg.in/yaml.v3"
)

// Lede, Practice, and WorkItem.Body are Markdown rendered at load time.
type Config struct {
	BaseURL      string
	Name         string
	Description  string
	Lede         template.HTML
	Practice     template.HTML
	SelectedWork []WorkItem
	Contact      Contact
}

type WorkItem struct {
	Title string
	Body  template.HTML
	URL   string
}

type Contact struct {
	GitHub   string `yaml:"github"`
	X        string `yaml:"x"`
	LinkedIn string `yaml:"linkedin"`
	Email    string `yaml:"email"`
}

type Essay struct {
	Slug        string
	Title       string
	Date        time.Time
	Description string
	BodyHTML    template.HTML
}

// slugPattern enforces lowercase kebab-case (hello-world).
var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

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
	if strings.TrimSpace(raw.BaseURL) == "" {
		return Config{}, fmt.Errorf("%s: base_url is required", path)
	}
	if strings.TrimSpace(raw.Name) == "" {
		return Config{}, fmt.Errorf("%s: name is required", path)
	}
	cfg := Config{
		BaseURL:     strings.TrimRight(raw.BaseURL, "/"),
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

func validateSlug(slug string) error {
	if !slugPattern.MatchString(slug) {
		return fmt.Errorf("invalid slug %q (want lowercase kebab-case, e.g. hello-world)", slug)
	}
	return nil
}

func parseEssay(slug string, data []byte) (Essay, error) {
	if err := validateSlug(slug); err != nil {
		return Essay{}, fmt.Errorf("essay %s: %w", slug, err)
	}
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
	if strings.TrimSpace(fm.Title) == "" {
		return Essay{}, fmt.Errorf("essay %s: title is required", slug)
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

// Markdown stays in goldmark's default safe mode (no html.WithUnsafe):
// raw HTML and javascript: URLs are stripped before we mark the result as template.HTML.
var md = goldmark.New(goldmark.WithExtensions(extension.GFM))

func markdown(src string) template.HTML {
	var buf bytes.Buffer
	// Convert only fails when the writer does, and a bytes.Buffer does not.
	_ = md.Convert([]byte(src), &buf)
	return template.HTML(buf.String())
}
