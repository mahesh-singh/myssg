package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
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

	if err := RenderIndexPage("templates/index.html", "output"); err != nil {
		fmt.Printf("error while rendering index page: %v \n", err)
		return
	}

	blogPosts, err := ParseBlogPosts("content/posts/")
	if err != nil {
		fmt.Printf("error while parsing blog posts: %v \n", err)
		return
	}

	if err := RenderBlogPosts(blogPosts, "templates/posts/post.html", "output/posts/"); err != nil {
		fmt.Printf("error while rendering blog posts: %v \n", err)
		return
	}

	if err := CopyStaticFiles("templates/static", "output/static"); err != nil {
		fmt.Printf("error while copy static: %v \n", err)
		return
	}

	if err := CopyStaticFiles("content/img", "output/img"); err != nil {
		fmt.Printf("error while copy static: %v \n", err)
		return
	}
}

func ParseBlogPosts(dir string) ([]Post, error) {
	var posts []Post

	files, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return nil, err
	}

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

	return posts, nil
}

func LoadMarkdownFile(filePath string) (*Post, error) {
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
	tmpl, err := template.ParseFiles("templates/base.html", "templates/partials/nav.html", templatePath)

	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}

	for _, post := range posts {
		var htmlContent strings.Builder
		if err := tmpl.ExecuteTemplate(&htmlContent, "base", post); err != nil {
			fmt.Printf("error rendering post %s: %v \n", post.Slug, err)
			continue
		}

		outputPath := filepath.Join(outputDir, post.Slug+".html")
		if err := os.WriteFile(outputPath, []byte(htmlContent.String()), 0644); err != nil {
			fmt.Printf("error saving file %s: %v \n", outputPath, err)
		}

	}

	err = RenderBlogPostsIndex(posts, "templates/posts/index.html", outputDir)
	if err != nil {
		return fmt.Errorf("error rendering post index: %v", err)
	}

	return nil
}

func RenderBlogPostsIndex(posts []Post, templatePath string, outputDir string) error {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/partials/nav.html", templatePath)
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}

	var htmlContent strings.Builder
	if err := tmpl.ExecuteTemplate(&htmlContent, "base", posts); err != nil {
		fmt.Printf("error rendering index of post %v \n", err)
	}

	outputPath := filepath.Join(outputDir, "index.html")

	if err := os.WriteFile(outputPath, []byte(htmlContent.String()), 0644); err != nil {
		fmt.Printf("error saving file %s: %v \n", outputPath, err)
	}
	return nil
}

func RenderIndexPage(templatePath string, outputDir string) error {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/partials/nav.html", templatePath)
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}

	var htmlContent strings.Builder
	if err := tmpl.ExecuteTemplate(&htmlContent, "base", nil); err != nil {
		fmt.Printf("error rendering index of post %v \n", err)
	}

	outputPath := filepath.Join(outputDir, "index.html")

	if err := os.WriteFile(outputPath, []byte(htmlContent.String()), 0644); err != nil {
		fmt.Printf("error saving file %s: %v \n", outputPath, err)
	}
	return nil
}

func CopyStaticFiles(srcDir, outputDir string) error {
	err := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(outputDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(destPath, os.ModePerm)
		}

		return CopyFile(path, destPath)
	})

	return err
}

func CopyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return nil
	}

	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return nil
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)

	return err

}
