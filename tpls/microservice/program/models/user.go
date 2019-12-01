package models

import (
	"github.com/jinzhu/gorm"
	"github.com/micro-kit/micro-common/db"
)

// User 表名-无前缀
type User struct {
	gorm.Model `json:"-"`
	Username   string `json:"username" gorm:"column:username"`
	Nickname   string `json:"nickname" gorm:"column:nickname"`
	Password   string `json:"password,omitempty" gorm:"column:password"`

	Token string `json:"token" gorm:"-"` // 登录时使用
}

// TableName 表名
func (User) TableName() string {
	return gorm.DefaultTableNameHandler(nil, "user")
}

// UserLogin 根据用户名密码查询数据库 - login使用
func (m *User) UserLogin(username, password string) error {
	return db.GetSlaveDB().Model(m).Where("username = ? and password = ?", username, password).Find(m).Error
}
