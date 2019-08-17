package fanfox

import (
	"fmt"
	"github.com/freddieptf/manga-scraper/pkg/scraper"
	"net/url"
	"strings"
	"testing"
)

// return the results and the search query
func getTestSearchResults(t *testing.T) (scraper.Manga, string) {
	testManga := "kengan"
	foxSource := &FoxManga{}
	results, err := foxSource.Search(testManga)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) <= 0 {
		t.Fatal("search results of length 0")
	}
	return results[0], testManga
}

func TestFoxGetMangaDetails(t *testing.T) {
	manga, _ := getTestSearchResults(t)
	doc, err := openFoxPage(fmt.Sprintf("%s/%s", foxURL, manga.MangaID))
	if err != nil {
		t.Fatal(err)
	}
	details := getMangaDetails(doc)
	if details.name != manga.MangaName {
		t.Error("manga details don't match")
	}
}

func testGetChapterUrlFromListing(t *testing.T) (chapterURL string) {
	manga, _ := getTestSearchResults(t)
	doc, err := openFoxPage(fmt.Sprintf("%s/%s", foxURL, manga.MangaID))
	if err != nil {
		t.Fatal(err)
	}
	chapterURL, chapterTitle := getChapterUrlFromListing("1", doc)
	if _, err := url.ParseRequestURI(chapterURL); err != nil {
		t.Fatalf("couldn't get chapterURL from listing = '%s' err: %v\n", chapterURL, err)
	}
	if chapterTitle == "" {
		t.Error("returned empty chapter title")
	}
	return
}

func testGetFoxChPageUrls(t *testing.T) (chapterPageUrls []string) {
	chapterURL := testGetChapterUrlFromListing(t)
	doc, err := openFoxPage(chapterURL)
	if err != nil {
		t.Fatal(err)
	}
	chapterPageUrls = getFoxChPageUrls(doc)
	if len(chapterPageUrls) <= 0 {
		t.Fatal("couldn't get the chapter page urls")
	}
	for _, chapterPageURL := range chapterPageUrls {
		if _, err := url.ParseRequestURI(chapterPageURL); err != nil {
			t.Errorf("chapterPageURL=%s : %v\n", chapterPageURL, err)
		}
	}
	return
}

func TestGetFoxChPageImgUrl(t *testing.T) {
	chapterPageUrls := testGetFoxChPageUrls(t)
	samplePageUrl := chapterPageUrls[0]
	imgUrl := getFoxChPageImgUrl(samplePageUrl)
	if _, err := url.ParseRequestURI(imgUrl); err != nil {
		t.Error(err)
	}
}

func TestFoxSearch(t *testing.T) {
	manga, query := getTestSearchResults(t)
	if !strings.Contains(strings.ToLower(manga.MangaName), query) {
		t.Error("results could be wrong")
	}
}

func TestFoxGetChapter(t *testing.T) {
	manga, _ := getTestSearchResults(t)
	foxSource := &FoxManga{}
	chapter, err := foxSource.GetChapter(manga.MangaID, "4")
	if err != nil {
		t.Error(err)
	}
	if len(chapter.ChapterPages) <= 0 {
		t.Error("no chapter pages returned")
	}
	for _, page := range chapter.ChapterPages {
		if _, err := url.ParseRequestURI(page.Url); err != nil {
			t.Errorf("invalid page url %s\n", page.Url)
		}
	}
}