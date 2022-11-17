package storage

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
	"golang.org/x/net/context"
	"io"
	"os"
	"path/filepath"
	"pirs.io/commons"
	pb "pirs.io/process-storage/grpc"
	"sort"
	"time"
)

var (
	log = commons.GetLoggerFor("git")
)

type GitClient struct {
	Context      context.Context
	RepoRootPath string
	Tenant       string
	ChunkSize    int64

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

// SaveFile saves given process
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
	var commitMessage UploadActionSummary
	newFileCreated := false
	// take actions based on process file existence
	if _, err := os.Stat(filePath); err == nil {
		// path/to/whatever exists
		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
		_, err = f.Write(file)
		if err != nil {
			log.Error().Msg(err.Error())
			return err
		}

	} else if errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does *not* exist
		f, err = os.Create(filePath)
		_, err = f.Write(file)
		if err != nil {
			log.Error().Msg(err.Error())
			return err
		}
		newFileCreated = true
	}
	commitMessage = createCommitMessage(
		*processId,
		commons.GetSingleValue(c.Context, commons.User),
		newFileCreated,
	)

	worktree, err := c.repo.Worktree()
	if err != nil {
		return err
	}
	_, err = c.commitAndTagFile(worktree, *processId, commitMessage)
	if err != nil {
		if err == git.ErrTagExists {
			log.Warn().Msg("Specified process version alredy exists, nothing will be updated")
		}
		log.Error().Msg(err.Error())
		return err
	}
	return nil
}

// DownloadProcess finds and returns save process file based in processId and version
func (c *GitClient) DownloadProcess(request *pb.ProcessDownloadRequest, w *io.PipeWriter) (*pb.ProcessMetadata, error) {
	log.Info().Msgf("Getting history for process: %s", request.ProcessId)
	processId, err := pb.ParseProcessId(request.ProcessId)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	processVersionCommit, err := c.getProcessCommitForTag(processId)
	if err != nil {
		log.Err(err)
		return nil, err
	}
	err = c.streamProcessFile(processId, processVersionCommit, w)
	if err != nil {
		log.Err(err)
	}
	// TODO use real values - finish missing
	return &pb.ProcessMetadata{
			ProcessId: processId.FullProcessIdWithVersionTag(),
			Filename:  processId.Process,
			Encoding:  0,
			Type:      0,
		},
		err
}

// GetProcessHistory finds upload action history for given processId
func (c *GitClient) GetProcessHistory(processId *pb.ProcessId) ([]UploadActionSummary, error) {
	fileHistory, err := c.getProcessCommitHistory(processId)
	if err != nil {
		return nil, err
	}
	var commits = make([]UploadActionSummary, 0)
	for _, commit := range fileHistory {
		c, err := parseCommitMessage(commit.Message)
		if err != nil {
			log.Error().Msg(err.Error())
			return nil, err
		}
		commits = append(commits, c)
	}
	if err != nil {
		return nil, err
	}
	return commits, err
}

// streamProcessFile
// processCommit is expected to contain only commit regarding specified processId
// target process file is streamed to pipeWriter
func (c *GitClient) streamProcessFile(processId *pb.ProcessId, processCommit *object.Commit, pipeWriter *io.PipeWriter) error {
	commitMessage, err := parseCommitMessage(processCommit.Message)
	if err != nil {
		log.Err(err)
		return err
	}
	processForVersion := struct {
		gitCommit     *object.Commit
		commitMessage *UploadActionSummary
	}{processCommit, &commitMessage}
	// if commit with version is found - checkout to that commit
	if processForVersion.gitCommit != nil && processForVersion.commitMessage != nil {
		workTree, err := c.repo.Worktree()
		if err != nil {
			log.Err(err)
			return err
		}
		err = workTree.Checkout(&git.CheckoutOptions{Hash: processForVersion.gitCommit.Hash})
		if err != nil {
			log.Err(err)
			return err
		}
		// after reading - checkout back to head
		defer func(name string) {
			head, err := c.repo.Head()
			if err != nil {
				log.Err(err)
			}
			err = workTree.Checkout(&git.CheckoutOptions{Hash: head.Hash()})
		}(*processId.ProcessWithinProject())

		// read desired version of file
		fd, err := os.Open(filepath.Join(c.tenantGitRepoPath, *processId.ProcessWithinProject()))
		return commons.StreamFileToPipe(fd, c.ChunkSize, pipeWriter)
	}
	return errors.New("process not found")
}

