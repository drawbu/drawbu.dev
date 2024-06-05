package blog

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"server/pkg/app"
	"server/pkg/components"
)

func Handler(serv *app.Server, w http.ResponseWriter, r *http.Request) error {
	if serv.Blog.RepoPath == "" {
		serv.Blog.GithubRepo = "https://github.com/drawbu/Notes"
		path, err := cloneRepo(serv.Blog.GithubRepo)
		if err != nil {
			return err
		}
		serv.Blog.RepoPath = path
		serv.Blog.Articles = getArticles(serv.Blog.RepoPath, serv.Blog.RepoPath)
	}

	if time.Since(serv.Blog.LastLookup).Hours() > 4 {
		fmt.Printf("Pulling from %s\n", serv.Blog.RepoPath)
		err := exec.Command("git", "-C", serv.Blog.RepoPath, "pull").Run()
		if err != nil {
			return err
		}
		serv.Blog.LastLookup = time.Now()
		serv.Blog.Articles = getArticles(serv.Blog.RepoPath, serv.Blog.RepoPath)
	}

	if r.URL.Path == "/blog" || r.URL.Path == "/blog/" {
		return components.Template(blog(serv.Blog.Articles)).Render(context.Background(), w)
	}

	a, err := findArticle(serv.Blog.Articles, r.URL.Path)
	if err != nil {
		return err
	}
	return components.Template(article(a)).Render(context.Background(), w)
}

func findArticle(articles []app.Article, path string) (app.Article, error) {
	for _, a := range articles {
		if a.URI == path {
			return a, nil
		}
	}
	return app.Article{}, errors.New("Article not found")
}

func cloneRepo(repo string) (string, error) {
	// Set path
	dirname := os.TempDir()
	if !strings.HasSuffix(dirname, "/") {
		dirname += "/"
	}
	dirname += "drawbu-blog"

	// Check if path exists
	_, err := os.Stat(dirname)
	if err == nil {
		err = os.RemoveAll(dirname)
		if err != nil {
			return "", err
		}
	}

	// Clone
	fmt.Printf("Cloning to %s\n", dirname)
	err = exec.Command("git", "clone", repo, dirname).Run()
	if err != nil {
		return "", errors.New("git: " + err.Error())
	}
	return dirname, nil
}

func getArticles(path string, basepath string) []app.Article {
	entries, err := os.ReadDir(path)
	if err != nil {
		return []app.Article{}
	}

	articles := []app.Article{}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		if entry.IsDir() {
			articles = append(articles, getArticles(path+"/"+entry.Name(), basepath)...)
			continue
		}
		if strings.HasSuffix(entry.Name(), ".md") {
			name := strings.TrimSuffix(entry.Name(), ".md")
			uri := "/blog" + strings.TrimPrefix(path, basepath) + "/" + name
			filepath := path + "/" + entry.Name()
			articles = append(articles, app.Article{Title: name, Path: filepath, URI: uri, Content: getContent(filepath)})
		}
	}
	return articles
}

func getContent(path string) []byte {
	file, err := os.ReadFile(path)
	if err != nil {
		return []byte{}
	}
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(file)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}
