package repository

import (
	"SimpleForum/internal/domain"
	"fmt"
)

func (serviceRepo *ServiceRepository) CreateUser(nickname, email, hashedPassword, role string) error {

	user := &domain.User{Nickname: nickname, Email: email, Password: hashedPassword, Role: role}

	err := serviceRepo.Repo.CreateUser(user)
	if err != nil {
		return fmt.Errorf("Service-CreateUser: %w", err)
	}
	return nil
}

func (serviceRepo *ServiceRepository) UpdateUser(user *domain.User) error {
	return nil
}

func (serviceRepo *ServiceRepository) DeleteUser(user *domain.User) error {
	return nil
}

func (serviceRepo *ServiceRepository) GetUserByID(id int64) (*domain.User, error) {
	return nil, nil
}

// func (serviceRepo *ServiceRepository) GetUserByEmail(email string) (*domain.User, error) {}

//func (serviceRepo *ServiceRepository) CheckUserByEmail(email string) (bool, error) {
//	isThereSuchEmail, err := serviceRepo.Repo.CheckUserByEmail(email)
//	if err != nil {
//		// Think about err Handling
//		return false, err
//	}
//	return isThereSuchEmail, nil
//
//}

func (serviceRepo *ServiceRepository) GetUserByEmail(email string) (*domain.User, error) {

	user, err := serviceRepo.Repo.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("Service-GetUserByEmail: %w", err)
	}
	return user, nil
}
