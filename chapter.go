package main

import "reflect"

type Chapter struct {
	Url   string   `json:"url"`
	Title string   `json:"title"`
	Pages []string `json:"pages"`
}

func NewChapter(url string) Chapter {
	return Chapter{Url: url}
}

func (c Chapter) WithTitle(title string) Chapter {
	c.Title = title
	return c
}

func (c Chapter) WithPages(pages []string) Chapter {
	c.Pages = pages
	return c
}

func (c Chapter) LessThan(other Chapter) bool {
	return reflect.DeepEqual(c, other)
}

func (c Chapter) EqualTo(other Chapter) bool {
	return reflect.DeepEqual(c, other)
}