// getProcessCommitForTag returns slice of all commits that affected given processId (process version is searched by tag)
func (c *GitClient) getProcessCommitForTag(processId *pb.ProcessId) (*object.Commit, error) {
	tag, err := c.repo.Tag(processId.FullProcessIdWithVersionTag())
	if err != nil {
		if err == git.ErrTagNotFound {
			log.Error().Msg("Specified version of process not found")
			return nil, err
		}
		log.Err(err)
		return nil, err
	}
	return c.repo.CommitObject(tag.Hash())
}

// getProcessCommitHistory returns slice of all commits that affected given processId
func (c *GitClient) getProcessCommitHistory(processId *pb.ProcessId) ([]*object.Commit, error) {
	// git log
	head, err := c.repo.Head()
	fileHistory, err := c.repo.Log(&git.LogOptions{
		From: head.Hash(),
	})
	fileCommits := commitIterToSlice(fileHistory)
	if len(fileCommits) <= 0 {
		log.Warn().Msgf("No history for process: %s", processId.ProcessWithinProject())
		return nil, err
	}
	// sort by commit time
	sort.Slice(fileCommits, func(i, j int) bool {
		return fileCommits[i].Author.When.After(fileCommits[j].Author.When)
	})
	// find commits asociated with requested processId
	commitsForProcess := lo.Filter(fileCommits, func(item *object.Commit, _ int) bool {
		message, err := parseCommitMessage(item.Message)
		if err != nil {
			log.Err(err)
			return false
		}
		allFileChanges := lo.Map(
			lo.Flatten([][]ProcessFile{message.AddedFiles, message.ModifiedFiled}),
			func(item ProcessFile, index int) string {
				return item.ProjectId + "/" + item.ProcessName
			},
		)
		return lo.Contains(allFileChanges, *processId.ProcessWithinProject())
	})
	return commitsForProcess, err
}

func (c *GitClient) createRepository() (*git.Repository, error) {
	return git.PlainInit(c.tenantGitRepoPath, false)
}

func (c *GitClient) commitAndTagFile(tree *git.Worktree, processId pb.ProcessId, commitMessage UploadActionSummary) (*string, error) {
	_, err := tree.Add(*processId.ProcessWithinProject())
	if err != nil {
		return nil, err
	}

	messageMarshalled, err := toml.Marshal(commitMessage)
	if err != nil {
		return nil, err
	}
	// TODO check after userinfo format for ctx will be implemented
	author := &object.Signature{
		Name:  commons.GetSingleValue(c.Context, commons.User),
		Email: commons.GetSingleValue(c.Context, commons.UserEmail),
		When:  time.Now(),
	}
	commit, err := tree.Commit(string(messageMarshalled), &git.CommitOptions{
		Author: author,
	})
	if err != nil {
		return nil, err
	}
	head, err := c.repo.Head()
	if err != nil {
		log.Err(err)
		return nil, err
	}
	versionTag := processId.FullProcessIdWithVersionTag()
	_, err = c.repo.CreateTag(versionTag, head.Hash(), nil)
	if err != nil {
		return nil, err
	}
	res := commit.String()
	return &res, nil
}

func commitIterToSlice(iter object.CommitIter) []*object.Commit {
	var res = make([]*object.Commit, 0)
	iter.ForEach(func(commit *object.Commit) error {
		res = append(res, commit)
		return nil
	})
	return res
}

func parseCommitMessage(message string) (UploadActionSummary, error) {
	var c UploadActionSummary
	err := toml.Unmarshal([]byte(message), &c)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	return c, nil
}

func createCommitMessage(
	processId pb.ProcessId,
	updatedBy string,
	newFileCreated bool) UploadActionSummary {
	newFile := ProcessFile{
		ProcessName: processId.Process,
		ProjectId:   processId.Project,
		Version:     processId.Version,
		LastUpdate:  time.Now().Unix(),
	}

	if newFileCreated {
		return UploadActionSummary{
			AddedFiles:    []ProcessFile{newFile},
			DeletedFiles:  nil,
			ModifiedFiled: nil,
			UpdatedBy:     updatedBy,
		}
	} else {
		return UploadActionSummary{
			AddedFiles:    nil,
			DeletedFiles:  nil,
			ModifiedFiled: []ProcessFile{newFile},
			UpdatedBy:     updatedBy,
		}
	}
}
