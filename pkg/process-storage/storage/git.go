package storage

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
	"golang.org/x/net/context"
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

// TODO refactor this method
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

// DownloadProcess finds and returns save process file based in processId and version
func (c *GitClient) DownloadProcess(request *pb.ProcessDownloadRequest) (*pb.ProcessMetadata, []byte, error) {
	log.Info().Msgf("Getting history for process: %s with version: %s", request.ProcessId, request.ProcessVersion)
	processId, err := pb.ParseProcessId(request.ProcessId)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	processCommits, err := c.getProcessCommitHistory(processId)
	file, err := c.getProcessFile(processId, request.ProcessVersion, processCommits)
	if err != nil {
		log.Err(err)
	}
	// TODO use real values - not values from request
	return &pb.ProcessMetadata{
			ProcessId: processId.FullProcessId(),
			Filename:  processId.Process,
			Version:   request.ProcessVersion,
			Encoding:  0,
			Type:      0,
		},
		file,
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

// TODO refactor this method
// getProcessFile returns file byte array of process with specified version
// processCommits is expected to contain only commits regarding specified processId
func (c *GitClient) getProcessFile(processId *pb.ProcessId, version string, processCommits []*object.Commit) ([]byte, error) {
	processForVersion := struct {
		gitCommit     *object.Commit
		commitMessage *UploadActionSummary
	}{}
	// search all commits for process
	for _, commit := range processCommits {
		cm, err := parseCommitMessage(commit.Message)
		if err != nil {
			log.Err(err)
		}
		// find process in commit metadata
		fileChangesInCommit := lo.Flatten([][]ProcessFile{cm.AddedFiles, cm.ModifiedFiled})
		// find the right version
		temp := lo.Filter(fileChangesInCommit, func(addedFile ProcessFile, i int) bool {
			return addedFile.Version == version
		})
		if len(temp) > 0 {
			processForVersion.gitCommit = commit
			processForVersion.commitMessage = &cm
		}
	}
	// if commit with version is found - checkout to that commit
	if processForVersion.gitCommit != nil && processForVersion.commitMessage != nil {
		workTree, err := c.repo.Worktree()
		if err != nil {
			log.Err(err)
			return nil, err
		}
		err = workTree.Checkout(&git.CheckoutOptions{Hash: processForVersion.gitCommit.Hash})
		if err != nil {
			log.Err(err)
			return nil, err
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
		f, err := os.Open(filepath.Join(c.tenantGitRepoPath, *processId.ProcessWithinProject()))
		stat, err := f.Stat()
		if err != nil {
			log.Err(err)
			return nil, err
		}
		// read into []byte
		res := make([]byte, stat.Size())
		if err != nil {
			log.Error().Msg(err.Error())
			return nil, err
		}
		_, err = f.Read(res)
		if err != nil {
			log.Err(err)
			return nil, err
		}
		return res, err
	}
	return nil, errors.New("process not found")
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

func (c *GitClient) commitFile(tree *git.Worktree, file string, commitMessage UploadActionSummary) (string, error) {
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

func createCommitMessageForNewProcess(
	processName string,
	project string,
	version string,
	updatedBy string,
	updated bool) UploadActionSummary {
	newFile := ProcessFile{
		ProcessName: processName,
		ProjectId:   project,
		Version:     version,
		LastUpdate:  time.Now().Unix(),
	}

	if updated {
		return UploadActionSummary{
			AddedFiles:    nil,
			DeletedFiles:  nil,
			ModifiedFiled: []ProcessFile{newFile},
			UpdatedBy:     updatedBy,
		}
	} else {
		return UploadActionSummary{
			AddedFiles:    []ProcessFile{newFile},
			DeletedFiles:  nil,
			ModifiedFiled: nil,
			UpdatedBy:     updatedBy,
		}
	}
}
