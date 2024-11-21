package controllers

import (
	"database/sql"
	"digibank/delivery/http/repository"
	"digibank/domain/entity"
	"digibank/infrastructure/encryptor"
	com "digibank/infrastructure/functions"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RegisterStruct struct {
	Base entity.UserBase `json:"base"`
	Info entity.UserInfo `json:"info"`
}

type LoginData struct {
	LoginStruct `json:"data"`
}

type LoginStruct struct {
	User entity.UserBase `json:"user"`
	Mac  string          `json:"mac_address"`
}

func Login(db *sql.DB, c *fiber.Ctx) error {

	defer com.PrintLog("======= Login Function END   =========")

	com.GenerateRandomID()

	com.PrintLog("======= Login Function START =========")

	var login LoginData

	if err := c.BodyParser(&login); err != nil {
		com.PrintLog(fmt.Sprintf("(LOGIN:0001)%s", err))
		return c.Status(501).JSON(fiber.Map{
			"status":  501,
			"message": "Error on Server Side",
		})
	}

	encrypted, err := encryptor.PasswordGenerator(login.User.Password)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(LOGIN:0002)%s", err))
		return c.Status(501).JSON(fiber.Map{
			"status":  501,
			"message": "Error on Server Side",
		})
	}

	rc, err := repository.LoginUser(db, &login.User, login.Mac, encrypted)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(LOGIN:0003)%s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  rc,
			"message": fmt.Sprintf("%s", err),
		})
	}

	res, err := encryptor.FieldGenerator(map[string]interface{}{
		"data":      login.User,
		"keepLogin": "1",
	})

	if err != nil {
		com.PrintLog(fmt.Sprintf("(LOGIN:0004)%s", err))
		return c.Status(401).JSON(fiber.Map{
			"status":  401,
			"message": err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": 200,
		"token":  res,
	})
}

func UpdateUser(db *sql.DB, c *fiber.Ctx) error {

	defer com.PrintLog("======= UPDATEUSER Function END   =========")

	com.GenerateRandomID()

	com.PrintLog("======= UPDATEUSER Function START =========")

	var userInfo entity.UserInfo
	var users string

	if err := c.BodyParser(&userInfo); err != nil {
		com.PrintLog(fmt.Sprintf("(UPDATEUSER:0001)%s", err))
		return c.Status(501).JSON(fiber.Map{
			"status":  501,
			"message": "Error on Server Side",
		})
	}

	token := c.Get("Authorization")

	encrypted, err := encryptor.VerifyField(token)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(UPDATEUSER:0002): %s", err))
		return c.Status(403).JSON(fiber.Map{
			"status":  http.StatusForbidden,
			"message": "Forbidden Login",
		})
	}

	rc, err := encryptor.Auth(encrypted, db, &userInfo.UserID, &users)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(UPDATEUSER:0003): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "Login Failed",
		})
	}

	rc, err = repository.UpdateUserInfo(db, &userInfo, users)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(UPDATEUSER:0004): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  rc,
			"message": fmt.Sprintf("%s", err),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Update successful",
	})
}

