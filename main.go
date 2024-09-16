package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/chromedp/chromedp"
)

const timeout = 60 * time.Second

func extractDomain(urlstr string) string {
	parsedURL, err := url.Parse(urlstr)
	if err != nil {
		log.Fatal("Failed to parse URL:", err)
	}
	return parsedURL.Hostname()
}

func normalizeAndValidateURL(urlString string) (string, bool) {
	if !strings.HasPrefix(urlString, "http://") && !strings.HasPrefix(urlString, "https://") {
		urlString = "https://" + urlString
	}
	_, err := url.ParseRequestURI(urlString)
	return urlString, err == nil
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Error: Valid URL is required")
		fmt.Println("Usage: murl <url>")
		fmt.Println("Example: murl https://example.com")
		os.Exit(1)
	}

	urlArg, validUrl := normalizeAndValidateURL(os.Args[1])

	if !validUrl {
		fmt.Println("Error: Invalid URL")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	allocatorContext, cancel := chromedp.NewRemoteAllocator(ctx, "ws://192.168.0.50:8000", chromedp.NoModifyURL)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocatorContext)
	defer cancel()

	var bodyHTML string
	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1280, 720),
		chromedp.Navigate(urlArg),
		chromedp.OuterHTML("body", &bodyHTML),
	)

	if err != nil {
		log.Fatal("Failed to fetch HTML:", err)
	}

	domain := extractDomain(urlArg)
	converter := md.NewConverter(domain, true, nil)
	markdown, err := converter.ConvertString(bodyHTML)

	if err != nil {
		log.Fatalf("Failed to convert HTML to Markdown: %v", err)
	}

	fmt.Println(markdown)
}
