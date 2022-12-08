package hangman

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"time"
)

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

func RandomPickLine(game *HangManData) *HangManData {
	var diff string
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	if game.Difficulty == "facile" {
		diff = "hangman/Words/words1.txt"
	}
	if game.Difficulty == "normal" {
		diff = "hangman/Words/words2.txt"
	}
	if game.Difficulty == "difficile" {
		diff = "hangman/Words/words3.txt"
	}
	if game.Difficulty == "tous les mots" {
		diff = "hangman/Words/words.txt"
	}
	file, err := os.Open(diff)
	var lines []string
	if err != nil {
		log.Fatal(err)
	}
	Scanner := bufio.NewScanner(file)
	Scanner.Split(bufio.ScanWords)
	for Scanner.Scan() {
		lines = append(lines, Scanner.Text())
	}
	if err := Scanner.Err(); err != nil {
		log.Fatal(err)
	}
	randomnum := r1.Intn(len(lines))
	game.ToFind = []rune(lines[randomnum])
	return game
}

func ModifyGameWord(game *HangManData, letter string) *HangManData { // permet de changer les _ du game.Word en lettres de toFind correspondantesS
	var indextofind []int
	for i := 0; i < len(game.ToFind); i++ {
		if rune(letter[0]) == game.ToFind[i] { // tableau qui stock les index où apparaît letter dans ToFind
			indextofind = append(indextofind, i)
		}
	}
	for j := 0; j < len(indextofind); j++ {
		game.Word[indextofind[j]] = game.ToFind[indextofind[j]]
	}
	game.StrWord = string(game.Word)
	return game
}

func RandomPickLetter(game *HangManData) []int {
	var randP []int
	var blacklist []int
	n := len(game.ToFind)/2 - 1
	for i := 0; i < n; i++ {
		nbUtils := RandomBlacklist(len(game.ToFind), blacklist)
		blacklist = append(blacklist, nbUtils)
		randP = append(randP, nbUtils)
	}
	return randP
}

func RandomBlacklist(max int, blacklisted []int) int {
	excluded := map[int]bool{}
	for _, x := range blacklisted {
		excluded[x] = true
	}
	for {
		n := rand.Intn(max)
		if !excluded[n] {
			return n
		}
	}
}

func ToHigher(s string) string { // Permet de mettre en MAJUSCULE une string
	modstr := []rune(s)
	for j := 0; j < len(s); j++ {
		if modstr[j] >= 97 && s[j] <= 122 {
			modstr[j] = modstr[j] - 32
		} else {
			continue
		}
	}
	s = string(modstr)
	return s
}

func ToLower(modstr []rune) []rune {
	var runetab []rune
	for j := 0; j < len(modstr); j++ {
		if modstr[j] <= 64 {
			j++
		} else {
			if modstr[j] >= 65 && modstr[j] <= 90 {
				runetab = append(runetab, modstr[j]+32)
			} else if modstr[j] >= 97 && modstr[j] <= 122 {
				runetab = append(runetab, modstr[j])
			} else {
				j++
			}
		}
	}
	return runetab
}
