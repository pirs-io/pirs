package storage

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
	"pirs.io/common"
	"text/template"
	"time"
)

const (
	repoIsEmpty = "remote repository is empty"
	remoteName  = "origin"
)

type OrganizationInfo struct {
	OrganizationId string
}

type GitClient struct {
	Context         context.Context
	Url             string
	Username        string
	Password        string
	TempRepoDirPath string
}

func (c *GitClient) InitializeStorage() error {
	auth := &http.BasicAuth{
		Username: c.Username,
		Password: c.Password,
	}
	r, err := c.clone(c.TempRepoDirPath, auth)
	initialFilePath, err := c.createInitialFile(c.TempRepoDirPath)
	worktree, _ := r.Worktree()

	if err != nil && err.Error() == repoIsEmpty {
		r, err = git.Init(r.Storer, worktree.Filesystem)
		_, err = r.CreateRemote(&config.RemoteConfig{
			Name: remoteName,
			URLs: []string{c.Url},
		})
		_, err = c.commitInitialFile(worktree, initialFilePath)
		err = c.push(r, auth)
		return err
	}
	//_, err = c.commitInitialFile(worktree, initialFilePath)
	//err = c.push(r, auth)
	return err
}

// Close cleanup cloned repository from FS/**
func (c *GitClient) Close() error {
	err := os.RemoveAll(c.TempRepoDirPath)
	if err != nil {
		return err
	}
	return nil
}

func (c *GitClient) clone(dir string, auth *http.BasicAuth) (*git.Repository, error) {
	return git.PlainClone(dir, false, &git.CloneOptions{
		Auth:     auth,
		URL:      c.Url,
		Progress: os.Stdout,
	})
}

func (c *GitClient) push(r *git.Repository, auth *http.BasicAuth) error {
	remote, err := r.Remote(remoteName)
	err = remote.Push(&git.PushOptions{
		RemoteName: remoteName,
		Auth:       auth,
	})
	if err != nil {
		return err
	}
	log.Info().Msg("Pushed to remote")
	return nil
}

func (c *GitClient) commitInitialFile(tree *git.Worktree, initialFilePath string) (string, error) {
	_, err := tree.Add(initialFilePath)
	if err != nil {
		return "", err
	}
	commit, err := tree.Commit("added initial file", &git.CommitOptions{
		Author: &object.Signature{
			Name:  common.GetSingleValue(c.Context, common.User),
			Email: common.GetSingleValue(c.Context, common.UserEmail),
			When:  time.Now(),
		},
	})
	if err != nil {
		return "", err
	}

	return commit.String(), nil
}

func (c *GitClient) createInitialFile(dir string) (string, error) {
	var temp *template.Template
	outFile, err := os.Create(filepath.Join(dir + "/README.md"))
	if err != nil {
		return "", nil
	}
	temp = template.Must(template.ParseFiles("goTemplates/README.md"))
	err = temp.Execute(outFile, OrganizationInfo{OrganizationId: "sk.dudak"})
	if err != nil {
		return "", err
	}
	err = outFile.Close()
	return filepath.Base(outFile.Name()), nil
}
