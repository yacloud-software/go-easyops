package ctx

import (
	"context"

	ge "golang.conradwood.net/apis/goeasyops"
)

// get authtags from context
func AuthTags(ctx context.Context) []string {
	ls := GetLocalState(ctx)
	return ls.AuthTags()
}

// true if context has given authtag
func HasAuthTag(ctx context.Context, tag string) bool {
	ls := GetLocalState(ctx)
	tags := ls.AuthTags()
	for _, t := range tags {
		if tag == t {
			return true
		}
	}
	return false
}

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

func IsDebug(ctx context.Context) bool {
	ls := GetLocalState(ctx)
	if ls == nil {
		return false
	}
	return ls.Debug()
}
func IsExperimentEnabled(ctx context.Context, name string) bool {
	ls := GetLocalState(ctx)
	if ls == nil {
		return false
	}
	for _, e := range ls.Experiments() {
		if e.Name == name {
			return true
		}
	}
	return false
}
func GetExperiments(ctx context.Context) []*ge.Experiment {
	ls := GetLocalState(ctx)
	if ls == nil {
		return nil
	}
	return ls.Experiments()
}
