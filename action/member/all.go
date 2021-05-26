package member

import (
	"context"

	"github.com/factly/mande-server/model"
	"github.com/factly/mande-server/util"
	"github.com/factly/x/middlewarex"
)

// All - to return all members
func All(ctx context.Context) (map[string]model.Member, error) {
	members := make(map[string]model.Member)

	organisationID, err := util.GetOrganisation(ctx)

	if err != nil {
		return members, err
	}

	userID, err := middlewarex.GetUser(ctx)

	if err != nil {
		return members, err
	}

	members = Mapper(organisationID, userID)

	return members, nil

}
