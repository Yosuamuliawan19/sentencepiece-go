SentencePiece Encoding in Golang

This is a simple SentencePiece-style subword tokenization implementation in Go. It trains a unigram-based model from a given corpus and encodes text into integer token IDs, handling unknown words (<unk>).

Features
• Train a unigram-based vocabulary from a text corpus
• Encode text into integer token IDs
• Handle unknown words with <unk> token
• Decode tokens back to human-readable text
• Prevent infinite loops by ensuring forward progress

Project Structure

📦 sentencepiece-go
├── 📄 corpus.txt # Sample corpus file
├── 📄 sentencepiece.go # Main implementation
├── 📄 README.md # Project documentation

Installation & Usage

1. Install Go

Ensure you have Go installed (version 1.18+ recommended).
Check with:

go version

2. Clone the Repository

git clone https://github.com/YOUR_USERNAME/sentencepiece-go.git
cd sentencepiece-go

3. Prepare a Corpus File

Create a corpus.txt file with sample sentences:

nano corpus.txt

Example content:

hello world
hello sentencepiece
hello encoding
world is beautiful
natural language processing is fun

Save (CTRL+X, then Y, then Enter).

4. Run the Tokenizer

go run sentencepiece.go

Sample Output

Trained Vocabulary: [{ID: 0, Piece: "hello", Probability: 0.214} {ID: 1, Piece: "world", Probability: 0.143} {ID: 2, Piece: "is", Probability: 0.143} {ID: 3, Piece: "naturallanguageprocessing", Probability: 0.071} {ID: 4, Piece: "funsentencepieceencodingbeautiful", Probability: 0.071} {ID: 5, Piece: "<unk>", Probability: 0.001}]
Encoded: [0 1 5]
Decoded: hello world <unk>

Explanation:
• "hello" → ID 0
• "world" → ID 1
• "unknownword" → <unk> (ID 5)

Code Overview

Training the Model

The TrainUnigramModel() function: 1. Reads the corpus and counts word frequencies. 2. Merges least probable subwords until vocabSize is reached. 3. Adds <unk> token for unknown words.

Encoding Text
• Finds the longest matching subword in the vocabulary.
• If no match is found, assigns <unk> and moves forward character by character to avoid infinite loops.

Decoding Tokens
• Converts integer token IDs back into words.
• Handles <unk> tokens properly.

Next Improvements
• Implement Byte-Pair Encoding (BPE) for better subword merging
• Save & Load pretrained vocab for reuse
• Implement subword splitting for finer granularity

Contributing

Contributions are welcome! Feel free to fork, improve, and submit a PR.

License

This project is open-source under the MIT License.
