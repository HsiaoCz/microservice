package service

import "errors"

type IUserService interface {
	GetName(userID int) string
	DeleteUser(userID int) error
}

type UserService struct{}

func (u UserService) GetName(userID int) string {
	if userID == 101 {
		return "bob"
	}
	return "alex"
}

func (u UserService) DeleteUser(userID int) error {
	if userID == 101 {
		return errors.New("无权限")
	}
	return nil
}
