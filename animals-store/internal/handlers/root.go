package handlers

import (
	"animals/internal/models"
	"animals/internal/services"
	"fmt"
	"log"
	"net/http"
)

type HandlerState struct {
	models.AppState
}

func (state *HandlerState) HandleRoot(writer http.ResponseWriter, req *http.Request) {
	rows, err := state.DB.Query("SELECT * FROM animals LIMIT 100")
	if err != nil {
		log.Printf("Error on selecting animals: %s", err)
		fmt.Fprintf(writer, "Error")
		return
	}
	defer rows.Close()

	var actualAnimals []models.Animal

	for rows.Next() {
		var id int
		var name string
		var breed string
		var age int
		var height int

		err = rows.Scan(&id, &name, &breed, &age, &height)
		if err != nil {
			log.Printf("Error on scanning an animal: %s", err)
			fmt.Fprintf(writer, "Error")
			return
		}

		animal := models.AnimalFromDb{
			Id:     id,
			Name:   name,
			Breed:  breed,
			Age:    age,
			Height: height,
		}

		if cat := models.TryCatFromDb(&animal); cat != nil {
			actualAnimals = append(actualAnimals, cat)
			log.Print(cat.Name)
		} else if dog := models.TryDogFromDb(&animal); dog != nil {
			actualAnimals = append(actualAnimals, dog)
			log.Print(dog.Name)
		} else {
			log.Printf("Error while parsing cats and dogs. This breed is not allowed: %s", animal.Breed)
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error in db rows: %s", err)
		fmt.Fprintf(writer, "Error")
		return
	}

	fmt.Fprint(writer, services.AnimalsListDocument(&state.AppState, actualAnimals))
}
