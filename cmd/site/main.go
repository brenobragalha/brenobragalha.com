package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

const homeEssays = 5

func main() {
	if err := build(".", "public"); err != nil {
		fmt.Fprintf(os.Stderr, "site: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("built site → public")
}

func build(root, out string) error {
	cfg, err := loadConfig(filepath.Join(root, "content", "site.yaml"))
	if err != nil {
		return err
	}
	essays, err := loadEssays(filepath.Join(root, "content", "writing"))
	if err != nil {
		return err
	}
	templates, err := parseTemplates(filepath.Join(root, "templates"))
	if err != nil {
		return err
	}

	tmp, err := os.MkdirTemp(filepath.Dir(out), ".site-*")
	if err != nil {
		return err
	}
	ok := false
	defer func() {
		if !ok {
			_ = os.RemoveAll(tmp)
		}
	}()

	if err := renderSite(root, tmp, cfg, essays, templates); err != nil {
		return err
	}
	if err := publish(tmp, out); err != nil {
		return err
	}
	ok = true
	return nil
}

// publish replaces out with tmp via a sibling swap (out → out.old → remove),
// restoring out.old if the final rename fails.
func publish(tmp, out string) error {
	old := out + ".old"
	_ = os.RemoveAll(old)
	if _, err := os.Stat(out); err == nil {
		if err := os.Rename(out, old); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.Rename(tmp, out); err != nil {
		if _, statErr := os.Stat(old); statErr == nil {
			_ = os.Rename(old, out)
		}
		return err
	}
	_ = os.RemoveAll(old)
	return nil
}

func renderSite(root, out string, cfg Config, essays []Essay, templates map[string]*template.Template) error {
	// Static first, so generated pages win over any colliding path.
	if err := os.CopyFS(out, os.DirFS(filepath.Join(root, "static"))); err != nil {
		return fmt.Errorf("copy static: %w", err)
	}

	recent := essays
	if len(recent) > homeEssays {
		recent = recent[:homeEssays]
	}
	pages := []struct {
		name string
		path string
		data page
	}{
		{"home", "index.html", page{
			Title:       cfg.Name,
			Description: cfg.Description,
			Canonical:   cfg.BaseURL + "/",
			Config:      cfg,
			Essays:      recent,
		}},
		{"archive", "writing/index.html", page{
			Title:       cfg.writingTitle(),
			Description: cfg.writingDescription(),
			Canonical:   cfg.writingURL(),
			Config:      cfg,
			Essays:      essays,
		}},
		{"notfound", "404.html", page{
			Title:       "Not found · " + cfg.Name,
			Description: "Page not found.",
			Config:      cfg,
		}},
	}
	for _, p := range pages {
		if err := writePage(templates[p.name], p.name, filepath.Join(out, p.path), p.data); err != nil {
			return err
		}
	}
	for i := range essays {
		essay := &essays[i]
		path := filepath.Join(out, "writing", essay.Slug, "index.html")
		if err := writePage(templates["essay"], "essay", path, page{
			Title:       essay.Title + " · " + cfg.Name,
			Description: essay.Description,
			Canonical:   cfg.essayURL(essay.Slug),
			Config:      cfg,
			Essay:       essay,
		}); err != nil {
			return err
		}
	}

	if err := writeFeed(filepath.Join(out, "feed.xml"), cfg, essays); err != nil {
		return err
	}
	return writeSitemap(filepath.Join(out, "sitemap.xml"), cfg, essays)
}

type page struct {
	Title       string
	Description string
	Canonical   string
	Config      Config
	Essays      []Essay
	Essay       *Essay
}

func parseTemplates(dir string) (map[string]*template.Template, error) {
	files := map[string]string{
		"home":     "home.html",
		"archive":  "archive.html",
		"essay":    "essay.html",
		"notfound": "404.html",
	}
	templates := make(map[string]*template.Template, len(files))
	for name, file := range files {
		t, err := template.New(name).ParseFiles(
			filepath.Join(dir, "base.html"),
			filepath.Join(dir, file),
		)
		if err != nil {
			return nil, fmt.Errorf("parse templates: %w", err)
		}
		templates[name] = t
	}
	return templates, nil
}

func writePage(t *template.Template, name, path string, data page) error {
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, name, data); err != nil {
		return fmt.Errorf("render %s: %w", name, err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}
