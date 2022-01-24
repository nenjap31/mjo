package outlet

import (
	"mjo/repository/outlet"
	"mjo/service/outlet/defined"
	"mjo/util/logger"
	"net/url"
)

const LOG_IDENTIFIER = "SERVICE_OUTLET"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

func NewOutlet(outlet defined.Outlet) *defined.Outlet {
	return &defined.Outlet{
		Id:           outlet.Id,
		MerchantId:       outlet.MerchantId,
		OutletName:      outlet.OutletName,
		CreatedAt:      outlet.CreatedAt,
		CreatedBy:      outlet.CreatedBy,
		UpdatedAt:     outlet.UpdatedAt,
		UpdatedBy:     outlet.UpdatedBy,
	}
}

type Service struct {
	repository outlet.IRepository
}

type IService interface {
	List(filters url.Values, limit int, offset int) ([]*defined.Outlet, error)
	Create(outlet defined.Outlet) (*defined.Outlet, error)
	FindById(id string) (*defined.Outlet, error)
	UpdateById(id string, outlet defined.Outlet) (*defined.Outlet, error)
	DeleteById(id string) error
}

func NewService(repository outlet.IRepository) IService {
	return &Service{
		repository: repository,
	}
}

func (service *Service) List(filters url.Values, limit int, offset int) ([]*defined.Outlet, error) {
	result, err := service.repository.List(filters, limit, offset)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}

func (service *Service) Create(outlet defined.Outlet) (*defined.Outlet, error) {
	repository, err := service.repository.Create(outlet)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result := NewOutlet(*repository)
	return result, nil
}

func (service *Service) FindById(id string) (*defined.Outlet, error) {
	repository, err := service.repository.FindById(id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result := NewOutlet(*repository)
	return result, nil
}

func (service *Service) UpdateById(id string, outlet defined.Outlet) (*defined.Outlet, error) {
	repository, err := service.repository.UpdateById(id, outlet)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result := NewOutlet(*repository)
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