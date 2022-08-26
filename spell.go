package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/eskriett/spell"
)

var (
	regex = regexp.MustCompile(`\w+`)
)

func load() {
	sp, err := spell.Load("dict.spell")
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("freaking.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Split(scanner.Text(), " ")
		for _, v := range words {
			if len(v) < 3 {
				continue
			}
			v = strings.ToLower(regex.FindString(v))
			suggestions, err := sp.Lookup(v)
			if err != nil {
				log.Println(err)
			}
			if entry, _ := sp.GetEntry(v); entry == nil {
				if len(suggestions) > 0 && suggestions.GetWords()[0] != v {
					fmt.Println("Suggestions for", v, ":", suggestions.String())
				}
			}
		}
	}
}

func loadFromDictionary() {
	sp := spell.New()

	file, err := os.Open("en-80k.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		word := line[0]
		freq, _ := strconv.Atoi(line[1])

		_, err := sp.AddEntry(spell.Entry{
			Word:      word,
			Frequency: uint64(freq),
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(sp.GetLongestWord())
	err = sp.Save("dict.spell")
	if err != nil {
		log.Fatal(err)
	}
}
