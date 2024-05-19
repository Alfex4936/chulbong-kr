// +build amd64

#include "textflag.h"

// func CheckProfanity(b *BadWordUtil, data []byte) bool
TEXT ·CheckProfanity(SB), NOSPLIT, $0-32
    // Load parameters
    MOVQ data+0(FP), SI         // SI = data pointer
    MOVQ data+8(FP), CX         // CX = len(data)
    MOVQ b+16(FP), DI           // DI = BadWordUtil pointer

    // Load the address of the bad words slice
    MOVQ (DI), BX               // BX = BadWordsListByte slice pointer
    MOVQ 8(DI), R8              // R8 = len(BadWordsListByte)

    // Initialize return value (no profanity found)
    XORQ AX, AX                 // AX = 0

check_next_word:
    // Check if we are done with all bad words
    TESTQ R8, R8
    JZ done

    // Load current bad word from struct
    MOVQ (BX), DX               // DX = BadWordsListByte[i] data pointer
    MOVQ 8(BX), R10             // R10 = len(BadWordsListByte[i])
    ADDQ $16, BX                // BX = next BadWordsListByte element
    DECQ R8                     // R8--

    // Check each position in data
    XORQ R11, R11               // R11 = j (index in data)
check_next_position:
    CMPQ R11, CX                // if j >= len(data)
    JAE check_next_word         //   continue to next bad word

    // Check if remaining data is less than bad word length
    MOVQ CX, R12                // R12 = len(data)
    SUBQ R11, R12               // R12 = len(data) - j
    CMPQ R12, R10               // if len(data) - j < len(BadWordsListByte[i])
    JB check_next_word          //   continue to next bad word

    // Compare substring of data with current bad word
    XORQ R12, R12               // R12 = k (index in BadWordsListByte[i])
compare_loop:
    CMPQ R12, R10               // if k >= len(BadWordsListByte[i])
    JAE found_profanity         //   found a match

    // Calculate effective addresses for data[j + k] and BadWordsListByte[i][k]
    MOVQ SI, R13                // R13 = SI (data pointer)
    ADDQ R11, R13               // R13 = data + j
    ADDQ R12, R13               // R13 = data + j + k

    MOVQ DX, R14                // R14 = DX (BadWordsListByte[i] pointer)
    ADDQ R12, R14               // R14 = BadWordsListByte[i] + k

    MOVB (R13), AL              // AL = data[j + k]
    MOVB (R14), BL              // BL = BadWordsListByte[i][k]
    CMPB AL, BL                 // if data[j + k] != BadWordsListByte[i][k]
    JNE no_match                //   continue to next position

    INCQ R12                    // k++
    JMP compare_loop

no_match:
    INCQ R11                    // j++
    JMP check_next_position

found_profanity:
    MOVQ $1, AX                 // AX = 1 (profanity found)

done:
    RET
