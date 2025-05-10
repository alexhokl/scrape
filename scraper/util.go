package scraper

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func containsCssClass(node *html.Node, class string) bool {
	if node.Type == html.ElementNode {
		for _, attr := range node.Attr {
			if attr.Key == "class" {
				for c := range strings.SplitSeq(attr.Val, " ") {
					if c == class {
						return true
					}
				}
			}
		}
	}
	return false
}

func printCssClasses(node *html.Node) {
	if node.Type == html.ElementNode {
		for _, attr := range node.Attr {
			if attr.Key == "class" {
				fmt.Println(attr.Val)
			}
		}
	}
}

func printAttributes(node *html.Node) {
	if node.Type == html.ElementNode {
		for _, attr := range node.Attr {
			fmt.Printf("%s: %s\n", attr.Key, attr.Val)
		}
	}
}
