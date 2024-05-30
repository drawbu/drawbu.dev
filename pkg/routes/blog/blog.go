package blog

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

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

	if serv.Blog.LastLookup >= 60*60 {
		fmt.Printf("Pulling from %s\n", serv.Blog.RepoPath)
		err := exec.Command("git", "pull", serv.Blog.RepoPath).Run()
		if err != nil {
			return err
		}
		serv.Blog.LastLookup = 0 // 1 hour
		serv.Blog.Articles = getArticles(serv.Blog.RepoPath, serv.Blog.RepoPath)
	} else {
		serv.Blog.LastLookup++
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
			articles = append(articles, app.Article{Title: name, Path: path + "/" + entry.Name(), URI: uri})
		}
	}
	return articles
}
