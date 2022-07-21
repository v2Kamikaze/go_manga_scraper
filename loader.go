package main

import (
	"fmt"
	"log"
	"os"
)

// HeavyStart carrega todos os mangás e todos os seus capítulos pela primeira vez.
// Demora algumas horas.
func HeavyStart(path string) {
	var mangaList []Manga

	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE, os.FileMode(0600))
	if err != nil {
		fmt.Println(err)
		return
	}
	log.SetOutput(file)

	defer file.Close()

	// Buscando a lista de mangás se não estiver salva.
	if mangaList, err = LoadList(path); err != nil {
		fmt.Println("Buscando a lista de mangás...")
		mangaList = MangaList("https://mangaschan.com/mangas/page/")
		if err := SaveMangaList(path, &mangaList); err != nil {
			log.Println(err.Error())
		}
	}

	// Para cada mangá, busca a lista de capítulos.
	for i := range mangaList {
		// Caso o mangá já esteja com a sua lista de capítulas adicionada, passamos para o próximo mangá.
		// Economizando tempo.
		if len(mangaList[i].Chapters) != 0 {
			continue
		}
		fmt.Printf("Buscando informações do mangá %s\n", mangaList[i].Title)
		m, err := GetMangaInfo(mangaList[i].Url)
		if err != nil {
			log.Printf("erro ao buscar informações do mangá '%s'. %+v", mangaList[i].Title, err)
		} else {
			mangaList[i] = m
			if err := SaveMangaList(path, &mangaList); err != nil {
				fmt.Println(err.Error())
			}
		}
	}

	// Para cada capítulo de cada mangá, busca as páginas desse capítulo.
	for i := range mangaList {
		for j := range mangaList[i].Chapters {
			// Caso o capítulo já tenha páginas adicionadas, passamos para o próximo capítulo.
			// Economizando tempo.
			if len(mangaList[i].Chapters[j].Pages) == 0 {
				fmt.Printf("Buscando %s do mangá %s\n", mangaList[i].Chapters[j].Title, mangaList[i].Title)
				c := GetChapter(mangaList[i].Chapters[j])

				if c.Title == "" || c.Url == "" || len(c.Pages) == 0 {
					log.Println(mangaList[i].Chapters[j])
				}
				mangaList[i].Chapters[j] = c
			}
		}

		if err := SaveMangaList(path, &mangaList); err != nil {
			fmt.Println(err.Error())
		}
	}
}
