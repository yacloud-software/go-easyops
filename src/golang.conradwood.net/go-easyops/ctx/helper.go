package ctx

import (
	"fmt"
	"golang.conradwood.net/apis/auth"
	"golang.conradwood.net/go-easyops/common"
)

func describeUser(su *auth.SignedUser) string {
	u := common.VerifySignedUser(su)
	if u == nil {
		return "[nouser]"
	}
	return fmt.Sprintf("[%s/%s]", u.ID, u.Email)
}
