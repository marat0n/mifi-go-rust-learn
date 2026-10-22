package models

type Animal interface {
	GetName() string
	GetBreed() string
	GetAge() int
	GetHeight() int

	GetImageUrl() string

	SetBreed(string)
	SetName(string)
	SetAge(int)
	SetHeight(int)
}

type AnimalFromDb struct {
	Id     int
	Name   string
	Breed  string
	Age    int
	Height int
}
