package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

func main() {
	fmt.Print("Enter filename:")
	var filename string

	fmt.Scanln(&filename)

	file, err := os.Open(filename)

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	defer file.Close()

	stopwords := map[string]bool{
		"the": true, "and": true, "in": true, "of": true, "to": true,
		"a": true, "an": true, "is": true, "it": true, "that": true,
		"for": true, "on": true, "with": true, "as": true, "at": true,
	}

	wordCount := make(map[string]int)
	scanner := bufio.NewScanner(file)

	re := regexp.MustCompile(`[A-Za-z0-9]+`)

	for scanner.Scan() {
		line := scanner.Text()

		words := re.FindAllString(line, -1)

		for _, w := range words {
			word := strings.ToLower(w)

			if stopwords[word] {
				continue
			}
			wordCount[word]++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	type pair struct {
		Word  string
		Count int
	}

	var freqList []pair
	for w, c := range wordCount {
		freqList = append(freqList, pair{Word: w, Count: c})
	}

	sort.Slice(freqList, func(i, j int) bool {
		return freqList[i].Count > freqList[j].Count
	})

	fmt.Println("\nWord frequencies:")

	for _, p := range freqList {
		fmt.Printf("%s: %d\n", p.Word, p.Count)
	}
}
