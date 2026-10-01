package animals

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
