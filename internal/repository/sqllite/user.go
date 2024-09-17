package sqllite

import "SimpleForum/internal/domain"

func (rp *Repository) CreateUser(user *domain.User) error {
	return nil
}

func (rp *Repository) UpdateUser(user *domain.User) error {
	return nil
}

func (rp *Repository) DeleteUser(user *domain.User) error {
	return nil
}

func (rp *Repository) GetUserByID(userId int) (domain.User, error) {
	return domain.User{}, nil
}
