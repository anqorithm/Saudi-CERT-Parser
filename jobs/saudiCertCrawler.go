package jobs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type Record struct {
	ID            int               `json:"id"`
	SeverityLevel string            `json:"severity_level"`
	Name          string            `json:"name"`
	ImageURL      string            `json:"image_url"`
	OriginalLink  string            `json:"original_link"`
	Details       map[string]string `json:"details"`
}

const (
	baseSecurityWarningsURL = "https://cert.gov.sa/en/security-warnings/?page="
	baseURL                 = "https://cert.gov.sa"
)

var (
	details         = make([]map[string]string, 0)
	severity_levels = make([]string, 0)
	bodies          = make([]string, 0)
	images          = make([]string, 0)
	links           = make([]string, 0)
)

func RunSaudiCertCrawler(fromPage, toPage int) {
	records := saudiCertCrawler(fromPage, toPage)
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	fileName := "./alerts/alerts.json"
	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing to file %s: %v", fileName, err)
	} else {
		fmt.Printf("Data written successfully to %s\n", fileName)
	}
}

func saudiCertCrawler(fromPage, toPage int) []Record {
	c := colly.NewCollector()
	var localHeaders, localBodies, localImages, localLinks []string
	uniqueLinks := make(map[string]bool)
	c.OnHTML(".card-header", func(e *colly.HTMLElement) {
		localHeaders = append(localHeaders, strings.TrimSpace(e.Text))
	})
	c.OnHTML(".cert-card-body-warning", func(e *colly.HTMLElement) {
		localBodies = append(localBodies, strings.TrimSpace(e.Text))
	})
	c.OnHTML(".security-alerts-cover-image", func(e *colly.HTMLElement) {
		localImages = append(localImages, baseURL+strings.TrimSpace(e.Attr("src")))
	})
	c.OnHTML(".col-md-6 a", func(e *colly.HTMLElement) {
		link := baseURL + strings.TrimSpace(e.Attr("href"))
		if !uniqueLinks[link] {
			uniqueLinks[link] = true
			localLinks = append(localLinks, link)
		}
	})
	c.OnScraped(func(r *colly.Response) {
		severity_levels = append(severity_levels, localHeaders...)
		bodies = append(bodies, localBodies...)
		images = append(images, localImages...)
		links = append(links, localLinks...)
		for _, link := range localLinks {
			processLinkFurther(link)
		}
	})
	for i := fromPage; i <= toPage; i++ {
		localHeaders = []string{}
		localBodies = []string{}
		localImages = []string{}
		localLinks = []string{}
		uniqueLinks = make(map[string]bool)
		fmt.Println(baseSecurityWarningsURL + strconv.Itoa(i))
		err := c.Visit(baseSecurityWarningsURL + strconv.Itoa(i))
		if err != nil {
			log.Printf("Error visiting page %d: %v\n", i, err)
			continue
		}
	}
	return convertToRecords()
}

func convertToRecords() []Record {
	fmt.Println(len(severity_levels))
	fmt.Println(len(links))
	fmt.Println(len(details))
	if len(severity_levels) != len(bodies) || len(severity_levels) != len(images) || len(severity_levels) != len(links) || len(severity_levels) != len(details) {
		log.Fatal("Slices are of different lengths")
	}
	records := make([]Record, len(severity_levels))
	for i := 0; i < len(severity_levels); i++ {
		records[i] = Record{
			ID:            i + 1,
			SeverityLevel: severity_levels[i],
			Name:          bodies[i],
			ImageURL:      images[i],
			OriginalLink:  links[i],
			Details:       details[i],
		}
	}
	return records
}

