package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
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

	markdownContent, err := LoadMarkdownFiles("content/hello.md")

	if err != nil {
		fmt.Printf("error while loading ghe markdown file: %v \n", err)
		return
	}

	content := ConvertMarkdownToHTML(markdownContent)

	p := post{
		Title:   "this is a title",
		Content: content,
		Slug:    "hello",
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

	outputPath := fmt.Sprintf("output/%s.html", p.Slug)
	err = SaveToFile(outputPath, htmlContent.String())
	if err != nil {
		fmt.Printf("error while saving the file: %v \n", err)
		return
	}

	fmt.Println((htmlContent.String()))
}

func LoadMarkdownFiles(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)

	if err != nil {
		fmt.Printf("error wile loading the markdown content from file: %v \n", err)
		return "", err
	}

	return string(content), err

}

func ConvertMarkdownToHTML(md string) template.HTML {
	var buf bytes.Buffer

	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		fmt.Printf("error while converting markdown into into html: %v \n", err)
		return ""
	}
	return template.HTML(buf.String())
}

func SaveToFile(filePath, content string) error {
	if err := os.MkdirAll("output", os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(filePath, []byte(content), 0644)
}
