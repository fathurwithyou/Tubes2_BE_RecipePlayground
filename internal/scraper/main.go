// internal/scraper/scraper.go
package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
)

type Element struct {
	Name    string     `json:"name"`
	Recipes [][]string `json:"recipes"`
	Tier    int        `json:"tier"`
}

type Result struct {
	Elements []Element `json:"elements"`
}

func Scrape(filename string) error {
	const url = "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http get: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return fmt.Errorf("parse HTML: %w", err)
	}

	tables := doc.Find("table.list-table.col-list.icon-hover")
	if tables.Length() < 3 {
		return fmt.Errorf("unexpected table count: %d", tables.Length())
	}

	var elements []Element
	tier := 0
	tables.Each(func(i int, tbl *goquery.Selection) {
		// skip the second table (i == 1)
		if i == 1 {
			return
		}

		tbl.Find("tr").Each(func(rowIdx int, row *goquery.Selection) {
			if rowIdx == 0 {
				return
			}
			tds := row.Find("td")
			if tds.Length() < 2 {
				return
			}

			var name string
			tds.Eq(0).Contents().EachWithBreak(func(_ int, node *goquery.Selection) bool {
				if goquery.NodeName(node) == "a" {
					if title, ok := node.Attr("title"); ok {
						name = title
						return false
					}
				}
				return true
			})
			if name == "" {
				return
			}

			var recipes [][]string
			tds.Eq(1).Find("ul > li").Each(func(_ int, li *goquery.Selection) {
				var combo []string
				li.Find("a").Each(func(_ int, a *goquery.Selection) {
					if t, ok := a.Attr("title"); ok {
						combo = append(combo, t)
					}
				})
				if len(combo) > 0 {
					recipes = append(recipes, combo)
				}
			})

			elements = append(elements, Element{
				Name:    name,
				Recipes: recipes,
				Tier:    tier,
			})
		})

		tier++
	})

	out := Result{Elements: elements}
	filePath := filepath.Join("..", "data", filename)
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}

	return nil
}
