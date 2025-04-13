package postgres

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"skillForce/internal/models"
	"skillForce/pkg/hash"
	"skillForce/pkg/logs"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func (d *Database) saveSession(ctx context.Context, userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
	})

	secretToken, err := token.SignedString([]byte(d.SESSION_SECRET))
	if err != nil {
		return "", err
	}

	_, err = d.conn.Exec("INSERT INTO sessions (user_id, token, expire) VALUES ($1, $2, $3)", userId, secretToken, time.Now().AddDate(1, 0, 0))
	if err != nil {
		return "", err
	}

	return secretToken, nil
}

// userExists - проверяет, существует ли пользователь с указанным email
func (d *Database) userExists(email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM usertable WHERE email = $1)"
	err := d.conn.QueryRow(query, email).Scan(&exists)
	return exists, err
}

func (d *Database) GetUserById(ctx context.Context, userId int) (*models.User, error) {
	var user models.User
	err := d.conn.QueryRow("SELECT email, name, hide_email FROM usertable WHERE id = $1", userId).Scan(&user.Email, &user.Name, &user.HideEmail)
	if err != nil {
		logs.PrintLog(ctx, "GetUserById", fmt.Sprintf("%+v", err))
		return nil, err
	}
	return &user, nil
}

func (d *Database) parseToken(ctx context.Context, token string) (map[string]interface{}, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(d.SESSION_SECRET), nil
	})

	if err != nil {
		logs.PrintLog(ctx, "parseToken", fmt.Sprintf("token parse error: %+v", err))
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		if exp, ok := claims["expire"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return nil, errors.New("token expired")
			}
		} else {
			return nil, errors.New("invalid or missing 'expire' field")
		}

		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RegisterUser - сохраняет нового пользователя в базе данных и создает сессию, тоже сохраняя её в базе
func (d *Database) RegisterUser(ctx context.Context, user *models.User) (string, error) {
	emailExists, err := d.userExists(user.Email)
	if err != nil {
		return "", err
	}
	if emailExists {
		return "", errors.New("email exists")
	}
	saltBase64 := base64.StdEncoding.EncodeToString(user.Salt)
	_, err = d.conn.Exec("INSERT INTO usertable (email, name, password, salt) VALUES ($1, $2, $3, $4)", user.Email, user.Name, user.Password, saltBase64)
	if err != nil {
		return "", err
	}

	logs.PrintLog(ctx, "RegisterUser", fmt.Sprintf("save user %+v in db", user))

	rows, err := d.conn.Query("SELECT id FROM usertable WHERE email = $1", user.Email)
	if err != nil {
		return "", err
	}
	for rows.Next() {
		err = rows.Scan(&user.Id)
		if err != nil {
			return "", err
		}
	}

	cookieValue, err := d.saveSession(ctx, user.Id)
	if err != nil {
		return "", err
	}

	logs.PrintLog(ctx, "RegisterUser", fmt.Sprintf("save session of user %+v in db", user))

	return cookieValue, nil
}

func (d *Database) ValidUser(ctx context.Context, user *models.User) (string, error) {
	emailExists, err := d.userExists(user.Email)
	if err != nil {
		return "", err
	}
	if emailExists {
		return "", errors.New("email exists")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
		"expire":   time.Now().Add(time.Hour).Unix(),
	})

	secretToken, err := token.SignedString([]byte(d.SESSION_SECRET))
	if err != nil {
		logs.PrintLog(ctx, "ValidUser", fmt.Sprintf("%+v", err))
		return "", err
	}

	logs.PrintLog(ctx, "ValidUser", fmt.Sprintf("create token for user %+v", user))
	return secretToken, nil
}

func (d *Database) GetUserByToken(ctx context.Context, token string) (*models.User, error) {
	var user models.User
	claims, err := d.parseToken(ctx, token)
	if err != nil {
		return nil, err
	}

	user.Name = claims["name"].(string)
	user.Email = claims["email"].(string)
	user.Password = claims["password"].(string)

	return &user, nil
}

