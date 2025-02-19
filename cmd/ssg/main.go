package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/yuin/goldmark"
)

type post struct {
	Title   string
	Date    string
	Tags    []string
	Content template.HTML
	Summary string
	Slug    string
	Draft   bool
}

type Metadata struct {
	Title string    `toml:"title"`
	Date  time.Time `toml:"date"`
	Tags  []string  `toml:"tags"`
	Slug  string    `toml:"slug"`
	Draft bool      `toml:"draft"`
}

func main() {
	fmt.Print("this is my static site generator")

	post, err := LoadMarkdownFiles("content/hello.md")

	if err != nil {
		fmt.Printf("error while loading ghe markdown file: %v \n", err)
		return
	}

	tmpl, err := template.ParseFiles("templates/base.html")
	if err != nil {
		fmt.Printf("error while parsing the template: %v \n", err)
		return
	}
	var htmlContent strings.Builder
	err = tmpl.Execute(&htmlContent, post)
	if err != nil {
		fmt.Printf("error while executing the template: %v \n", err)
		return
	}

	outputPath := fmt.Sprintf("output/%s.html", post.Slug)
	err = SaveToFile(outputPath, htmlContent.String())
	if err != nil {
		fmt.Printf("error while saving the file: %v \n", err)
		return
	}

	fmt.Println((htmlContent.String()))
}

func LoadMarkdownFiles(filePath string) (*post, error) {
	content, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Printf("error wile loading the markdown content from file: %v \n", err)
		return nil, err
	}

	text := string(content)
	metadata, markdown, err := ExtractMetadata(text)

	if err != nil {
		return nil, err
	}

	htmlContent := ConvertMarkdownToHTML(markdown)

	return &post{
		Title:   metadata.Title,
		Date:    metadata.Date.String(),
		Tags:    metadata.Tags,
		Content: htmlContent,
		Slug:    metadata.Slug,
		Draft:   metadata.Draft,
	}, nil

}

func ConvertMarkdownToHTML(md string) template.HTML {
	var buf bytes.Buffer

	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		fmt.Printf("error while converting markdown into into html: %v \n", err)
		return ""
	}
	return template.HTML(buf.String())
}

func ExtractMetadata(text string) (Metadata, string, error) {
	re := regexp.MustCompile(`(?s)^\+\+\+\n(.*?)\n\+\+\+\n(.*)`)
	matches := re.FindStringSubmatch(text)

	if len(matches) < 3 {
		return Metadata{}, "", fmt.Errorf("failed to extract metadata")
	}

	tomlData := matches[1]
	markdownContent := matches[2]

	var metadata Metadata

	err := toml.Unmarshal([]byte(tomlData), &metadata)

	if err != nil {
		return Metadata{}, "", fmt.Errorf("error parsing TOML metadata %v \n", err)
	}

	return metadata, markdownContent, nil
}

func SaveToFile(filePath, content string) error {
	if err := os.MkdirAll("output", os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(filePath, []byte(content), 0644)
}
