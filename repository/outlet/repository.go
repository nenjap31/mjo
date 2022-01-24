package outlet

import (
	"database/sql"
	querybuilder "mjo/repository/util/query"
	"mjo/service/outlet/defined"
	"mjo/util/logger"
	"net/url"
	"time"

	"gorm.io/gorm"
)

const LOG_IDENTIFIER = "REPOSITORY_OUTLET"
const TABLE_NAME = "outlets"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

type Outlet struct {
	Id           uint           `gorm:"column:id;primaryKey"`
	MerchantId       uint `gorm:"column:merchant_id;default:0"`
	OutletName     sql.NullString `gorm:"column:outlet_name;size:40"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	CreatedBy    uint           `gorm:"column:created_by;default:0"`
	UpdatedAt    time.Time      `gorm:"column:update_at;autoUpdateTime"`
	UpdatedBy    uint           `gorm:"column:updated_by;default:0"`
}

type Tabler interface {
	TableName() string
}

func (Outlet) TableName() string {
	return TABLE_NAME
}

func Migrate(db gorm.DB) error {
	if !db.Migrator().HasTable(Outlet{}) {
		err := db.Migrator().CreateTable(Outlet{})
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

func NewOutlet(outlet defined.Outlet) *Outlet {
	return &Outlet{
		Id:           outlet.Id,
		MerchantId:       outlet.MerchantId,
		OutletName:        outlet.OutletName,
		CreatedAt:      outlet.CreatedAt,
		CreatedBy:      outlet.CreatedBy,
		UpdatedAt:     outlet.UpdatedAt,
		UpdatedBy:     outlet.UpdatedBy,
	}
}

func (outlet *Outlet) Map() defined.Outlet {
	var data defined.Outlet
	data.Id = outlet.Id
	data.MerchantId = outlet.MerchantId
	data.OutletName = outlet.OutletName
	data.CreatedAt = outlet.CreatedAt
	data.CreatedBy = outlet.CreatedBy
	data.UpdatedAt = outlet.UpdatedAt
	data.UpdatedBy = outlet.UpdatedBy
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
	List(filters url.Values, limit int, offset int) ([]*defined.Outlet, error)
	Create(outlet defined.Outlet) (*defined.Outlet, error)
	FindById(id string) (*defined.Outlet, error)
	UpdateById(id string, outlet defined.Outlet) (*defined.Outlet, error)
	DeleteById(id string) error
}

func (repository Repository) List(filters url.Values, limit int, offset int) ([]*defined.Outlet, error) {
	outlets := []*Outlet{}
	querybuilder.GormFilterBuilder(repository.db, filters, limit, offset).Find(&outlets)
	result := []*defined.Outlet{}
	for _, data := range outlets {
		newData := data.Map()
		result = append(result, &newData)
	}
	return result, nil
}

func (repository Repository) Create(outlet defined.Outlet) (*defined.Outlet, error) {
	newOutlet := NewOutlet(outlet)
	inserted := repository.db.Create(&newOutlet)
	if inserted.RowsAffected == 0 {
		return nil, inserted.Error
	}
	result := newOutlet.Map()
	return &result, nil
}

func (repository Repository) FindById(id string) (*defined.Outlet, error) {
	outlet := Outlet{}
	data := repository.db.Find(&outlet, id)
	if data.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	result := outlet.Map()
	return &result, nil
}

func (repository Repository) UpdateById(id string, outlet defined.Outlet) (*defined.Outlet, error) {
	newOutlet := Outlet{}
	data := repository.db.Find(&newOutlet, id)
	if data.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	newOutlet.OutletName = outlet.OutletName
	newOutlet.MerchantId = outlet.MerchantId
	newOutlet.UpdatedBy = outlet.UpdatedBy
	err := repository.db.Save(&newOutlet)
	if err.Error != nil {
		return nil, err.Error
	}
	result := newOutlet.Map()
	return &result, nil
}

func (repository Repository) DeleteById(id string) error {
	newOutlet := Outlet{}
	data := repository.db.Find(&newOutlet, id)
	if data.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	repository.db.Delete(&newOutlet)
	return nil
}