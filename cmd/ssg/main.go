package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/yuin/goldmark"
)

type Post struct {
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
	blogPosts, err := ParseBlogPosts("content/posts/")
	if err != nil {
		fmt.Printf("error while parsing blog posts: %v \n", err)
		return
	}

	if err := RenderBlogPosts(blogPosts, "templates/posts/post.html", "output/posts/"); err != nil {
		fmt.Printf("error while rendering blog posts: %v \n", err)
		return
	}
}

func ParseBlogPosts(dir string) ([]Post, error) {
	var posts []Post

	files, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return nil, err
	}
	fmt.Printf("filepath loaded count %d \n", len(files))
	for _, file := range files {
		post, err := LoadMarkdownFile(file)

		if err != nil {
			fmt.Printf("error while processing file %s: %v \n", file, err)
			continue
		}
		if !post.Draft { // Skip drafts
			posts = append(posts, *post)
		}
	}
	fmt.Printf("generate posts count %d \n", len(posts))
	return posts, nil
}

func LoadMarkdownFile(filePath string) (*Post, error) {
	content, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Printf("error wile loading the markdown content from file: %v \n", err)
		return nil, err
	}

	fmt.Printf("markdown loaded from path %s, len of %d \n", filePath, len(content))

	text := string(content)
	metadata, markdown, err := ExtractMetadata(text)

	if err != nil {
		return nil, err
	}

	fmt.Printf("meta data parsed for  %s \n", metadata.Title)
	htmlContent := ConvertMarkdownToHTML(markdown)

	fmt.Printf("html generated for %s, len of %d \n", filePath, len(htmlContent))

	return &Post{
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
		return Metadata{}, "", fmt.Errorf("error parsing TOML metadata %v", err)
	}

	return metadata, markdownContent, nil
}

func RenderBlogPosts(posts []Post, templatePath, outputDir string) error {
	tmpl, err := template.ParseFiles("templates/base.html", templatePath)
	fmt.Printf("loaded template %s \n", templatePath)
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}
	fmt.Printf("loaded template %s \n", templatePath)

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}

	for _, post := range posts {
		var htmlContent strings.Builder
		if err := tmpl.ExecuteTemplate(&htmlContent, "base", post); err != nil {
			fmt.Printf("error rendering post %s: %v \n", post.Slug, err)
			continue
		}
		fmt.Printf("template executed for post %v of content %s \n", post, &htmlContent)
		outputPath := filepath.Join(outputDir, post.Slug+".html")
		if err := os.WriteFile(outputPath, []byte(htmlContent.String()), 0644); err != nil {
			fmt.Printf("error saving file %s: %v \n", outputPath, err)
		}
		fmt.Printf("output generated %s \n", outputPath)
	}

	return nil
}
