package repository

import (
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/models"
	"gorm.io/gorm"
)

type LoginRepositoryInterface interface {
	CreateUser(u models.User) (interface{}, error)
	FetchHashedPass(name string) (string, error)
}

type loginRepository struct {
	DB *gorm.DB
}

func NewLoginRepository(db *gorm.DB) LoginRepositoryInterface {
	return &loginRepository{db}
}

func (loginRepo *loginRepository) CreateUser(u models.User) (interface{}, error) {
	result := loginRepo.DB.Create(&u)
	return result, result.Error
}

func (loginRepo *loginRepository) FetchHashedPass(name string) (string, error) {
	type Result struct {
		Username string
		Password string
	}
	var r Result
	result := loginRepo.DB.Table("users").Where("username = ?", name).Scan(&r)
	return r.Password, result.Error
}
