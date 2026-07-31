package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBuild builds the real site into a temporary directory.
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
