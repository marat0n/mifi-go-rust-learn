package models

import (
	"fmt"
	"learn-3/internal/services"
)

type CatBreed = string

const (
	CatBreedEgyptian CatBreed = "Egyptian"
	CatBreedBritain  CatBreed = "Britain"
	CatBreedNA       CatBreed = "CatNA"
)

type Cat struct {
	Name   string
	Breed  CatBreed
	Age    int
	Height int
}

func (cat *Cat) GetName() string {
	return cat.Name
}

func (cat *Cat) GetBreed() string {
	breed := cat.Breed
	if breed == CatBreedNA {
		breed = "NA"
	}
	return "Cat: " + breed
}

func (cat *Cat) GetAge() int {
	return cat.Age
}

func (cat *Cat) GetHeight() int {
	return cat.Height
}

func (_ *Cat) GetImageUrl() string {
	digit := services.Random(0, 100)
	return fmt.Sprintf(
		"https://cataas.com/cat?%d",
		digit,
	)
}

func (cat *Cat) SetBreed(breed string) {
	cat.Breed = breed
}

func (cat *Cat) SetName(name string) {
	cat.Name = name
}

func (cat *Cat) SetAge(age int) {
	cat.Age = age
}

func (cat *Cat) SetHeight(height int) {
	cat.Height = height
}
