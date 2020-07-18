package main

import (
	"github.com/imonke/monkebase"
	"github.com/imonke/monketype"

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
	monkebase.Connect(os.Getenv("MONKEBASE_CONNECTION"))

	var result int = main.Run()
	monkebase.EmptyTable(monkebase.USER_TABLE)
	os.Exit(result)
}

// honestly fuck these tests they're so ugly and gross idc anymore
func Test_feedAll(test *testing.T) {
	var query map[string]int
	var queries []map[string]int = []map[string]int{
		map[string]int{
			"offset": 0,
			"size":   10,
		},
		map[string]int{
			"offset": 10,
			"size":   20,
		},
	}

	monkebase.EmptyTable(monkebase.CONTENT_TABLE)
	seed(40)
	defer monkebase.EmptyTable(monkebase.CONTENT_TABLE)

	var request *http.Request
	var fetched []monketype.Content
	var code, size int
	var r_map map[string]interface{}
	var err error

	for _, query = range queries {
		request = blank.WithContext(
			context.WithValue(
				blank.Context(),
				"parsed_query",
				query,
			),
		)

		if code, r_map, err = feedAll(request); err != nil {
			test.Fatal(err)
		}

		if code != 200 {
			test.Errorf("got code %d", code)
		}

		size = r_map["size"].(int)
		fetched = r_map["content"].([]monketype.Content)

		if size != query["size"] {
			test.Errorf("wanted size mismatch! have: %d, want: %d", size, query["size"])
		}

		if len(fetched) != query["size"] {
			test.Errorf("actual size mismatch! have: %d, want: %d", len(fetched), query["size"])
		}

		sequenceOK(test, fetched)
	}
}

func Test_feedAll_some(test *testing.T) {
	var population = 20
	var query map[string]int = map[string]int{
		"offset": 10,
		"size":   20,
	}

	var projected = population - query["offset"]

	monkebase.EmptyTable(monkebase.CONTENT_TABLE)
	seed(population)
	defer monkebase.EmptyTable(monkebase.CONTENT_TABLE)

	var request *http.Request = blank.WithContext(
		context.WithValue(
			blank.Context(),
			"parsed_query",
			query,
		),
	)

	var code int
	var r_map map[string]interface{}
	var err error
	if code, r_map, err = feedAll(request); err != nil {
		test.Fatal(err)
	}

	if code != 200 {
		test.Errorf("got code %d", code)
	}

	var size int = r_map["size"].(int)
	var fetched []monketype.Content = r_map["content"].([]monketype.Content)

	if size != projected {
		test.Errorf("wanted size mismatch! have: %d, want: %d", projected, size)
	}

	if len(fetched) != projected {
		test.Errorf("actual size mismatch! have: %d, want: %d", len(fetched), projected)
	}

	sequenceOK(test, fetched)
}
