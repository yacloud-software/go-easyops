package ctx

import (
	"context"
)

// get requestid from context
func GetRequestID(ctx context.Context) string {
	ls := GetLocalState(ctx)
	return ls.RequestID()
}

// get sessionid from context or "" (empty string) if none
func GetSessionID(ctx context.Context) string {
	sess := GetLocalState(ctx).Session()
	if sess == nil {
		return ""
	}
	return sess.SessionID
}

// get organisationid from context.Session or "" (empty string) if none
func GetOrganisationID(ctx context.Context) string {
	sess := GetLocalState(ctx).Session()
	if sess == nil {
		return ""
	}
	org := sess.Organisation
	if org == nil {
		return ""
	}
	return org.ID
}
