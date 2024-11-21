package repository

import (
	"database/sql"
	"digibank/domain/entity"
	com "digibank/infrastructure/functions"
	"fmt"
	"net/http"
	"strconv"
)

func ComparePin(db *sql.DB, userBase entity.UserBase, encryptedPin string) (rc int, err error) {
	defer com.PrintLog(" ==================== COMPARE PIN END ===================")
	com.PrintLog("================ COMPARE PIN START ==================")

	var pin string

	sqls := "SELECT pin from user_base WHERE username = $1 AND status = 0"

	errors := db.QueryRow(sqls, userBase.Username).Scan(&pin)

	if errors != nil {
		com.PrintLog(fmt.Sprintf("(LOGINUSER:9) : %s", err))
		return 401, fmt.Errorf("errors (LOGINUSER:9): wrong password or username")
	}

	if pin != encryptedPin {
		com.PrintLog("(COMPAREPIN:0001) WRONG PIN!")
		return http.StatusUnauthorized, fmt.Errorf("wrong pin")
	}

	com.PrintLog("PIN OK!")

	return http.StatusOK, nil
}

func CreateAccount(db *sql.DB, userBase entity.UserBase, account entity.AccountBase) (rc int, err error) {

	defer com.PrintLog("======================== CREATE ACCOUNT END ==========================")

	com.PrintLog("======================== CREATE ACCOUNT START ==========================")

	com.PrintLog(fmt.Sprintf("userBase: %v", userBase))
	com.PrintLog(fmt.Sprintf("accountBase : %v", account))

	tx, err := db.Begin()

	if err != nil {
		com.PrintLog(fmt.Sprintf("(CREATEACCOUNT:0001) %s", err))
		return 501, fmt.Errorf("error on server side")
	}

	sqls := `
			INSERT INTO account_base (
				account_number, user_id, account_type, 
				bs_balance, current_balance, joint_with, 
				created_at, updated_at, status
			)
			VALUES ($1, $2, 1, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0)
			`

	res, err := tx.Exec(sqls, account.AccountNumber, userBase.ID,
		0, 0, nil)

	if err != nil {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("(CREATEACCOUNT:0002): %s", err))
		return 501, fmt.Errorf("error on server side")
	}

	cnt, err := res.RowsAffected()

	if err != nil || cnt < 1 {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
		if err != nil {
			com.PrintLog(fmt.Sprintf("(CREATEACCOUNT:0003): %s", err))
		}
		return 501, fmt.Errorf("error on server side")
	}

	tx.Commit()

	com.PrintLog("CREATE ACCOUNT OKE")

	return 200, nil
}

func GetAllAccount(db *sql.DB, userBase entity.UserBase, account *[]entity.AccountBase) (rc int, err error) {

	defer com.PrintLog("======================== GET ALL ACCOUNT END ==========================")

	com.PrintLog("======================== GET ALL ACCOUNT START ==========================")

	com.PrintLog(fmt.Sprintf("userBase: %v", userBase))
	com.PrintLog(fmt.Sprintf("accountBase : %v", account))

	var cnt int = 0

	sqls := "SELECT count(1) from account_base WHERE user_id = $1 AND status = 0"

	err = db.QueryRow(sqls, userBase.ID).Scan(&cnt)

	if err != nil || cnt == 0 {
		com.PrintLog(fmt.Sprintf("(GETACCOUNTMUtATION:0010) %s", err))
		return 404, fmt.Errorf("account not found")
	}

	sqls = "SELECT account_number, account_type, bs_balance ,current_balance, created_at from account_base WHERE user_id = $1 AND status = 0"

	res, errors := db.Query(sqls, userBase.ID)

	if errors != nil {
		if errors == sql.ErrNoRows {
			com.PrintLog(fmt.Sprintf("(GETALLACCOUNT:0001): %s", errors))
			return 401, fmt.Errorf("account not found")
		}
		return 501, fmt.Errorf("error on server side")
	}

	defer res.Close()
	for res.Next() {
		var acct entity.AccountBase
		errs := res.Scan(&acct.AccountNumber, &acct.AccountType, &acct.BsBalance, &acct.CurrentBalance, &acct.CreatedAt)
		if errs != nil {
			com.PrintLog(fmt.Sprintf("(GETALLACCOUNT:0002) %s", err))
			return 501, fmt.Errorf("error on server side")
		}

		*account = append(*account, acct)
	}

	return 200, nil
}

func GetAccount(db *sql.DB, userBase entity.UserBase, account *entity.AccountBase) (rc int, err error) {
	defer com.PrintLog("======================== GET ACCOUNT END ==========================")

	com.PrintLog("======================== GET ACCOUNT START ==========================")

	com.PrintLog(fmt.Sprintf("userBase: %v", userBase))
	com.PrintLog(fmt.Sprintf("accountBase : %v", account))

	sqls := "SELECT account_type, bs_balance ,current_balance, created_at from account_base WHERE user_id = $1 AND account_number = $2 AND status = 0"

	errors := db.QueryRow(sqls, userBase.ID, account.AccountNumber).Scan(&account.AccountType, &account.BsBalance, &account.CurrentBalance, &account.CreatedAt)

	if errors != nil {
		com.PrintLog(fmt.Sprintf("(GETACCOUNT:0001) %s", errors))
		return 404, fmt.Errorf("account not found")
	}

	return 200, nil
}

