package json2struct

import (
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/c9s/inflect"
	"github.com/go-openapi/swag"
)

type ObjectType uint

const (
	Invalid ObjectType = iota
	Value
	Hash
	Array
)

const (
	rootID            = "$"
	DefaultStructName = "data"
)

var debugMode bool
var option Options

func SetDebug(v bool) {
	debugMode = v
}

type Options struct {
	UseOmitempty   bool
	UseShortStruct bool
	UseLocal       bool
	UseExample     bool
	Name           string
	Prefix         string
	Suffix         string
}

func Parse(reader io.Reader, opt Options) (string, error) {
	var input interface{}
	if err := json.NewDecoder(reader).Decode(&input); err != nil {
		return "", err
	}
	if opt.Name == "" {
		opt.Name = DefaultStructName
	}
	opt.Prefix = swag.ToGoName(opt.Prefix)
	opt.Suffix = swag.ToGoName(opt.Suffix)
	option = opt
	walker := NewWalker(input)
	walker.start()
	if debugMode {
		b, _ := json.MarshalIndent(walker.structure, "", "  ")
		walker.logln(string(b))
	}
	return walker.output(), nil
}

type Walker struct {
	root      interface{}
	nest      int
	structure *Structure
}

func NewWalker(root interface{}) *Walker {
	return &Walker{
		root: root,
	}
}

func (w *Walker) logln(a ...interface{}) {
	if debugMode {
		prefix := ""
		if w.nest > 1 {
			prefix = strings.Repeat("\t", w.nest-1)
		}
		a = append([]interface{}{prefix}, a...)
		log.Println(a...)
	}
}

func (w *Walker) start() {
	w.walk(rootID, option.Name, w.root, nil)
}

func (w *Walker) output() string {
	return strings.Join(w.structure.Output(), "\n\n")
}

func (w *Walker) walk(spath, name string, data interface{}, parent *Structure) {
	if debugMode {
		w.nest++
		defer func() {
			w.nest--
		}()
	}

	if name != "" {
		spath = fmt.Sprintf("%s.%s", spath, name)
	}

	switch getType(data) {
	case Value:
		v := reflect.ValueOf(data)
		kind := v.Kind()
		if kind == reflect.Float64 {
			kind = getNumberKind(v.Float())
		}
		w.logln(name, kind)
		parent.AddPropety(name, kind, v.Interface(), nil)
	case Array:
		spath = fmt.Sprintf("%s[]", spath)
		w.logln(name)
		w.logln("[")
		list, _ := data.([]interface{})
		item := inflect.Singularize(name)
		current := NewStructure(spath, item)
		for _, val := range list {
			w.walk(spath, item, val, current)
		}
		if parent == nil {
			w.structure = current
		} else {
			parent.AddPropety(name, reflect.Array, list, current)
		}
		w.logln("]")
	case Hash:
		current := NewStructure(spath, name)
		w.logln(name)
		w.logln("{")
		h, _ := data.(map[string]interface{})
		for key, val := range h {
			w.walk(spath, key, val, current)
		}
		w.logln("}")
		if parent == nil {
			w.structure = current
		} else {
			parent.AddPropety(name, reflect.Map, nil, current)
		}
	case Invalid:
		parent.AddPropety(name, reflect.Interface, nil, nil)
	}
	return
}

type Structure struct {
	ID    string
	Name  string
	Props Props
}

type Props []Propety

type Propety struct {
	Name  string
	Kind  reflect.Kind
	Value interface{}
	Refs  *Structure `json:",omitempty"`
}

func NewStructure(spath, name string) *Structure {
	if !option.UseShortStruct {
		name = SpathToName(spath, name)
	}
	name = fmt.Sprintf("%s%s%s", option.Prefix, swag.ToGoName(name), option.Suffix)
	if option.UseLocal {
		name = swag.ToVarName(name)
	}
	return &Structure{
		ID:    spath,
		Name:  name,
		Props: make([]Propety, 0, 8),
	}
}

func SpathToName(spath, name string) string {
	args := strings.Split(spath, ".")
	result := make([]string, 0, len(args))
	for _, v := range args {
		if v == rootID {
			continue
		}
		if strings.HasSuffix(v, "[]") {
			continue
		}
		result = append(result, swag.ToGoName(v))
	}
	if strings.HasSuffix(spath, "[]") {
		result = append(result, swag.ToGoName(name))
	}
	return strings.Join(result, "")
}

