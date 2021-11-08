package appinfo

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
