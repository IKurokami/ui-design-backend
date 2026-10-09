package main

import "time"

type ChapterSeed struct {
	Number        int
	Title         string
	Content       string
	Illustrations []string
}

type StorySeed struct {
	ID          string
	Title       string
	Description string
	CoverURL    string
	Category    string
	Tags        string
	Status      string
	CreatedAt   time.Time
	Chapters    []ChapterSeed
}
