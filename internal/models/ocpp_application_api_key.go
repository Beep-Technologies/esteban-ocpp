// Code generated by SQLBoiler 4.6.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"strconv"

	"github.com/Beep-Technologies/beepbeep3-iam/pkg/constants"
)

// OcppApplicationAPIKey is an object representing the database table.
type OcppApplicationAPIKey struct {
	APIKeyHash    string `gorm:"column:api_key_hash"`
	ApplicationID string `gorm:"column:application_id"`
	Description   string `gorm:"column:description"`
	IsActive      bool   `gorm:"column:is_active"`
}

func (OcppApplicationAPIKey) TableName() string {
	_ = strconv.Quote("")
	return constants.SCHEMA + "ocpp_application_api_key"
}