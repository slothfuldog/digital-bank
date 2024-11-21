package controllers

import (
	"database/sql"
	"digibank/delivery/http/repository"
	"digibank/domain/entity"
	"digibank/infrastructure/encryptor"
	com "digibank/infrastructure/functions"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type OverbookTran struct {
	FromAcctNo string  `json:"from_account"`
	ToAcctNo   string  `json:"to_account"`
	TrxAmt     float64 `json:"trx_amt"`
}

func AddBalance(db *sql.DB, c *fiber.Ctx) error {

	defer com.PrintLog("=================== ADD BALANCE END =========================")

	com.GenerateRandomID()

	com.PrintLog("=================== ADD BALANCE START =========================")

	var transaction entity.LedgerTransaction
	var userBase entity.UserBase

	if err := c.BodyParser(&transaction); err != nil {
		com.PrintLog(fmt.Sprintf("(ADDBALANCE:0001): %s", err))
		return c.Status(403).JSON(fiber.Map{
			"status":  501,
			"message": "Error on server side",
		})
	}

	token := c.Get("Authorization")

	encrypted, err := encryptor.VerifyField(token)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(ADDBALANCE:0002): %s", err))
		return c.Status(403).JSON(fiber.Map{
			"status":  http.StatusForbidden,
			"message": "Forbidden Login",
		})
	}

	rc, err := encryptor.Auth(encrypted, db, &userBase.ID, &userBase.Username)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(ADDBALANCE:0003): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "Login Failed",
		})
	}

	rc, err = repository.TopUpBalance(db, userBase, &transaction)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(ADDBALANCE:0004): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  rc,
			"message": fmt.Sprintf("%s", err),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"data":    transaction,
		"message": "Balance successfully added",
	})
}

func Overbooking(db *sql.DB, c *fiber.Ctx) error {

	defer com.PrintLog("=================== ADD BALANCE END =========================")

	com.GenerateRandomID()

	com.PrintLog("=================== ADD BALANCE START =========================")

	var transactionFrom entity.LedgerTransaction
	var transactionTo entity.LedgerTransaction
	var userBase entity.UserBase
	var body OverbookTran

	if err := c.BodyParser(&body); err != nil {
		com.PrintLog(fmt.Sprintf("(OVERBOOKING:0001): %s", err))
		return c.Status(403).JSON(fiber.Map{
			"status":  501,
			"message": "Error on server side",
		})
	}

	token := c.Get("Authorization")

	encrypted, err := encryptor.VerifyField(token)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(OVERBOOKING:0002): %s", err))
		return c.Status(403).JSON(fiber.Map{
			"status":  http.StatusForbidden,
			"message": "Forbidden Login",
		})
	}

	rc, err := encryptor.Auth(encrypted, db, &userBase.ID, &userBase.Username)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(OVERBOOKING:0003): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "Login Failed",
		})
	}

	transactionFrom.AccountNumber = body.FromAcctNo
	transactionFrom.TrxAmt = body.TrxAmt
	transactionTo.AccountNumber = body.ToAcctNo
	transactionTo.TrxAmt = body.TrxAmt

	rc, err = repository.Overbook(db, userBase, &transactionFrom, &transactionTo)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(OVERBOOKING:0004): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  rc,
			"message": fmt.Sprintf("%s", err),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": http.StatusOK,
		"data": map[string]interface{}{
			"transactionFrom": transactionFrom,
			"transactionTo":   transactionTo.AccountNumber,
		},
		"message": "Overbook successfully executed",
	})
}

