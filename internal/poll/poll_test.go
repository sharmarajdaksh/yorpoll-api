package poll

import (
	"testing"
)

var tests = []struct {
	title string
}{
	{
		title: "testTitle 1",
	},
	{
		title: "hjsjadkjsdbajk saashdiha dsadbnjka sd",
	},
	{
		title: "`jhasdhkasnd ;;ad 90u21h  321097u 213921 39u 213ih;n21nnmn 213",
	},
}

func TestNewOptionTitle(t *testing.T) {

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			opt := NewOption(tt.title)
			if opt.Title != tt.title {
				t.Errorf("Option title does not match")
			}
		})
	}
}

func TestNewOptionUniqueIDs(t *testing.T) {

	seen := []string{}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			opt := NewOption(tt.title)
			for _, id := range seen {
				if opt.ID == id {
					t.Errorf("Option ID is not unique")
				}
			}
		})
	}
}
