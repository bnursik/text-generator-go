package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Prefix []string

func (p Prefix) String() string {
	return strings.Join(p, " ")
}

func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

type Chain struct {
	chain     map[string][]string
	prefixLen int
}

func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string][]string), prefixLen}
}

func (c *Chain) Build(r io.Reader, prefix string) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	found := false

	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}

		if prefix != " " && strings.Join(p, " ") == prefix {
			found = true
		}

		key := p.String()
		c.chain[key] = append(c.chain[key], s)
		p.Shift(s)
	}

	if prefix != " " && !found {
		fmt.Println("Given prefix is not found in the original text")
		os.Exit(1)
	}
}

func (c *Chain) Generate(n int, prefix string) string {
	p := make(Prefix, c.prefixLen)
	var words []string
	if prefix != " " {
		words = append(words, prefix)
		p = strings.Fields(prefix)

		for i := 0; i < n-len(p); i++ {
			suffixes := c.chain[p.String()]
			if len(suffixes) == 0 {
				break
			}

			nextWord := suffixes[rand.Intn(len(suffixes))]
			words = append(words, nextWord)
			p.Shift(nextWord)
		}
	} else {
		for i := 0; i < n; i++ {
			suffixes := c.chain[p.String()]
			if len(suffixes) == 0 {
				break
			}

			nextWord := suffixes[rand.Intn(len(suffixes))]
			words = append(words, nextWord)
			p.Shift(nextWord)
		}
	}

	return strings.Join(words, " ")
}

func main() {
	numWords := flag.Int("w", 100, "number of words to print")
	prefixLen := flag.Int("l", 2, "number of words in the prefix")
	prefix := flag.String("p", " ", "starting prefix")
	help := flag.Bool("help", false, "usage text")

	flag.Parse()

	rand.Seed(time.Now().UnixNano()) // updates the random seed everytime programm runs

	if *help {
		usage := `Markov Chain text generator.

Usage:
	markovchain [-w <N>] [-p <S>] [-l <N>]
	markovchain --help

Options:
	--help  Show this screen.
	-w N    Number of maximum words
	-p S    Starting prefix
	-l N    Prefix length`

		fmt.Println(usage)
		os.Exit(1)
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 { // to check if there is data piped in
		fmt.Println("Error: no input text")
		os.Exit(1)
	}

	if (*numWords <= 0) || (*numWords > 10000) {
		fmt.Println("Number of words to print should be in range [1:10,000]")
		os.Exit(1)
	}

	if (*prefixLen <= 0) || (*prefixLen > 5) {
		fmt.Println("Prefix length should be in range [1:5]")
		os.Exit(1)
	}

	if *prefix != " " && len(strings.Fields(*prefix)) != *prefixLen {
		fmt.Println("Length of given prefix must be equal to the prefix length(2 by default)")
		os.Exit(1)
	}

	if *prefix != " " && len(strings.Fields(*prefix)) > *numWords {
		fmt.Println("Length of given prefix must NOT be more than the number of words (100 by default)")
		os.Exit(1)
	}

	c := NewChain(*prefixLen)
	c.Build(os.Stdin, *prefix)
	text := c.Generate(*numWords, *prefix)
	fmt.Println(text)
}
