package storage

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pelletier/go-toml/v2"
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
	processId, err := pb.ParseProcessId(processMetadata.ProcessId)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	err = os.MkdirAll(filepath.Join(c.tenantGitRepoPath, processId.Project), os.ModePerm)
	filePath := filepath.Join(c.tenantGitRepoPath, processId.Project, processId.Process)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	var f *os.File
	var commitMessage CommitMessage
	if _, err := os.Stat(filePath); err == nil {
		// path/to/whatever exists
		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
		_, err = f.Write(file)
		if err != nil {
			log.Error().Msg(err.Error())
			return err
		}
		commitMessage = createCommitMessageForNewProcess(
			processId.Process,
			processId.Project,
			processMetadata.Version,
			commons.GetSingleValue(c.Context, commons.User),
			true,
		)
	} else if errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does *not* exist
		f, err = os.Create(filePath)
		_, err = f.Write(file)
		if err != nil {
			log.Error().Msg(err.Error())
			return err
		}
		commitMessage = createCommitMessageForNewProcess(
			processId.Process,
			processId.Project,
			processMetadata.Version,
			commons.GetSingleValue(c.Context, commons.User),
			false,
		)
	}

	worktree, err := c.repo.Worktree()
	if err != nil {
		return err
	}
	_, err = c.commitFile(worktree, filepath.Join(processId.Project, processId.Process), commitMessage)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	return nil
}

func (c *GitClient) GetProcessHistory(processId *pb.ProcessId) []CommitMessage {
	processPath := filepath.Join(processId.Project, processId.Process)
	fileHistory, err := c.repo.Log(&git.LogOptions{
		FileName: &processPath,
	})
	if err != nil {
		return nil
	}
	var commits = make([]CommitMessage, 0)
	err = fileHistory.ForEach(func(commit *object.Commit) error {
		var c CommitMessage
		err := toml.Unmarshal([]byte(commit.Message), &c)
		if err != nil {
			log.Error().Msg(err.Error())
			return err
		}
		commits = append(commits, c)
		return nil
	})
	if err != nil {
		return nil
	}

	return commits
}

func (c *GitClient) createRepository() (*git.Repository, error) {
	return git.PlainInit(c.tenantGitRepoPath, false)
}

func (c *GitClient) commitFile(tree *git.Worktree, file string, commitMessage CommitMessage) (string, error) {
	_, err := tree.Add(file)
	if err != nil {
		return "", err
	}

	b, err := toml.Marshal(commitMessage)
	if err != nil {
		return "", err
	}
	commit, err := tree.Commit(string(b), &git.CommitOptions{
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

func createCommitMessageForNewProcess(
	processName string,
	project string,
	version string,
	updatedBy string,
	updated bool) CommitMessage {
	newFile := ProcessFile{
		ProcessName: processName,
		ProjectId:   project,
		Version:     version,
		LastUpdate:  time.Now().Unix(),
	}

	if updated {
		return CommitMessage{
			AddedFiles:    nil,
			DeletedFiles:  nil,
			ModifiedFiled: []ProcessFile{newFile},
			UpdatedBy:     updatedBy,
		}
	} else {
		return CommitMessage{
			AddedFiles:    []ProcessFile{newFile},
			DeletedFiles:  nil,
			ModifiedFiled: nil,
			UpdatedBy:     updatedBy,
		}
	}
}
