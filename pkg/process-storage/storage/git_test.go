package storage

import (
	"context"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/samber/lo"
	"io"
	"pirs.io/commons"
	pb "pirs.io/commons/parsers"
	"pirs.io/process-storage/grpc"
	"testing"
)

var (
	repoRootPath       = "/"
	tenant             = "tenant"
	chunkSize    int64 = 1024
)

type testData struct {
	name         string
	repoRootPath string
	tenant       string
	chunkSize    int64
	fs           billy.Filesystem
	dotGitFS     billy.Filesystem
	ctx          context.Context

	files []struct {
		processId string
		file      []byte
	}

	want string
}

func TestGitClient_InitializeStorage(t *testing.T) {
	tests := []testData{
		{name: "no git repo is created", repoRootPath: repoRootPath, tenant: tenant, chunkSize: chunkSize, fs: memfs.New(), dotGitFS: memfs.New()},
		{
			name:         "open existing repo",
			repoRootPath: repoRootPath,
			tenant:       tenant,
			chunkSize:    chunkSize,
			fs: func() billy.Filesystem {
				fs := memfs.New()
				lruCache := cache.NewObjectLRU(2048)
				storage := filesystem.NewStorage(fs, lruCache)
				_, _ = git.Init(storage, fs)
				return fs
			}(),
			dotGitFS: memfs.New(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gitClient := newGitClient(test)
			err := gitClient.InitializeStorage()
			if err != nil {
				if err != nil {
					t.Errorf("Failed to initialize storage")
				}
			}
		})
	}
}

func TestGitClient_SaveFile_ProcessNotExistsInSystem(t *testing.T) {
	test := testData{
		name:         "upload new process",
		repoRootPath: repoRootPath,
		tenant:       tenant,
		chunkSize:    chunkSize,
		fs:           memfs.New(),
		dotGitFS: func() billy.Filesystem {
			return memfs.New()
		}(),
		ctx: func() context.Context {
			ctx := context.Background()
			context.WithValue(ctx, commons.User, "user")
			return ctx
		}(),
		files: []struct {
			processId string
			file      []byte
		}{{processId: "organization.tenant.project.process:1", file: make([]byte, 64)}},
		want: "tenant/project/process",
	}
	gitClient := newGitClient(test)
	err := gitClient.InitializeStorage()
	if err != nil {
		t.Errorf("failed to initialize storage")
		t.Fatalf(err.Error())
	}
	for _, file := range test.files {
		err = gitClient.SaveFile(&grpc.ProcessMetadata{
			ProcessId:     file.processId,
			Filename:      "testProcess.xml",
			Encoding:      0,
			Type:          0,
			FileSizeBytes: 128,
		}, file.file)
	}
	checkProcessFileExistence(t, err, test)
}

func TestGitClient_SaveFile_UploadingDuplicate(t *testing.T) {
	test := testData{
		name:         "upload new process",
		repoRootPath: repoRootPath,
		tenant:       tenant,
		chunkSize:    chunkSize,
		fs:           memfs.New(),
		dotGitFS: func() billy.Filesystem {
			return memfs.New()
		}(),
		ctx: func() context.Context {
			ctx := context.Background()
			context.WithValue(ctx, commons.User, "user")
			return ctx
		}(),
		files: []struct {
			processId string
			file      []byte
		}{
			{processId: "organization.tenant.project.process:1", file: make([]byte, 64)},
			{processId: "organization.tenant.project.process:1", file: make([]byte, 64)},
		},
		want: "tenant/project/process",
	}
	gitClient := newGitClient(test)
	err := gitClient.InitializeStorage()
	if err != nil {
		t.Errorf("failed to initialize storage")
		t.Fatalf(err.Error())
	}
	for i, file := range test.files {
		err = gitClient.SaveFile(&grpc.ProcessMetadata{
			ProcessId:     file.processId,
			Filename:      "testProcess.xml",
			Encoding:      0,
			Type:          0,
			FileSizeBytes: 128,
		}, file.file)
		// first save
		if i == 0 {
			checkProcessFileExistence(t, err, test)
		}
		if i == 1 {
			checkProcessFileExistence(t, err, test, ErrProcessAlreadyExists)
			if err == nil {
				t.Fatalf("System didn't recognized process duplicate!")
			}
			if err != ErrProcessAlreadyExists {
				t.Fatalf("expected: %s, got: %s", ErrProcessAlreadyExists, err)
			}
		}
	}
}

