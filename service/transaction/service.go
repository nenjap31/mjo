package transaction

import (
	authServiceDefined "mjo/service/auth/defined"
	"mjo/repository/transaction"
	"mjo/service/transaction/defined"
	"mjo/util/logger"
	"net/url"
	"github.com/dgrijalva/jwt-go"
	"errors"
)

const LOG_IDENTIFIER = "SERVICE_TRANSACTION"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

func NewTransaction(transaction defined.Transaction) *defined.Transaction {
	return &defined.Transaction{
		Id:           transaction.Id,
		MerchantId:       transaction.MerchantId,
		OutletId:      transaction.OutletId,
		BillTotal:      transaction.BillTotal,
		CreatedAt:      transaction.CreatedAt,
		CreatedBy:      transaction.CreatedBy,
		UpdatedAt:     transaction.UpdatedAt,
		UpdatedBy:     transaction.UpdatedBy,
	}
}

type Service struct {
	repository transaction.IRepository
}

type IService interface {
	List(filters url.Values, limit int, offset int) ([]*defined.Transaction, error)
	Create(transaction defined.Transaction) (*defined.Transaction, error)
	FindById(id string) (*defined.Transaction, error)
	UpdateById(id string, transaction defined.Transaction) (*defined.Transaction, error)
	DeleteById(id string) error
	MonthlyReport(token *jwt.Token, stardate string, enddate string, limit int, offset int) ([]*defined.ReportMerchant, error)
	MonthlyOutletReport(token *jwt.Token, stardate string, enddate string, limit int, offset int) ([]*defined.ReportOutlet, error)
}

func NewService(repository transaction.IRepository) IService {
	return &Service{
		repository: repository,
	}
}

func (service *Service) List(filters url.Values, limit int, offset int) ([]*defined.Transaction, error) {
	result, err := service.repository.List(filters, limit, offset)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}

func (service *Service) Create(transaction defined.Transaction) (*defined.Transaction, error) {
	repository, err := service.repository.Create(transaction)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result := NewTransaction(*repository)
	return result, nil
}

func (service *Service) FindById(id string) (*defined.Transaction, error) {
	repository, err := service.repository.FindById(id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result := NewTransaction(*repository)
	return result, nil
}

func (service *Service) UpdateById(id string, transaction defined.Transaction) (*defined.Transaction, error) {
	repository, err := service.repository.UpdateById(id, transaction)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result := NewTransaction(*repository)
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

func (service *Service) MonthlyReport(token *jwt.Token, startdate string, enddate string, limit int, offset int) ([]*defined.ReportMerchant, error) {
	claims, success := token.Claims.(*authServiceDefined.AccessTokenClaims)
	if !success {
		return nil, errors.New("internal server error")
	}

	result, err := service.repository.MonthlyReport(claims.UserId, startdate, enddate, limit, offset)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}

func (service *Service) MonthlyOutletReport(token *jwt.Token, startdate string, enddate string, limit int, offset int) ([]*defined.ReportOutlet, error) {
	claims, success := token.Claims.(*authServiceDefined.AccessTokenClaims)
	if !success {
		return nil, errors.New("internal server error")
	}

	result, err := service.repository.MonthlyOutletReport(claims.UserId, startdate, enddate, limit, offset)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}