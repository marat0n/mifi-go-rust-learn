package main

import (
	"fmt"
	a "learn/internal/animals"
	randomgenerator "learn/internal/randomGenerator"
)

// Спрашиваем пользователя: Добавить ещё животное?
// Да:
//		Породу
//		Имя
//		Возраст
//		Высоту
// Нет:
//		Заканчиваем спрашивать и выводим результат

func main() {
	fmt.Println(randomgenerator.Random1(50, 77))
	fmt.Println(randomgenerator.Random2(50, 77))
}

func zoo() {
	animals := []a.Animal{}

	for {
		fmt.Println("Добавить ещё животное?")

		userDecision := ""
		_, err := fmt.Scan(&userDecision)
		if err != nil {
			fmt.Printf("Случилась ошибка: %s\n", err)
			break
		}

		if userDecision == "Да" {
			var animal a.Animal

			breed := ""
			_, err := fmt.Scan(&breed)
			if err != nil {
				fmt.Printf("Случилась ошибка: %s\n", err)
				break
			}

			switch breed {
			case a.CatBreedEgyptian:
				animal = &a.Cat{}
			case a.CatBreedBritain:
				animal = &a.Cat{}
			case a.CatBreedNA:
				animal = &a.Cat{}
			case a.DogBreedCollie:
				animal = &a.Dog{}
			case a.DogBreedGerman:
				animal = &a.Dog{}
			case a.DogBreedNA:
				animal = &a.Dog{}
			}

			if animal == nil {
				fmt.Println("Переменная пуста!")
			}
			animal.SetBreed(breed)

			name := ""
			_, err = fmt.Scan(&name)
			if err != nil {
				fmt.Printf("Случилась ошибка: %s\n", err)
				break
			}

			animal.SetName(name)

			age := 0
			_, err = fmt.Scan(&name)
			if err != nil {
				fmt.Printf("Случилась ошибка: %s\n", err)
				break
			}

			animal.SetAge(age)

			height := 0
			_, err = fmt.Scan(&name)
			if err != nil {
				fmt.Printf("Случилась ошибка: %s\n", err)
				break
			}

			animal.SetHeight(height)

			animals = append(animals, animal)
		} else if userDecision == "Нет" {
			fmt.Println("Пока пока!")
			break
		}
	}

	fmt.Println()

	for index, animal := range animals {
		fmt.Printf("%d: %s\n", index, animal.GetName())
	}
}