func TestGitClient_DownloadProcess_ExistingProcess(t *testing.T) {
	test := testData{
		name:         "download existing process",
		repoRootPath: repoRootPath,
		tenant:       tenant,
		chunkSize:    chunkSize,
		fs:           memfs.New(),
		dotGitFS: func() billy.Filesystem {
			return memfs.New()
		}(),
		ctx: func() context.Context {
			ctx := context.Background()
			context.WithValue(ctx, commons.User, "user")
			return ctx
		}(),
		files: []struct {
			processId string
			file      []byte
		}{{processId: "organization.tenant.project.process:1", file: make([]byte, 64)}},
		want: "organization.tenant.project.process:1",
	}
	gitClient := newGitClient(test)
	err := gitClient.InitializeStorage()
	if err != nil {
		t.Errorf("failed to initialize storage")
		t.Fatalf(err.Error())
	}
	for _, file := range test.files {
		r, w := io.Pipe()
		err = gitClient.SaveFile(&grpc.ProcessMetadata{
			ProcessId:     file.processId,
			Filename:      "testProcess.xml",
			Encoding:      0,
			Type:          0,
			FileSizeBytes: 128,
		}, file.file)
		process, err := gitClient.DownloadProcess(&grpc.ProcessDownloadRequest{
			ProcessId: test.want,
		}, w)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if process.ProcessId != test.want {
			t.Fatalf("want: %s, got: %s", test.want, process.ProcessId)
		}
		file := make([]byte, 64)
		_, err = r.Read(file)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}
}

func TestGitClient_DownloadProcess_NonExistingProcess(t *testing.T) {
	test := testData{
		name:         "download non-existing process",
		repoRootPath: repoRootPath,
		tenant:       tenant,
		chunkSize:    chunkSize,
		fs:           memfs.New(),
		dotGitFS: func() billy.Filesystem {
			return memfs.New()
		}(),
		ctx: func() context.Context {
			ctx := context.Background()
			context.WithValue(ctx, commons.User, "user")
			return ctx
		}(),
		files: []struct {
			processId string
			file      []byte
		}{{processId: "organization.tenant.project.process:1", file: make([]byte, 64)}},
		want: "organization.tenant.project.process:2",
	}
	gitClient := newGitClient(test)
	err := gitClient.InitializeStorage()
	if err != nil {
		t.Errorf("failed to initialize storage")
		t.Fatalf(err.Error())
	}
	for _, file := range test.files {
		err = gitClient.SaveFile(&grpc.ProcessMetadata{
			ProcessId:     file.processId,
			Filename:      "testProcess.xml",
			Encoding:      0,
			Type:          0,
			FileSizeBytes: 128,
		}, file.file)
	}
	_, w := io.Pipe()
	_, err = gitClient.DownloadProcess(&grpc.ProcessDownloadRequest{
		ProcessId: test.want,
	}, w)
	if err != ErrProcessNotFound {
		t.Fatalf("want: %s, got: %s", test.want, err)
	}
}

