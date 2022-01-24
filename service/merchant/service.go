package merchant

import (
	"mjo/repository/merchant"
	"mjo/service/merchant/defined"
	"mjo/util/logger"
	"net/url"
)

const LOG_IDENTIFIER = "SERVICE_MERCHANT"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

func NewMerchant(merchant defined.Merchant) *defined.Merchant {
	return &defined.Merchant{
		Id:           merchant.Id,
		UserId:       merchant.UserId,
		MerchantName:      merchant.MerchantName,
		CreatedAt:      merchant.CreatedAt,
		CreatedBy:      merchant.CreatedBy,
		UpdatedAt:     merchant.UpdatedAt,
		UpdatedBy:     merchant.UpdatedBy,
	}
}

type Service struct {
	repository merchant.IRepository
}

type IService interface {
	List(filters url.Values, limit int, offset int) ([]*defined.Merchant, error)
	Create(merchant defined.Merchant) (*defined.Merchant, error)
	FindById(id string) (*defined.Merchant, error)
	UpdateById(id string, merchant defined.Merchant) (*defined.Merchant, error)
	DeleteById(id string) error
}

func NewService(repository merchant.IRepository) IService {
	return &Service{
		repository: repository,
	}
}

func (service *Service) List(filters url.Values, limit int, offset int) ([]*defined.Merchant, error) {
	result, err := service.repository.List(filters, limit, offset)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}

func (service *Service) Create(merchant defined.Merchant) (*defined.Merchant, error) {
	repository, err := service.repository.Create(merchant)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result := NewMerchant(*repository)
	return result, nil
}

func (service *Service) FindById(id string) (*defined.Merchant, error) {
	repository, err := service.repository.FindById(id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result := NewMerchant(*repository)
	return result, nil
}

func (service *Service) UpdateById(id string, merchant defined.Merchant) (*defined.Merchant, error) {
	repository, err := service.repository.UpdateById(id, merchant)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result := NewMerchant(*repository)
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