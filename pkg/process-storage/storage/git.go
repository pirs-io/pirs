package storage

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
	"pirs.io/commons"
	pb "pirs.io/process-storage/grpc"
	"time"
)

var (
	log = commons.GetLoggerFor("git")
)

type OrganizationInfo struct {
	OrganizationId string
}

type GitClient struct {
	Context      context.Context
	RepoRootPath string
	Username     string
	Password     string

	repo *git.Repository
	auth *http.BasicAuth
}

func (c *GitClient) InitializeStorage() error {
	auth := &http.BasicAuth{
		Username: c.Username,
		Password: c.Password,
	}
	r, err := git.PlainOpen(c.RepoRootPath)
	c.repo = r
	c.auth = auth

	return err
}

func (c *GitClient) SaveFile(processMetadata *pb.ProcessMetadata, file []byte) {
	f, err := os.Create(filepath.FromSlash(c.RepoRootPath + "/" + processMetadata.Filename))
	_, err = f.Write(file)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	worktree, err := c.repo.Worktree()
	_, err = c.commitFile(worktree, processMetadata.Filename, "added process: "+processMetadata.ProcessId)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}

func (c *GitClient) commitFile(tree *git.Worktree, initialFilePath string, commitMessage string) (string, error) {
	_, err := tree.Add(initialFilePath)
	if err != nil {
		return "", err
	}
	commit, err := tree.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  commons.GetSingleValue(c.Context, commons.User),
			Email: commons.GetSingleValue(c.Context, commons.UserEmail),
			When:  time.Now(),
		},
	})
	if err != nil {
		return "", err
	}

	return commit.String(), nil
}
