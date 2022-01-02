package auth

import "testing"

type testEnforce struct {
	claim  ApplicationClaim
	group  ApiGroups
	action ApiAction
}

func TestEnforce(t *testing.T) {
	testContexts := []testEnforce{
		{
			// case want to read & can read
			claim: ApplicationClaim{
				UserAuthorization: map[ApiGroups][]ApiAction{
					UsersApiGroup: {ReadAction},
				},
			},
			group:  UsersApiGroup,
			action: ReadAction,
		},
		{
			// case want to read & is admin
			claim: ApplicationClaim{
				UserAuthorization: map[ApiGroups][]ApiAction{
					UsersApiGroup: {AllAction},
				},
			},
			group:  UsersApiGroup,
			action: ReadAction,
		},
		{
			// case want to read on self action & is admin
			claim: ApplicationClaim{
				UserAuthorization: map[ApiGroups][]ApiAction{
					UsersApiGroup: {AllAction},
				},
			},
			group:  UsersApiGroup,
			action: ReadSelfAction,
		},
		{
			// case want to read on self action & has elevated
			claim: ApplicationClaim{
				UserAuthorization: map[ApiGroups][]ApiAction{
					UsersApiGroup: {ReadAction},
				},
			},
			group:  UsersApiGroup,
			action: ReadSelfAction,
		},
		{
			// case want to read action & can read only on self
			claim: ApplicationClaim{
				UserAuthorization: map[ApiGroups][]ApiAction{
					UsersApiGroup: {ReadSelfAction},
				},
			},
			group:  UsersApiGroup,
			action: ReadAction,
		},
		{
			// case want to read action & has no action
			claim: ApplicationClaim{
				UserAuthorization: map[ApiGroups][]ApiAction{
					UsersApiGroup: {},
				},
			},
			group:  UsersApiGroup,
			action: ReadAction,
		},
	}
	expected := []bool{
		true,
		true,
		true,
		true,
		false,
		false,
	}

	for i, context := range testContexts {
		result, _ := enforce(context.claim, context.group, context.action, nil, []SelfActionHandler{})
		if result != expected[i] {
			t.Errorf("expected %v, got %v at iteration %v", expected[i], result, i+1)
		}
	}
}
