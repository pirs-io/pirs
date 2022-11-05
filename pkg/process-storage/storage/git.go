package storage

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
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

type GitClient struct {
	Context      context.Context
	RepoRootPath string
	Tenant       string

	repo              *git.Repository
	tenantGitRepoPath string
}

func (c *GitClient) InitializeStorage() error {
	c.tenantGitRepoPath = filepath.Join(c.RepoRootPath, c.Tenant)
	r, err := git.PlainOpen(c.tenantGitRepoPath)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			log.Info().Msgf("Git storage folder is empty, creating repository for tenant: %s", c.Tenant)
			r, err = c.createRepository()
		}
	}
	c.repo = r
	return err
}

func (c *GitClient) SaveFile(processMetadata *pb.ProcessMetadata, file []byte) error {
	filePath := filepath.Join(c.tenantGitRepoPath, processMetadata.Filename)
	f, err := os.Create(filePath)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	_, err = f.Write(file)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	worktree, err := c.repo.Worktree()
	_, err = c.commitFile(worktree, processMetadata.Filename, "added process: "+processMetadata.ProcessId)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	return nil
}

func (c *GitClient) createRepository() (*git.Repository, error) {
	return git.PlainInit(c.tenantGitRepoPath, false)
}

func (c *GitClient) commitFile(tree *git.Worktree, filePath string, commitMessage string) (string, error) {
	_, err := tree.Add(filePath)
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
