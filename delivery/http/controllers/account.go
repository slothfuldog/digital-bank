package controllers

import (
	"database/sql"
	"digibank/delivery/http/repository"
	"digibank/domain/entity"
	"digibank/infrastructure/encryptor"
	com "digibank/infrastructure/functions"
	"fmt"
	"net/http"
	_ "time"

	"github.com/gofiber/fiber/v2"
)

func CreateAccount(db *sql.DB, c *fiber.Ctx) error {

	defer com.PrintLog("================= Create Account END ==================")

	com.GenerateRandomID()

	com.PrintLog("================= Create Account START ==================")

	var userBase entity.UserBase

	var account entity.AccountBase

	token := c.Get("Authorization")

	encrypted, err := encryptor.VerifyField(token)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(CREATEACCOUNT:0002): %s", err))
		return c.Status(403).JSON(fiber.Map{
			"status":  http.StatusForbidden,
			"message": "Forbidden Login",
		})
	}

	rc, err := encryptor.Auth(encrypted, db, &userBase.ID, &userBase.Username)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(CREATEACCOUNT:0003): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "Login Failed",
		})
	}

	com.GenerateAccountNumber(&account.AccountNumber)

	com.PrintLog(fmt.Sprintf("ACCOUNT NUMBER : %s", account.AccountNumber))

	rc, err = repository.CreateAccount(db, userBase, account)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(CREATEACCOUNT:0004) %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  rc,
			"message": fmt.Sprintf("%s", err),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Successfully create account",
	})
}

func CheckPin(db *sql.DB, c *fiber.Ctx) error {

	defer com.PrintLog("================= CHECK PIN START =====================")

	com.GenerateRandomID()

	com.PrintLog("================= CHECK PIN START =====================")

	var pin struct {
		Pin string `json:"pin"`
	}

	var userBase entity.UserBase

	if err := c.BodyParser(&pin); err != nil {
		com.PrintLog(fmt.Sprintf("(CHECKPIN:0001)%s", err))
		return c.Status(501).JSON(fiber.Map{
			"status":  501,
			"message": "Error on Server Side",
		})
	}

	encryptedPin, err := encryptor.GeneratePin(pin.Pin)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(CHECKPIN:0002)%s", err))
		return c.Status(501).JSON(fiber.Map{
			"status":  501,
			"message": "Error on Server Side",
		})
	}

	token := c.Get("Authorization")

	encrypted, err := encryptor.VerifyField(token)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(CHEKPIN:0002): %s", err))
		return c.Status(403).JSON(fiber.Map{
			"status":  http.StatusForbidden,
			"message": "Forbidden Login",
		})
	}

	rc, err := encryptor.Auth(encrypted, db, &userBase.ID, &userBase.Username)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(CHEKPIN:0003): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "Login Failed",
		})
	}

	rc, err = repository.ComparePin(db, userBase, encryptedPin)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(CHEKPIN:0004): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  rc,
			"message": fmt.Sprintf("%s", err),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "PIN authorized",
	})
}

func GetAllAcct(db *sql.DB, c *fiber.Ctx) error {

	defer com.PrintLog("================= GETALLACCT START =====================")

	com.GenerateRandomID()

	com.PrintLog("================= GETALLACCT START =====================")

	var acct []entity.AccountBase
	var userBase entity.UserBase

	token := c.Get("Authorization")

	encrypted, err := encryptor.VerifyField(token)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(GETALLACCT:0001): %s", err))
		return c.Status(403).JSON(fiber.Map{
			"status":  http.StatusForbidden,
			"message": "Forbidden Login",
		})
	}

	rc, err := encryptor.Auth(encrypted, db, &userBase.ID, &userBase.Username)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(GETALLACCT:0002): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "Login Failed",
		})
	}

	rc, err = repository.GetAllAccount(db, userBase, &acct)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(GETALLACCT:0003): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  rc,
			"message": fmt.Sprintf("%s", err),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": 200,
		"data":   acct,
	})
}

func GetAcct(db *sql.DB, c *fiber.Ctx) error {
	defer com.PrintLog("================= GETACCT START =====================")

	com.GenerateRandomID()

	com.PrintLog("================= GETACCT START =====================")

	var acct entity.AccountBase
	var userBase entity.UserBase

	token := c.Get("Authorization")

	accountNo := c.Params("accountno")

	encrypted, err := encryptor.VerifyField(token)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(GETACCT:0001): %s", err))
		return c.Status(403).JSON(fiber.Map{
			"status":  http.StatusForbidden,
			"message": "Forbidden Login",
		})
	}

	rc, err := encryptor.Auth(encrypted, db, &userBase.ID, &userBase.Username)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(GETACCT:0002): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "Login Failed",
		})
	}

	acct.AccountNumber = accountNo

	rc, err = repository.GetAccount(db, userBase, &acct)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(GETACCT:0003): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  rc,
			"message": fmt.Sprintf("%s", err),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": 200,
		"data":   acct,
	})
}

func CloseAcct(db *sql.DB, c *fiber.Ctx) error {

	defer com.PrintLog("================= GETACCT START =====================")

	com.GenerateRandomID()

	com.PrintLog("================= GETACCT START =====================")

	var acct entity.AccountBase
	var userBase entity.UserBase

	token := c.Get("Authorization")

	accountNo := c.Params("accountno")

	encrypted, err := encryptor.VerifyField(token)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(GETACCT:0001): %s", err))
		return c.Status(403).JSON(fiber.Map{
			"status":  http.StatusForbidden,
			"message": "Forbidden Login",
		})
	}

	rc, err := encryptor.Auth(encrypted, db, &userBase.ID, &userBase.Username)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(GETACCT:0002): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "Login Failed",
		})
	}

	acct.AccountNumber = accountNo

	rc, err = repository.CloseAccount(db, userBase, &acct)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(GETACCT:0003): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  rc,
			"message": fmt.Sprintf("%s", err),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "Successfully delete account",
	})
}