// GetUserByCookie - получение пользователя по cookie
func (d *Database) GetUserByCookie(ctx context.Context, cookieValue string) (*models.UserProfile, error) {
	var userProfile models.UserProfile
	err := d.conn.QueryRow("SELECT u.id, u.email, u.name, COALESCE(u.bio, ''), u.avatar_src, u.hide_email FROM usertable u JOIN sessions s ON u.id = s.user_id WHERE s.token = $1 AND s.expire > NOW();",
		cookieValue).Scan(&userProfile.Id, &userProfile.Email, &userProfile.Name, &userProfile.Bio, &userProfile.AvatarSrc, &userProfile.HideEmail)
	if err != nil {
		logs.PrintLog(ctx, "GetUserByCookie", fmt.Sprintf("error in GetUserByCookie %+v", err))
		return nil, err
	}
	return &userProfile, err
}

// AuthenticateUser - проверяет есть ли пользователь с указанным email и паролем в базе данных, елси есть - возвращает его id и сохраняет сессию в базе
func (d *Database) AuthenticateUser(ctx context.Context, email string, password string) (string, error) {
	var id int
	emailExists, err := d.userExists(email)
	if err != nil {
		return "", err
	}
	if !emailExists {
		return "", errors.New("email or password incorrect")
	}
	var passwordFromDB string
	var salt string
	err2 := d.conn.QueryRow("SELECT id, password, salt FROM usertable WHERE email = $1", email).Scan(&id, &passwordFromDB, &salt)
	if err2 != nil {
		return "", err2
	}
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return "", err
	}

	if !hash.CheckPassword(password, passwordFromDB, saltBytes) {
		return "", errors.New("email or password incorrect")
	}

	logs.PrintLog(ctx, "AuthenticateUser", fmt.Sprintf("login user with email %+v in db", email))

	cookieValue, err := d.saveSession(ctx, id)
	if err != nil {
		return "", err
	}

	logs.PrintLog(ctx, "AuthenticateUser", fmt.Sprintf("save session of user with email %+v in db", email))
	return cookieValue, nil
}

// LogoutUser - удаляет сессию пользователя из базы данных
func (d *Database) LogoutUser(ctx context.Context, userId int) error {
	_, err := d.conn.Exec("DELETE FROM sessions WHERE user_id = $1", userId)
	if err != nil {
		return err
	}
	logs.PrintLog(ctx, "LogoutUser", fmt.Sprintf("logout user with id %+v in db", userId))
	return err
}

func (d *Database) UpdateProfile(ctx context.Context, userId int, userProfile *models.UserProfile) error {
	_, err := d.conn.Exec("UPDATE usertable SET email = $1, name = $2, bio = $3, hide_email = $4 WHERE id = $5",
		userProfile.Email, userProfile.Name, userProfile.Bio, userProfile.HideEmail, userId)
	if err != nil {
		return err
	}
	logs.PrintLog(ctx, "UpdateProfile", fmt.Sprintf("update profile %+v of user with id %+v in db", userProfile, userId))
	return err
}

func (d *Database) UpdateProfilePhoto(ctx context.Context, photoUrl string, userId int) (string, error) {
	_, err := d.conn.Exec("UPDATE usertable SET avatar_src = $1 WHERE id = $2", photoUrl, userId)
	if err != nil {
		return "", err
	}

	logs.PrintLog(ctx, "UpdateProfilePhoto", fmt.Sprintf("update profile photo to %+v of user with id %+v in db", photoUrl, userId))

	return photoUrl, nil
}

func (d *Database) DeleteProfilePhoto(ctx context.Context, userId int) error {
	delaultAvatar := "http://217.16.21.64:8006/avatars/default_avatar.png"
	_, err := d.conn.Exec("UPDATE usertable SET avatar_src = $1 WHERE id = $2", delaultAvatar, userId)
	if err != nil {
		logs.PrintLog(ctx, "DeleteProfilePhoto", fmt.Sprintf("%+v", err))
		return err
	}

	logs.PrintLog(ctx, "DeleteProfilePhoto", fmt.Sprintf("update profile photo to default of user with id %+v in db", userId))

	return nil
}