func processLinkFurther(link string) {
	keys := []string{"warning_date", "warning_number", "targeted_sector"}
	values := []string{}
	var description, threats, bestPractice string
	var affectedProducts, threatList, recommendationLinks []string
	c := colly.NewCollector()
	index := 0
	c.OnHTML(".col-7 p", func(e *colly.HTMLElement) {
		index++
		if index == 2 {
			return
		}
		value := strings.TrimSpace(e.Text)
		values = append(values, value)
	})
	c.OnHTML(".cert-body.cert-gray-70.m-3", func(e *colly.HTMLElement) {
		e.ForEach("p", func(index int, el *colly.HTMLElement) {
			if el.Text == "Description:" {
				description = el.DOM.Next().Text()
				el.DOM.Next().Next().Children().Each(func(i int, s *goquery.Selection) {
					if s.Nodes[0].Data == "li" {
						affectedProduct := strings.TrimSpace(s.Text())
						affectedProducts = append(affectedProducts, affectedProduct)
					}
				})
			}
			if el.Text == "Threats:" {
				threats = el.DOM.Next().Text()
				el.DOM.Next().Next().Children().Each(func(i int, s *goquery.Selection) {
					if s.Nodes[0].Data == "li" {
						threat := strings.TrimSpace(s.Text())
						threatList = append(threatList, threat)
					}
				})
			}
			if el.Text == "Best practice and Recommendations:" {
				bestPractice = el.DOM.Next().Text()
				el.DOM.Next().Next().Children().Each(func(i int, s *goquery.Selection) {
					if s.Nodes[0].Data == "li" {
						link := strings.TrimSpace(s.Text())
						recommendationLinks = append(recommendationLinks, link)
					} else if s.Nodes[0].Data == "a" {
						link := strings.TrimSpace(s.Text())
						recommendationLinks = append(recommendationLinks, link)
					}
				})
			}
		})
		e.ForEach("strong", func(index int, el *colly.HTMLElement) {
			if el.Text == "Description:" {
				descriptionNode := el.DOM.Next().Children().First()
				description = strings.TrimSpace(descriptionNode.Text())
				if description == "" {
					descriptionNode = descriptionNode.Next()
					description = strings.TrimSpace(descriptionNode.Text())
				}
				el.DOM.Next().Children().Last().Each(func(i int, s *goquery.Selection) {
					s.Children().Each(func(i int, s *goquery.Selection) {
						affectedProduct := strings.TrimSpace(s.Text())
						affectedProducts = append(affectedProducts, affectedProduct)
					})
				})
			}
			if el.Text == "Threats:" {
				threats = el.DOM.Next().Children().First().Text()
				if threats == "" {
					threats = el.DOM.Next().Children().Next().Text()
				}
				if el.DOM.Next().Children().Length() == 2 {
					el.DOM.Next().Children().Last().Each(func(i int, s *goquery.Selection) {
						s.Children().Each(func(i int, s *goquery.Selection) {
							threat := strings.TrimSpace(s.Text())
							threatList = append(threatList, threat)
						})
					})
				} else {
					el.DOM.Next().Children().Next().Next().Each(func(i int, s *goquery.Selection) {
						s.Children().Each(func(i int, s *goquery.Selection) {
							threat := strings.TrimSpace(s.Text())
							threatList = append(threatList, threat)
						})
					})
				}
			}
			if el.Text == "Best practice and Recommendations:" {
				bestPractice = el.DOM.Next().Children().First().Text()
				if bestPractice != "" {
					bestPractice = el.DOM.Next().Children().Next().Text()
				}
				el.DOM.Next().Children().Last().Each(func(i int, s *goquery.Selection) {
					if s.Is("ul") {
						el.DOM.Next().Children().Last().Each(func(i int, s *goquery.Selection) {
							s.Children().Each(func(i int, s *goquery.Selection) {
								link := strings.TrimSpace(s.Text())
								recommendationLinks = append(recommendationLinks, link)
							})
						})
					} else if s.Is("p") {
						el.DOM.Next().Children().Each(func(i int, s *goquery.Selection) {
							s.Children().Each(func(i int, s *goquery.Selection) {
								link := strings.TrimSpace(s.Text())
								recommendationLinks = append(recommendationLinks, link)
							})
						})
					}
				})
			}
		})
	})
	err := c.Visit(link)
	if err != nil {
		log.Printf("Error visiting %s: %v", link, err)
		return
	}
	dataMap := make(map[string]string)
	for i, key := range keys {
		if i < len(values) {
			dataMap[key] = values[i]
		}
	}
	re := regexp.MustCompile(`\s+`)
	dataMap["description"] = description
	dataMap["affected_products"] = re.ReplaceAllString(strings.Join(affectedProducts, "#"), " ")
	dataMap["threats"] = threats
	dataMap["threat_list"] = re.ReplaceAllString(strings.Join(threatList, "#"), " ")
	dataMap["best_practice"] = bestPractice
	dataMap["recommendation_links"] = re.ReplaceAllString(strings.Join(recommendationLinks, "#"), " ")
	details = append(details, dataMap)
}
