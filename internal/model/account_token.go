package model

import "time"

const AccountTokenTable string = "account_token"

type AccountToken struct {
	Id                int       `gorm:"column:id"`
	Openid            string    `gorm:"column:openid"`
	SessionKey        string    `gorm:"column:session_key"`
	InternalToken     string    `gorm:"column:internal_token"`
	InternalTokenTime time.Time `gorm:"column:internal_token_time"`
	CTime             time.Time `gorm:"column:c_time"`
	MTime             time.Time `gorm:"column:m_time"`
	PersonID          string    `gorm:"column:person_id"`
}
