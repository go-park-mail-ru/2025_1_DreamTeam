package postgres

import (
	"context"
	"database/sql"
	"fmt"
	usermodels "skillForce/internal/models/user"
	"skillForce/pkg/logs"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestGetUserById_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	userId := 1

	expectedUser := usermodels.User{
		Email:     "test@example.com",
		Name:      "John Doe",
		HideEmail: true,
	}

	mock.ExpectQuery("SELECT email, name, hide_email FROM usertable WHERE id = \\$1").
		WithArgs(userId).
		WillReturnRows(sqlmock.NewRows([]string{"email", "name", "hide_email"}).
			AddRow(expectedUser.Email, expectedUser.Name, expectedUser.HideEmail))

	user, err := database.GetUserById(ctx, userId)
	require.NoError(t, err)
	require.Equal(t, &expectedUser, user)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserById_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	userId := 999

	mock.ExpectQuery("SELECT email, name, hide_email FROM usertable WHERE id = \\$1").
		WithArgs(userId).
		WillReturnError(sql.ErrNoRows)

	user, err := database.GetUserById(ctx, userId)
	require.Error(t, err)
	require.Nil(t, user)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserById_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	userId := 42

	mock.ExpectQuery("SELECT email, name, hide_email FROM usertable WHERE id = \\$1").
		WithArgs(userId).
		WillReturnError(fmt.Errorf("some db error"))

	user, err := database.GetUserById(ctx, userId)
	require.Error(t, err)
	require.Nil(t, user)

	require.NoError(t, mock.ExpectationsWereMet())
}
