package models

import (
	"animals/internal"
	"fmt"
	"strings"
)

type CatBreed = string

const (
	CatBreedEgyptian CatBreed = "CatEgyptian"
	CatBreedBritain  CatBreed = "CatBritain"
	CatBreedNA       CatBreed = "CatNA"
)

type Cat struct {
	Id     int
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
	digit := internal.Random(0, 100)
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

func TryCatFromDb(animal *AnimalFromDb) *Cat {
	if strings.HasPrefix(animal.Breed, "Cat") {
		return &Cat{
			Id:     animal.Id,
			Name:   animal.Name,
			Breed:  animal.Breed,
			Age:    animal.Age,
			Height: animal.Height,
		}
	}

	return nil
}
