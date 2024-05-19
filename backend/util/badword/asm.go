//go:build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

func main() {
	Package("github.com/Alfex4936/chulbong-kr/util")

	TEXT("CheckProfanity", NOSPLIT, "func(b *BadWordUtil, data []byte) bool")
	Doc("CheckProfanity checks if the input data contains any bad words using SIMD instructions.")

	// Load parameters
	b := Dereference(Param("b"))
	data := Load(Param("data").Base(), GP64())
	dataLen := Load(Param("data").Len(), GP64())

	// Load BadWordsListByte from the struct
	badWordsList := b.Field("BadWordsListByte")
	badWords := Load(badWordsList.Base(), GP64())
	badWordsLen := Load(badWordsList.Len(), GP64())

	// Initialize return value (no profanity found)
	found := GP64()
	XORQ(found, found) // found = 0

	Label("check_next_word")
	// Check if we are done with all bad words
	TESTQ(badWordsLen, badWordsLen)
	JZ(LabelRef("done"))

	// Load current bad word from struct
	currentWord := GP64()
	currentWordLen := GP64()
	MOVQ(Mem{Base: badWords}, currentWord)
	MOVQ(Mem{Base: badWords, Disp: 8}, currentWordLen)
	ADDQ(U32(16), badWords)
	DECQ(badWordsLen)

	// Check each position in data
	position := GP64()
	XORQ(position, position) // position = 0
	
	Label("check_next_position")
	CMPQ(position, dataLen)
	JAE(LabelRef("check_next_word"))

	// Load substring of data
	X0 := XMM()
	MOVOU(Mem{Base: data, Index: position, Scale: 1}, X0)

	// Compare with current bad word
	X1 := XMM()
	MOVOU(Mem{Base: currentWord}, X1)
	PCMPEQB(X1, X0)   // Compare bytes
	R11 := GP32()     // PMOVMSKB extracts 32 bits, use GP32 register
	PMOVMSKB(X1, R11) // Extract comparison result

	// Check if any byte matched
	TESTL(R11, R11) // TESTL because R11 is GP32
	JNZ(LabelRef("found_profanity"))

	INCQ(position)
	JMP(LabelRef("check_next_position"))

	Label("found_profanity")
	MOVQ(U32(1), found)

	Label("done")
	Store(found, ReturnIndex(0)) // Store to return index 0
	RET()
	Generate()
}
