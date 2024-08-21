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

func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)

	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}

		key := p.String()
		c.chain[key] = append(c.chain[key], s)
		p.Shift(s)
	}
}

func (c *Chain) Generate(n int) string {
	p := make(Prefix, c.prefixLen)
	var words []string
	for i := 0; i < n; i++ {
		suffixes := c.chain[p.String()]
		if len(suffixes) == 0 {
			break
		}

		nextWord := suffixes[rand.Intn(len(suffixes))]
		words = append(words, nextWord)
		p.Shift(nextWord)
	}

	return strings.Join(words, " ")
}

func main() {
	numWords := flag.Int("w", 100, "number of words to print")
	prefixLen := flag.Int("l", 2, "number of words in the prefix")

	flag.Parse()

	rand.Seed(time.Now().UnixNano()) //updates the random seed everytime programm runs

	c := NewChain(*prefixLen)
	c.Build(os.Stdin)
	text := c.Generate(*numWords)
	fmt.Println(text)
}