func TestGitClient_GetProcessHistory_existingProcess(t *testing.T) {
	test := testData{
		name:         "upload new process",
		repoRootPath: repoRootPath,
		tenant:       tenant,
		chunkSize:    chunkSize,
		fs:           memfs.New(),
		dotGitFS: func() billy.Filesystem {
			return memfs.New()
		}(),
		ctx: func() context.Context {
			ctx := context.Background()
			context.WithValue(ctx, commons.User, "user")
			return ctx
		}(),
		files: []struct {
			processId string
			file      []byte
		}{
			{processId: "organization.tenant.project.process:1", file: make([]byte, 64)},
			{processId: "organization.tenant.project.process:2", file: make([]byte, 64)},
		},
	}
	gitClient := newGitClient(test)
	err := gitClient.InitializeStorage()
	if err != nil {
		t.Errorf("failed to initialize storage")
		t.Fatalf(err.Error())
	}
	for _, file := range test.files {
		err = gitClient.SaveFile(&grpc.ProcessMetadata{
			ProcessId:     file.processId,
			Filename:      "testProcess.xml",
			Encoding:      0,
			Type:          0,
			FileSizeBytes: 64,
		}, file.file)
	}
	history, err := gitClient.GetProcessHistory(&pb.ProcessId{
		Organization: "organization",
		Tenant:       "tenant",
		Project:      "project",
		Process:      "process",
		Version:      2,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	v1 := history[0].Version
	if v1 != 1 {
		t.Fatalf("wanted version: 1, got: %d", v1)
	}
	v2 := history[1].Version
	if v2 != 2 {
		t.Fatalf("wanted version: 2, got: %d", v2)
	}
}

func TestGitClient_GetProcessHistory_emptyRepo(t *testing.T) {
	test := testData{
		name:         "upload new process",
		repoRootPath: repoRootPath,
		tenant:       tenant,
		chunkSize:    chunkSize,
		fs:           memfs.New(),
		dotGitFS: func() billy.Filesystem {
			return memfs.New()
		}(),
		ctx: func() context.Context {
			ctx := context.Background()
			context.WithValue(ctx, commons.User, "user")
			return ctx
		}(),
		files: []struct {
			processId string
			file      []byte
		}{},
	}
	gitClient := newGitClient(test)
	err := gitClient.InitializeStorage()
	if err != nil {
		t.Errorf("failed to initialize storage")
		t.Fatalf(err.Error())
	}
	for _, file := range test.files {
		err = gitClient.SaveFile(&grpc.ProcessMetadata{
			ProcessId:     file.processId,
			Filename:      "testProcess.xml",
			Encoding:      0,
			Type:          0,
			FileSizeBytes: 64,
		}, file.file)
	}
	history, err := gitClient.GetProcessHistory(&pb.ProcessId{
		Organization: "organization",
		Tenant:       "tenant",
		Project:      "project",
		Process:      "process",
		Version:      2,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(history) != 0 {
		t.Fatalf("history present in empty repo!")
	}
}

func TestGitClient_GetProcessHistory_nonExistingProcess(t *testing.T) {
	test := testData{
		name:         "upload new process",
		repoRootPath: repoRootPath,
		tenant:       tenant,
		chunkSize:    chunkSize,
		fs:           memfs.New(),
		dotGitFS: func() billy.Filesystem {
			return memfs.New()
		}(),
		ctx: func() context.Context {
			ctx := context.Background()
			context.WithValue(ctx, commons.User, "user")
			return ctx
		}(),
		files: []struct {
			processId string
			file      []byte
		}{
			{processId: "organization.tenant.project.process:1", file: make([]byte, 64)},
			{processId: "organization.tenant.project.process:2", file: make([]byte, 64)},
		},
	}
	gitClient := newGitClient(test)
	err := gitClient.InitializeStorage()
	if err != nil {
		t.Errorf("failed to initialize storage")
		t.Fatalf(err.Error())
	}
	for _, file := range test.files {
		err = gitClient.SaveFile(&grpc.ProcessMetadata{
			ProcessId:     file.processId,
			Filename:      "testProcess.xml",
			Encoding:      0,
			Type:          0,
			FileSizeBytes: 64,
		}, file.file)
	}
	history, err := gitClient.GetProcessHistory(&pb.ProcessId{
		Organization: "organization",
		Tenant:       "tenant",
		Project:      "project",
		Process:      "non-existing-process",
		Version:      2,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(history) != 0 {
		t.Fatalf("found history for bad process!")
	}
}

func newGitClient(test testData) *GitClient {
	return &GitClient{
		Context:      test.ctx,
		RepoRootPath: test.repoRootPath,
		Tenant:       test.tenant,
		ChunkSize:    test.chunkSize,
		DataFS:       test.fs,
		DotGitFS:     test.dotGitFS,
	}
}

func checkProcessFileExistence(t *testing.T, err error, test testData, errorsToIgnore ...error) {
	errorMessages := lo.Map(errorsToIgnore, func(item error, index int) string {
		return item.Error()
	})
	if err != nil && !lo.Contains(errorMessages, err.Error()) {
		t.Fatalf(err.Error())
	}
	processFile, err := test.fs.Open(test.want)
	if err != nil {
		t.Fatalf(err.Error(), test.want)
	}
	if processFile == nil {
		t.Fatalf("process not saved! want: %s", test.want)
	}
}