func GetUserMutation(db *sql.DB, c *fiber.Ctx) error {
	defer com.PrintLog("==================== GET USER MUTATION END ====================")

	com.GenerateRandomID()

	com.PrintLog("==================== GET USER MUTATION START ====================")

	var ledger []entity.LedgerTransaction
	var userBase entity.UserBase
	var account entity.AccountBase

	accountNumber := c.Params("accountnumber")

	account.AccountNumber = accountNumber

	token := c.Get("Authorization")

	encrypted, err := encryptor.VerifyField(token)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(GETUSERMUTATION:0002): %s", err))
		return c.Status(403).JSON(fiber.Map{
			"status":  http.StatusForbidden,
			"message": "Forbidden Login",
		})
	}

	rc, err := encryptor.Auth(encrypted, db, &userBase.ID, &userBase.Username)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(GETUSERMUTATION:0003): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "Login Failed",
		})
	}

	rc, err = repository.GetAccountMutation(db, userBase, account, &ledger)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(GETUSERMUTATION:0004): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  rc,
			"message": fmt.Sprintf("%s", err),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": http.StatusOK,
		"data":   ledger,
	})
}

func TimeDepositSimulation(db *sql.DB, c *fiber.Ctx) error {
	defer com.PrintLog("==================== TIME DEPOSIT SIMULATION END ====================")

	com.GenerateRandomID()

	com.PrintLog("==================== TIME DPOSIT SIMULATION START ====================")

	var body struct {
		TrxAmt   float64 `json:"trx_amt"`
		FromDate string  `json:"from_date"`
		ToDate   string  `json:"to_date"`
	}

	var result float64
	var interest_rate float64 = 5.5
	var totalDays int
	var strTemp string
	layout := "2006-01-02"
	var userBase entity.UserBase

	if err := c.BodyParser(&body); err != nil {
		com.PrintLog(fmt.Sprintf("(TIMEDEPOSIMU:0001): %s", err))
		return c.Status(403).JSON(fiber.Map{
			"status":  501,
			"message": "Error on server side",
		})
	}

	com.PrintLog(fmt.Sprintf("BODY     =   %v", body))

	token := c.Get("Authorization")

	encrypted, err := encryptor.VerifyField(token)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(TIMEDEPOSIMU:0002): %s", err))
		return c.Status(403).JSON(fiber.Map{
			"status":  http.StatusForbidden,
			"message": "Forbidden Login",
		})
	}

	rc, err := encryptor.Auth(encrypted, db, &userBase.ID, &userBase.Username)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(TIMEDEPOSIMU:0003): %s", err))
		return c.Status(rc).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "Login Failed",
		})
	}

	parsedTime, err := time.Parse(layout, body.FromDate)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(TIMEDEPOSIMU:00010): %s", err))
		return c.Status(501).JSON(fiber.Map{
			"status":  501,
			"message": "Error on server side",
		})
	}

	parsedTime2, err := time.Parse(layout, body.ToDate)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(TIMEDEPOSIMU:00011): %s", err))
		return c.Status(501).JSON(fiber.Map{
			"status":  501,
			"message": "Error on server side",
		})
	}

	com.PrintLog(fmt.Sprintf("PARSEDTIME 1   =   %v", parsedTime))
	com.PrintLog(fmt.Sprintf("PARSEDTIME 2   =   %v", parsedTime2))

	duration := parsedTime2.Sub(parsedTime)

	com.PrintLog(fmt.Sprintf("DURATION     =    %v", duration))

	totalDays = int(duration.Hours() / 24)

	if totalDays <= 10 {
		com.PrintLog("(TIMEDEPOSIM:0030): totalDays must be more than 10")
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"status":  http.StatusForbidden,
			"message": "total days must more than 10 days",
		})
	}

	com.PrintLog(fmt.Sprintf("TOTAL DAYS   =   %d", totalDays))

	result = ((body.TrxAmt * interest_rate) / 100) / 365

	result = result * float64(totalDays)

	strTemp = fmt.Sprintf("%.2f", result)

	res, err := strconv.ParseFloat(strTemp, 64)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(TIMEDEPOSIMU:0004): %s", err))
		return c.Status(501).JSON(fiber.Map{
			"status":  501,
			"message": "Error on server side",
		})
	}

	result = res

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": http.StatusOK,
		"result": result,
	})
}
