package storage

type Provider int

const (
	GitStorageProvider Provider = 1
)

type CommitMessage struct {
	AddedFiles    []ProcessFile
	DeletedFiles  []ProcessFile
	ModifiedFiled []ProcessFile
	UpdatedBy     string
}

type ProcessFile struct {
	ProcessName string
	ProjectId   string
	Version     string
	LastUpdate  int64
}
