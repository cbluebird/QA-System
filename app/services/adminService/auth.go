package adminService

import (
	"QA-System/app/models"
	"QA-System/app/utils"
	"QA-System/config/database"
)

func GetAdminByUsername(username string) (*models.User, error) {
	var user models.User
	result := database.DB.Model(&models.User{}).Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if user.Password != "" {
		aesDecryptPassword(&user)
	}
	return &user, result.Error
}

func GetAdminByID(id int) (*models.User, error) {
	user := models.User{}
	result := database.DB.Model(&models.User{}).Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	aesDecryptPassword(&user)
	return &user, nil
}

func IsAdminExist(username string) error {
	var user models.User
	result := database.DB.Model(models.User{}).Where("username = ?", username).First(&user)
	return result.Error
}

func CreateAdmin(user models.User) error {
	aesEncryptPassword(&user)
	result := database.DB.Model(models.User{}).Create(&user)
	return result.Error
}

func aesDecryptPassword(user *models.User) {
	user.Password = utils.AesDecrypt(user.Password)
}

func aesEncryptPassword(user *models.User) {
	user.Password = utils.AesEncrypt(user.Password)
}
