package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yudppp/json2struct"
)

var (
	debug     = flag.Bool("debug", false, "Set debug mode")
	omitempty = flag.Bool("omitempty", false, "Set omitempty mode")
	short     = flag.Bool("short", false, "Set short struct name mode")
	local     = flag.Bool("local", false, "Use local struct mode")
	example   = flag.Bool("example", false, "Use example tag mode")
	prefix    = flag.String("prefix", "", "Set struct name prefix")
	suffix    = flag.String("suffix", "", "Set struct name suffix")
	name      = flag.String("name", json2struct.DefaultStructName, "Set struct name")
)

func main() {
	flag.Parse()
	json2struct.SetDebug(*debug)
	opt := json2struct.Options{
		UseOmitempty:   *omitempty,
		UseShortStruct: *short,
		UseLocal:       *local,
		UseExample:     *example,
		Prefix:         *prefix,
		Suffix:         *suffix,
		Name:           strings.ToLower(*name),
	}
	parsed, err := json2struct.Parse(os.Stdin, opt)
	if err != nil {
		panic(err)
	}
	fmt.Println(parsed)
}
