package extractor

import (
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type EntityTags struct {
	Tags map[string]bool
}

func (e *EntityTags) Contains(tag string) bool {
	return e.Tags[tag]
}

// Fetch webpage content
func fetchURL(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the HTML content
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// Recursive function to extract text content from HTML nodes
func extractText(n *html.Node, sb *strings.Builder, excludedTags *EntityTags) {
	// Skip the node if it's an element that should be excluded
	if n.Type == html.ElementNode && excludedTags.Contains(n.Data) {
		return
	}

	// If it's a text node, clean and append the text content
	if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if len(text) > 0 && !strings.HasPrefix(text, "<iframe") && !strings.HasPrefix(text, "<img") {
			sb.WriteString(text + " ")
		}
	}

	// Traverse the child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, sb, excludedTags)
	}
}

func GetTextFromURL(url string) (string, error) {
	// Fetch the HTML document
	doc, err := fetchURL(url)
	if err != nil {
		return "", err
	}

	excludedTags := &EntityTags{
		Tags: map[string]bool{
			"script": true,
			"style":  true,
			"head":   true,
			"nav":    true,
			"footer": true,
			"img":    true,
			"iframe": true,
		},
	}

	// Extract the text content
	var sb strings.Builder
	extractText(doc, &sb, excludedTags)

	// Return the cleaned-up text
	return sb.String(), nil
}
