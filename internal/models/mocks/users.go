package mocks

import (
	"fmt"
	"github.com/martinheinrich2/goUebernachter/internal/models"
	"time"
)

type UserModel struct{}

func (m *UserModel) Insert(lastname, firstname, email, job, room, password string) error {
	fmt.Println("mocks/users.go Insert method")
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) Update(id int, lastname, firstname, email, job, room string, admin bool, active bool) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (m *UserModel) MailExists(email string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *UserModel) GetUserId(email string) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (m *UserModel) SelectUserByJob(job string) ([]models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *UserModel) All() ([]models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *UserModel) AllNames() ([]models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *UserModel) ResetPassword(id int, newPassword string) error {
	//TODO implement me
	panic("implement me")
}

func (m *UserModel) Get(id int) (models.User, error) {
	fmt.Println("this is the mocks/users.go Get method")
	if id == 1 {
		u := models.User{
			ID:        1,
			LastName:  "Cooper",
			FirstName: "Alice",
			Email:     "alice@example.com",
			JobTitle:  "Social worker",
			Room:      "765",
			Created:   time.Now(),
		}

		return u, nil
	}

	return models.User{}, models.ErrNoRecord
}

func (m *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	if id == 1 {
		if currentPassword != "pa$$word" {
			return models.ErrInvalidCredentials
		}

		return nil
	}

	return models.ErrNoRecord
}
