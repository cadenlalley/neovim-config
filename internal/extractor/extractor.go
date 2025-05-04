package extractor

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

var (
	ErrRequestBlocked = errors.New("request blocked")
)

var excludedTags = []string{
	"script",
	"style",
	"head",
	"nav",
	"footer",
	"img",
	"iframe",
	"svg",
	"form",
	"button",
	"input",
	"select",
	"textarea",
	"link",
	"picture",
	"video",
	"audio",
	"source",
	"canvas",
	"noscript",
	"figure",
}

// Fetch webpage content
func fetchURL(inputURL string) (*html.Node, error) {
	// Get domain of URL
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return nil, err
	}
	host := parsedURL.Hostname()

	req, err := http.NewRequest("GET", inputURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	req.Header.Set("Host", host)
	req.Header.Set("Accept", "text/html")
	req.Header.Set("Accept-Language", "en-US")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if resp.StatusCode == 401 || resp.StatusCode == 403 {
			return nil, ErrRequestBlocked
		}
		return nil, fmt.Errorf("unexpected response from source domain '%d'", resp.StatusCode)
	}

	// Parse the HTML content
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func GetTextFromURL(url string) (string, error) {
	// Fetch the HTML document
	doc, err := fetchURL(url)
	if err != nil {
		return "", err
	}

	gqDoc := goquery.NewDocumentFromNode(doc)
	gqDoc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		s.SetAttr("href", "")
	})

	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(
				commonmark.WithStrongDelimiter("__"),
			),
		),
	)

	for _, tag := range excludedTags {
		conv.Register.TagType(tag, converter.TagTypeRemove, 0)
	}

	md, err := conv.ConvertNode(gqDoc.Get(0))
	if err != nil {
		return "", err
	}

	return string(md), nil
}
