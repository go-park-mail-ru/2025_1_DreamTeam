package models

import "skillForce/internal/models/dto"

type User struct {
	Id        int
	Name      string
	Email     string
	Password  string
	Salt      []byte
	HideEmail bool
}

type UserProfile struct {
	Id        int
	Name      string
	Email     string
	Bio       string
	AvatarSrc string
	HideEmail bool
}

func NewUser(dto dto.UserDTO) *User {
	return &User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
	}
}

func NewUserProfile(dto dto.UserProfileDTO) *UserProfile {
	return &UserProfile{
		Name:      dto.Name,
		Email:     dto.Email,
		Bio:       dto.Bio,
		AvatarSrc: dto.AvatarSrc,
		HideEmail: dto.HideEmail,
	}
}
