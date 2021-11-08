package cmdline

// this file is automatically updated by the build server

const (
	BUILD_NUMBER        = 1         // replaceme
	BUILD_DESCRIPTION   = "not set" //replaceme
	BUILD_TIMESTAMP     = 1         //replaceme
	BUILD_REPOSITORY_ID = 1         // replaceme
	BUILD_REPOSITORY    = "not set" // replaceme
	BUILD_COMMIT        = "not set" // replaceme
)

type AppVersionInfo struct {
	Number         uint64
	Description    string
	Timestamp      uint64
	RepositoryID   uint64
	RepositoryName string
	CommitID       string
}

func RegisterAppInfo(avi *AppVersionInfo) {
}
