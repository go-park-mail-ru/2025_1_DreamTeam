package postgres

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	usermodels "skillForce/internal/models/user"
	"skillForce/pkg/hash"
	"skillForce/pkg/logs"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveSession_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	secret := "super-secret"
	database := &Database{
		conn:           db,
		SESSION_SECRET: secret,
	}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	userId := 123

	mock.ExpectExec("INSERT INTO sessions").
		WithArgs(userId, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	token, err := database.saveSession(ctx, userId)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveSession_JWTSigningError(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	database := &Database{
		conn:           db,
		SESSION_SECRET: string([]byte{}),
	}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})

	database.SESSION_SECRET = string([]byte{})
	token, err := database.saveSession(ctx, 1)
	require.Error(t, err)
	require.Empty(t, token)
}

func TestSaveSession_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	database := &Database{
		conn:           db,
		SESSION_SECRET: "valid-secret",
	}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	userId := 123

	mock.ExpectExec("INSERT INTO sessions").
		WithArgs(userId, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(fmt.Errorf("insert error"))

	token, err := database.saveSession(ctx, userId)
	require.Error(t, err)
	require.Empty(t, token)
	require.EqualError(t, err, "insert error")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserExists_True(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	database := &Database{conn: db}
	email := "test@example.com"

	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM usertable WHERE email = \\$1\\)").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := database.userExists(email)
	require.NoError(t, err)
	require.True(t, exists)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserExists_False(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	database := &Database{conn: db}
	email := "notfound@example.com"

	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM usertable WHERE email = \\$1\\)").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	exists, err := database.userExists(email)
	require.NoError(t, err)
	require.False(t, exists)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserExists_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	database := &Database{conn: db}
	email := "fail@example.com"

	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM usertable WHERE email = \\$1\\)").
		WithArgs(email).
		WillReturnError(fmt.Errorf("query error"))

	exists, err := database.userExists(email)
	require.Error(t, err)
	require.False(t, exists)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserById_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

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
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

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
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

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

func TestParseToken_Success(t *testing.T) {
	secret := "mysecret"
	database := &Database{SESSION_SECRET: secret}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 1,
		"expire":  float64(time.Now().Add(time.Hour).Unix()),
	})
	signedToken, err := token.SignedString([]byte(secret))
	require.NoError(t, err)

	claims, err := database.parseToken(ctx, signedToken)
	require.NoError(t, err)
	require.Equal(t, float64(1), claims["user_id"])
}

func TestParseToken_Expired(t *testing.T) {
	secret := "mysecret"
	database := &Database{SESSION_SECRET: secret}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 1,
		"expire":  float64(time.Now().Add(-time.Hour).Unix()),
	})
	signedToken, err := token.SignedString([]byte(secret))
	require.NoError(t, err)

	claims, err := database.parseToken(ctx, signedToken)
	require.Error(t, err)
	require.Contains(t, err.Error(), "token expired")
	require.Nil(t, claims)
}

func TestParseToken_MissingExpire(t *testing.T) {
	secret := "mysecret"
	database := &Database{SESSION_SECRET: secret}
	ctx := context.Background()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 1,
	})
	signedToken, err := token.SignedString([]byte(secret))
	require.NoError(t, err)

	claims, err := database.parseToken(ctx, signedToken)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid or missing 'expire'")
	require.Nil(t, claims)
}

func TestParseToken_MalformedToken(t *testing.T) {
	secret := "mysecret"
	database := &Database{SESSION_SECRET: secret}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})

	claims, err := database.parseToken(ctx, "not.a.jwt.token")
	require.Error(t, err)
	require.Nil(t, claims)
}

func TestRegisterUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	user := &usermodels.User{
		Email:    "test@example.com",
		Name:     "Test",
		Password: "hashedpassword",
		Salt:     []byte("randomsalt"),
	}

	database := &Database{conn: db, SESSION_SECRET: "secret"}

	mock.ExpectQuery("SELECT EXISTS.*FROM usertable").
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectExec("INSERT INTO usertable").
		WithArgs(user.Email, user.Name, user.Password, base64.StdEncoding.EncodeToString(user.Salt)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("SELECT id FROM usertable").
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))

	mock.ExpectExec("INSERT INTO sessions").
		WithArgs(123, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	token, err := database.RegisterUser(ctx, user)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Equal(t, 123, user.Id)
}

func TestRegisterUser_EmailExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	database := &Database{conn: db, SESSION_SECRET: "secret"}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})

	user := &usermodels.User{Email: "test@example.com"}

	mock.ExpectQuery("SELECT EXISTS.*FROM usertable").
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	token, err := database.RegisterUser(ctx, user)
	require.Error(t, err)
	require.Equal(t, "email exists", err.Error())
	require.Empty(t, token)
}

