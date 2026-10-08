package models

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
	return "Cat: " + cat.Breed
}

func (cat *Cat) GetAge() int {
	return cat.Age
}

func (cat *Cat) GetHeight() int {
	return cat.Height
}

func (_ *Cat) GetImageUrl() string {
	return "https://cataas.com/cat"
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
