package storage

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
	"pirs.io/common"
	pb "pirs.io/process-storage/grpc"
	"text/template"
	"time"
)

const (
	repoIsEmpty = "remote repository is empty"
	remoteName  = "origin"
)

var (
	log = common.GetLoggerFor("git")
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

	repo *git.Repository
	auth *http.BasicAuth
}

func (c *GitClient) InitializeStorage() error {
	auth := &http.BasicAuth{
		Username: c.Username,
		Password: c.Password,
	}
	r, err := c.clone(c.TempRepoDirPath, auth)
	worktree, _ := r.Worktree()

	c.repo = r
	c.auth = auth

	if err != nil && err.Error() == repoIsEmpty {
		initialFilePath, err := c.createInitialFile()
		r, err = git.Init(r.Storer, worktree.Filesystem)
		_, err = r.CreateRemote(&config.RemoteConfig{
			Name: remoteName,
			URLs: []string{c.Url},
		})
		_, err = c.commitFile(worktree, initialFilePath, "Commited initial file")
		err = c.push(r, auth)
		return err
	}
	return err
}

func (c *GitClient) SaveFile(
	processMetadata *pb.ProcessMetadata,
	file []byte) {
	f, err := os.Create(filepath.FromSlash(c.TempRepoDirPath + "/" + processMetadata.Filename))
	_, err = f.Write(file)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	worktree, err := c.repo.Worktree()
	_, err = c.commitFile(worktree, processMetadata.Filename, "added process: "+processMetadata.ProcessId)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	err = c.push(c.repo, c.auth)
	if err != nil {
		log.Error().Msg(err.Error())
	}
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

func (c *GitClient) commitFile(tree *git.Worktree, initialFilePath string, commitMessage string) (string, error) {
	_, err := tree.Add(initialFilePath)
	if err != nil {
		return "", err
	}
	commit, err := tree.Commit(commitMessage, &git.CommitOptions{
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

func (c *GitClient) createInitialFile() (string, error) {
	var temp *template.Template
	outFile, err := os.Create(filepath.Join(c.TempRepoDirPath + "/README.md"))
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
