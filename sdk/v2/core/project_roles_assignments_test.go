package core

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brigadecore/brigade/sdk/v2/authz"
	rmTesting "github.com/brigadecore/brigade/sdk/v2/internal/restmachinery/testing" // nolint: lll
	libAuthz "github.com/brigadecore/brigade/sdk/v2/lib/authz"
	"github.com/stretchr/testify/require"
)

func TestNewProjectRoleAssignmentsClient(t *testing.T) {
	client := NewProjectRoleAssignmentsClient(
		rmTesting.TestAPIAddress,
		rmTesting.TestAPIToken,
		nil,
	)
	require.IsType(t, &projectRoleAssignmentsClient{}, client)
	rmTesting.RequireBaseClient(
		t,
		client.(*projectRoleAssignmentsClient).BaseClient,
	)
}

func TestProjectRoleAssignmentsClientGrant(t *testing.T) {
	testProjectRoleAssignment := ProjectRoleAssignment{
		Role:      libAuthz.Role("ceo"),
		ProjectID: "bluebook",
		Principal: libAuthz.PrincipalReference{
			Type: authz.PrincipalTypeUser,
			ID:   "tony@starkindustries.com",
		},
	}
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				defer r.Body.Close()
				require.Equal(t, http.MethodPost, r.Method)
				require.Equal(t, "/v2/project-role-assignments", r.URL.Path)
				bodyBytes, err := ioutil.ReadAll(r.Body)
				require.NoError(t, err)
				projectRoleAssignment := ProjectRoleAssignment{}
				err = json.Unmarshal(bodyBytes, &projectRoleAssignment)
				require.NoError(t, err)
				require.Equal(t, testProjectRoleAssignment, projectRoleAssignment)
				w.WriteHeader(http.StatusOK)
			},
		),
	)
	defer server.Close()
	client := NewProjectRoleAssignmentsClient(
		server.URL,
		rmTesting.TestAPIToken,
		nil,
	)
	err := client.Grant(context.Background(), testProjectRoleAssignment)
	require.NoError(t, err)
}

func TestProjectRoleAssignmentsClientRevoke(t *testing.T) {
	testProjectRoleAssignment := ProjectRoleAssignment{
		Role:      libAuthz.Role("ceo"),
		ProjectID: "bluebook",
		Principal: libAuthz.PrincipalReference{
			Type: authz.PrincipalTypeUser,
			ID:   "tony@starkindustries.com",
		},
	}
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodDelete, r.Method)
				require.Equal(t, "/v2/project-role-assignments", r.URL.Path)
				require.Equal(
					t,
					testProjectRoleAssignment.Role,
					libAuthz.Role(r.URL.Query().Get("role")),
				)
				require.Equal(
					t,
					testProjectRoleAssignment.ProjectID,
					r.URL.Query().Get("projectID"),
				)
				require.Equal(
					t,
					testProjectRoleAssignment.Principal.Type,
					libAuthz.PrincipalType(r.URL.Query().Get("principalType")),
				)
				require.Equal(
					t,
					testProjectRoleAssignment.Principal.ID,
					r.URL.Query().Get("principalID"),
				)
				w.WriteHeader(http.StatusOK)
			},
		),
	)
	defer server.Close()
	client := NewProjectRoleAssignmentsClient(
		server.URL,
		rmTesting.TestAPIToken,
		nil,
	)
	err := client.Revoke(context.Background(), testProjectRoleAssignment)
	require.NoError(t, err)
}
