package authremote

import (
	"context"
	"strconv"

	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/errors"
)

/*
returns the partitionid for the user in this context.
a partition is either:
  - specific to a particular user
  - specific to an organisation

the information in the context determines which partition will be returned
a context without user information will always use partition 0

	Note: currently this is a bit of a stub. it only resolves the userid to a partition and does not consider the organisation

	PartitionIDs start from 100 upwards. this is meant to make it easier for tools to support custom partitions, such as "any user" or "no user" or so
*/
func PartitionID(ctx context.Context) (uint64, error) {
	u := auth.GetUser(ctx)
	if u == nil {
		// no user -> partition #0
		return 0, nil
	}
	id, err := strconv.ParseUint(u.ID, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err)
	}
	return id + 100, nil // just adding some random offset so it is guaranteed to be NOT the same as the userid
}
