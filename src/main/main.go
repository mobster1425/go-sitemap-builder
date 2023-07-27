package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"

	//"os"
	"strings"

	//"feyin/sitemap-builder/link"
	"golang.org/x/net/html"
)

//to run
//go run main.go -domain="https://example.com" -depth=3

// Sitemap represents the XML sitemap structure.
type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []URL    `xml:"url"`
}

// URL represents the URL structure in the sitemap.
type URL struct {
	Loc string `xml:"loc"`
}

func main() {
	// Parse command line flags to get the website domain and max depth
	domain := flag.String("domain", "", "Website domain to build sitemap for")
	maxDepth := flag.Int("depth", 3, "Maximum depth to follow links")
	flag.Parse()

	if *domain == "" {
		fmt.Println("Please provide the website domain using -domain flag")
		return
	}

	// Start building the sitemap from the root URL
	sitemap, err := buildSitemap(*domain, *maxDepth)
	if err != nil {
		fmt.Println("Error building sitemap:", err)
		return
	}

	// Output the sitemap in XML format
	xmlData, err := xml.MarshalIndent(sitemap, "", "  ")
	if err != nil {
		fmt.Println("Error encoding sitemap to XML:", err)
		return
	}

	// Print the XML sitemap to stdout
	fmt.Println(xml.Header + string(xmlData))
}

func buildSitemap(domain string, maxDepth int) (*Sitemap, error) {
	//creating an instance of the siemap struct using the addr so any changes made in this function reflects
	//in the original struct
	sitemap := &Sitemap{}
	visited := make(map[string]bool)
	//  a visited map to keep track of visited URLs, and a stack to store URLs to be visited.
	stack := []string{domain}

	for len(stack) > 0 {
		//get the top of stack
		//url rep the current url been processed
		url := stack[len(stack)-1]
		//This pops the last URL from the stack, removing it from consideration for future processing. I
		stack = stack[:len(stack)-1]

		//This checks if the current url has already been visited (by checking the visited map) or if
		// its depth (number of slashes in the URL minus 2) exceeds
		// the maxDepth value. If either of these conditions is true, the URL is skipped, and the loop continues to the next iteration.
		if visited[url] || len(strings.Split(url, "/"))-2 > maxDepth {
			continue
		}

		// Fetch the HTML content of the URL
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("error fetching URL %s: %v", url, err)
		}
		//close the http get , this will exec after the surroounding function
		//This ensures that the HTTP response body is always closed after the function execution,
		// regardless of whether there was an error or not.
		//Closing the response body is essential to release resources associated with the HTTP request and to avoid leaking resources.
		defer resp.Body.Close()

		// Parse the HTML content to find links
		links, err := findLinks(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error parsing HTML for URL %s: %v", url, err)
		}

		// Add the current URL to the sitemap
		sitemap.URLs = append(sitemap.URLs, URL{Loc: url})
		visited[url] = true

		// Add new links to the stack
		for _, link := range links {
			//checking if the domain name and the link been passed starts with the same prefix
			//whcih means we are not accessing an outside domain, we are only dealingwith an inside the website domain
			if isSameDomain(domain, link) {
				//then append to the stack
				stack = append(stack, link)
			}
		}
	}

	return sitemap, nil
}

/*
e findLinks function takes an io.Reader (HTML content) as input, parses it using the

	html.Parse function, and calls the linkNodes and buildLink functions to find and extract the links from the parsed HTML.
*/
func findLinks(r io.Reader) ([]string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := linkNodes(doc)
	var links []string
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}
	return links, nil
}

// The buildLink function extracts the URL
// from the <a href="..."> tag in the HTML node by iterating through its attributes and finding the href attribute.
func buildLink(n *html.Node) string {
	var ret string
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret = attr.Val
			break
		}
	}
	return ret
}

// The linkNodes function takes an HTML node as input
// and returns a slice of all <a> tags within that node. It uses recursion to traverse the HTML tree and collect all the <a> tags.
func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {

		/*
						In the code ret = append(ret, linkNodes(c)...), ... is called the "ellipsis" or "variadic" operator in Go.
						 It is used when calling a variadic function with multiple arguments of the same type.

						In this case, linkNodes(c) returns a slice of *html.Node, and ret is also a slice of *html.Node.
						 When using append to add the elements of one slice to another, you need to use the ... operator to "unpack" the elements
						 from the source slice and add them individually to the destination slice.

						Without the ... operator, append would treat linkNodes(c) as a single element of type []*html.Node, which is not the desired behavior.

			For example, let's say linkNodes(c) returns the slice []*html.Node{node1, node2, node3}, and ret is currently []*html.Node{node4}.

			After executing ret = append(ret, linkNodes(c)...), ret will become []*html.Node{node4, node1, node2, node3}.
			The elements of sourceSlice (node1, node2, and node3) are added individually to the ret slice.


		*/
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}

/*
The isSameDomain function checks if a given link is from the same domain as the specified domain.

	It uses the strings.HasPrefix function to determine if the link starts with the domain name.
	If it does, it is considered to be from the same domain.
*/
func isSameDomain(domain, link string) bool {
	return strings.HasPrefix(link, domain)
}
