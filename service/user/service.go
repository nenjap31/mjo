package user

import (
	"mjo/repository/user"
	"mjo/service/user/defined"
	"mjo/service/util/generatepwd"
	"mjo/util/logger"
	"net/url"
)

const LOG_IDENTIFIER = "SERVICE_USER"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

func NewUser(user defined.User) *defined.User {
	return &defined.User{
		Id:           user.Id,
		Name:       user.Name,
		UserName:      user.UserName,
		Password:     user.Password,
		CreatedAt:      user.CreatedAt,
		CreatedBy:      user.CreatedBy,
		UpdatedAt:     user.UpdatedAt,
		UpdatedBy:     user.UpdatedBy,
	}
}

type Service struct {
	repository user.IRepository
}

type IService interface {
	List(filters url.Values, limit int, offset int) ([]*defined.User, error)
	Create(user defined.User) (*defined.User, error)
	FindById(id string) (*defined.User, error)
	UpdateById(id string, user defined.User) (*defined.User, error)
	DeleteById(id string) error
}

func NewService(repository user.IRepository) IService {
	return &Service{
		repository: repository,
	}
}

func (service *Service) List(filters url.Values, limit int, offset int) ([]*defined.User, error) {
	result, err := service.repository.List(filters, limit, offset)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}

func (service *Service) Create(user defined.User) (*defined.User, error) {
	password, err := generatepwd.GeneratePassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = password
	repository, err := service.repository.Create(user)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result := NewUser(*repository)
	return result, nil
}

func (service *Service) FindById(id string) (*defined.User, error) {
	repository, err := service.repository.FindById(id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result := NewUser(*repository)
	return result, nil
}

func (service *Service) UpdateById(id string, user defined.User) (*defined.User, error) {
	repository, err := service.repository.UpdateById(id, user)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result := NewUser(*repository)
	return result, nil
}

func (service *Service) DeleteById(id string) error {
	err := service.repository.DeleteById(id)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}