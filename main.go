package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

// Subword represents a tokenized subword with an ID.
type Subword struct {
	ID          int
	Piece       string
	Probability float64
}

// LoadCorpus reads a corpus from a file and returns it as a slice of strings.
func LoadCorpus(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var corpus []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		corpus = append(corpus, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return corpus, nil
}

// TrainUnigramModel trains a simple unigram-based subword model.
func TrainUnigramModel(corpus []string, vocabSize int) []Subword {
	wordCounts := make(map[string]int)

	// Count occurrences of words in the corpus
	for _, sentence := range corpus {
		words := strings.Split(sentence, " ")
		for _, word := range words {
			wordCounts[word]++
		}
	}

	// Initialize vocabulary with words and probabilities
	var vocab []Subword
	for word, count := range wordCounts {
		vocab = append(vocab, Subword{Piece: word, Probability: float64(count)})
	}

	// Normalize probabilities
	total := 0.0
	for _, v := range vocab {
		total += v.Probability
	}
	for i := range vocab {
		vocab[i].Probability /= total
	}

	// Iteratively merge least probable subwords
	rand.Seed(time.Now().UnixNano())
	for len(vocab) > vocabSize-1 { // Reserve space for <unk>
		// Sort by probability (ascending)
		sort.Slice(vocab, func(i, j int) bool {
			return vocab[i].Probability < vocab[j].Probability
		})

		// Merge the two least probable subwords
		merged := vocab[0].Piece + vocab[1].Piece
		mergedProbability := (vocab[0].Probability + vocab[1].Probability) / 2

		// Replace the first two with the merged subword
		vocab = append(vocab[2:], Subword{Piece: merged, Probability: mergedProbability})
	}

	// Add <unk> token to handle unknown words
	vocab = append(vocab, Subword{Piece: "<unk>", Probability: 0.001})

	// Assign integer IDs to subwords
	for i := range vocab {
		vocab[i].ID = i // Assign ID based on index
	}

	return vocab
}

// EncodeText tokenizes input text into integer token IDs, handling unknown words.
func EncodeText(text string, vocab []Subword) []int {
	tokens := []int{}
	remaining := text

	// Map subword pieces to their IDs for fast lookup
	pieceToID := make(map[string]int)
	unkID := -1
	for _, subword := range vocab {
		pieceToID[subword.Piece] = subword.ID
		if subword.Piece == "<unk>" {
			unkID = subword.ID
		}
	}

	for len(remaining) > 0 {
		match := ""
		matchID := unkID // Default to <unk>
		matchLength := 1 // Default to consuming one character for unknowns

		// Find the longest matching subword
		for _, subword := range vocab {
			if strings.HasPrefix(remaining, subword.Piece) {
				if len(subword.Piece) > len(match) {
					match = subword.Piece
					matchID = subword.ID
					matchLength = len(subword.Piece)
				}
			}
		}

		tokens = append(tokens, matchID)
		remaining = remaining[matchLength:] // Move forward by the length of the match
	}

	return tokens
}

// DecodeTokens converts token IDs back to text.
func DecodeTokens(tokenIDs []int, vocab []Subword) string {
	idToSubword := make(map[int]string)
	for _, subword := range vocab {
		idToSubword[subword.ID] = subword.Piece
	}

	var words []string
	for _, id := range tokenIDs {
		if word, found := idToSubword[id]; found {
			words = append(words, word)
		}
	}

	return strings.Join(words, "")
}

func main() {
	// Load corpus from a file
	corpus, err := LoadCorpus("corpus.txt")
	if err != nil {
		log.Fatalf("Error loading corpus: %v", err)
	}

	// Train a unigram model with a vocab size of 6 (including <unk>)
	vocab := TrainUnigramModel(corpus, 1000)
	fmt.Println("Trained Vocabulary:", vocab)

	// Encode a test sentence (includes unknown words)
	testSentence := "Merry Christmas, Marmee! Many of them! Thank you for our books; we read some, and mean to every day, they cried, in chorus."
	encoded := EncodeText(testSentence, vocab)
	fmt.Println("Encoded:", encoded)

	// Decode back to text
	decoded := DecodeTokens(encoded, vocab)
	fmt.Println("Decoded:", decoded)
}
