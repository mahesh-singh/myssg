package main

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/yuin/goldmark"
)

type post struct {
	Title   string
	Date    string
	Tags    []string
	Content template.HTML
	Summary string
	Slug    string
}

func main() {
	fmt.Print("this is my static site generator")
	p := post{
		Title:   "this is a title",
		Content: ConvertMarkdownToHTML("this is content \n **this is string**"),
	}

	tmpl, err := template.ParseFiles("templates/base.html")
	if err != nil {
		fmt.Printf("error while parsing the template: %v \n", err)
		return
	}
	var htmlContent strings.Builder
	err = tmpl.Execute(&htmlContent, p)

	if err != nil {
		fmt.Printf("error while executing the template: %v \n", err)
		return
	}

	fmt.Println((htmlContent.String()))
}

func ConvertMarkdownToHTML(md string) template.HTML {
	var buf bytes.Buffer

	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		fmt.Printf("error while converting markdown into into html: %v \n", err)
		return ""
	}
	return template.HTML(buf.String())
}
