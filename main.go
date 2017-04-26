package main

import (
	"bufio"
	"os"
	"runtime/debug"
)

// we only handle a-z
const firstLetter int = 'a'
const lastLetter int = 'z'
const nLetters int = lastLetter - firstLetter + 1

// node is part of a trie that holds one class of anagrams
type node struct {
	// anagrams holds the byte representation of the anagrams joined by spaces
	anagrams []byte
	// children points to all anagram classes that have the same prefix of
	// ordered characters as this one, along with another letter.
	// that letter, offset from "a", indexes this array
	children [nLetters]*node
}

// makeNode creates a blank node
func makeNode() *node {
	return &node{}
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// turn of garbage collection, since we memory isn't an issue here
	debug.SetGCPercent(-1)

	// defer profile.Start(profile.ProfilePath(".")).Stop()
	process(os.Args[1], os.Args[2])
}

func process(inputFilename, outputFilname string) {
	file, err := os.Open(inputFilename)
	handleErr(err)

	scanner := bufio.NewScanner(file)
	n := makeNode()
	for scanner.Scan() {
		// copy bytes, because it reuses the bytes
		n.add(append([]byte{}, scanner.Bytes()...))
	}
	file.Close()

	file, err = os.Create(outputFilname)
	handleErr(err)
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	n.write(writer)
}

// add adds a word to the trie, by first finding the node where it should
// below, using it's sorted letters, then adding it there
func (n *node) add(word []byte) {
	sorted := sort(word)
	n = n.search(sorted)
	n.anagrams = concatWords(n.anagrams, word)
}

// search returns the the node corresponding to the sorted letters,
// assuming n is the root
func (n *node) search(sorted [nLetters]int) *node {
	for i, nChars := range sorted {
		for ; nChars != 0; nChars-- {
			n = n.child(i)
		}
	}
	return n
}

// sort returns an array holding the number of each of the characters
// in a word. The first item in the array holds the number of "a"s, the second
// "b"s, etc.
func sort(word []byte) (sorted [nLetters]int) {
	for _, r := range word {
		sorted[int(r)-firstLetter]++
	}
	return
}

// child returns the child node associated with the ith characters
// creating it if it doesn't exist
func (n *node) child(i int) (childNode *node) {
	childNode = n.children[i]
	if childNode == nil {
		childNode = makeNode()
		n.children[i] = childNode
	}
	return
}

// write prints out all anagram classes on seperate lines
// It uses a depth first traversal of the trie, printing each set of anagrams
// on a line, if there are word on that node
func (n *node) write(writer *bufio.Writer) {
	if len(n.anagrams) > 0 {
		writer.Write(n.anagrams)
		writer.WriteRune('\n')
	}

	for _, childN := range n.children {
		if childN != nil {
			childN.write(writer)
		}
	}
}

// concatWords som existing words, with another word, with a space in the
// middle. if there are not existing words, just returns the new word
func concatWords(words, word []byte) []byte {
	if len(words) == 0 {
		return word
	}
	return append(append(words, byte(' ')), word...)
}
