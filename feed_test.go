package main

import (
	"git.gastrodon.io/imonke/monkebase"
	"git.gastrodon.io/imonke/monketype"

	"context"
	"net/http"
	"net/url"
	"os"
	"testing"
)

type querySet struct {
	URL    string
	Size   int
	Offset int
	Code   int
	OK     bool
}

const (
	nick  = "zero"
	email = "mail@imonke.io"
)

var (
	blank *http.Request  = new(http.Request)
	user  monketype.User = monketype.NewUser(nick, "", email)
)

func seed(size int) {
	for size != 0 {
		monkebase.WriteContent(monketype.NewContent("", user.ID, "png", nil, true, false).Map())
		size--
	}
}

func urlMustParse(it string) (parsed *url.URL) {
	var err error
	if parsed, err = url.Parse(it); err != nil {
		panic(err)
	}

	return
}

func sequenceOK(test *testing.T, content []monketype.Content) {
	var index int
	var it monketype.Content

	for index, it = range content[1:] {
		if content[index].Created < it.Created {
			test.Errorf(
				"feed is out of order! %s at %d -> %s at %d",
				content[index].ID,
				content[index].Created,
				it.ID,
				it.Created,
			)
		}
	}
}

func TestMain(main *testing.M) {
	monkebase.Connect(os.Getenv("DATABASE_CONNECTION"))
	seed(100)

	var result int = main.Run()
	monkebase.EmptyTable(monkebase.CONTENT_TABLE)
	os.Exit(result)
}

func Test_feedAll(test *testing.T) {
	var targets [5]int = [5]int{10, 20, 30, 40, 50}
	var content []monketype.Content
	var request *http.Request
	var code, size, target int
	var r_map map[string]interface{}
	var err error
	for _, target = range targets {
		request = blank.WithContext(
			context.WithValue(
				blank.Context(),
				"query",
				map[string]interface{}{"size": target, "before": ""},
			),
		)

		if code, r_map, err = feedAll(request); err != nil {
			test.Fatal(err)
		}

		if code != 200 {
			test.Errorf("bad code %d", code)
		}

		size = r_map["size"].(map[string]int)["content"]
		content = r_map["content"].([]monketype.Content)

		if size != target {
			test.Errorf("bad reported size %d, want: %d", size, target)
		}

		if len(content) != target {
			test.Errorf("bad actual size %d, want: %d", len(content), target)
		}

		sequenceOK(test, content)
	}
}

func Test_feedAll_after(test *testing.T) {
	var target, offset int = 50, 11
	var first, second []monketype.Content
	var request *http.Request
	var r_map map[string]interface{}
	var err error

	request = blank.WithContext(
		context.WithValue(
			blank.Context(),
			"query",
			map[string]interface{}{"size": target, "before": ""},
		),
	)

	if _, r_map, err = feedAll(request); err != nil {
		test.Fatal(err)
	}

	first = r_map["content"].([]monketype.Content)

	request = blank.WithContext(
		context.WithValue(
			blank.Context(),
			"query",
			map[string]interface{}{"size": target, "before": first[offset].ID},
		),
	)

	if _, r_map, err = feedAll(request); err != nil {
		test.Fatal(err)
	}

	second = r_map["content"].([]monketype.Content)

	sequenceOK(test, first)
	sequenceOK(test, second)

	var index int
	var content monketype.Content
	for index, content = range first[offset+1:] {
		if content.ID != second[index].ID {
			test.Errorf("IDs are not aligned at %d: %s != %s", index, content.ID, second[index].ID)
		}
	}
}
