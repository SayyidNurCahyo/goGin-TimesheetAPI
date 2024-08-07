package impl

import (
	"timesheet-app/config"
	"timesheet-app/entity"
	"timesheet-app/helper"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type TimeSheetRepository struct{}

func NewTimeSheetRepository() *TimeSheetRepository {
	return &TimeSheetRepository{}
}

func (TimeSheetRepository) CreateTimeSheet(timesheet entity.TimeSheet) (*entity.TimeSheet, error) {
	err := config.DB.Create(&timesheet).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return &timesheet, nil
}

func (TimeSheetRepository) UpdateTimeSheet(ts entity.TimeSheet) (*entity.TimeSheet, error) {
	err := config.DB.Transaction(func(db *gorm.DB) error {
		err := db.Save(&ts).Error
		if err != nil {
			log.Error().Msg(err.Error())
			return err
		}
		err = db.Save(&ts.TimeSheetDetails).Error
		if err != nil {
			log.Error().Msg(err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return &ts, nil
}

func (TimeSheetRepository) DeleteTimeSheet(id string) error {
	return config.DB.Model(&entity.TimeSheet{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

func (TimeSheetRepository) GetStatusTimeSheetByID(id string) (*entity.StatusTimeSheet, error) {
	var status *entity.StatusTimeSheet
	err := config.DB.Where("id = ?", id).First(&status).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return status, nil
}

func (TimeSheetRepository) GetStatusTimeSheetByName(name string) (*entity.StatusTimeSheet, error) {
	var status *entity.StatusTimeSheet
	err := config.DB.Where("status_name = ?", name).First(&status).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return status, nil
}

func (TimeSheetRepository) GetTimeSheetByID(id string) (*entity.TimeSheet, error) {
	var ts entity.TimeSheet
	err := config.DB.Preload("TimeSheetDetails").First(&ts, "id = ?", id).Error
	return &ts, err
}

func (TimeSheetRepository) GetAllTimeSheets(spec []func(db *gorm.DB) *gorm.DB) (*[]entity.TimeSheet, string, error) {
	var timeSheets []entity.TimeSheet
	db := config.DB.Scopes(spec...).Preload("TimeSheetDetails").Find(&timeSheets)
	totalRows := helper.GetTotalRows(db)
	return &timeSheets, totalRows, db.Error
}

func (TimeSheetRepository) ApproveManagerTimeSheet(id string, userID string) error {
	var status *entity.StatusTimeSheet
	err := config.DB.Where("status_name = ?", "accepted").First(&status).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	return config.DB.Model(&entity.TimeSheet{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"confirmed_manager_by": userID,
			"status_time_sheet_id": status.ID,
		}).Error
}

func (TimeSheetRepository) RejectManagerTimeSheet(id string, userID string) error {
	var status *entity.StatusTimeSheet
	err := config.DB.Where("status_name = ?", "denied").First(&status).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	return config.DB.Model(&entity.TimeSheet{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"confirmed_manager_by": userID,
			"status_time_sheet_id": status.ID,
		}).Error
}

func (TimeSheetRepository) ApproveBenefitTimeSheet(id string, userID string) error {
	var status *entity.StatusTimeSheet
	err := config.DB.Where("status_name = ?", "approved").First(&status).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	return config.DB.Model(&entity.TimeSheet{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"confirmed_benefit_by": userID,
			"status_time_sheet_id": status.ID,
		}).Error
}

func (TimeSheetRepository) RejectBenefitTimeSheet(id string, userID string) error {
	var status *entity.StatusTimeSheet
	err := config.DB.Where("status_name = ?", "rejected").First(&status).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	return config.DB.Model(&entity.TimeSheet{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"confirmed_benefit_by": userID,
			"status_time_sheet_id": status.ID,
		}).Error
}

func (t TimeSheetRepository) UpdateTimeSheetStatus(id string) error {
	var ts entity.TimeSheet
	err := config.DB.Preload("TimeSheetDetails").First(&ts, "id = ?", id).Error
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	status, err := t.GetStatusTimeSheetByName("pending")
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	ts.StatusTimeSheetID = status.ID

	_, err = t.UpdateTimeSheet(ts)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	return nil
}
