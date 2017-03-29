package main

import (
	"bufio"
	"os"
	"runtime/debug"
)

const firstLetter int = 'a'
const lastLetter int = 'z'
const nLetters int = lastLetter - firstLetter + 1

type anagram struct {
	word []byte
	next *anagram
}

type node struct {
	anagrams *anagram
	children [nLetters]*node
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
	debug.SetGCPercent(-1)
	// defer profile.Start(profile.ProfilePath(".")).Stop()
	process(os.Args[1], os.Args[2])
}

func process(inputFilename, outputFilname string) {
	n := makeNode()
	file, err := os.Open(inputFilename)
	handleErr(err)

	stat, err := file.Stat()
	handleErr(err)
	// use reader size equal to filesize, so that bytes are not ovewritten when
	// they are read
	reader := bufio.NewReaderSize(file, int(stat.Size()))
	// use reader instead of scanner so we can set buffer size and don't
	// have to copy bytes
	for line, isPrefix, err := reader.ReadLine(); len(line) > 0; line, isPrefix, err = reader.ReadLine() {
		handleErr(err)
		if isPrefix {
			panic("is prefix!")
		}
		n.add(line)
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
	n.anagrams = &anagram{
		word: word,
		next: n.anagrams,
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
	firstAnagram := n.anagrams
	if firstAnagram != nil {
		writer.Write(firstAnagram.word)
		for a := firstAnagram.next; a != nil; a = a.next {
			writer.WriteRune(' ')
			writer.Write(a.word)
		}

		writer.WriteRune('\n')
	}

	for _, childN := range n.children {
		if childN != nil {
			childN.write(writer)
		}
	}
}
