package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const otherWord = "*"

var transforms = loadTransformsFromFile()

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		t := transforms[rand.Intn(len(transforms))]
		fmt.Println(strings.Replace(t, otherWord, s.Text(), -1))
	}
}

func loadTransformsFromFile() []string {
	f, err := os.Open("transforms.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "File could not read: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	content := make([]string, 0, 100)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		content = append(content, scanner.Text())
	}
	if serr := scanner.Err(); serr != nil {
		fmt.Fprintf(os.Stderr, "File scan error: %v\n", err)
	}

	return content
}
