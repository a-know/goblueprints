package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/a-know/goblueprints/chapter4/thesaurus"
)

func main() {
	apiKey := os.Getenv("BHT_APIKEY")
	thesaurus := &thesaurus.BigHuge{APIKey: apiKey}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := s.Text()
		syns, err := thesaurus.Synonyms(word)
		if err != nil {
			log.Fatalf("Failed to search synonym of %q : %v\n", word, err)
		}
		if len(syns) == 0 {
			log.Fatalf("There are no synonyms with %q\n", word)
		}
		for _, syn := range syns {
			fmt.Println(syn)
		}
	}
}
