package scrape

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/zikster3262/shared-lib/chapter"
	"github.com/zikster3262/shared-lib/img"
	"github.com/zikster3262/shared-lib/page"
	"github.com/zikster3262/shared-lib/source"
)

func ScapeSource(mp source.SourceSQL) (m []page.Page) {
	// Request the HTML page.
	res, err := http.Get(mp.Manga_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(mp.Home_Pattern).Each(func(i int, s *goquery.Selection) {

		// Fetch page address
		v, _ := s.Attr("href")

		// Fetch page title
		t, _ := s.Attr("title")

		// create Page
		mn := page.Page{
			Url:             v,
			Title:           t,
			Source_Id:       mp.Id,
			Append:          mp.Append,
			Chapter_Pattern: mp.Chapter_Pattern,
			Page_Pattern:    mp.Page_Pattern,
		}

		if mp.Append {
			mn.Url = mp.Manga_URL + v
		}

		m = append(m, mn)

	})

	return m
}

func ScapePage(p page.PageSQL, page_id, sid int64) (m []page.PageSQL) {
	// Request the HTML page.
	res, err := http.Get(p.Url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(p.Page_Pattern).Each(func(i int, s *goquery.Selection) {

		href, _ := s.Attr("href")

		mn := page.PageSQL{
			Id:              page_id,
			Url:             href,
			Title:           p.Title,
			Source_Id:       int(sid),
			Chapter_Pattern: p.Chapter_Pattern,
			Append:          p.Append,
		}

		m = append(m, mn)

	})

	return m
}

func ScapeChapter(cha chapter.Chapter) (images []img.Image) {
	// Request the HTML page.
	res, err := http.Get(cha.Url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(cha.Chapter_Pattern).Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("src")

		img := img.Image{
			Title:   cha.Title,
			Url:     href,
			Chapter: cha.Url,
		}
		images = append(images, img)

	})
	return images
}
