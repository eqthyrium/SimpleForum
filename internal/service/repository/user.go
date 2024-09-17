package repository

import "SimpleForum/internal/domain"

func (userObject *ServiceRepository) CreateUser(user *domain.User) error {
	return nil
}

func (userObject *ServiceRepository) UpdateUser(user *domain.User) error {
	return nil
}

func (userObject *ServiceRepository) DeleteUser(user *domain.User) error {
	return nil
}

func (userObject *ServiceRepository) GetUserByID(id int64) (*domain.User, error) {
	return nil, nil
}