func Register(db *sql.DB, c *fiber.Ctx) error {

	defer com.PrintLog("======= Register Function END =========")

	com.GenerateRandomID() //Generate Global Id

	com.PrintLog("======= Register Function START =========")

	var register RegisterStruct
	var userBase entity.UserBase
	var userInfo entity.UserInfo

	if err := c.BodyParser(&register); err != nil {
		com.PrintLog(fmt.Sprintf("(REGISTER:0001)%s", err)) //Only show actual error on Log
		return c.Status(501).JSON(fiber.Map{
			"status":  501,
			"message": "Error on Server Side",
		})
	}

	userBase = register.Base
	userBase.UserType = 1
	com.PrintLog(userBase.Pin)

	if userBase.Pin == "" || len(userBase.Pin) != 6 {
		com.PrintLog("USER PIN CANNOT BE NULL AND LENGTH MUST BE MORE THAN 6")
		return c.Status(409).JSON(fiber.Map{
			"status":  409,
			"message": "User PIN must be filled and length must be more than 6",
		})
	}

	encryptedPin, err := encryptor.GeneratePin(userBase.Pin)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(Register:0052): %s", err))
		return c.Status(501).JSON(fiber.Map{
			"status":  501,
			"message": "error on server side",
		})
	}
	userBase.Pin = encryptedPin
	userInfo = register.Info
	userInfo.CreatedBy = userBase.Username
	userInfo.UpdatedBy = userBase.Username

	if userBase.Username == "" || userBase.Password == "" {
		com.PrintLog("(REGISTER:0000) Username or Password must be filled") //Only show actual error on Log
		return c.Status(401).JSON(fiber.Map{
			"status":  401,
			"message": "Username or Password must be filled",
		})
	}

	encrypted, err := encryptor.PasswordGenerator(userBase.Password)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(REGISTER:0002)%s", err)) //Only show actual error on Log
		return c.Status(501).JSON(fiber.Map{
			"status":  501,
			"message": "Error on Server Side",
		})
	}

	userBase.Password = encrypted
	userBase.CreatedAt = time.Now()

	rc, err := repository.RegisterUser(db, &userBase, &userInfo)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(REGISTER:0004)%s", err)) //Only show actual error on Log
		return c.Status(rc).JSON(fiber.Map{
			"status":  rc,
			"message": fmt.Sprintf("%s", err),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Successfuly register " + userBase.Username,
	})
}

func DeactiveUserAccount(db *sql.DB, c *fiber.Ctx) error {

	defer com.PrintLog("============ DEACTIVEUSERACCOUNT END ===============")

	com.GenerateRandomID()

	com.PrintLog("============ DEACTIVEUSERACCOUNT START ===============")

	var userBase entity.UserBase

	if err := c.BodyParser(&userBase); err != nil {
		com.PrintLog(fmt.Sprintf("(DEACTIVEUSERACCOUNT:0001)%s", err))
		return c.Status(501).JSON(fiber.Map{
			"status":  501,
			"message": "Error on Server Side",
		})
	}

	token := c.Get("Authorization")

	encrypted, err := encryptor.VerifyField(token)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(DACTIVEUSERACCOUNT:0002): %s", err))
		return c.Status(403).JSON(fiber.Map{
			"status":  http.StatusForbidden,
			"message": "Forbidden Login",
		})
	}

	rc, err := encryptor.Auth(encrypted, db, &userBase.ID, &userBase.Username)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(DEACTIVEUSERACCOUNT:0003): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "Login Failed",
		})
	}

	encryptedPass, err := encryptor.PasswordGenerator(userBase.Password)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(DACTIVEUSERACCOUNT:0002)%s", err))
		return c.Status(501).JSON(fiber.Map{
			"status":  501,
			"message": "Error on Server Side",
		})
	}

	rc, err = repository.DeactiveUser(db, &userBase, encryptedPass)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(DEACTIVATEUSERACCOUNT:0003) %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  rc,
			"message": fmt.Sprintf("%s", err),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Account successfully deleted",
	})
}

func GetUser(db *sql.DB, c *fiber.Ctx) error {

	defer com.PrintLog("==================== GET USER START ======================")

	com.PrintLog("==================== GET USER START ======================")

	var userInfo entity.UserInfo
	var username string

	token := c.Get("Authorization")

	encrypted, err := encryptor.VerifyField(token)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(GETUSER:0001): %s", err))
		return c.Status(403).JSON(fiber.Map{
			"status":  http.StatusForbidden,
			"message": "Forbidden Login",
		})
	}

	rc, err := encryptor.Auth(encrypted, db, &userInfo.UserID, &username)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(GETUSER:0002): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "Login Failed",
		})
	}

	rc, err = repository.GetUserInfo(db, &userInfo)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(GETUSER:0003): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  rc,
			"message": "Login Failed",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": 200,
		"data":   userInfo,
	})
}

func Welcome(db *sql.DB, c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "WELCOME",
	})
}
