package main

import (
	"fmt"
	"learn-3/internal/models"
	"net/http"
	"strings"
)

// 1. DONE: Interface

// 2. DONE: Struct injecting

// 3. TODO: JSON serialization / deserialization

// Любая программа это: данные и поведение

// ООП:
// 1. Инкапсуляция
// 2. Наследование
// 3. Полиморфизм
//
// Большинство ООП языков сразу работают с ссылками и классами

// Class vs Struct
// Class: это уже ссылка на данные
//			+ дополнительные методы для работы с данными
//			например:
//				1. сравнение ссылок,
//				2. расчёт хэша объекта,
//				3. сравнение хэшей объектов
//
// Struct: это только значение

func main() {
	app := models.CreateAppState("Animals Shop")

	http.HandleFunc("/list", func(writer http.ResponseWriter, req *http.Request) {
		var animalsRender strings.Builder

		for index, animal := range app.Animals {
			fmt.Fprintf(
				&animalsRender,
				`<div class="animal">
					%d. %s
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
				index,
				animal.GetName(),
				animal.GetBreed(),
				animal.GetHeight(),
				animal.GetAge(),
				animal.GetImageUrl(),
			)
		}

		fmt.Fprintf(
			writer,
			` <!DOCTYPE html>
				<html>
				<head>
				<title>%s</title>
				<style>
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
					<div class="animals-wrapper">
						%s
					</div>
				</body>
				</html>
			`,
			app.Title,
			animalsRender.String(),
		)
	})

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(
			w,
			` <!DOCTYPE html>
				<html>
				<head>
				<title>%s</title>
				</head>
				<body>
					%s
				</body>
				</html>
			`,
			app.Title,
			"Hello, World!",
		)
	})

	server := &http.Server{
		Addr: ":4000",
	}

	server.ListenAndServe()
}
