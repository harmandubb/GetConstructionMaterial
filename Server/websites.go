package server

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
)

func FindEmailsOnPage(page string) ([]string, error) {
	c := colly.NewCollector()

	var mailtoLinks []string

	// OnHTML callback for mailto links
	c.OnHTML("a[href^='mailto:']:not(header a[href^='mailto:'])", func(e *colly.HTMLElement) {
		mailtoLink := e.Attr("href")
		mailtoLinks = append(mailtoLinks, mailtoLink)
	})

	// OnHTML callback for specific elements
	c.OnHTML("body", func(e *colly.HTMLElement) {
		// Load the HTML of the current element
		node, err := htmlquery.Parse(strings.NewReader(string(e.Response.Body)))
		if err != nil {
			log.Fatal(err)
		}

		// XPath query to select nodes that directly contain 'email'
		expr := "//p/text() | //div/text() | //span/text() | //a/text()"
		nodes, err := htmlquery.QueryAll(node, expr)
		if err != nil {
			log.Fatal(err)
		}

		// Iterate over nodes and check if they contain 'email'
		for _, n := range nodes {
			if strings.Contains(strings.ToLower(n.Data), "email") {
				mailtoLinks = append(mailtoLinks, extractEmailFromString(n.Data))
				// fmt.Printf("Text: %s\n", n.Data)
			}
		}
	})

	// Handle visiting the page
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping
	err := c.Visit(page) // Replace with the URL you want to scrape
	if err != nil {
		return mailtoLinks, err
	}

	return mailtoLinks, nil

}

func extractEmailFromString(str string) string {
	// Regular expression for matching an email address
	re := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

	// FindString will return the first match found
	email := re.FindString(str)

	return email
}

// isValidWebsite checks if the given string is a valid website URL
func isValidWebsite(u string) bool {
	parsedURL, err := url.ParseRequestURI(u)
	if err != nil {
		return false
	}

	// Check for valid scheme and host
	return parsedURL.Scheme != "" && parsedURL.Host != ""
}

func FindContactURLOnPage(page string) ([]string, error) {
	c := colly.NewCollector()

	// Define a slice to store contact links
	var contactLinks []string

	// OnHTML callback for links
	c.OnHTML("a", func(e *colly.HTMLElement) {
		linkText := e.Text
		if strings.Contains(strings.ToLower(linkText), "contact us") {
			contactLink := e.Attr("href")

			//check if it is an absolute website link or a relative link
			if !isValidWebsite(contactLink) {
				contactLink = fmt.Sprintf("%s%s", page, contactLink)
			}
			contactLinks = append(contactLinks, contactLink)
		}
	})

	// Handle visiting the page
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping
	err := c.Visit(page) // Replace with the URL you want to scrape
	if err != nil {
		return contactLinks, err
	}

	return contactLinks, nil
}

type FormDetails struct {
	InputFields []InputFieldDetails
	SubmitID    string
}

type InputFieldDetails struct {
	name string
	id   string
}

func FindContactFormInputs(page string) (FormDetails, error) {
	c := colly.NewCollector()

	formDetails := FormDetails{
		SubmitID: "",
	}

	c.OnHTML("main form", func(e *colly.HTMLElement) {

		e.ForEach("input", func(_ int, el *colly.HTMLElement) {
			if el.Attr("type") == "submit" {
				formDetails.SubmitID = el.Attr("id")

			} else {
				formDetails.InputFields = append(formDetails.InputFields,
					InputFieldDetails{
						name: el.Attr("name"),
						id:   el.Attr("id"),
					})
			}
		})

		e.ForEach("textarea", func(_ int, el *colly.HTMLElement) {

			formDetails.InputFields = append(formDetails.InputFields,
				InputFieldDetails{
					name: el.Attr("name"),
					id:   el.Attr("id"),
				})

		})

		if formDetails.SubmitID == "" {
			e.ForEach("button", func(_ int, el *colly.HTMLElement) {
				if el.Attr("type") == "submit" {
					formDetails.SubmitID = el.Attr("id")
				}
			})

		}

	})

	// Handle visiting the page
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping
	err := c.Visit(page) // Replace with the URL you want to scrape
	if err != nil {
		log.Fatal(err)
	}

	return formDetails, nil

}
