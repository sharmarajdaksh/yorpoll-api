package poll

import (
	"reflect"
	"testing"
	"time"
)

var optionTests = []struct {
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

var pollTests = []struct {
	title       string
	description string
	expiry      int64
}{
	{
		title:       "first poll",
		description: "generic descr for the first poll",
		expiry:      999999,
	},
	{
		title:       "second poll with a silghtly long title",
		description: "generic description for the second poll",
		expiry:      19827,
	},
	{
		title:       "3poll",
		description: "generic description for the first poll",
		expiry:      16843,
	},
}

func TestNewOption(t *testing.T) {
	for _, tt := range optionTests {
		t.Run(tt.title, func(t *testing.T) {
			opt := NewOption(tt.title)
			if opt.Title != tt.title {
				t.Errorf("Option title does not match")
			}
			if opt.Votes != 0 {
				t.Errorf("Option initialized with non-zero votes")
			}
		})
	}
}

func TestNewPoll(t *testing.T) {
	var opts = []Option{}
	for _, opt := range optionTests {
		opts = append(opts, NewOption(opt.title))
	}

	for _, tt := range pollTests {
		t.Run(tt.title, func(t *testing.T) {
			p := New(tt.title, tt.description, opts, tt.expiry)
			if p.Title != tt.title {
				t.Errorf("Poll title does not match")
			}
			if p.Description != tt.description {
				t.Errorf("Poll title does not match")
			}
			if p.Created > time.Now().Unix() {
				t.Errorf("Created timestamp is in the future")
			}
			if p.Created != p.Modified {
				t.Errorf("Created and modified timestamp do not match for new poll")
			}
			if p.Expiry != tt.expiry {
				t.Errorf("Poll expiry does not match")
			}
			if !reflect.DeepEqual(p.Options, opts) {
				t.Errorf("Options do not match")
			}
		})
	}
}
