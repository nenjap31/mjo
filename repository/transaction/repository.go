package transaction

import (
	querybuilder "mjo/repository/util/query"
	"mjo/service/transaction/defined"
	"mjo/util/logger"
	"net/url"
	"time"

	"gorm.io/gorm"
)

const LOG_IDENTIFIER = "REPOSITORY_TRANSACTIONS"
const TABLE_NAME = "transactions"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

type Transaction struct {
	Id           uint           `gorm:"column:id;primaryKey"`
	MerchantId       uint `gorm:"column:merchant_id;default:0"`
	OutletId       uint `gorm:"column:outlet_id;default:0"`
	BillTotal       float64 `gorm:"column:bill_total;default:0"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	CreatedBy    uint           `gorm:"column:created_by;default:0"`
	UpdatedAt    time.Time      `gorm:"column:update_at;autoUpdateTime"`
	UpdatedBy    uint           `gorm:"column:updated_by;default:0"`
}

type MonthlyTrxReport struct {
	MerchantName       string `gorm:"column:merchant_name;"`
	Date       string `gorm:"column:date;"`
	Omzet       float64 `gorm:"column:omzet;default:0"`
}

type MonthlyTrxOutletReport struct {
	MerchantName       string `gorm:"column:merchant_name;"`
	OutletName       string `gorm:"column:outlet_name;"`
	Date       string `gorm:"column:date;"`
	Omzet       float64 `gorm:"column:omzet;default:0"`
}

type Tabler interface {
	TableName() string
}

func (Transaction) TableName() string {
	return TABLE_NAME
}

func Migrate(db gorm.DB) error {
	if !db.Migrator().HasTable(Transaction{}) {
		err := db.Migrator().CreateTable(Transaction{})
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

func NewTransaction(transaction defined.Transaction) *Transaction {
	return &Transaction{
		Id:           transaction.Id,
		MerchantId:       transaction.MerchantId,
		OutletId:        transaction.OutletId,
		BillTotal:        transaction.BillTotal,
		CreatedAt:      transaction.CreatedAt,
		CreatedBy:      transaction.CreatedBy,
		UpdatedAt:     transaction.UpdatedAt,
		UpdatedBy:     transaction.UpdatedBy,
	}
}

func (transaction *Transaction) Map() defined.Transaction {
	var data defined.Transaction
	data.Id = transaction.Id
	data.MerchantId = transaction.MerchantId
	data.OutletId = transaction.OutletId
	data.BillTotal = transaction.BillTotal
	data.CreatedAt = transaction.CreatedAt
	data.CreatedBy = transaction.CreatedBy
	data.UpdatedAt = transaction.UpdatedAt
	data.UpdatedBy = transaction.UpdatedBy
	return data
}

func (monthlyreport *MonthlyTrxReport) Map() defined.ReportMerchant {
	var data defined.ReportMerchant
	data.MerchantName = monthlyreport.MerchantName
	data.Date = monthlyreport.Date
	data.Omzet = monthlyreport.Omzet
	return data
}

func (monthlyoutletreport *MonthlyTrxOutletReport) Map() defined.ReportOutlet {
	var data defined.ReportOutlet
	data.MerchantName = monthlyoutletreport.MerchantName
	data.OutletName = monthlyoutletreport.OutletName
	data.Date = monthlyoutletreport.Date
	data.Omzet = monthlyoutletreport.Omzet
	return data
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) (*Repository, error) {
	repository := Repository{db}
	return &repository, nil
}

type IRepository interface {
	List(filters url.Values, limit int, offset int) ([]*defined.Transaction, error)
	Create(transaction defined.Transaction) (*defined.Transaction, error)
	FindById(id string) (*defined.Transaction, error)
	MonthlyReport(merchant_id uint,startdate string, enddate string, limit int, offset int)([]*defined.ReportMerchant, error)
	MonthlyOutletReport(merchant_id uint,startdate string, enddate string, limit int, offset int)([]*defined.ReportOutlet, error)
	UpdateById(id string, transaction defined.Transaction) (*defined.Transaction, error)
	DeleteById(id string) error
}

func (repository Repository) List(filters url.Values, limit int, offset int) ([]*defined.Transaction, error) {
	transactions := []*Transaction{}
	querybuilder.GormFilterBuilder(repository.db, filters, limit, offset).Find(&transactions)
	result := []*defined.Transaction{}
	for _, data := range transactions {
		newData := data.Map()
		result = append(result, &newData)
	}
	return result, nil
}

func (repository Repository) Create(transaction defined.Transaction) (*defined.Transaction, error) {
	newTransaction := NewTransaction(transaction)
	inserted := repository.db.Create(&newTransaction)
	if inserted.RowsAffected == 0 {
		return nil, inserted.Error
	}
	result := newTransaction.Map()
	return &result, nil
}

func (repository Repository) FindById(id string) (*defined.Transaction, error) {
	transaction := Transaction{}
	data := repository.db.Find(&transaction, id)
	if data.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	result := transaction.Map()
	return &result, nil
}

func (repository Repository) UpdateById(id string, transaction defined.Transaction) (*defined.Transaction, error) {
	newTransaction := Transaction{}
	data := repository.db.Find(&newTransaction, id)
	if data.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	newTransaction.OutletId = transaction.OutletId
	newTransaction.MerchantId = transaction.MerchantId
	newTransaction.BillTotal = transaction.BillTotal
	newTransaction.UpdatedBy = transaction.UpdatedBy
	err := repository.db.Save(&newTransaction)
	if err.Error != nil {
		return nil, err.Error
	}
	result := newTransaction.Map()
	return &result, nil
}

func (repository Repository) DeleteById(id string) error {
	newTransaction := Transaction{}
	data := repository.db.Find(&newTransaction, id)
	if data.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	repository.db.Delete(&newTransaction)
	return nil
}

func (repository Repository) MonthlyReport(merchant_id uint,startdate string, enddate string, limit int, offset int) ([]*defined.ReportMerchant, error) {
	monthlyreport := []*MonthlyTrxReport{}
	_ = repository.db.Raw("select merchants.merchant_name,DATE_FORMAT(transactions.created_at,'%Y-%m-%d')as date,sum(bill_total)as omzet FROM transactions LEFT JOIN merchants on merchants.id = transactions.merchant_id WHERE transactions.merchant_id=? AND transactions.created_at between ? and ? GROUP BY merchant_id, transactions.created_at order by transactions.created_at LIMIT ? OFFSET ?", merchant_id, startdate, enddate, limit, offset).Scan(&monthlyreport)
	result := []*defined.ReportMerchant{}
	for _, data := range monthlyreport {
		newData := data.Map()
		result = append(result, &newData)
	}
	return result, nil
}

func (repository Repository) MonthlyOutletReport(merchant_id uint,startdate string, enddate string, limit int, offset int) ([]*defined.ReportOutlet, error) {
	monthlyoutletreport := []*MonthlyTrxOutletReport{}
	_ = repository.db.Raw("select merchants.merchant_name,outlets.outlet_name,DATE_FORMAT(transactions.created_at,'%Y-%m-%d')as date,sum(bill_total)as omzet FROM transactions LEFT JOIN merchants on merchants.id = transactions.merchant_id LEFT JOIN outlets on outlets.id = transactions.outlet_id WHERE outlet.merchant_id=? AND transactions.created_at between ? and ? GROUP BY transactions.outlet_id, transactions.created_at order by transactions.created_at LIMIT ? OFFSET ?", merchant_id, startdate, enddate, limit, offset).Scan(&monthlyoutletreport)
	result := []*defined.ReportOutlet{}
	for _, data := range monthlyoutletreport {
		newData := data.Map()
		result = append(result, &newData)
	}
	return result, nil
}