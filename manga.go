package main

type Manga struct {
	Title    string    `json:"title"`
	Url      string    `json:"url"`
	Chapters []Chapter `json:"chapters"`
}

func NewManga(url string) Manga {
	return Manga{Url: url}
}

func (m Manga) WithTitle(title string) Manga {
	m.Title = title
	return m
}

func (m Manga) AddChapter(chapter Chapter) Manga {
	m.Chapters = append(m.Chapters, chapter)
	return m
}

func (m Manga) AddChapters(chapters []Chapter) Manga {
	m.Chapters = chapters
	return m
}
