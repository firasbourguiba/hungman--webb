package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func Random(difficulty string) string {
	file := ""
	word := ""
	var index int
	if difficulty == "EASY" {
		index = rand.Intn(37) + 1
		file = "words.txt"
	} else if difficulty == "NORMAL" {
		index = rand.Intn(23) + 1
		file = "words2.txt"
	} else if difficulty == "HARD" {
		index = rand.Intn(24) + 1
		file = "words3.txt"
	} else {
		fmt.Println("Erreur")
		os.Exit(1)
	}
	fileOpen, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	i := 0
	scanner := bufio.NewScanner(fileOpen)
	for scanner.Scan() {
		i++
		if i == index {
			word = scanner.Text()
		}
	}
	return word
}