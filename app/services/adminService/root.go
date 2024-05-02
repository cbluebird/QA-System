package adminService

import (
	"QA-System/app/models"
	"QA-System/config/database"
)

func GetUserByName(username string) (models.User, error) {
	var user models.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	return user, err
}

func CreatePermission(id int, surveyID int) error {
	err := database.DB.Create(&models.Manage{UserID: id, SurveyID: surveyID}).Error
	return err
}

func DeletePermission(id int, surveyID int) error {
	err := database.DB.Where("user_id = ? AND survey_id = ?", id, surveyID).Delete(&models.Manage{}).Error
	return err
}

func CheckPermission(id int, surveyID int) error {
	var manage models.Manage
	err := database.DB.Where("user_id = ? AND survey_id = ?", id, surveyID).First(&manage).Error
	return err
}