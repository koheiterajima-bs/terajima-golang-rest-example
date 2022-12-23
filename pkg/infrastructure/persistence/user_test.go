package persistence

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/taaaaakahiro/golang-rest-example/pkg/domain/entity"
	testfixtures "github.com/taaaaakahiro/golang-rest-example/test_fixtures"
)

const userTable = "users"

func TestUserRepository_GetUser(t *testing.T) {
	db, _ := sql.Open("mysql", mysqlDsn)

	// CleanUp
	if err := testfixtures.TruncateTables(db, []string{userTable}); err != nil {
		t.Errorf("truncate error: %s\n", err.Error())
	}
	t.Cleanup(func() {
		if err := testfixtures.TruncateTables(db, []string{userTable}); err != nil {
			t.Errorf("truncate error: %s\n", err.Error())
		}
	})

	// Fixture
	users := []struct {
		Id   string
		Name string
	}{
		{Id: "1", Name: "user1"},
		{Id: "2", Name: "user2"},
	}

	for _, user := range users {
		if err := testfixtures.InsertTable(db, "users", interface{}(user)); err != nil {
			t.Errorf("insert error: %s\n", err.Error())
		}
	}

	// TestCase
	tests := []struct {
		name    string
		userID  string
		want    *entity.User
		wantErr error
	}{
		{
			name:   "ok",
			userID: "1",
			want: &entity.User{
				ID: "1", Name: "user1",
			},
			wantErr: nil,
		},
		{
			name:   "ok",
			userID: "2",
			want: &entity.User{
				ID: "2", Name: "user2",
			},
			wantErr: nil,
		},
		{
			name:    "notExistUserId",
			userID:  "999",
			want:    nil,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		c := context.Background()
		t.Run(tt.name, func(t *testing.T) {
			got, err := userRepo.GetUser(c, tt.userID)
			if diff := cmp.Diff(tt.wantErr, err); len(diff) != 0 {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want, got); len(diff) != 0 {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}

}

func TestUserRepository_CreateUser(t *testing.T) {
	db, _ := sql.Open("mysql", mysqlDsn)

	// CleanUp
	if err := testfixtures.TruncateTables(db, []string{userTable}); err != nil {
		t.Errorf("truncate error: %s\n", err.Error())
	}
	t.Cleanup(func() {
		if err := testfixtures.TruncateTables(db, []string{userTable}); err != nil {
			t.Errorf("truncate error: %s\n", err.Error())
		}
	})

	tests := []struct {
		name    string
		args    string
		want    *entity.User
		wantErr error
	}{
		{
			name: "ok",
			args: "createUser1",
			want: &entity.User{
				ID:   "1",
				Name: "createUser1",
			},
			wantErr: nil,
		},
		{
			name: "ok",
			args: "createUser2",
			want: &entity.User{
				ID:   "2",
				Name: "createUser2",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		c := context.Background()
		t.Run(tt.name, func(t *testing.T) {
			id, err := userRepo.CreateUser(c, tt.args)
			var got *entity.User
			if err == nil {
				stmtOut, err := db.Prepare("select id, name from users where id = ?")
				assert.NoError(t, err)
				var user entity.User
				err = stmtOut.QueryRow(*id).Scan(&user.ID, &user.Name)
				assert.NoError(t, err)
				got = &user
			}

			if diff := cmp.Diff(tt.wantErr, err); len(diff) != 0 {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want, got); len(diff) != 0 {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUserRepository_UpdateUser(t *testing.T) {
	db, _ := sql.Open("mysql", mysqlDsn)

	// Fixture
	users := []struct {
		Id   string
		Name string
	}{
		{Id: "1", Name: "user1"},
		{Id: "2", Name: "user2"},
	}

	type args struct {
		userID string
		name   string
	}

	// TestCase
	tests := []struct {
		name    string
		args    args
		want    *entity.User
		wantErr error
	}{
		{
			name: "ok",
			args: args{
				userID: "1",
				name:   "updateUserId1",
			},
			want: &entity.User{
				ID: "1", Name: "updateUserId1",
			},
			wantErr: nil,
		},
		{
			name: "ok",
			args: args{
				userID: "2",
				name:   "updateUserId2",
			},
			want: &entity.User{
				ID: "2", Name: "updateUserId2",
			},
			wantErr: nil,
		},
		{
			name: "notExistUserId",
			args: args{
				userID: "999",
				name:   "updateUserId1",
			},
			want:    &entity.User{},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		if err := testfixtures.TruncateTables(db, []string{userTable}); err != nil {
			t.Errorf("truncate error: %s\n", err.Error())
		}

		t.Cleanup(func() {
			if err := testfixtures.TruncateTables(db, []string{userTable}); err != nil {
				t.Errorf("truncate error: %s\n", err.Error())
			}
		})

		for _, user := range users {
			if err := testfixtures.InsertTable(db, "users", interface{}(user)); err != nil {
				t.Errorf("insert error: %s\n", err.Error())
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			c := context.Background()
			err := userRepo.UpdateUser(c, tt.args.userID, tt.args.name)
			var got *entity.User
			if err == nil {
				stmtOut, _ := db.Prepare("select id, name from users where id = ?")
				var user entity.User
				_ = stmtOut.QueryRowContext(c, tt.args.userID).Scan(&user.ID, &user.Name)
				got = &user
			}

			if diff := cmp.Diff(tt.wantErr, err); len(diff) != 0 {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want, got); len(diff) != 0 {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUserRepository_DeleteUser(t *testing.T) {
	db, _ := sql.Open("mysql", mysqlDsn)

	// Fixture
	users := []struct {
		Id   string
		Name string
	}{
		{Id: "1", Name: "user1"},
		{Id: "2", Name: "user2"},
	}

	// TestCase
	tests := []struct {
		name    string
		userID  string
		want    *entity.User
		wantErr error
	}{
		{
			name:    "ok",
			userID:  "1",
			want:    &entity.User{},
			wantErr: nil,
		},
		{
			name:    "ok",
			userID:  "2",
			want:    &entity.User{},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		if err := testfixtures.TruncateTables(db, []string{userTable}); err != nil {
			t.Errorf("truncate error: %s\n", err.Error())
		}

		t.Cleanup(func() {
			if err := testfixtures.TruncateTables(db, []string{userTable}); err != nil {
				t.Errorf("truncate error: %s\n", err.Error())
			}
		})

		for _, user := range users {
			if err := testfixtures.InsertTable(db, "users", interface{}(user)); err != nil {
				t.Errorf("insert error: %s\n", err.Error())
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			c := context.Background()
			err := userRepo.DeleteUser(c, tt.userID)
			var got *entity.User
			if err == nil {
				stmtOut, _ := db.Prepare("select id, name from users where id = ?")
				var user entity.User
				_ = stmtOut.QueryRowContext(c, tt.userID).Scan(&user.ID)
				got = &user
			}

			if diff := cmp.Diff(tt.wantErr, err); len(diff) != 0 {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want, got); len(diff) != 0 {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