func (s *Structure) AddPropety(name string, kind reflect.Kind, val interface{}, refs *Structure) {
	for i, prop := range s.Props {
		if prop.Name != name {
			continue
		}
		// float64 > int
		if prop.Kind == reflect.Float64 && kind == reflect.Int ||
			prop.Kind == reflect.Int && kind == reflect.Float64 {
			s.Props[i] = Propety{Name: name, Kind: reflect.Float64}
			return
		}
		// other kinds -> interface
		if prop.Kind != kind {
			s.Props[i] = Propety{Name: name, Kind: reflect.Interface}
			return
		}
		// merge map pops
		if kind == reflect.Map {
			if refs == nil || prop.Refs == nil {
				return
			}
			for _, p := range refs.Props {
				prop.Refs.AddPropety(p.Name, p.Kind, val, p.Refs)
			}
		}
		return
	}
	prop := Propety{Name: name, Kind: kind, Value: val}
	if refs != nil {
		prop.Refs = refs
	}
	s.Props = append(s.Props, prop)
	sort.Sort(s.Props)
}

func (s *Structure) Output() []string {
	refs := s.Refs()
	result := make([]string, 0, 8)
	if !strings.HasSuffix(s.ID, "[]") {
		result = append(result, s.String())
	}
	for _, ref := range refs {
		result = append(result, ref.Output()...)
	}
	return result
}

func (s *Structure) String() string {
	props := make([]string, len(s.Props))
	for i, prop := range s.Props {
		props[i] = prop.String()
	}
	str := fmt.Sprintf("type %s struct{\n%v\n}", s.Name, strings.Join(props, "\n"))

	formated, err := format.Source([]byte(str))
	if err != nil {
		return str
	}
	return string(formated)
}

func (s *Structure) Refs() []*Structure {
	refs := make([]*Structure, 0, len(s.Props))
	for _, v := range s.Props {
		if v.Refs != nil {
			switch v.Kind {
			case reflect.Map:
				refs = append(refs, v.Refs)
			case reflect.Array:
				refs = append(refs, v.Refs.Refs()...)
			}
		}
	}
	return refs
}

func (p *Propety) String() string {
	kind := "interface{}"
	isStruct := false
	switch p.Kind {
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		kind = p.Kind.String()
	case reflect.Slice, reflect.Array:
		kind = "[]interface{}"
		if p.Refs.Name != "" && len(p.Refs.Props) != 0 {
			if len(p.Refs.Refs()) == 0 && len(p.Refs.Props) == 1 {
				if p.Refs.Props[0].Kind != reflect.Interface {
					kind = fmt.Sprintf("[]%s", p.Refs.Props[0].Kind)
				}
			} else {
				kind = fmt.Sprintf("[]%s", p.Refs.Name)
				isStruct = true
			}
		}

	case reflect.Map:
		if p.Refs.Name != "" {
			kind = p.Refs.Name
			isStruct = true
			if option.UseOmitempty {
				kind = fmt.Sprintf("*%s", kind)
			}
		}
	}
	jsonOption := ""
	if option.UseOmitempty {
		jsonOption = ",omitempty"
	}
	exampleOption := ""
	if option.UseExample && p.Value != nil && !isStruct {
		list, ok := p.Value.([]interface{})
		if ok {
			strs := make([]string, len(list))
			for i, v := range list {
				strs[i] = fmt.Sprint(v)
			}
			p.Value = strings.Join(strs, ",")
		}
		if p.Value != "" {
			exampleOption = fmt.Sprintf(" example:\"%v\"", p.Value)
		}
	}
	propName := swag.ToGoName(p.Name)
	if option.UseLocal {
		propName = swag.ToVarName(propName)
	}
	return fmt.Sprintf("\t%s %s `json:\"%s%s\"%s`", propName, kind, p.Name, jsonOption, exampleOption)
}

func (p Props) Len() int {
	return len(p)
}

func (p Props) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Props) Less(i, j int) bool {
	return p[i].Name < p[j].Name
}

func getType(data interface{}) ObjectType {
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return Value
	case reflect.Slice, reflect.Array:
		return Array
	case reflect.Map:
		return Hash
	default:
		return Invalid
	}
}

func getNumberKind(f float64) reflect.Kind {
	decimals := 10000
	shift := float64(decimals) * f
	num := int(shift)
	if num%decimals == 0 {
		return reflect.Int
	}
	return reflect.Float64
}
