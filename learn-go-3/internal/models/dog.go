package models

import "learn-3/internal/services"

type DogBreed = string

const (
	DogBreedGerman DogBreed = "German"
	DogBreedCollie DogBreed = "Collie"
	DogBreedNA     DogBreed = "DogNA"
)

var dogsImages = []string{
	"https://images.dog.ceo/breeds/komondor/n02105505_3389.jpg",
	"https://images.dog.ceo/breeds/leonberg/n02111129_1832.jpg",
	"https://images.dog.ceo/breeds/maltese/n02085936_2927.jpg",
	"https://images.dog.ceo/breeds/cavapoo/doggo3.jpg",
	"https://images.dog.ceo/breeds/deerhound-scottish/n02092002_6003.jpg",
}

type Dog struct {
	Name   string
	Breed  DogBreed
	Age    int
	Height int
}

func (dog *Dog) GetName() string {
	return dog.Name
}

func (dog *Dog) GetBreed() string {
	breed := dog.Breed
	if breed == DogBreedNA {
		breed = "NA"
	}
	return "Dog: " + breed
}

func (dog *Dog) GetAge() int {
	return dog.Age
}

func (dog *Dog) GetHeight() int {
	return dog.Height
}

func (_ *Dog) GetImageUrl() string {
	randomNumber := services.Random(0, 5)
	return dogsImages[randomNumber]
}

func (dog *Dog) SetBreed(breed string) {
	dog.Breed = breed
}

func (dog *Dog) SetName(name string) {
	dog.Name = name
}

func (dog *Dog) SetAge(age int) {
	dog.Age = age
}

func (dog *Dog) SetHeight(height int) {
	dog.Height = height
}
