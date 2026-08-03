package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuild(t *testing.T) {
	root := fixtureRoot(t)
	out := filepath.Join(t.TempDir(), "public")
	if err := build(root, out); err != nil {
		t.Fatal(err)
	}
	// Rebuild into the same out to exercise the sibling-swap publish path.
	if err := build(root, out); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{
		"index.html", "writing/index.html", "writing/hello-world/index.html",
		"404.html", "feed.xml", "sitemap.xml", "robots.txt", "css/site.css",
	} {
		if _, err := os.Stat(filepath.Join(out, name)); err != nil {
			t.Errorf("missing %s: %v", name, err)
		}
	}

	essay, err := os.ReadFile(filepath.Join(out, "writing/hello-world/index.html"))
	if err != nil {
		t.Fatal(err)
	}
	essayStr := string(essay)
	if !strings.Contains(essayStr, "<strong>emphasis</strong>") {
		t.Error("essay page is missing rendered Markdown")
	}
	if !strings.Contains(essayStr, `rel="canonical" href="https://example.com/writing/hello-world/"`) {
		t.Error("essay page is missing the canonical link")
	}

	feed, err := os.ReadFile(filepath.Join(out, "feed.xml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(feed), "<description>A short fixture essay.</description>") {
		t.Error("feed should use the essay description, not the full HTML body")
	}
}

// fixtureRoot builds an isolated site tree: content + stub static/ from
// testdata/site, and the real templates/ (so template edits stay covered).
func fixtureRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.CopyFS(root, os.DirFS(filepath.Join("testdata", "site"))); err != nil {
		t.Fatalf("copy fixture: %v", err)
	}
	templatesDir := filepath.Join("..", "..", "templates")
	if err := os.CopyFS(filepath.Join(root, "templates"), os.DirFS(templatesDir)); err != nil {
		t.Fatalf("copy templates: %v", err)
	}
	return root
}
