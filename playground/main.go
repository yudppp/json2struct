package main

import (
	"strings"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
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
		{"name": "k8s"}
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
}

// Render implements the vecty.Component interface.
func (p *PageView) Render() *vecty.HTML {
	return elem.Body(
		elem.Header(
			prop.Class("header"),
			vecty.Text("json2struct playground"),
		),
		elem.Div(
			prop.Class("wrapper"),
			elem.Div(
				prop.Class("column input"),
				elem.Div(
					vecty.Tag("label",
						vecty.Text("input json"),
					),
					elem.TextArea(
						prop.Class("u-full-width"),
						vecty.Text(p.Input),
						event.Input(func(e *vecty.Event) {
							p.Input = e.Target.Get("value").String()
							vecty.Rerender(p)
						}),
					),
				),
				elem.Div(
					vecty.Tag("label",
						vecty.Text("struct name"),
					),
					elem.Input(
						prop.Value(p.StructName),
						prop.Type(prop.TypeText),
						event.Input(func(e *vecty.Event) {
							p.StructName = e.Target.Get("value").String()
							vecty.Rerender(p)
						}),
					),
				),
				elem.Div(
					vecty.Tag("label",
						vecty.Text("struct prefix name"),
					),
					elem.Input(
						prop.Value(p.Prefix),
						prop.Type(prop.TypeText),
						event.Input(func(e *vecty.Event) {
							p.Prefix = e.Target.Get("value").String()
							vecty.Rerender(p)
						}),
					),
				),
				elem.Div(
					vecty.Tag("label",
						vecty.Text("struct suffix name"),
					),
					elem.Input(
						prop.Value(p.Suffix),
						prop.Type(prop.TypeText),
						event.Input(func(e *vecty.Event) {
							p.Suffix = e.Target.Get("value").String()
							vecty.Rerender(p)
						}),
					),
				),
				elem.Div(
					vecty.Tag("label",
						vecty.Text("omitempty mode"),
					),
					elem.Input(
						prop.Class("toggle"),
						prop.Type(prop.TypeCheckbox),
						prop.Checked(p.UseOmitempty),
						event.Change(func(e *vecty.Event) {
							p.UseOmitempty = e.Target.Get("checked").Bool()
							vecty.Rerender(p)
						}),
					),
				),
				elem.Div(
					vecty.Tag("label",
						vecty.Text("short mode"),
					),
					elem.Input(
						prop.Class("toggle"),
						prop.Type(prop.TypeCheckbox),
						prop.Checked(p.UseShortStruct),
						event.Change(func(e *vecty.Event) {
							p.UseShortStruct = e.Target.Get("checked").Bool()
							vecty.Rerender(p)
						}),
					),
				),
				elem.Div(
					vecty.Tag("label",
						vecty.Text("local mode"),
					),
					elem.Input(
						prop.Class("toggle"),
						prop.Type(prop.TypeCheckbox),
						prop.Checked(p.UseLocal),
						event.Change(func(e *vecty.Event) {
							p.UseLocal = e.Target.Get("checked").Bool()
							vecty.Rerender(p)
						}),
					),
				),
			),
			elem.Div(
				prop.Class("column output"),
				vecty.Tag("label",
					vecty.Text("output struct"),
				),

				vecty.Tag("pre",
					elem.Code(
						&StructObject{
							Input: p.Input,
							Option: json2struct.Options{
								Name:           p.StructName,
								Prefix:         p.Prefix,
								Suffix:         p.Suffix,
								UseShortStruct: p.UseShortStruct,
								UseLocal:       p.UseLocal,
								UseOmitempty:   p.UseOmitempty,
							},
						},
					),
				),
			),
		),
		elem.Footer(
			prop.Class("footer"),
			vecty.Text("Used by "),
			elem.Anchor(
				prop.Href("https://github.com/yudppp/json2struct"),
				vecty.Text("yudppp/json2struct"),
			),
		),
	)

}

// StructObject is output values.
type StructObject struct {
	vecty.Core
	Input          string
	Option         json2struct.Options
	StructName     string
	Prefix         string
	Suffix         string
	UseShortStruct bool
	UseLocal       bool
	UseOmitempty   bool
}

// Render implements the vecty.Component interface.
func (m *StructObject) Render() *vecty.HTML {

	output, err := json2struct.Parse(strings.NewReader(m.Input), m.Option)
	if err != nil {
		return elem.Div(
			vecty.Text("invalid"),
		)
	}
	return elem.Div(
		vecty.Text(output),
	)
}
