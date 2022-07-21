package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// MangaList retorna uma lista contento todos os mangás presentes no site.
// Retorna uma lista vazia caso ocorra algum erro.
// Qualquer erro será escrito no logger.
func MangaList(url string) []Manga {

	var mangas []Manga

	res, err := http.Get(url + "/1")

	if err != nil {
		log.Printf("a página '%s' não está disponível <MangaList>. %+v\n", url+"/1", err)
		return []Manga{}
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("status code %d em '%s' <MangaList>", res.StatusCode, url)
		return []Manga{}
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Printf("erro ao criar reader da página '%s' em <MangaList>. %+v\n", url+"/1", err.Error())
		return []Manga{}
	}

	numPages := NumPagesOfSite(doc)

	if numPages == -1 {
		log.Println("não foi possível encontrar o número de páginas com mangás do site. <MangaList>")
		return []Manga{}
	}

	fmt.Println("[x] Página 1 OK!")
	mangas = append(mangas, FindMangaList(doc)...)

	for i := 2; i <= numPages; i++ {
		ml := GetOnePageList(url + strconv.Itoa(i))
		if ml != nil {
			fmt.Printf("[x] Página %d OK!.\n", i)
			mangas = append(mangas, ml...)
		} else {
			fmt.Printf("[!] Página %d não retornou nenhum mangá.\n", i)
		}
	}

	return mangas

}

// GetOnePageList retorna uma lista com os mangás encontrados nessa página.
func GetOnePageList(url string) []Manga {
	res, err := http.Get(url)

	if res.StatusCode != http.StatusOK || err != nil {
		return []Manga{}
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return []Manga{}
	}

	return FindMangaList(doc)
}

// GetMangaInfo returna um mangá com suas informações básicas
// preenchidas e um erro caso ocorra.
func GetMangaInfo(url string) (Manga, error) {

	res, err := http.Get(url)

	if err != nil {
		return NewManga(url), err
	}

	if res.StatusCode != http.StatusOK {
		return NewManga(url), fmt.Errorf("status code %d em '%s' <GetMangaInfo>", res.StatusCode, url)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return NewManga(url), fmt.Errorf("erro ao criar reader em <GetMangaInfo>. %+v", err.Error())
	}

	title := GetTitle(doc)
	chapList := GetChaptersList(doc)

	return NewManga(url).WithTitle(title).AddChapters(chapList), nil
}

// GetTitle retorna o título do mangá presente na página.
func GetTitle(doc *goquery.Document) string {
	title := doc.Find("div.seriestuheader > h1").Text()
	return title
}

// GetChaptersList retorna a lista de capítulos presentes na página.
func GetChaptersList(doc *goquery.Document) []Chapter {

	var chapList []Chapter

	doc.Find("#chapterlist > ul > li> div > div").Each(func(i int, s *goquery.Selection) {
		title := s.Find("a > span.chapternum").Text()
		url, exists := s.Find("a").Attr("href")
		if exists {
			chapList = append(chapList, NewChapter(url).WithTitle(title))
		}
	})

	return chapList
}

// GetChapter retorna um capítulo de um mangá com todas as suas informações preenchidas.
func GetChapter(chapter Chapter) Chapter {
	res, err := http.Get(chapter.Url)

	if err != nil {
		log.Printf("não foi possível buscar o capítulo '%s' de '%s' <GetChapter>. %+v\n", chapter.Title, chapter.Url, err)
		return Chapter{}
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("status code %d em '%s' <GetChapter>", res.StatusCode, chapter.Url)
		return Chapter{}
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Println("erro ao criar reader em <GetChapter>. ", err.Error())
	}

	noscriptText := FindNoScriptTag(doc)

	reader := strings.NewReader(noscriptText)

	doc, err = goquery.NewDocumentFromReader(reader)

	if err != nil {
		log.Println("erro ao criar reader <noscript> em <GetChapter>. ", err.Error())
		return chapter
	}

	pages := GetChapterPages(doc)

	return chapter.WithPages(pages)
}
