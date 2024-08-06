package repository

import (
	"final-project-enigma/dto/request"
	"final-project-enigma/entity"
	"gorm.io/gorm"
)

type TimeSheetRepository interface {
	CreateTimeSheet(ts entity.TimeSheet) (*entity.TimeSheet, error)
	UpdateTimeSheet(ts entity.TimeSheet) (*entity.TimeSheet, error)
	DeleteTimeSheet(id string) error
	GetStatusTimeSheetByID(id string) (*entity.StatusTimeSheet, error)
	GetStatusTimeSheetByName(name string) (*entity.StatusTimeSheet, error)
	GetTimeSheetByID(id string) (*entity.TimeSheet, error)
	GetAllTimeSheets(spec []func(db *gorm.DB) *gorm.DB) (*[]entity.TimeSheet, string, error)
	ApproveManagerTimeSheet(id string, userID string) error
	RejectManagerTimeSheet(id string, userID string) error
	ApproveBenefitTimeSheet(id string, userID string) error
	RejectBenefitTimeSheet(id string, userID string) error
	UpdateTimeSheetStatus(req request.UpdateTimeSheetStatusRequest) error
}
