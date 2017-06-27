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
		Prefix:         *prefix,
		Suffix:         *suffix,
		Name:           strings.ToLower(*name),
	}
	parsed := json2struct.Parse(os.Stdin, opt)
	fmt.Println(parsed)
}
