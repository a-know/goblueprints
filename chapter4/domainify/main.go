package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

// var tlds = []string{"com", "net"}
type strslice []string

func (s *strslice) String() string {
	return fmt.Sprintf("%v", tlds)
}

func (s *strslice) Set(v string) error {
	*s = append(*s, v)
	return nil
}

var tlds strslice

const allowedChars = "abcdefghijklmnopqrstuvwxyz0123456789-"

func main() {
	flag.Var(&tlds, "d", "Domain strings")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		text := strings.ToLower(s.Text())
		var newText []rune
		for _, r := range text {
			if unicode.IsSpace(r) {
				r = '-'
			}
			if !strings.ContainsRune(allowedChars, r) {
				continue
			}
			newText = append(newText, r)
		}
		fmt.Println(string(newText) + "." + tlds[rand.Intn(len(tlds))])
	}
}
