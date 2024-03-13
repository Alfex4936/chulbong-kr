package utils

import (
	"log"
	"testing"
)

func TestCheckForBadWords(t *testing.T) {
	if err := LoadBadWords("../badwords.txt"); err != nil {
		log.Fatalf("Failed to load bad words: %v", err)
	}

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "contains bad word",
			input: "시잇ㅄ발",
			want:  true,
		},
		{
			name:  "contains no bad word",
			input: "좋아요 굿!!!",
			want:  false,
		},
		{
			name:  "contains multiple bad words",
			input: "뭐라노 ㅅㅂ",
			want:  true,
		},
		{
			name:  "contains no bad word 2",
			input: "ㅋㅋㅋㅋㅋㅋ너무좋아요 근데 이게 뭐라고 ~~",
			want:  false,
		},
		{
			name:  "contains bad word with punctuation",
			input: "아닠ㅋㅋ병신?",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckForBadWords(tt.input)
			if err != nil {
				t.Errorf("CheckForBadWords() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("CheckForBadWords() got = %v, want %v", got, tt.want)
			}
		})
	}
}