func TestRegisterUser_InsertFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	database := &Database{conn: db, SESSION_SECRET: "secret"}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})

	user := &usermodels.User{
		Email:    "test@example.com",
		Name:     "Test",
		Password: "hashedpassword",
		Salt:     []byte("randomsalt"),
	}

	mock.ExpectQuery("SELECT EXISTS.*FROM usertable").
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectExec("INSERT INTO usertable").
		WithArgs(user.Email, user.Name, user.Password, base64.StdEncoding.EncodeToString(user.Salt)).
		WillReturnError(errors.New("insert failed"))

	token, err := database.RegisterUser(ctx, user)
	require.Error(t, err)
	require.Contains(t, err.Error(), "insert failed")
	require.Empty(t, token)
}

func TestRegisterUser_GetIDFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	database := &Database{conn: db, SESSION_SECRET: "secret"}
	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	user := &usermodels.User{
		Email:    "test@example.com",
		Name:     "Test",
		Password: "hashedpassword",
		Salt:     []byte("randomsalt"),
	}

	mock.ExpectQuery("SELECT EXISTS.*FROM usertable").
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectExec("INSERT INTO usertable").
		WithArgs(user.Email, user.Name, user.Password, base64.StdEncoding.EncodeToString(user.Salt)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("SELECT id FROM usertable").
		WithArgs(user.Email).
		WillReturnError(errors.New("query id failed"))

	token, err := database.RegisterUser(ctx, user)
	require.Error(t, err)
	require.Contains(t, err.Error(), "query id failed")
	require.Empty(t, token)
}

func TestValidUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	database := &Database{
		conn:           db,
		SESSION_SECRET: "supersecret",
	}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})

	user := &usermodels.User{
		Email:    "new@example.com",
		Name:     "Alice",
		Password: "securehash",
	}

	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM usertable WHERE email = \\$1\\)").
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	token, err := database.ValidUser(ctx, user)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestValidUser_EmailExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	database := &Database{
		conn:           db,
		SESSION_SECRET: "supersecret",
	}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})

	user := &usermodels.User{
		Email: "existing@example.com",
	}

	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM usertable WHERE email = \\$1\\)").
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	token, err := database.ValidUser(ctx, user)
	require.Error(t, err)
	require.Equal(t, "email exists", err.Error())
	require.Empty(t, token)
}

func TestValidUser_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	database := &Database{
		conn:           db,
		SESSION_SECRET: "supersecret",
	}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})

	user := &usermodels.User{
		Email: "fail@example.com",
	}

	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM usertable WHERE email = \\$1\\)").
		WithArgs(user.Email).
		WillReturnError(errors.New("db is down"))

	token, err := database.ValidUser(ctx, user)
	require.Error(t, err)
	require.Contains(t, err.Error(), "db is down")
	require.Empty(t, token)
}

func TestGetUserByToken_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	database := &Database{
		SESSION_SECRET: "mysecret",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":     "Alice",
		"email":    "alice@example.com",
		"password": "hashedpassword",
		"expire":   time.Now().Add(time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(database.SESSION_SECRET))
	require.NoError(t, err)

	user, err := database.GetUserByToken(ctx, tokenStr)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, "Alice", user.Name)
	require.Equal(t, "alice@example.com", user.Email)
	require.Equal(t, "hashedpassword", user.Password)
}

func TestGetUserByToken_InvalidToken(t *testing.T) {
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	database := &Database{
		SESSION_SECRET: "mysecret",
	}

	invalidToken := "not.a.real.token"

	user, err := database.GetUserByToken(ctx, invalidToken)
	require.Error(t, err)
	require.Nil(t, user)
}

// func TestGetUserByCookie_Success(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	require.NoError(t, err)
// 	defer func() {

// 	database := &Database{conn: db}
// 	ctx := context.Background()
// 	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
// 		Data: make([]*logs.LogString, 0),
// 	})
// 	cookieValue := "valid_cookie_token"

// 	expected := usermodels.UserProfile{
// 		Id:        1,
// 		Email:     "user@example.com",
// 		Name:      "Test User",
// 		Bio:       "Developer",
// 		AvatarSrc: "avatar.jpg",
// 		HideEmail: false,
// 	}

// 	rows := sqlmock.NewRows([]string{"id", "email", "name", "bio", "avatar_src", "hide_email"}).
// 		AddRow(expected.Id, expected.Email, expected.Name, expected.Bio, expected.AvatarSrc, expected.HideEmail)

// 	mock.ExpectQuery(`SELECT u.id, u.email, u.name, COALESCE\(u.bio, ''\), u.avatar_src, u.hide_email FROM usertable u JOIN sessions s ON u.id = s.user_id WHERE s.token = \$1 AND s.expire > NOW\(\);`).
// 		WithArgs(cookieValue).
// 		WillReturnRows(rows)

