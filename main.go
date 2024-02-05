package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/dchooyc/film"
)

type Genres struct {
	Genres []string `json:"genres"`
}

func main() {
	input, err := film.RetrieveFilms("output.json")
	if err != nil {
		fmt.Println("retrieve films failed: ", err)
		return
	}
	fmt.Println("Length of all movies: ", len(input.Films))
	genreToMovies := sortByGenre(input.Films)
	createJsons(genreToMovies)
}

func createJsons(genreToFilms map[string][]film.Film) {
	genres := Genres{}
	err := os.Mkdir("./jsons", 0755)
	if err != nil {
		fmt.Println("failed creating dir: ", err)
		return
	}

	for genre := range genreToFilms {
		films := film.Films{Films: genreToFilms[genre]}
		filename := "jsons/" + genre + ".json"

		err := createJsonFilms(filename, films)
		if err != nil {
			fmt.Println(genre, err)
			continue
		}

		genres.Genres = append(genres.Genres, genre)
	}

	err = createJsonGenres("genres.json", genres)
	if err != nil {
		fmt.Println("failed creating genres json: ", err)
	}
}

func createJsonFilms(filename string, films film.Films) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed creating json: %w", err)
	}

	sort.Slice(films.Films, func(i, j int) bool {
		return films.Films[i].AudienceScore > films.Films[j].AudienceScore
	})

	jsonData, err := json.Marshal(films)
	if err != nil {
		return fmt.Errorf("failed marshalling json: %w", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed writing json: %w", err)
	}

	fmt.Println("Created: ", filename, " with length: ", len(films.Films))
	return nil
}

func createJsonGenres(filename string, genres Genres) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed creating json: %w", err)
	}

	sort.Strings(genres.Genres)

	jsonData, err := json.Marshal(genres)
	if err != nil {
		return fmt.Errorf("failed marshalling json: %w", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed writing json: %w", err)
	}

	return nil
}

func sortByGenre(films []film.Film) map[string][]film.Film {
	genreToMovies := make(map[string][]film.Film)

	for i := 0; i < len(films); i++ {
		f := films[i]
		genres := strings.Split(f.Genre, "/")
		for _, genre := range genres {
			genreToMovies[genre] = append(genreToMovies[genre], f)
		}
	}

	return genreToMovies
}
