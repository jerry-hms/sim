package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"sim/app/global/variable"
	"time"
)

type BaseModel struct {
	*gorm.DB  `gorm:"-" json:"-"`
	Id        uint64    `gorm:"column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func ConnDb() *gorm.DB {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		variable.ConfigYml.GetString("mysql.username"),
		variable.ConfigYml.GetString("mysql.password"),
		variable.ConfigYml.GetString("mysql.host"),
		variable.ConfigYml.GetString("mysql.port"),
		variable.ConfigYml.GetString("mysql.database"),
		variable.ConfigYml.GetString("mysql.charset"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   variable.ConfigYml.GetString("mysql.prefix"),
			SingularTable: true,
		},
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	return db
}
