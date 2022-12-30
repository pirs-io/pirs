package storage

type Provider int

const (
	GitStorageProvider Provider = 1
)

type UploadActionSummary struct {
	AddedFiles    []ProcessFile
	DeletedFiles  []ProcessFile
	ModifiedFiled []ProcessFile
	UpdatedBy     string
}

type ProcessFile struct {
	ProcessName string
	ProjectId   string
	Version     int64
	LastUpdate  int64
}
