package service

import (
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/models"
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/repository"
)

type LoginServiceInterface interface {
	CreateUserService(user models.User) (interface{}, error)
	GetPasswordService(username string) (string, error)
}

type loginService struct{}

var (
	loginRepository repository.LoginRepositoryInterface
)

func NewLoginService(repository repository.LoginRepositoryInterface) LoginServiceInterface {
	loginRepository = repository
	return &loginService{}
}

func (lS *loginService) CreateUserService(user models.User) (interface{}, error) {
	return loginRepository.CreateUser(user)
}

func (lS *loginService) GetPasswordService(username string) (string, error) {
	return loginRepository.FetchHashedPass(username)
}
