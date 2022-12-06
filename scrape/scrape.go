package scrape

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/zikster3262/shared-lib/chapter"
	"github.com/zikster3262/shared-lib/img"
	"github.com/zikster3262/shared-lib/page"
	"github.com/zikster3262/shared-lib/source"
	"github.com/zikster3262/shared-lib/utils"
)

func ScapeSource(mangaSource source.SQL) []page.Page {
	var pages []page.Page
	// Request the HTML page.
	res, err := http.Get(mangaSource.MangaURL)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	// Find the review items
	doc.Find(mangaSource.HomePattern).Each(func(i int, s *goquery.Selection) {
		// Fetch page address
		hrefAttr, _ := s.Attr("href")

		// Fetch page title
		t, _ := s.Attr("title")

		// create Page
		mangaPage := page.Page{
			URL:            hrefAttr,
			Title:          t,
			SourceID:       mangaSource.ID,
			Append:         mangaSource.Append,
			ChapterPattern: mangaSource.ChapterPattern,
			PagePattern:    mangaSource.PagePattern,
		}

		if mangaSource.Append {
			mangaPage.URL = mangaSource.MangaURL + hrefAttr
		}

		pages = append(pages, mangaPage)

	})

	return pages
}

func ScapePage(pgs page.SQL, pageID, sid int64) []page.SQL {
	var pages []page.SQL
	// Request the HTML page.
	res, err := http.Get(pgs.URL)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	// Find the review items
	doc.Find(pgs.PagePattern).Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")

		mangaPage := page.SQL{
			ID:             pageID,
			Title:          pgs.Title,
			URL:            href,
			SourceID:       int(sid),
			PagePattern:    "",
			ChapterPattern: pgs.ChapterPattern,
			DateAdded:      sql.NullTime{},
			Append:         pgs.Append,
		}

		pages = append(pages, mangaPage)

	})

	return pages
}

func ScapeChapter(cha chapter.Chapter) []img.Image {
	var images []img.Image
	// Request the HTML page.
	res, err := http.Get(cha.URL)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	chapter := utils.GetIDFromChapterURL(cha.URL)

	// Find the review items
	doc.Find(cha.ChapterPattern).Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("src")

		img := img.Image{
			Title:    cha.Title,
			URL:      href,
			Chapter:  chapter,
			Filename: utils.GetFileName(href),
		}
		images = append(images, img)

	})

	return images
}