// 	userProfile, err := database.GetUserByCookie(ctx, cookieValue)
// 	require.NoError(t, err)
// 	require.Equal(t, &expected, userProfile)
// }

func TestGetUserByCookie_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	cookieValue := "token"

	mock.ExpectQuery(`SELECT u.id, u.email, u.name, COALESCE\(u.bio, ''\), u.avatar_src, u.hide_email FROM usertable u JOIN sessions s ON u.id = s.user_id WHERE s.token = \$1 AND s.expire > NOW\(\);`).
		WithArgs(cookieValue).
		WillReturnError(sql.ErrNoRows)

	userProfile, err := database.GetUserByCookie(ctx, cookieValue)
	require.Error(t, err)
	require.Nil(t, userProfile)
}

func TestLogoutUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	userId := 42

	mock.ExpectExec(`DELETE FROM sessions WHERE user_id = \$1`).
		WithArgs(userId).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = database.LogoutUser(ctx, userId)
	require.NoError(t, err)
}

func TestLogoutUser_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	userId := 42

	mock.ExpectExec(`DELETE FROM sessions WHERE user_id = \$1`).
		WithArgs(userId).
		WillReturnError(fmt.Errorf("some db error"))

	err = database.LogoutUser(ctx, userId)
	require.Error(t, err)
	require.EqualError(t, err, "some db error")
}

func TestUpdateProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	database := &Database{conn: db}

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	userId := 1
	profile := &usermodels.UserProfile{
		Email:     "new@example.com",
		Name:      "New Name",
		Bio:       "Updated bio",
		HideEmail: true,
	}

	mock.ExpectExec("UPDATE usertable SET email = \\$1, name = \\$2, bio = \\$3, hide_email = \\$4 WHERE id = \\$5").
		WithArgs(profile.Email, profile.Name, profile.Bio, profile.HideEmail, userId).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = database.UpdateProfile(ctx, userId, profile)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProfilePhoto(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	database := &Database{conn: db}

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	userId := 42
	photoUrl := "https://example.com/avatar.jpg"

	mock.ExpectExec("UPDATE usertable SET avatar_src = \\$1 WHERE id = \\$2").
		WithArgs(photoUrl, userId).
		WillReturnResult(sqlmock.NewResult(0, 1))

	result, err := database.UpdateProfilePhoto(ctx, photoUrl, userId)
	require.NoError(t, err)
	require.Equal(t, photoUrl, result)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteProfilePhoto(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	database := &Database{conn: db}

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	userId := 42
	defaultAvatar := "http://217.16.21.64:8006/avatars/default_avatar.png"

	mock.ExpectExec("UPDATE usertable SET avatar_src = \\$1 WHERE id = \\$2").
		WithArgs(defaultAvatar, userId).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = database.DeleteProfilePhoto(ctx, userId)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthenticateUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	email := "test@example.com"
	password := "password123"
	salt := []byte("somesalt")
	hashedPassword := hash.HashPassword(password, salt)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM usertable WHERE email = \$1\)`).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectQuery(`SELECT id, password, salt FROM usertable WHERE email = \$1`).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "password", "salt"}).
			AddRow(1, hashedPassword, saltBase64))

	mock.ExpectExec(`INSERT INTO sessions \(user_id, token, expire\)`).
		WithArgs(1, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	token, err := database.AuthenticateUser(ctx, email, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGetUserByCookie_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	database := &Database{conn: db}

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	cookie := "testcookie"

	mock.ExpectQuery(`SELECT u.id, u.email, u.name, COALESCE\(u.bio, ''\), u.avatar_src, u.hide_email, u.role FROM usertable u JOIN sessions s ON u.id = s.user_id WHERE s.token = \$1 AND s.expire > NOW\(\)`).
		WithArgs(cookie).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "bio", "avatar_src", "hide_email", "role"}).
			AddRow(1, "user@example.com", "User", "Bio", "/avatars/user.png", false, "admin"))

	userProfile, err := database.GetUserByCookie(ctx, cookie)
	assert.NoError(t, err)
	assert.Equal(t, "user@example.com", userProfile.Email)
	assert.True(t, userProfile.IsAdmin)
}

func TestDeleteProfilePhoto_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	database := &Database{conn: db}

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	userId := 1
	defaultAvatar := "http://217.16.21.64:8006/avatars/default_avatar.png"

	mock.ExpectExec(`UPDATE usertable SET avatar_src = \$1 WHERE id = \$2`).
		WithArgs(defaultAvatar, userId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = database.DeleteProfilePhoto(ctx, userId)
	assert.NoError(t, err)
}
