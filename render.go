package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"path"
	"strings"
)

const (
	Views = "./views"
)

type Layout struct {
	Page template.HTML
}

type Renderer struct {
	Files map[string]string
}

func NewRenderer() (r Renderer, err error) {
	var c []byte
	r.Files = make(map[string]string)

	files, err := ioutil.ReadDir(Views)
	if err != nil {
		return
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".tmpl") {
			continue
		}

		c, err = ioutil.ReadFile(path.Join(Views, file.Name()))
		if err != nil {
			return
		}

		r.Files[strings.Split(file.Name(), ".")[0]] = string(c)
	}

	return
}

func (r Renderer) Render(tmpl string, ctx interface{}) (body []byte, err error) {
	var layoutBuf, pageBuf bytes.Buffer

	page, err := template.New("page").Parse(r.Files[tmpl])
	if err != nil {
		return
	}

	err = page.Execute(&pageBuf, ctx)
	if err != nil {
		return
	}

	l := Layout{
		Page: template.HTML(pageBuf.String()),
	}

	layout, err := template.New("layout").Parse(r.Files["layout"])
	if err != nil {
		return
	}

	err = layout.Execute(&layoutBuf, l)
	if err != nil {
		return
	}

	body = layoutBuf.Bytes()

	return
}
