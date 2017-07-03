# json2struct

Translate to Golang struct from json

## Playground

Try playground [https://yudppp.github.io/json2struct](https://yudppp.github.io/json2struct)

## CLI

### Install

```
$ go get github.com/yudppp/json2struct/cmd/json2struct
``` 

### How to use

```
$ echo '{"url": "http://blog.yudppp.com", "text": "Hello:)", "status": 1, "categories": [{"name": "k8s"}]}' |  json2struct -name=blog
type Blog struct {
	Categories []BlogCategory `json:"categories"`
	Status     int            `json:"status"`
	Text       string         `json:"text"`
	URL        string         `json:"url"`
}

type BlogCategory struct {
	Name string `json:"name"`
}
```

#### options

| option | description |
|:-----------|:-----------|
| name | Set struct name (default "data") |
| prefix | Set struct name prefix |
| suffix | Set struct name suffix |
| short | Set short struct name mode |
| local | Use local struct mode |
| omitempty | Set omitempty mode |
| example | Use example tag (https://github.com/yudppp/structs)|
