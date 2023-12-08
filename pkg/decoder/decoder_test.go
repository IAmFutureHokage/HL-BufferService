package decoder

import (
	"reflect"
	"testing"
	// types "github.com/IAmFutureHokage/HL-Coder/pkg/types"
)

func TestParseString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "Normal case",
			input: "word1= word2= word3=",
			want:  []string{"word1", "word2", "word3"},
		},
		{
			name:  "No equals sign",
			input: "word1 word2 word3",
			want:  []string{"word1", "word2", "word3"},
		},
		{
			name:  "Mixed case",
			input: "word1= word2 word3=",
			want:  []string{"word1", "word2", "word3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseString(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitSequence(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "Normal case",
			input: "block1 block2 922xx block3 block4",
			want:  []string{"block1 block2", "block1 xx081 block3 block4"},
		},
		{
			name:  "No 922 prefix",
			input: "block1 block2 block3",
			want:  []string{"block1 block2 block3"},
		},
		{
			name:  "Multiple 922 prefixes",
			input: "block1 922xx block2 922yy block3",
			want:  []string{"block1 xx081 block2", "block1 yy081 block3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitSequence(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitSequence() = %v, want %v", got, tt.want)
			}
		})
	}
}
