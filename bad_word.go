package main

import (
	"strings"
)

var listOfBadWords = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

func BadWorkReplacement(message string) string {
	wordsInMessage := strings.Fields(message)
	newMessage := ""
	for i, word := range wordsInMessage {
		for _, badWord := range listOfBadWords {
			if strings.ToLower(word) == strings.ToLower(badWord) {
				word = "****"
			}
		}

		newMessage += word
		if i < len(wordsInMessage)-1 {
			newMessage += " "
		}
	}

	return newMessage
}
