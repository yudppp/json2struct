package json2struct

import (
	"strings"
	"testing"
)

type TestCase struct {
	Input       string
	InputOption Options
	// MEMO: replace backquote to singlequote
	Expected string
}

var testCases = []TestCase{
	TestCase{
		Input: `{"text": "hello"}`,
		Expected: `type Data struct {
	Text string 'json:"text"'
}`,
	},
	TestCase{
		Input:       `[{"id": 123}]`,
		InputOption: Options{Name: "Categories"},
		Expected: `type Category struct {
	ID int 'json:"id"'
}`,
	},
	TestCase{
		Input:       `[{"name": 123},{"name": 3.14}]`,
		InputOption: Options{Name: "Tags"},
		Expected: `type Tag struct {
	Name float64 'json:"name"'
}`,
	},
	TestCase{
		Input:       `[{"name": 123},{"name": "stringer"}]`,
		InputOption: Options{Name: "Tags"},
		Expected: `type Tag struct {
	Name interface{} 'json:"name"'
}`,
	},
	TestCase{
		Input: `{"nest": {"text": "hello"}}`,
		Expected: `type Data struct {
	Nest DataNest 'json:"nest"'
}

type DataNest struct {
	Text string 'json:"text"'
}`,
	},
	TestCase{
		Input:       `{"nest": {"text": "hello"}}`,
		InputOption: Options{UseShortStruct: true},
		Expected: `type Data struct {
	Nest Nest 'json:"nest"'
}

type Nest struct {
	Text string 'json:"text"'
}`,
	},
	TestCase{
		Input:       `{"nest": {"text": "hello"}}`,
		InputOption: Options{UseLocal: true},
		Expected: `type data struct {
	nest dataNest 'json:"nest"'
}

type dataNest struct {
	text string 'json:"text"'
}`,
	},
	TestCase{
		Input:       `{"nest": {"text": "hello"}}`,
		InputOption: Options{UseOmitempty: true},
		Expected: `type Data struct {
	Nest *DataNest 'json:"nest,omitempty"'
}

type DataNest struct {
	Text string 'json:"text,omitempty"'
}`,
	},
	TestCase{
		Input:       `{"nest": {"text": "hello"}}`,
		InputOption: Options{Prefix: "input"},
		Expected: `type InputData struct {
	Nest InputDataNest 'json:"nest"'
}

type InputDataNest struct {
	Text string 'json:"text"'
}`,
	},
	TestCase{
		Input:       `{"nest": {"text": "hello"}}`,
		InputOption: Options{Suffix: "result"},
		Expected: `type DataResult struct {
	Nest DataNestResult 'json:"nest"'
}

type DataNestResult struct {
	Text string 'json:"text"'
}`,
	},
	TestCase{
		Input: `{"categories": [{}]}`,
		Expected: `type Data struct {
	Categories []DataCategory 'json:"categories"'
}

type DataCategory struct {
}`,
	},
	TestCase{
		Input: `{"categories": []}`,
		Expected: `type Data struct {
	Categories []interface{} 'json:"categories"'
}`,
	},
	TestCase{
		Input: `{"categories": [1]}`,
		Expected: `type Data struct {
	Categories []int 'json:"categories"'
}`,
	},
	TestCase{
		Input: `{"categories": [1,"abc"]}`,
		Expected: `type Data struct {
	Categories []interface{} 'json:"categories"'
}`,
	},
	TestCase{
		Input: `{"categories": null}`,
		Expected: `type Data struct {
	Categories interface{} 'json:"categories"'
}`,
	},
	TestCase{
		Input: `{"post": {"status": 1, "accept_comment": true, "title": "hello world", "tags": [1,2,4]}, "categories": [{"name": "aws", "num": 123}, {"name": 123, "num": 3.14}], "url": "http://blog.yudppp.com", "profile_image_url": "http://blog.yudppp.com/img/profile.gif", "comments": []}`,
		Expected: `type Data struct {
	Categories      []DataCategory 'json:"categories"'
	Comments        []interface{}  'json:"comments"'
	Post            DataPost       'json:"post"'
	ProfileImageURL string         'json:"profile_image_url"'
	URL             string         'json:"url"'
}

type DataCategory struct {
	Name interface{} 'json:"name"'
	Num  float64     'json:"num"'
}

type DataPost struct {
	AcceptComment bool   'json:"accept_comment"'
	Status        int    'json:"status"'
	Tags          []int  'json:"tags"'
	Title         string 'json:"title"'
}`,
	},
	TestCase{
		Input:       `[{"post": {"status": 1, "accept_comment": true, "title": "hello world", "tags": [1,2,4]}, "categories": [{"name": "aws", "num": 123}, {"name": 123, "num": 3.14}], "url": "http://blog.yudppp.com", "profile_image_url": "http://blog.yudppp.com/img/profile.gif", "comments": []}]`,
		InputOption: Options{UseShortStruct: true, Name: "json"},
		Expected: `type JSON struct {
	Categories      []Category    'json:"categories"'
	Comments        []interface{} 'json:"comments"'
	Post            Post          'json:"post"'
	ProfileImageURL string        'json:"profile_image_url"'
	URL             string        'json:"url"'
}

type Category struct {
	Name interface{} 'json:"name"'
	Num  float64     'json:"num"'
}

type Post struct {
	AcceptComment bool   'json:"accept_comment"'
	Status        int    'json:"status"'
	Tags          []int  'json:"tags"'
	Title         string 'json:"title"'
}`,
	},
}

func TestParse(t *testing.T) {
	for _, v := range testCases {
		expected := strings.Replace(v.Expected, "'", "`", -1)
		actual, _ := Parse(strings.NewReader(v.Input), v.InputOption)
		if actual != expected {
			t.Errorf("\ninput:\n%v\ngot:\n%v\nwant:\n%v", v.Input, actual, expected)
		}
	}
}
