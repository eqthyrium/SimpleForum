package repository

import (
	"SimpleForum/internal/domain"
	"fmt"
)

func (serviceRepo *ServiceRepository) CreateUser(nickname, memberIdentity, hashedPassword, role string) error {

	user := &domain.User{Nickname: nickname, MemberIdentity: memberIdentity, Password: hashedPassword, Role: role}

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

// func (serviceRepo *ServiceRepository) GetUserByEmail(memberIdentity string) (*domain.User, error) {}

//func (serviceRepo *ServiceRepository) CheckUserByEmail(memberIdentity string) (bool, error) {
//	isThereSuchEmail, err := serviceRepo.Repo.CheckUserByEmail(memberIdentity)
//	if err != nil {
//		// Think about err Handling
//		return false, err
//	}
//	return isThereSuchEmail, nil
//
//}

func (serviceRepo *ServiceRepository) GetUserByEmail(memberIdentity string) (*domain.User, error) {

	user, err := serviceRepo.Repo.GetUserByEmail(memberIdentity)
	if err != nil {
		return nil, fmt.Errorf("Service-GetUserByEmail: %w", err)
	}
	return user, nil
}
