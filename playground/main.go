package main

import (
	"strings"
	"time"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/hexops/vecty/prop"
	"github.com/yudppp/json2struct"
)

func main() {
	vecty.SetTitle("json2struct playground")
	vecty.AddStylesheet("https://cdnjs.cloudflare.com/ajax/libs/skeleton/2.0.4/skeleton.min.css")
	vecty.AddStylesheet("app.css")
	vecty.RenderBody(&PageView{
		Input: `{
	"text": "Hello",
	"categories": [
		{"id": 1,"name": "Golang"}
	]
}`,
		StructName: "blog",
	})

}

// PageView is our main page component.
type PageView struct {
	vecty.Core
	Input          string
	StructName     string
	Prefix         string
	Suffix         string
	UseShortStruct bool
	UseLocal       bool
	UseOmitempty   bool
	UseExampleTag  bool
	lastTimeKey    int64
}

// Render implements the vecty.Component interface.
func (p *PageView) Render() vecty.ComponentOrHTML {
	return elem.Body(
		elem.Header(
			vecty.Markup(
				vecty.Class("header"),
			),
			vecty.Text("json2struct playground"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("wrapper"),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("col", "input"),
				),
				elem.Div(
					vecty.Tag("label",
						vecty.Text("input json"),
					),
					elem.TextArea(
						vecty.Markup(
							vecty.Class("u-full-width"),
							event.Input(func(e *vecty.Event) {
								p.Input = e.Target.Get("value").String()
								p.Rerender()
							}),
						),
						vecty.Text(p.Input),
					),
				),
				elem.Div(
					vecty.Tag("label",
						vecty.Text("struct name"),
					),
					elem.Input(
						vecty.Markup(
							prop.Value(p.StructName),
							prop.Type(prop.TypeText),
							event.Input(func(e *vecty.Event) {
								p.StructName = e.Target.Get("value").String()
								p.Rerender()
							}),
						),
					),
				),
				elem.Div(
					vecty.Tag("label",
						vecty.Text("struct prefix name"),
					),
					elem.Input(
						vecty.Markup(
							prop.Value(p.Prefix),
							prop.Type(prop.TypeText),
							event.Input(func(e *vecty.Event) {
								p.Prefix = e.Target.Get("value").String()
								p.Rerender()
							}),
						),
					),
				),
				elem.Div(
					vecty.Tag("label",
						vecty.Text("struct suffix name"),
					),
					elem.Input(
						vecty.Markup(
							prop.Value(p.Suffix),
							prop.Type(prop.TypeText),
							event.Input(func(e *vecty.Event) {
								p.Suffix = e.Target.Get("value").String()
								p.Rerender()
							}),
						),
					),
				),
				elem.Div(
					vecty.Tag("label",
						vecty.Text("omitempty mode"),
					),
					elem.Input(
						vecty.Markup(
							vecty.Class("toggle"),
							prop.Type(prop.TypeCheckbox),
							prop.Checked(p.UseOmitempty),
							event.Change(func(e *vecty.Event) {
								p.UseOmitempty = e.Target.Get("checked").Bool()
								p.Rerender()
							}),
						),
					),
				),
				elem.Div(
					vecty.Tag("label",
						vecty.Text("short mode"),
					),
					elem.Input(
						vecty.Markup(
							vecty.Class("toggle"),
							prop.Type(prop.TypeCheckbox),
							prop.Checked(p.UseShortStruct),
							event.Change(func(e *vecty.Event) {
								p.UseShortStruct = e.Target.Get("checked").Bool()
								p.Rerender()
							}),
						),
					),
				),
				elem.Div(
					vecty.Tag("label",
						vecty.Text("local mode"),
					),
					elem.Input(
						vecty.Markup(
							vecty.Class("toggle"),
							prop.Type(prop.TypeCheckbox),
							prop.Checked(p.UseLocal),
							event.Change(func(e *vecty.Event) {
								p.UseLocal = e.Target.Get("checked").Bool()
								p.Rerender()
							}),
						),
					),
				),
				elem.Div(
					vecty.Tag("label",
						vecty.Text("example tag mode"),
					),
					elem.Input(
						vecty.Markup(
							vecty.Class("toggle"),
							prop.Type(prop.TypeCheckbox),
							prop.Checked(p.UseExampleTag),
							event.Change(func(e *vecty.Event) {
								p.UseExampleTag = e.Target.Get("checked").Bool()
								p.Rerender()
							}),
						),
					),
				),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("col", "output"),
				),
				vecty.Tag("label",
					vecty.Text("output struct"),
				),

				vecty.Tag("pre",
					elem.Code(
						vecty.Text(p.renderStructOutput()),
					),
				),
			),
		),
		elem.Footer(
			vecty.Markup(
				vecty.Class("footer"),
			),
			vecty.Text("Used by "),
			elem.Anchor(
				vecty.Markup(
					prop.Href("https://github.com/yudppp/json2struct"),
				),
				vecty.Text("yudppp/json2struct"),
			),
		),
	)

}

// Rerender is rerender and debounce
func (p *PageView) Rerender() {
	timeKey := time.Now().UnixNano()
	p.lastTimeKey = timeKey
	go func() {
		time.Sleep(800 * time.Millisecond)
		if timeKey == p.lastTimeKey {
			vecty.Rerender(p)
		}
	}()
}

// renderStructOutput renders the struct output as text
func (p *PageView) renderStructOutput() string {
	opt := json2struct.Options{
		Name:           p.StructName,
		Prefix:         p.Prefix,
		Suffix:         p.Suffix,
		UseShortStruct: p.UseShortStruct,
		UseLocal:       p.UseLocal,
		UseOmitempty:   p.UseOmitempty,
		UseExample:     p.UseExampleTag,
	}
	output, err := json2struct.Parse(strings.NewReader(p.Input), opt)
	if err != nil {
		return "invalid JSON"
	}
	return output
}
