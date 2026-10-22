package services

import (
	"animals/internal/models"
	"fmt"
	"strings"
)

func AnimalCard(animal models.Animal) string {
	return fmt.Sprintf(
		`<div class="animal">
			%s
			<br />

			<b>Порода</b>
			%s
			<br />

			<b>Высота</b>
			%d
			<br />

			<b>Возраст</b>
			%d
			<br />

			<img src="%s" width="450px" />
		</div>`,
		animal.GetName(),
		animal.GetBreed(),
		animal.GetHeight(),
		animal.GetAge(),
		animal.GetImageUrl(),
	)
}

func AnimalsListDocument(appState *models.AppState, animals []models.Animal) string {
	var animalsRender strings.Builder

	for _, animal := range animals {
		animalsRender.WriteString(AnimalCard(animal))
	}

	return fmt.Sprintf(
		` <!DOCTYPE html>
			<html>
			<head>
			<title>%s</title>
			<style>
				html {
					background: white;
					filter: invert(1) hue-rotate(180deg);
				}

				img {
					filter: invert(1) hue-rotate(180deg);
				}

				.animal {
					border: 1px solid  black;
					border-radius: 8px;
					min-width: 500px;
					font-size: 20pt;
					padding: 10px 20px;
				}

				.animals-wrapper {
					display: flex;
					flex-wrap: wrap;
					gap: 14px;
				}
			</style>
			</head>
			<body>
				<form action="POST">

				</form>

				<div class="animals-wrapper">
					%s
				</div>
			</body>
			</html>
		`,
		appState.Title,
		animalsRender.String(),
	)
}
