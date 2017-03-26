package main

import (
	"bufio"
	"os"

	"github.com/pkg/profile"
)

const firstLetter int = 'a'
const lastLetter int = 'z'
const nLetters int = lastLetter - firstLetter + 1

type anagram struct {
	word string
	next *anagram
}

type node struct {
	firstAnagram string
	nextAnagram  *anagram
	children     [nLetters]*node
}

func makeNode() *node {
	return &node{}
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	defer profile.Start(profile.ProfilePath(".")).Stop()
	process(os.Args[1], os.Args[2])
}

func process(inputFilename, outputFilname string) {
	n := makeNode()
	file, err := os.Open(inputFilename)
	handleErr(err)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		n.add(scanner.Text())
	}
	file.Close()

	file, err = os.Create(outputFilname)
	handleErr(err)
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	n.write(writer)
}

func (n *node) add(word string) {
	sorted := sort(word)
	n = n.search(sorted)
	n.addValue(word)
}

func (n *node) addValue(word string) {
	if n.firstAnagram == "" {
		n.firstAnagram = word
	} else {
		n.nextAnagram = &anagram{
			word: word,
			next: n.nextAnagram,
		}
	}
}

func (n *node) search(sorted [nLetters]int) *node {
	for i, nChars := range sorted {
		for ; nChars != 0; nChars-- {
			childNode := n.children[i]
			if childNode == nil {
				childNode = makeNode()
				n.children[i] = childNode
			}
			n = childNode
		}
	}
	return n
}

func sort(word string) (sorted [nLetters]int) {
	for _, r := range word {
		sorted[int(r)-firstLetter]++
	}
	return
}

func (n *node) write(writer *bufio.Writer) {
	if n.firstAnagram != "" {
		writer.WriteString(n.firstAnagram)
		a := n.nextAnagram
		if a != nil {
			for ; a.next != nil; a = a.next {
				writer.WriteString(a.word)
				writer.WriteRune(' ')
			}
			if a.word != "" {
				writer.WriteString(a.word)
			}
		}
		writer.WriteRune('\n')
	}

	for _, childN := range n.children {
		if childN != nil {
			childN.write(writer)
		}
	}
}
