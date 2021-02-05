package vault_plugin_database_redisenterprise

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	dbtesting "github.com/hashicorp/vault/sdk/database/dbplugin/v5/testing"
)



func TestRedisEnterpriseDB_NewUser_Without_Database(t *testing.T) {

	database := ""

	db := setupRedisEnterpriseDB(t, database, enableACL)

	createRequest := dbplugin.NewUserRequest{
		UsernameConfig: dbplugin.UsernameMetadata{
			DisplayName: "tester_new_without_db",
			RoleName:    "test",
		},
		Statements: dbplugin.Statements{
			Commands: []string{`{"role":"DB Member"}`},
		},
		Password:   "testing",
		Expiration: time.Now().Add(time.Minute),
	}

	if enableACL {
		createRequest.Statements.Commands = []string{`{"acl":"Not Dangerous"}`}
	}

	res := dbtesting.AssertNewUser(t, db, createRequest)

	assertUserExists(t, url, username, password, res.Username)

	teardownUserFromDatabase(t, db, res.Username)
}

func TestRedisEnterpriseDB_NewUser_With_Database(t *testing.T) {

	db := setupRedisEnterpriseDB(t, database, enableACL)

	createRequest := dbplugin.NewUserRequest{
		UsernameConfig: dbplugin.UsernameMetadata{
			DisplayName: "tester_new_with_db",
			RoleName:    "test",
		},
		Statements: dbplugin.Statements{
			Commands: []string{`{"role":"DB Member"}`},
		},
		Password:   "testing",
		Expiration: time.Now().Add(time.Minute),
	}

	if enableACL {
		createRequest.Statements.Commands = []string{`{"acl":"Not Dangerous"}`}
	}

	res := dbtesting.AssertNewUser(t, db, createRequest)

	assertUserExists(t, url, username, password, res.Username)

	teardownUserFromDatabase(t, db, res.Username)
}

func TestRedisEnterpriseDB_NewUser_With_Database_With_ACL(t *testing.T) {

	enableACL := true

	db := setupRedisEnterpriseDB(t, database, enableACL)

	createRequest := dbplugin.NewUserRequest{
		UsernameConfig: dbplugin.UsernameMetadata{
			DisplayName: "tester_new_with_db_with_acl",
			RoleName:    "test",
		},
		Statements: dbplugin.Statements{
			Commands: []string{`{"role":"DB Member"}`},
		},
		Password:   "testing",
		Expiration: time.Now().Add(time.Minute),
	}

	if enableACL {
		createRequest.Statements.Commands = []string{`{"acl":"Not Dangerous"}`}
	}

	res := dbtesting.AssertNewUser(t, db, createRequest)

	assertUserExists(t, url, username, password, res.Username)

	teardownUserFromDatabase(t, db, res.Username)
}

func TestRedisEnterpriseDB_NewUser_Detect_Cluster_Errors(t *testing.T) {

	database := ""

	db := setupRedisEnterpriseDB(t, database, enableACL)

	for _, spec := range [][]string{{"", ""}, {"", "Not Dangerous"}} {
		createReq := newUserRequest(spec[0], spec[1])

		ctx, cancel := context.WithTimeout(context.Background(), context_timeout)
		defer cancel()

		_, err := db.NewUser(ctx, createReq)
		if err == nil {
			t.Fatalf("Failed to detect NewUser (cluster) error with role (%s) and acl (%s)", spec[0], spec[1])
		}
	}
}

func TestRedisEnterpriseDB_NewUser_Detect_Errors_With_Database_Without_ACL(t *testing.T) {

	db := setupRedisEnterpriseDB(t, database, enableACL)

	for _, spec := range [][]string{{"", ""}, {"", "Not Dangerous"}, {"garbage", ""}} {
		createReq := newUserRequest(spec[0], spec[1])

		ctx, cancel := context.WithTimeout(context.Background(), context_timeout)
		defer cancel()

		_, err := db.NewUser(ctx, createReq)
		if err == nil {
			t.Fatalf("Failed to detect NewUser (database, no acl_only) error with role (%s) and acl (%s)", spec[0], spec[1])
		}
	}
}

func TestRedisEnterpriseDB_NewUser_Detect_Errors_With_Database_With_ACL(t *testing.T) {

	enableACL := true

	db := setupRedisEnterpriseDB(t, database, enableACL)

	for _, spec := range [][]string{{"", ""}, {"", "garbage"}} {
		createReq := newUserRequest(spec[0], spec[1])

		ctx, cancel := context.WithTimeout(context.Background(), context_timeout)
		defer cancel()

		_, err := db.NewUser(ctx, createReq)
		if err == nil {
			t.Fatalf("Failed to detect NewUser (database, acl_only) error with role (%s) and acl (%s)", spec[0], spec[1])
		}

	}
}
