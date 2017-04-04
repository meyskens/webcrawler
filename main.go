package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var (
	mainURL          string
	urlsAdded        = map[string]bool{}
	urlsToScan       = []string{}
	urlInfos         = []urlInfo{}
	schemeRegex      = regexp.MustCompile(`^https*://`)
	urlFragmentRegex = regexp.MustCompile(`\#.*$`)                         // to filter out #fragments
	nonHTTPLinkRegex = regexp.MustCompile(`^(mailto|tel|sms|javascript):`) // most common ones
)

// urlInfo contains the data that needs to returned per crawled url
type urlInfo struct {
	URL    string   `json:"url"`
	Assets []string `json:"assets"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No URL specified")
		return
	}
	mainURL = os.Args[1]
	if mainURL[len(mainURL)-1] != '/' { // add tailing / if not present
		mainURL += "/"
	}
	if !schemeRegex.MatchString(mainURL) {
		mainURL = "http://" + mainURL
	}

	addURLToScan(mainURL)

	for len(urlsToScan) > 0 {
		url := urlsToScan[0]
		crawlURL(url)
		urlsToScan = urlsToScan[1:]
	}
	jsonOutput, _ := json.MarshalIndent(urlInfos, "", "  ") // tabs vs spaces RIGHT HERE
	fmt.Println(string(jsonOutput))
}

// crawlURL craws a url and does the necessary actions
func crawlURL(url string) {
	assets := []string{}

	resp, err := http.Get(url)
	if err != nil {
		return // don't show errors for now
	}

	url = resp.Request.URL.String() // maybe we got redirected
	urlsAdded[url] = true           // also add redirected URL
	if !isInSameDomain(url) {
		return // we shouldn't scan this URL
	}

	toknizer := html.NewTokenizer(resp.Body)

L:
	for {
		tokenType := toknizer.Next()
		if tokenType == html.ErrorToken {
			// end of HTML page
			break L
		}

		token := toknizer.Token()
		if token.Data == "a" {
			scanLink(url, token)
		} else { // I could pass script, img, link,... but what will they add next?
			if asset := scanForAsset(url, token); asset != "" {
				assets = append(assets, asset)
			}
		}
	}

	info := urlInfo{
		URL:    url,
		Assets: assets,
	}
	urlInfos = append(urlInfos, info)
}

// scanLink looks at an anchor to list any URLs needed to be scanned
func scanLink(url string, token html.Token) {
	for _, attribute := range token.Attr {
		if attribute.Key == "href" {
			if schemeRegex.MatchString(attribute.Val) {
				// absolute URL only need to verify for subdomain change
				if isInSameDomain(attribute.Val) {
					addURLToScan(attribute.Val)
				}
			} else {
				// make the URL absolute
				url, err := makeLinkAbsolute(url, attribute.Val)
				if err == nil {
					addURLToScan(url)
				}
			}
		}
	}
}

// addURLToScan adds an URL to urlsToScan for main() to pick up
func addURLToScan(url string) {
	url = urlFragmentRegex.ReplaceAllString(url, "") // filter out fragments as they are ignored by the server

	// get a clean url to not have duplicates
	cleanURL := strings.Trim(url, "/")                    // remove trailing /
	cleanURL = schemeRegex.ReplaceAllString(cleanURL, "") // remove scheme

	if _, exists := urlsAdded[cleanURL]; !exists {
		urlsToScan = append(urlsToScan, url)
		urlsAdded[cleanURL] = true
	}
}

// scanForAsset scans elements for resources that are loaded
func scanForAsset(url string, token html.Token) string {
	assetURL := ""
	rel := ""

	for _, attribute := range token.Attr {
		if attribute.Key == "src" || (token.Data == "link" && attribute.Key == "href") {
			assetURL = attribute.Val
		}
		if attribute.Key == "rel" {
			rel = attribute.Val
		}
	}

	if rel != "" && rel != "stylesheet" && rel != "icon" { // <link> elements are tricky...
		return ""
	}
	if schemeRegex.MatchString(assetURL) { // absolute
		return assetURL
	}
	if assetURL, err := makeLinkAbsolute(url, assetURL); err == nil { // make the URL absolute
		return assetURL
	}

	return ""
}

// isInSameDomain checks if the given URL is in the same domain as the one we're crawling
func isInSameDomain(url string) bool {
	// remove the schema and tailing /
	mainDomain := strings.Trim(schemeRegex.ReplaceAllString(mainURL, ""), "/")
	url = string(schemeRegex.ReplaceAll([]byte(url), []byte("")))

	urlParts := strings.Split(url, "/")
	return urlParts[0] == mainDomain
}

// makeLinkAbsolute turns relative URLs into absolute ones
func makeLinkAbsolute(scanURL, link string) (string, error) {
	if link == "" {
		return "", fmt.Errorf("Link is empty")
	}
	if nonHTTPLinkRegex.MatchString(link) {
		return "", fmt.Errorf("Link is not to HTTP")
	}
	switch {
	case len(link) > 2 && link[:2] == "//":
		return schemeRegex.FindString(scanURL) + link[2:], nil
	case link[0] == '/':
		return mainURL + link[1:], nil
	case len(link) > 2 && link[:2] == "./":
		return getDirectory(scanURL, 0) + link[2:], nil
	case len(link) > 2 && link[:2] == "~/":
		return mainURL + link[2:], nil
	case len(link) > 3 && link[:3] == "../":
		return getDirectory(scanURL, 1) + link[3:], nil
	}

	// no prefix
	return getDirectory(scanURL, 0) + link, nil
}

// getDirectory gives the directory a file is located in, or go back a certain offset
func getDirectory(url string, offset int) string {
	urlParts := strings.Split(url, "/")
	newURL := strings.Join(urlParts[:len(urlParts)-(1+offset)], "/")
	if newURL[len(newURL)-1] != '/' {
		newURL += "/" // add tailing slash if needed
	}
	return newURL
}
