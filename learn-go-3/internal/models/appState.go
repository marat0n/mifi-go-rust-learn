package models

// Инкапсуляция   ->  In capsule

type TitleState struct {
	Title string
}

func (self *TitleState) SetTitle(newTitle string) {
	self.Title = newTitle
}

type AnimalsState struct {
	Animals []Animal
}

type AppState struct {
	TitleState
	AnimalsState
}

func CreateAppState(title string) AppState {
	title = "Title: " + title
	return AppState{
		TitleState: TitleState{
			Title: title,
		},
		AnimalsState: AnimalsState{
			Animals: []Animal{
				&Cat{
					Name:   "Буся",
					Breed:  CatBreedEgyptian,
					Age:    14,
					Height: 25,
				},
				&Cat{
					Name:   "Дуся",
					Breed:  CatBreedBritain,
					Age:    12,
					Height: 23,
				},
				&Dog{
					Name:   "Шарик",
					Breed:  DogBreedCollie,
					Age:    16,
					Height: 48,
				},
				&Dog{
					Name:   "Гарик",
					Breed:  DogBreedGerman,
					Age:    8,
					Height: 62,
				},
				&Cat{
					Name:   "Гуся",
					Breed:  CatBreedNA,
					Age:    3,
					Height: 31,
				},
			},
		},
	}
}
