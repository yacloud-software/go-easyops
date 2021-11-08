package appinfo

type AppVersionInfo struct {
	Number         uint64
	Description    string
	Timestamp      int64
	RepositoryID   uint64
	RepositoryName string
	CommitID       string
}

var (
	AppInfo *AppVersionInfo
)

func RegisterAppInfo(avi *AppVersionInfo) {
}
