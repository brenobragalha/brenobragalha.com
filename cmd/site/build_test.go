package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuild(t *testing.T) {
	out := filepath.Join(t.TempDir(), "public")
	if err := build("../..", out); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"index.html", "writing/index.html", "404.html",
		"feed.xml", "sitemap.xml", "robots.txt", "css/site.css",
	} {
		if _, err := os.Stat(filepath.Join(out, name)); err != nil {
			t.Errorf("missing %s: %v", name, err)
		}
	}

	home, err := os.ReadFile(filepath.Join(out, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(home), "Breno Bragalha") {
		t.Error("home page is missing the site name")
	}
	if !strings.Contains(string(home), `rel="alternate" type="application/rss+xml"`) {
		t.Error("home page is missing the RSS alternate link")
	}

	notFound, err := os.ReadFile(filepath.Join(out, "404.html"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(notFound), `rel="canonical"`) {
		t.Error("404 page should omit the canonical link")
	}
}

func TestParseEssay(t *testing.T) {
	essay, err := parseEssay("hello-world", []byte("---\ntitle: Hello\ndate: 2026-07-29\n---\n\nBody **text**.\n"))
	if err != nil {
		t.Fatal(err)
	}
	if essay.Title != "Hello" || essay.Date.Format("2006-01-02") != "2026-07-29" {
		t.Errorf("unexpected front matter: %+v", essay)
	}
	if essay.Description != "Hello" {
		t.Errorf("description should fall back to the title, got %q", essay.Description)
	}
	if !strings.Contains(string(essay.BodyHTML), "<strong>text</strong>") {
		t.Errorf("body was not rendered as Markdown: %q", essay.BodyHTML)
	}
}

func TestParseEssayRejectsBadSlug(t *testing.T) {
	body := []byte("---\ntitle: Hello\ndate: 2026-07-29\n---\n\nBody.\n")
	for _, slug := range []string{"", ".", "..", "Hello"} {
		if _, err := parseEssay(slug, body); err == nil {
			t.Errorf("expected reject for slug %q", slug)
		}
	}
}

func TestParseEssayRequiresTitle(t *testing.T) {
	_, err := parseEssay("hello-world", []byte("---\ntitle: \"\"\ndate: 2026-07-29\n---\n\nBody.\n"))
	if err == nil {
		t.Fatal("expected empty title to be rejected")
	}
}

func TestLoadConfigValidates(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "site.yaml")
	if err := os.WriteFile(path, []byte("base_url: https://example.com/\nname: Test\ndescription: d\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := loadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BaseURL != "https://example.com" {
		t.Errorf("expected trailing slash trimmed, got %q", cfg.BaseURL)
	}

	if err := os.WriteFile(path, []byte("base_url: \"\"\nname: Test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := loadConfig(path); err == nil {
		t.Fatal("expected empty base_url to be rejected")
	}

	if err := os.WriteFile(path, []byte("base_url: https://example.com\nname: \"\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := loadConfig(path); err == nil {
		t.Fatal("expected empty name to be rejected")
	}
}
