package main

import (
	"fmt"
	"hangman/hangman"
	"html/template"
	"log"
	"net/http"
)

/* NOTRE STRUCT
type HangManData struct {
	Word       []rune // Word composed of '_', ex: H_ll_
	ToFind     []rune // Final word chosen by the program at the beginning. It is the word to find
	Attempts   int    // Number of attempts left
	History    string
	Difficulty string
	UserEntry  string
	StrWord    string
	StrToFind  string
}
*/

var gamee hangman.HangManData
var game = &gamee

func Home(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./pages/index.html", "./templates/footer.html", "./templates/header.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func Informations(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./pages/informations.html", "./templates/footer.html", "./templates/header.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func Hangman(game *hangman.HangManData) *hangman.HangManData {
	var histobool bool = false
	letter := game.UserEntry
	if letter != "" {
		for j := 0; j < len(game.ToFind); j++ {
			ct := 0
			if rune(letter[0]) == game.ToFind[j] && len(letter) == 1 {
				game = hangman.ModifyGameWord(game, string(letter)) // func qui transforme la string Word en mettant en place la lettre
				break
			} else if len(letter) > 1 && string(letter) != string(game.ToFind) && j == len(game.ToFind)-1 && ct < 1 {
				game.Attempts -= 2
				ct++
				break
			} else if j == len(game.ToFind)-1 { // si la lettre testée est fausse alors on l'ajoute à la string Historique pour avoir un suivi et ne pas réessayer les même lettres
				for i := 0; i < len(game.History); i++ { // Permet de ne pas avoir une lettre qui apparaît plusieurs fois si elle est testée plusieurs fois
					if string(letter[0]) == string(game.History[i]) {
						histobool = true
					}
				}
				if !histobool {
					game.History += hangman.ToHigher(string(letter[0]))
					game.Attempts--
				}
				histobool = false
			}
		}
	}
	count := 0
	for g := 0; g < len(game.Word); g++ {
		if game.Word[g] != '_' {
			count++
		}
		if count == len(game.Word) || letter == string(game.ToFind) { // si tous les caractères de game.Word sont autre chose que des underscore
			game.Word = game.ToFind
			game.StrWord = game.StrToFind
		}
	}
	count = 0
	return game
}

func Game(w http.ResponseWriter, r *http.Request, game *hangman.HangManData) {
	template, err := template.ParseFiles("./pages/game.html", "./templates/footer.html", "./templates/header.html", "./templates/lose.html", "./templates/win.html", "./templates/diff.html")
	if err != nil {
		log.Fatal(err)
	}
	if game.Difficulty == "" {
		if r.FormValue("diff") == "Facile" {
			game.Difficulty = "facile"
			game = InitStruct(game)
		}
		if r.FormValue("diff") == "Normal" {
			game.Difficulty = "normal"
			game = InitStruct(game)
		}
		if r.FormValue("diff") == "Difficile" {
			game.Difficulty = "difficile"
			game = InitStruct(game)
		}
		if r.FormValue("diff") == "Tous les mots" {
			game.Difficulty = "tous les mots"
			game = InitStruct(game)
		}
	}
	if r.FormValue("rejouer") == "Rejouer" { // Boutons rejouer
		game.Difficulty = ""
	}
	if game.Difficulty != "" {
		game.UserEntry = string(hangman.ToLower([]rune(r.FormValue("UserEntry"))))
		fmt.Println(game.UserEntry)
		if game.UserEntry != "" {
			game = Hangman(game)
		}
	}
	template.Execute(w, game)
}

func InitStruct(game *hangman.HangManData) *hangman.HangManData {
	game.Word = nil
	game = hangman.RandomPickLine(game)
	T := hangman.RandomPickLetter(game)
	fmt.Println(string(game.ToFind)) // TEST
	for i := 0; i < len(game.ToFind); i++ {
		game.Word = append(game.Word, '_')
	}
	for j := 0; j < len(T); j++ {
		game.Word[T[j]] = game.ToFind[T[j]]
	}
	game.Attempts = 10
	game.UserEntry = ""
	game.History = ""
	game.StrWord = string(game.Word)
	game.StrToFind = string(game.ToFind)
	return game
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/informations", Informations)
	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		Game(w, r, game)
	})
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/Ressources/", http.StripPrefix("/Ressources/", http.FileServer(http.Dir("./Ressources"))))
	fmt.Println("Server running on port :8080")
	http.ListenAndServe(":8080", nil)
}
