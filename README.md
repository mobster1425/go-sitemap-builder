# go-sitemap-builder

The Sitemap Builder is a command-line tool that generates a sitemap for a given domain. A sitemap is a map of all the pages within a specific domain, used by search engines and other tools to inform them of all the pages on the domain.

# How to Use
Clone the repository:

``` bash
Copy code
git clone https://github.com/your-username/sitemap-builder.git
Navigate to the project directory:
```

bash
Copy code
cd sitemap-builder
Build the project:

go
Copy code
go build
Run the Sitemap Builder with the desired domain and depth:

bash
Copy code
./sitemap-builder -domain="https://example.com" -depth=3
Replace https://example.com with the domain you want to generate the sitemap for, and 3 with the maximum depth you want to follow when building the sitemap.

The Sitemap Builder will generate an XML sitemap and print it to the console. You can also redirect the output to a file if needed:

bash
Copy code
./sitemap-builder -domain="https://example.com" -depth=3 > sitemap.xml
Dependencies
The Sitemap Builder uses the following third-party libraries:

encoding/xml: For generating XML output.
flag: For parsing command-line flags.
fmt: For printing messages.
io: For reading and writing data.
net/http: For sending HTTP requests and fetching HTML pages.
os: For interacting with the operating system.
strings: For string manipulation.
The Sitemap Builder also uses the go-html-link-parser/link package (created in a previous exercise) for parsing HTML pages and extracting links.

# How it Works
The Sitemap Builder starts with the given domain as the root page and recursively follows links to other pages on the same domain up to the specified depth. It uses the link parser package to parse HTML pages and extract links.

To avoid infinite loops and cycles, the Sitemap Builder keeps track of visited pages using a map. If a link points to a different domain or exceeds the maximum depth, it is skipped.

The final output is an XML sitemap in the following format:

xml
Copy code
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/</loc>
  </url>
  <url>
    <loc>https://example.com/page1</loc>
  </url>
  <!-- More URLs -->
</urlset>
# Bonus: Depth Limit
As a bonus feature, the Sitemap Builder can accept a -depth flag to limit the maximum depth of links to follow. If a page is more than the specified depth away from the root domain, it will not be included in the sitemap.