func CloseAccount(db *sql.DB, userBase entity.UserBase, account *entity.AccountBase) (rc int, err error) {
	defer com.PrintLog("======================== CLOSE ACCOUNT END ==========================")

	com.PrintLog("======================== CLOSE ACCOUNT START ==========================")

	com.PrintLog(fmt.Sprintf("userBase: %v", userBase))
	com.PrintLog(fmt.Sprintf("accountBase : %v", account))

	sqls := `
			UPDATE account_base
			SET status = 1, updated_at = CURRENT_TIMESTAMP
			WHERE account_number = $1
			AND status = 0
			`

	res, err := db.Exec(sqls, account.AccountNumber)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(CLOSEACCOUNT:0001): %s", err))
		return 501, fmt.Errorf("error on server side")
	}

	cnt, err := res.RowsAffected()

	if err != nil || cnt < 1 {
		com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
		if err != nil {
			com.PrintLog(fmt.Sprintf("(CLOSEACCOUNT:0002): %s", err))
		}
		return 501, fmt.Errorf("error on server side")
	}

	return 200, nil
}

func CreateTimeDepositAcct(db *sql.DB, userBase entity.UserBase, account *entity.AccountBase, ledger *entity.LedgerTransaction, accrual entity.AccrualBase) (rc int, err error) {

	defer com.PrintLog("======================== CREATE TIME DEPOSIT END ==========================")

	com.PrintLog("======================== CREATE TIME DEPOSIT START ==========================")

	var interest_rate float64 = 5.5
	var totalDays int
	var strTemp string

	tx, err := db.Begin()

	if err != nil {
		com.PrintLog(fmt.Sprintf("(CREATETIMEDEPOACC:0001) %s", err))
		return 501, fmt.Errorf("error on server side")
	}

	sqls := `
			SELECT current_balance
			FROM   account_base
			WHERE  account_number = $1
			AND    user_id = $2;
	`

	err = tx.QueryRow(sqls, account.AccountNumber, userBase.ID).Scan(&account.CurrentBalance)

	if ledger.TrxAmt > account.CurrentBalance {
		tx.Rollback()
		com.PrintLog("(CREATETIMEDEPOACC:0002): Balance less than trx amount")
		return http.StatusForbidden, fmt.Errorf("balance less than trx amount")
	}

	duration := accrual.ToDate.Sub(accrual.FromDate)

	com.PrintLog(fmt.Sprintf("DURATION     =    %v", duration))

	totalDays = int(duration.Hours() / 24)

	com.PrintLog(fmt.Sprintf("TOTAL DAYS   =   %d", totalDays))

	accrual.TotalAccrual = ((ledger.TrxAmt * interest_rate) / 100) / 365

	accrual.TotalAccrual = accrual.TotalAccrual * float64(totalDays)

	strTemp = fmt.Sprintf("%.2f", accrual.TotalAccrual)

	res, _ := strconv.ParseFloat(strTemp, 64)

	accrual.TotalAccrual = res

	sqls = `
			INSERT INTO account_base (
				account_number, user_id, account_type, 
				bs_balance, current_balance, joint_with, 
				created_at, updated_at, status
			)
			VALUES ($1, $2, 2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0)
`

	result, err := tx.Exec(sqls, account.AccountNumber, userBase.ID,
		0, 0, nil)

	if err != nil {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("(CREATETIMEDEPO:0002): %s", err))
		return 501, fmt.Errorf("error on server side")
	}

	cnt, err := result.RowsAffected()

	if err != nil || cnt < 1 {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
		if err != nil {
			com.PrintLog(fmt.Sprintf("(CREATETIMEDEPO:0003): %s", err))
		}
		return 501, fmt.Errorf("error on server side")
	}

	sqls = `
	INSERT INTO accrual_base (
		id, account_number, total_history, base_amount, interest_rate, 
		total_accrual, current_accrued, from_date, to_date, 
		status, created_at, updated_at
	) 
	VALUES (
		(SELECT count(1) + 1 FROM accrual_base), $1, 0, $2, 5.5, $3, 
		0, $4, $5, 0, 
		CURRENT_DATE, CURRENT_DATE
	)
`

	result, err = tx.Exec(sqls, ledger.AccountNumber, ledger.TrxAmt, accrual.TotalAccrual, accrual.FromDate, accrual.ToDate)

	if err != nil {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("(CREATETIMEDEPO:0001): %s", err))
		return 501, fmt.Errorf("error on server side")
	}

	cnt, err = result.RowsAffected()

	if err != nil || cnt < 1 {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
		return 501, fmt.Errorf("error on server side")
	}

	return 200, nil
}
