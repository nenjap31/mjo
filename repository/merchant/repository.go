package merchant

import (
	"database/sql"
	querybuilder "mjo/repository/util/query"
	"mjo/service/merchant/defined"
	"mjo/util/logger"
	"net/url"
	"time"

	"gorm.io/gorm"
)

const LOG_IDENTIFIER = "REPOSITORY_MERCHANT"
const TABLE_NAME = "merchants"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

type Merchant struct {
	Id           uint           `gorm:"column:id;primaryKey"`
	UserId       uint `gorm:"column:user_id;default:0"`
	MerchantName     sql.NullString `gorm:"column:merchant_name;size:40"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	CreatedBy    uint           `gorm:"column:created_by;default:0"`
	UpdatedAt    time.Time      `gorm:"column:update_at;autoUpdateTime"`
	UpdatedBy    uint           `gorm:"column:updated_by;default:0"`
}

type Tabler interface {
	TableName() string
}

func (Merchant) TableName() string {
	return TABLE_NAME
}

func Migrate(db gorm.DB) error {
	if !db.Migrator().HasTable(Merchant{}) {
		err := db.Migrator().CreateTable(Merchant{})
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

func NewMerchant(merchant defined.Merchant) *Merchant {
	return &Merchant{
		Id:           merchant.Id,
		UserId:       merchant.UserId,
		MerchantName:        merchant.MerchantName,
		CreatedAt:      merchant.CreatedAt,
		CreatedBy:      merchant.CreatedBy,
		UpdatedAt:     merchant.UpdatedAt,
		UpdatedBy:     merchant.UpdatedBy,
	}
}

func (merchant *Merchant) Map() defined.Merchant {
	var data defined.Merchant
	data.Id = merchant.Id
	data.UserId = merchant.UserId
	data.MerchantName = merchant.MerchantName
	data.CreatedAt = merchant.CreatedAt
	data.CreatedBy = merchant.CreatedBy
	data.UpdatedAt = merchant.UpdatedAt
	data.UpdatedBy = merchant.UpdatedBy
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
	List(filters url.Values, limit int, offset int) ([]*defined.Merchant, error)
	Create(merchant defined.Merchant) (*defined.Merchant, error)
	FindById(id string) (*defined.Merchant, error)
	UpdateById(id string, merchant defined.Merchant) (*defined.Merchant, error)
	DeleteById(id string) error
}

func (repository Repository) List(filters url.Values, limit int, offset int) ([]*defined.Merchant, error) {
	merchants := []*Merchant{}
	querybuilder.GormFilterBuilder(repository.db, filters, limit, offset).Find(&merchants)
	result := []*defined.Merchant{}
	for _, data := range merchants {
		newData := data.Map()
		result = append(result, &newData)
	}
	return result, nil
}

func (repository Repository) Create(merchant defined.Merchant) (*defined.Merchant, error) {
	newMerchant := NewMerchant(merchant)
	inserted := repository.db.Create(&newMerchant)
	if inserted.RowsAffected == 0 {
		return nil, inserted.Error
	}
	result := newMerchant.Map()
	return &result, nil
}

func (repository Repository) FindById(id string) (*defined.Merchant, error) {
	merchant := Merchant{}
	data := repository.db.Find(&merchant, id)
	if data.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	result := merchant.Map()
	return &result, nil
}

func (repository Repository) UpdateById(id string, merchant defined.Merchant) (*defined.Merchant, error) {
	newMerchant := Merchant{}
	data := repository.db.Find(&newMerchant, id)
	if data.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	newMerchant.MerchantName = merchant.MerchantName
	newMerchant.UserId = merchant.UserId
	newMerchant.UpdatedBy = merchant.UpdatedBy
	err := repository.db.Save(&newMerchant)
	if err.Error != nil {
		return nil, err.Error
	}
	result := newMerchant.Map()
	return &result, nil
}

func (repository Repository) DeleteById(id string) error {
	newMerchant := Merchant{}
	data := repository.db.Find(&newMerchant, id)
	if data.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	repository.db.Delete(&newMerchant)
	return nil
}