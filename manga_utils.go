package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// NumPagesOfSite retorna o número de páginas que possuem mangás.
// Se o elemento não for encontrado, retorna -1.
func NumPagesOfSite(doc *goquery.Document) int {
	numPagesStr := doc.Find("#content > div.wrapper > div.postbody > div > div.page > div.pagination > a:nth-child(5)").Text()
	numPages, err := strconv.Atoi(numPagesStr)

	if err != nil {
		return -1
	}

	return numPages
}

// FindMangaList retorna uma lista com os mangás contidos na página.
func FindMangaList(doc *goquery.Document) []Manga {

	var mangas []Manga

	doc.Find("#content > div.wrapper > div.postbody > div > div.page > div.listupd.cp > div > div > a").Each(func(i int, s *goquery.Selection) {
		title := s.Find("div.bigor > div.tt").Text()
		title = strings.TrimPrefix(title, "\n")
		url, exists := s.Attr("href")
		if exists {
			m := NewManga(url).WithTitle(title)
			mangas = append(mangas, m)
		}
	})

	return mangas
}

// FindNoScriptTag retorna o conteúdo da tag <noscript> para a criação de um novo Reader.
func FindNoScriptTag(doc *goquery.Document) string {
	return doc.Find("noscript").Text()
}

// GetChapterPages retorna as páginas de um capítulo.
func GetChapterPages(doc *goquery.Document) []string {

	var pages []string

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		pageUrl, exists := s.Attr("href")
		if exists {
			pages = append(pages, pageUrl)
		}
	})

	if pages == nil {
		doc.Find("img").Each(func(i int, s *goquery.Selection) {
			pageUrl, exists := s.Attr("src")
			if exists {
				pages = append(pages, pageUrl)
			}
		})
	}

	return pages
}

// SaveMangaList salva a lista de mangás no diretório passado.
// Retorna um erro caso não seja possível.
func SaveMangaList(path string, mangaList *[]Manga) error {
	file, err := os.Create(path)

	if err != nil {
		return fmt.Errorf("erro ao criar arquivo '%s'. %+v", path, err)
	}

	defer file.Close()

	enconder := json.NewEncoder(file)
	enconder.SetEscapeHTML(true)
	enconder.SetIndent("", "  ")

	if err := enconder.Encode(mangaList); err != nil {
		return fmt.Errorf("erro ao salvar lista de mangás. %+v", err)
	}

	return nil
}

// LoadList carrega a lista de mangás contida no diretório passado em memória.
// Retorna a lista e nil caso tudo ocorra bem.
// Retorna uma lista vazia e um erro caso contrário.
func LoadList(path string) ([]Manga, error) {

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return []Manga{}, err
	}

	file, err := os.Open(path)

	if err != nil {
		return []Manga{}, err
	}

	defer file.Close()

	var mangas []Manga

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&mangas); err != nil {
		return []Manga{}, err
	}

	return mangas, nil
}
