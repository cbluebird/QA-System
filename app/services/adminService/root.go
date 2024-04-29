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