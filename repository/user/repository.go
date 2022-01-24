package user

import (
	"database/sql"
	querybuilder "mjo/repository/util/query"
	"mjo/service/user/defined"
	"mjo/util/logger"
	"net/url"
	"time"

	"gorm.io/gorm"
)

const LOG_IDENTIFIER = "REPOSITORY_USER"
const TABLE_NAME = "users"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

type User struct {
	Id           uint           `gorm:"column:id;primaryKey"`
	Name         sql.NullString `gorm:"column:name;size:45"`
	UserName     sql.NullString `gorm:"column:user_name;size:45"`
	Password     string         `gorm:"column:password;size:255"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	CreatedBy    uint           `gorm:"column:created_by;default:0"`
	UpdatedAt    time.Time      `gorm:"column:update_at;autoUpdateTime"`
	UpdatedBy    uint           `gorm:"column:updated_by;default:0"`
}

type Tabler interface {
	TableName() string
}

func (User) TableName() string {
	return TABLE_NAME
}

func Migrate(db gorm.DB) error {
	if !db.Migrator().HasTable(User{}) {
		err := db.Migrator().CreateTable(User{})
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

func NewUser(user defined.User) *User {
	return &User{
		Id:           user.Id,
		Name:       user.Name,
		UserName:        user.UserName,
		Password:     user.Password,
		CreatedAt:      user.CreatedAt,
		CreatedBy:      user.CreatedBy,
		UpdatedAt:     user.UpdatedAt,
		UpdatedBy:     user.UpdatedBy,
	}
}

func (user *User) Map() defined.User {
	var data defined.User
	data.Id = user.Id
	data.Name = user.Name
	data.UserName = user.UserName
	data.Password = user.Password
	data.CreatedAt = user.CreatedAt
	data.CreatedBy = user.CreatedBy
	data.UpdatedAt = user.UpdatedAt
	data.UpdatedBy = user.UpdatedBy
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
	List(filters url.Values, limit int, offset int) ([]*defined.User, error)
	Create(user defined.User) (*defined.User, error)
	FindById(id string) (*defined.User, error)
	FindByUsername(username string) (*defined.User, error)
	UpdateById(id string, user defined.User) (*defined.User, error)
	DeleteById(id string) error
}

func (repository Repository) List(filters url.Values, limit int, offset int) ([]*defined.User, error) {
	users := []*User{}
	querybuilder.GormFilterBuilder(repository.db, filters, limit, offset).Find(&users)
	result := []*defined.User{}
	for _, data := range users {
		newData := data.Map()
		result = append(result, &newData)
	}
	return result, nil
}

func (repository Repository) Create(user defined.User) (*defined.User, error) {
	newUser := NewUser(user)
	inserted := repository.db.Create(&newUser)
	if inserted.RowsAffected == 0 {
		return nil, inserted.Error
	}
	result := newUser.Map()
	return &result, nil
}

func (repository Repository) FindById(id string) (*defined.User, error) {
	user := User{}
	data := repository.db.Find(&user, id)
	if data.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	result := user.Map()
	return &result, nil
}

func (repository Repository) FindByUsername(username string) (*defined.User, error) {
	user := User{}
	data := repository.db.Where("user_name = ?", username).First(&user)
	if data.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	result := user.Map()
	return &result, nil
}

func (repository Repository) UpdateById(id string, user defined.User) (*defined.User, error) {
	newUser := User{}
	data := repository.db.Find(&newUser, id)
	if data.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	newUser.Name = user.Name
	newUser.UserName = user.UserName
	newUser.Password = user.Password
	newUser.UpdatedBy = user.UpdatedBy
	err := repository.db.Save(&newUser)
	if err.Error != nil {
		return nil, err.Error
	}
	result := newUser.Map()
	return &result, nil
}

func (repository Repository) DeleteById(id string) error {
	newUser := User{}
	data := repository.db.Find(&newUser, id)
	if data.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	repository.db.Delete(&newUser)
	return nil
}