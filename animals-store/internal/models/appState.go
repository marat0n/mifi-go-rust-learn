package models

import sqlite "gosqlite.org"

type TitleState struct {
	Title string
}

func (self *TitleState) SetTitle(newTitle string) {
	self.Title = newTitle
}

type DbState struct {
	DB *sqlite.DB
}

type AppState struct {
	TitleState
	DbState
}

func CreateAppState(title string, db *sqlite.DB) AppState {
	title = "Title: " + title
	return AppState{
		TitleState{Title: title},
		DbState{DB: db},
	}
}
