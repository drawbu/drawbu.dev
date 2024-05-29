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

func getArticles(repo string) (string, error) {
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

	fmt.Println(repo, dirname)
	cmd := exec.Command("git", "git", "clone", repo, dirname)
	if cmd == nil {
		return "", errors.New("error running git")
	}
	err = cmd.Run()
	if err != nil {
		return "", errors.New("error cloning repo: " + err.Error())
	}
	return dirname, nil
}

func Handler(serv *app.Server, w http.ResponseWriter, r *http.Request) error {
	if serv.Blog.RepoPath == "" {
		serv.Blog.GithubRepo = "https://github.com/drawbu/Notes"
		path, err := getArticles(serv.Blog.GithubRepo)
		if err != nil {
			return err
		}
		serv.Blog.RepoPath = path
	}

	_, article_path, found := strings.Cut(r.RequestURI, "/blog/")

	if !found || article_path == "" {
		return components.Template(blog()).Render(context.Background(), w)
	}
	return components.Template(article()).Render(context.Background(), w)
}
