package storage

import (
	"errors"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"os"
	"path/filepath"
	"pirs.io/commons"
	"pirs.io/commons/parsers"
	pb "pirs.io/process-storage/grpc"
	"sort"
	"strings"
	"time"
)

var (
	log = commons.GetLoggerFor("git")
)

var (
	ErrProcessNotFound      commons.ErrorResponse = status.New(codes.InvalidArgument, "process not found").Err()
	ErrProcessAlreadyExists commons.ErrorResponse = status.New(codes.InvalidArgument, "process with specified version already exists").Err()
)

type GitClient struct {
	Context      context.Context
	RepoRootPath string
	Tenant       string
	ChunkSize    int64
	DataFS       billy.Filesystem
	DotGitFS     billy.Filesystem

	repo              *git.Repository
	tenantGitRepoPath string
	lruCache          *cache.ObjectLRU
	storage           *filesystem.Storage
}

func (c *GitClient) InitializeStorage() error {
	c.tenantGitRepoPath = filepath.Join(c.RepoRootPath, c.Tenant)

	// if not specified create osFs with .git folder
	if c.DotGitFS == nil {
		c.DotGitFS = osfs.New(filepath.Join(c.DataFS.Root(), ".git"))
	}
	c.lruCache = cache.NewObjectLRU(2048)
	c.storage = filesystem.NewStorage(c.DotGitFS, c.lruCache)

	r, err := git.Open(c.storage, c.DataFS)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			log.Info().Msgf("Git storage folder is empty, creating repository for tenant: %s", c.Tenant)
			r, err = c.createRepository()
			if err != nil {
				log.Err(err)
				return err
			}
		}
	}
	c.repo = r
	return err
}

// SaveFile saves given process
func (c *GitClient) SaveFile(processMetadata *pb.ProcessMetadata, file []byte) error {
	processId, err := parsers.ParseProcessId(processMetadata.ProcessId)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	filePath := filepath.Join(processId.Tenant, processId.Project, processId.Process)

	var commitMessage UploadActionSummary
	newFileCreated := false
	// take actions based on process file existence
	if _, err := c.DataFS.Stat(filePath); err == nil {
		// path/to/whatever exists
		f, err := c.DataFS.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
		_, err = f.Write(file)
		if err != nil {
			log.Error().Msg(err.Error())
			return err
		}

	} else if errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does *not* exist
		f, err := c.DataFS.Create(filePath)
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
			log.Warn().Msg("Specified process version already exists, nothing will be updated")
			return ErrProcessAlreadyExists
		}
		log.Error().Msg(err.Error())
		return err
	}
	log.Info().Msgf("process %s successfully saved", processId.FullProcessIdWithVersionTag())
	return nil
}

// DownloadProcess finds and returns save process file based in processId and version
func (c *GitClient) DownloadProcess(request *pb.ProcessDownloadRequest, w *io.PipeWriter) (*pb.ProcessMetadata, error) {
	log.Info().Msgf("Getting history for process: %s", request.ProcessId)
	processId, err := parsers.ParseProcessId(request.ProcessId)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	processVersionCommit, err := c.getProcessCommitForTag(processId)
	if err != nil {
		log.Err(err)
		return nil, err
	}
	err = c.streamProcessFile(processId, processVersionCommit, w)
	if err != nil {
		log.Err(err)
		return nil, err
	}
	// TODO use real values - finish missing
	return &pb.ProcessMetadata{
			ProcessId: processId.FullProcessIdWithVersionTag(),
			Filename:  processId.Process,
			Encoding:  0,
			Type:      0,
		},
		nil
}

// GetProcessHistory finds upload action history for given processId
func (c *GitClient) GetProcessHistory(processId *parsers.ProcessId) ([]ProcessFile, error) {
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
	res := lo.FlatMap(commits, func(item UploadActionSummary, index int) []ProcessFile {
		return lo.Flatten([][]ProcessFile{item.ModifiedFiled, item.AddedFiles})
	})
	return res, err
}

// streamProcessFile
// processCommit is expected to contain only commit regarding specified processId
// target process file is streamed to pipeWriter
func (c *GitClient) streamProcessFile(processId *parsers.ProcessId, processCommit *object.Commit, pipeWriter *io.PipeWriter) error {
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
		defer func() {
			head, err := c.repo.Head()
			if err != nil {
				log.Err(err)
			}
			err = workTree.Checkout(&git.CheckoutOptions{Hash: head.Hash()})
		}()

		// read desired version of file
		fd, err := c.DataFS.Open(*processId.ProcessWithinProject())
		go func() {
			err := commons.StreamFileToPipe(fd, c.ChunkSize, pipeWriter)
			if err != nil {
				log.Err(err)
			}
		}()
	}
	return nil
}

// getProcessCommitForTag returns slice of all commits that affected given processId (process version is searched by tag)
func (c *GitClient) getProcessCommitForTag(processId *parsers.ProcessId) (*object.Commit, error) {
	tag, err := c.repo.Tag(processId.FullProcessIdWithVersionTag())
	if err != nil {
		if err == git.ErrTagNotFound {
			log.Error().Msg("Specified version of process not found")
			return nil, ErrProcessNotFound
		}
		log.Err(err)
		return nil, err
	}
	return c.repo.CommitObject(tag.Hash())
}

// getProcessCommitHistory returns slice of all commits that affected given processId
func (c *GitClient) getProcessCommitHistory(processId *parsers.ProcessId) ([]*object.Commit, error) {
	// git log
	tags, err := c.repo.Tags()
	tagsForProcess := make([]plumbing.Hash, 0)
	err = tags.ForEach(func(reference *plumbing.Reference) error {
		if strings.Contains(reference.String(), processId.FullProcessIdWithoutVersionTag()) {
			tagsForProcess = append(tagsForProcess, reference.Hash())
		}
		log.Info().Msg(reference.String())
		return nil
	})
	if err != nil {
		return nil, err
	}

	fileCommits := lo.FilterMap(tagsForProcess, func(item plumbing.Hash, index int) (*object.Commit, bool) {
		res, e := c.repo.CommitObject(item)
		return res, e == nil
	})

	if len(fileCommits) <= 0 {
		log.Warn().Msgf("No history for process: %s", *processId.ProcessWithinProject())
		return nil, err
	}
	// sort by commit time
	sort.Slice(fileCommits, func(i, j int) bool {
		return fileCommits[i].Author.When.After(fileCommits[j].Author.When)
	})
	return fileCommits, err
}

func (c *GitClient) createRepository() (*git.Repository, error) {
	return git.Init(c.storage, c.DataFS)
}

func (c *GitClient) commitAndTagFile(tree *git.Worktree, processId parsers.ProcessId, commitMessage UploadActionSummary) (*string, error) {
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

func parseCommitMessage(message string) (UploadActionSummary, error) {
	var c UploadActionSummary
	err := toml.Unmarshal([]byte(message), &c)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	return c, nil
}

func createCommitMessage(
	processId parsers.ProcessId,
	updatedBy string,
	newFileCreated bool) UploadActionSummary {
	newFile := ProcessFile{
		ProcessName: processId.Process,
		Project:     processId.Project,
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
