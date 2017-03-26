package main

import (
	"bufio"
	"os"
)

const firstLetter int = 'a'
const lastLetter int = 'z'
const nLetters int = lastLetter - firstLetter + 1

type anagram struct {
	word []byte
	next *anagram
}

type node struct {
	firstAnagram []byte
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
	process(os.Args[1], os.Args[2])
}

func process(inputFilename, outputFilname string) {
	n := makeNode()
	file, err := os.Open(inputFilename)
	handleErr(err)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// must copy bytes, because can overwrite
		b := scanner.Bytes()
		bNew := make([]byte, len(b))
		copy(bNew, b)
		n.add(bNew)
	}
	file.Close()

	file, err = os.Create(outputFilname)
	handleErr(err)
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	n.write(writer)
}

func (n *node) add(word []byte) {
	sorted := sort(word)
	n = n.search(sorted)
	n.addValue(word)
}

func (n *node) addValue(word []byte) {
	if len(n.firstAnagram) == 0 {
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

func sort(word []byte) (sorted [nLetters]int) {
	for _, r := range word {
		sorted[int(r)-firstLetter]++
	}
	return
}

func (n *node) write(writer *bufio.Writer) {
	if len(n.firstAnagram) > 0 {
		writer.Write(n.firstAnagram)
		for a := n.nextAnagram; a != nil; a = a.next {
			if len(a.word) > 0 {
				writer.WriteRune(' ')
				writer.Write(a.word)
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
