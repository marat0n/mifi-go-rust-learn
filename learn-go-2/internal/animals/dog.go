package animals

type DogBreed = string

const (
	DogBreedGerman DogBreed = "German"
	DogBreedCollie DogBreed = "Collie"
	DogBreedNA     DogBreed = "DogNA"
)

type Dog struct {
	Name   string
	Breed  DogBreed
	Age    int
	Height int
}

func (dog *Dog) GetName() string {
	return dog.Name
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
