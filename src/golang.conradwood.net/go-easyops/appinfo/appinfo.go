/*
This package stores data about the build process.
During build variables are set, such as file location, build number and time.
These can be retrieved through this package.
*/
package appinfo

import (
	"fmt"
	"strconv"
)

type AppVersionInfo struct {
	Number         uint64
	Description    string
	Timestamp      int64
	RepositoryID   uint64
	RepositoryName string
	CommitID       string
	GitURL         string
	ArtefactID     uint64
}

var (
	appInfo           *AppVersionInfo // set by init() of some modules (overrides)
	OldAppInfo        *AppVersionInfo // set by cmdline
	LD_Number         string          // set by linker flags
	LD_Description    string
	LD_Timestamp      string
	LD_RepositoryID   string
	LD_ArtefactID     string
	LD_RepositoryName string
	LD_CommitID       string
	LD_GitURL         string
)

func RegisterAppInfo(avi *AppVersionInfo) {
	appInfo = avi
}
func AppInfo() *AppVersionInfo {
	if appInfo != nil {
		return appInfo
	}
	if LD_Number != "" {
		a := &AppVersionInfo{
			Number:         required_number(LD_Number),
			Description:    LD_Description,
			Timestamp:      int64(required_number(LD_Timestamp)),
			RepositoryID:   required_number(LD_RepositoryID),
			ArtefactID:     required_number(LD_ArtefactID),
			RepositoryName: LD_RepositoryName,
			CommitID:       LD_CommitID,
			GitURL:         LD_GitURL,
		}
		appInfo = a
		return a
	}
	if OldAppInfo != nil {
		return OldAppInfo
	}
	a := &AppVersionInfo{RepositoryName: "notset"}
	return a
}
func required_number(num string) uint64 {
	if num == "" {
		return 0
	}

	n, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		fmt.Printf("Not a number: \"%s\" (%s)\n", num, err)
		panic("invalid linker flags")
	}
	return uint64(n)
}
