package repository

import (
	"database/sql"
	"digibank/domain/entity"
	com "digibank/infrastructure/functions"
	"fmt"
	"net/http"
)

func TopUpBalance(db *sql.DB, userBase entity.UserBase, ledger *entity.LedgerTransaction) (rc int, err error) {

	defer com.PrintLog("======================== TOPUPBALANCE END ================================")

	com.PrintLog("======================== TOPUPBALANCE START ================================")

	com.PrintLog(fmt.Sprintf("userBase: %v", userBase))
	com.PrintLog(fmt.Sprintf("ledger : %v", ledger))

	var baseCnt int

	tx, err := db.Begin()

	if err != nil {
		com.PrintLog(fmt.Sprintf("(TOPUPBALANCE:0001) %s", err))
		return 501, fmt.Errorf("error on server side")
	}

	sqls := `
		SELECT current_balance 
		FROM   account_base
		WHERE  account_number = $1
		AND    user_id = $2
	`

	err = tx.QueryRow(sqls, ledger.AccountNumber, userBase.ID).Scan(&ledger.BeforeBalance)

	if err != nil {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("(TOPUPBALANCE:0002): %s", err))
		return 501, fmt.Errorf("error on server side")
	}

	ledger.AfterBalance = ledger.BeforeBalance + ledger.TrxAmt

	com.PrintLog(fmt.Sprintf("Current Balance = %.2f", ledger.BeforeBalance))
	com.PrintLog(fmt.Sprintf("After Balance = %.2f", ledger.AfterBalance))

	sqls = `
	INSERT INTO ledger_transaction (
	    id, account_number, before_balance, trx_amt, 
	    after_balance, global_id, remark, 
	    created_at, created_by
	) 
	VALUES (
	    (SELECT count(1) + 1 FROM ledger_transaction),$1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, $7
	)
`

	res, err := tx.Exec(sqls, ledger.AccountNumber, ledger.BeforeBalance, ledger.TrxAmt, ledger.AfterBalance, com.GlobalId, "Top Up Balance", userBase.ID)

	if err != nil {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("(TOPUPBALANCE:0004): %s", err))
		return 501, fmt.Errorf("error on server side")
	}

	cnt, err := res.RowsAffected()

	if err != nil || cnt < 1 {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
		if err != nil {
			com.PrintLog(fmt.Sprintf("(TOPUPBALANCE:0005): %s", err))
		}
		return 501, fmt.Errorf("error on server side")
	}

	sqls = `
		SELECT count(1)
		FROM   ledger_base
		WHERE  account_code IN (
		'10908213',
		'24033123'
		)
		AND reference_no = $1
	`

	err = tx.QueryRow(sqls, ledger.AccountNumber).Scan(&baseCnt)

	if err != nil {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("(TOPUPBALANCE:0006): %s", err))
		return 501, fmt.Errorf("error on server side")
	}

	if baseCnt != 2 {
		sqls = `INSERT INTO ledger_base (
    			id, reference_no,account_code, trx_amount, global_id, created_at, updated_at
				) 
				VALUES 
    				((SELECT count(1) + 1 FROM ledger_base),$1,'10908213', $2, $3, CURRENT_DATE, CURRENT_DATE)
				ON CONFLICT (account_code,reference_no) 
				DO NOTHING;`

		res, err := tx.Exec(sqls, ledger.AccountNumber, ledger.TrxAmt, com.GlobalId)

		if err != nil {
			tx.Rollback()
			com.PrintLog(fmt.Sprintf("(TOPUPBALANCE:0080): %s", err))
			return 501, fmt.Errorf("error on server side")
		}

		cnt, err := res.RowsAffected()

		if err != nil || cnt < 1 {
			tx.Rollback()
			com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
			if err != nil {
				com.PrintLog(fmt.Sprintf("(TOPUPBALANCE:0081): %s", err))
			}
			return 501, fmt.Errorf("error on server side")
		}

		sqls = `INSERT INTO ledger_base (
    			id, reference_no,account_code, trx_amount, global_id, created_at, updated_at
				) 
				VALUES 
    				((SELECT count(1) + 1 FROM ledger_base),$1,'24033123', $2, $3, CURRENT_DATE, CURRENT_DATE)
				ON CONFLICT (account_code, reference_no) 
				DO NOTHING;`

		res, err = tx.Exec(sqls, ledger.AccountNumber, ledger.TrxAmt, com.GlobalId)

		if err != nil {
			tx.Rollback()
			com.PrintLog(fmt.Sprintf("(TOPUPBALANCE:0082): %s", err))
			return 501, fmt.Errorf("error on server side")
		}

		cnt, err = res.RowsAffected()

		if err != nil || cnt < 1 {
			tx.Rollback()
			com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
			if err != nil {
				com.PrintLog(fmt.Sprintf("(TOPUPBALANCE:0083): %s", err))
			}
			return 501, fmt.Errorf("error on server side")
		}
	} else {
		sqls = `UPDATE ledger_base 
        SET updated_at = CURRENT_DATE, trx_amount = trx_amount + $1
        WHERE (account_code = '10908213' OR account_code = '24033123')
		AND   reference_no = $2`

		res, err := tx.Exec(sqls, ledger.TrxAmt, ledger.AccountNumber)

		if err != nil {
			tx.Rollback()
			com.PrintLog(fmt.Sprintf("(TOPUPBALANCE:0009): %s", err))
			return 501, fmt.Errorf("error on server side")
		}

		cnt, err := res.RowsAffected()

		if err != nil || cnt < 1 {
			tx.Rollback()
			com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
			if err != nil {
				com.PrintLog(fmt.Sprintf("(UPDATELEDGER:0010): %s", err))
			}
			return 501, fmt.Errorf("error on server side")
		}

	}

	sqls = `
		    UPDATE ledger_master
		    SET
		        account_type = CASE
		            WHEN account_code LIKE '1%' OR account_code LIKE '2%' THEN 'B'
		            WHEN account_code LIKE '3%' OR account_code LIKE '4%' OR account_code LIKE '5%' THEN 'P'
		        END,
		        debit_amt = CASE
		            WHEN account_code LIKE '1%' OR account_code LIKE '5%' THEN debit_amt + $1
	            ELSE debit_amt
		        END,
		        credit_amt = CASE
		            WHEN account_code LIKE '2%' OR account_code LIKE '3%' OR account_code LIKE '4%' THEN credit_amt + $1
		            ELSE credit_amt
		        END,
		        total_amt = total_amt + $1,
		        updated_at = CURRENT_TIMESTAMP
		    WHERE account_code IN(
		 		'10908213',
		 		'24033123'
		 	);
			`

	res, err = tx.Exec(sqls, ledger.TrxAmt)
	if err != nil {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("(TOPUPBALANCE:0011): %s", err))
		return 501, fmt.Errorf("error on server side")
	}

	cnt, err = res.RowsAffected()
	if err != nil || cnt < 1 {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
		if err != nil {
			com.PrintLog(fmt.Sprintf("(TOPUPBALANCE:0012): %s", err))
		}
		return 501, fmt.Errorf("error on server side")
	}

	sqls = `UPDATE account_base 
        SET updated_at = CURRENT_TIMESTAMP, bs_balance = $1, current_balance = $2
        WHERE account_number = $3`

	res, err = tx.Exec(sqls, ledger.AfterBalance, ledger.AfterBalance, ledger.AccountNumber)

	if err != nil {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("(TOPUPBALANCE:0085): %s", err))
		return 501, fmt.Errorf("error on server side")
	}

	cnt, err = res.RowsAffected()

	if err != nil || cnt < 1 {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
		if err != nil {
			com.PrintLog(fmt.Sprintf("(TOPUPBALANC:0086): %s", err))
		}
		return 501, fmt.Errorf("error on server side")
	}

	tx.Commit()

	return http.StatusOK, nil
}

func Overbook(db *sql.DB, userBase entity.UserBase, ledgerFrom *entity.LedgerTransaction, ledgerTo *entity.LedgerTransaction) (rc int, err error) {

	defer com.PrintLog("============================ OVERBOOK END ==========================")

	com.PrintLog("============================ OVERBOOK START ==========================")

	com.PrintLog(fmt.Sprintf("Ledger From  =  %v", ledgerFrom))
	com.PrintLog(fmt.Sprintf("Ledger To    =  %v", ledgerTo))

	var toCnt int

	tx, err := db.Begin()

	if err != nil {
		com.PrintLog(fmt.Sprintf("(TOPUPBALANCE:0001) %s", err))
		return 501, fmt.Errorf("error on server side")
	}

	sqls := `
		SELECT current_balance 
		FROM   account_base
		WHERE  account_number = $1
		AND    user_id = $2
	`

	err = tx.QueryRow(sqls, ledgerFrom.AccountNumber, userBase.ID).Scan(&ledgerFrom.BeforeBalance)

	if err != nil {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("(TOPUPBALANCE:0002): %s", err))
		return 501, fmt.Errorf("error on server side")
	}

	if ledgerFrom.BeforeBalance < ledgerFrom.TrxAmt {
		tx.Rollback()
		com.PrintLog("LEDGER FROM BALANCE LESS THAN TRX AMT")
		com.PrintLog(fmt.Sprintf("Current Balance = %.2f, trx Amount = %.2f", ledgerFrom.BeforeBalance, ledgerFrom.TrxAmt))
		return http.StatusForbidden, fmt.Errorf("balance less than transaction amount")
	}

	sqls = `
		SELECT current_balance 
		FROM   account_base
		WHERE  account_number = $1
	`

	err = tx.QueryRow(sqls, ledgerTo.AccountNumber).Scan(&ledgerTo.BeforeBalance)

	if err != nil {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("(OVERBOOK:0003): %s", err))
		return 501, fmt.Errorf("error on server side")
	}

	ledgerFrom.AfterBalance = ledgerFrom.BeforeBalance - ledgerFrom.TrxAmt

	ledgerTo.AfterBalance = ledgerTo.BeforeBalance + ledgerTo.TrxAmt

	sqls = `
	INSERT INTO ledger_transaction (
	    id, account_number, before_balance, trx_amt, 
	    after_balance, global_id, remark, 
	    created_at, created_by
	) 
	VALUES (
	    (SELECT count(1) + 1 FROM ledger_transaction),$1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, $7
	)
`

	res, err := tx.Exec(sqls, ledgerFrom.AccountNumber, ledgerFrom.BeforeBalance, ledgerFrom.TrxAmt, ledgerFrom.AfterBalance, com.GlobalId, "Top Up Balance", userBase.ID)

	if err != nil {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("(OVERBOOK:0004): %s", err))
		return 501, fmt.Errorf("error on server side")
	}

	cnt, err := res.RowsAffected()

	if err != nil || cnt < 1 {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
		if err != nil {
			com.PrintLog(fmt.Sprintf("(OVERBOOK:0005): %s", err))
		}
		return 501, fmt.Errorf("error on server side")
	}

	sqls = `
	INSERT INTO ledger_transaction (
	    id, account_number, before_balance, trx_amt, 
	    after_balance, global_id, remark, 
	    created_at, created_by
	) 
	VALUES (
	    (SELECT count(1) + 1 FROM ledger_transaction),$1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, $7
	)
`

	res, err = tx.Exec(sqls, ledgerTo.AccountNumber, ledgerTo.BeforeBalance, ledgerTo.TrxAmt, ledgerTo.AfterBalance, com.GlobalId, "Top Up Balance", userBase.ID)

	if err != nil {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("(OVERBOOK:0006): %s", err))
		return 501, fmt.Errorf("error on server side")
	}

	cnt, err = res.RowsAffected()

	if err != nil || cnt < 1 {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
		if err != nil {
			com.PrintLog(fmt.Sprintf("(TOPUPBALANCE:0007): %s", err))
		}
		return 501, fmt.Errorf("error on server side")
	}

	sqls = `UPDATE ledger_base 
        SET updated_at = CURRENT_DATE, trx_amount = trx_amount - $1
        WHERE (account_code = '10908213' OR account_code = '24033123')
		AND   reference_no = $2`

	res, err = tx.Exec(sqls, ledgerFrom.TrxAmt, ledgerFrom.AccountNumber)

	if err != nil {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("(OVERBOOK:0008): %s", err))
		return 501, fmt.Errorf("error on server side")
	}

	cnt, err = res.RowsAffected()

	if err != nil || cnt < 1 {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
		if err != nil {
			com.PrintLog(fmt.Sprintf("(OVERBOOK:0009): %s", err))
		}
		return 501, fmt.Errorf("error on server side")
	}

	sqls = `
		SELECT count(1)
		FROM   ledger_base
		WHERE  account_code IN (
		'10908213',
		'24033123'
		)
		AND reference_no = $1
	`

	err = tx.QueryRow(sqls, ledgerTo.AccountNumber).Scan(&toCnt)

	if toCnt != 2 {
		sqls = `INSERT INTO ledger_base (
    			id, reference_no,account_code, trx_amount, global_id, created_at, updated_at
				) 
				VALUES 
    				((SELECT count(1) + 1 FROM ledger_base),$1,'10908213', $2, $3, CURRENT_DATE, CURRENT_DATE)
				ON CONFLICT (account_code,reference_no) 
				DO NOTHING;`

		res, err := tx.Exec(sqls, ledgerTo.AccountNumber, ledgerTo.TrxAmt, com.GlobalId)

		if err != nil {
			tx.Rollback()
			com.PrintLog(fmt.Sprintf("(OVERBOOK:0080): %s", err))
			return 501, fmt.Errorf("error on server side")
		}

		cnt, err := res.RowsAffected()

		if err != nil || cnt < 1 {
			tx.Rollback()
			com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
			if err != nil {
				com.PrintLog(fmt.Sprintf("(OVERBOOK:0081): %s", err))
			}
			return 501, fmt.Errorf("error on server side")
		}

		sqls = `INSERT INTO ledger_base (
    			id, reference_no,account_code, trx_amount, global_id, created_at, updated_at
				) 
				VALUES 
    				((SELECT count(1) + 1 FROM ledger_base),$1,'24033123', $2, $3, CURRENT_DATE, CURRENT_DATE)
				ON CONFLICT (account_code, reference_no) 
				DO NOTHING;`

		res, err = tx.Exec(sqls, ledgerTo.AccountNumber, ledgerTo.TrxAmt, com.GlobalId)

		if err != nil {
			tx.Rollback()
			com.PrintLog(fmt.Sprintf("(OVERBOOK:0013): %s", err))
			return 501, fmt.Errorf("error on server side")
		}

		cnt, err = res.RowsAffected()

		if err != nil || cnt < 1 {
			tx.Rollback()
			com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
			if err != nil {
				com.PrintLog(fmt.Sprintf("(OVERBOOK:0014): %s", err))
			}
			return 501, fmt.Errorf("error on server side")
		}
	} else {
		sqls = `UPDATE ledger_base 
        SET updated_at = CURRENT_DATE, trx_amount = trx_amount + $1
        WHERE (account_code = '10908213' OR account_code = '24033123')
		AND   reference_no = $2`

		res, err := tx.Exec(sqls, ledgerTo.TrxAmt, ledgerTo.AccountNumber)

		if err != nil {
			tx.Rollback()
			com.PrintLog(fmt.Sprintf("(OVERBOOK:0016): %s", err))
			return 501, fmt.Errorf("error on server side")
		}

		cnt, err := res.RowsAffected()

		if err != nil || cnt < 1 {
			tx.Rollback()
			com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
			if err != nil {
				com.PrintLog(fmt.Sprintf("(OVERBOOK:0019): %s", err))
			}
			return 501, fmt.Errorf("error on server side")
		}

	}

	sqls = `UPDATE account_base 
        SET updated_at = CURRENT_TIMESTAMP, bs_balance = $1, current_balance = $2
        WHERE account_number = $3`

	res, err = tx.Exec(sqls, ledgerFrom.AfterBalance, ledgerFrom.AfterBalance, ledgerFrom.AccountNumber)

	if err != nil {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("(OVERBOOK:0085): %s", err))
		return 501, fmt.Errorf("error on server side")
	}

	cnt, err = res.RowsAffected()

	if err != nil || cnt < 1 {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
		if err != nil {
			com.PrintLog(fmt.Sprintf("(OVERBOOK:0086): %s", err))
		}
		return 501, fmt.Errorf("error on server side")
	}

	sqls = `UPDATE account_base 
        SET updated_at = CURRENT_TIMESTAMP, bs_balance = $1, current_balance = $2
        WHERE account_number = $3`

	res, err = tx.Exec(sqls, ledgerTo.AfterBalance, ledgerTo.AfterBalance, ledgerTo.AccountNumber)

	if err != nil {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("(OVERBOOK:0085): %s", err))
		return 501, fmt.Errorf("error on server side")
	}

	cnt, err = res.RowsAffected()

	if err != nil || cnt < 1 {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
		if err != nil {
			com.PrintLog(fmt.Sprintf("(OVERBOOK:0086): %s", err))
		}
		return 501, fmt.Errorf("error on server side")
	}

	tx.Commit()

	return 200, nil
}

func GetAccountMutation(db *sql.DB, userBase entity.UserBase, account entity.AccountBase, ledger *[]entity.LedgerTransaction) (rc int, err error) {

	defer com.PrintLog("====================== GETACCOUNTUTATION END =======================")

	com.PrintLog("====================== GETACCOUNTUTATION START =======================")

	com.PrintLog(fmt.Sprintf("userBase     =     %v", userBase))
	com.PrintLog(fmt.Sprintf("account      =     %v", account))

	var cnt int = 0

	sqls := `
			SELECT     count(1)
			FROM       ledger_transaction a 
			INNER JOIN account_base b ON a.account_number = b.account_number
			WHERE a.account_number = $1
			AND   b.user_id = $2
	`

	err = db.QueryRow(sqls, account.AccountNumber, userBase.ID).Scan(&cnt)

	if err != nil || cnt == 0 {
		com.PrintLog(fmt.Sprintf("(GETACCOUNTMUtATION:0010) %s", err))
		return 404, fmt.Errorf("account not found")
	}

	sqls = `
			SELECT     a.account_number, a.trx_amt, a.after_balance, a.created_at 
			FROM       ledger_transaction a 
			INNER JOIN account_base b ON a.account_number = b.account_number
			WHERE a.account_number = $1
			AND   b.user_id = $2
	`

	res, err := db.Query(sqls, account.AccountNumber, userBase.ID)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(GETACCOUNTMUTATION:0001): %s", err))
		if err == sql.ErrNoRows {
			com.PrintLog(fmt.Sprintf("(GETACCOUNTMUTATION:0001): %s", err))
			return 401, fmt.Errorf("account not found")
		}
		return 501, fmt.Errorf("error on server side")
	}

	defer res.Close()

	for res.Next() {
		var ledgerTran entity.LedgerTransaction
		errs := res.Scan(&ledgerTran.AccountNumber, &ledgerTran.TrxAmt, &ledgerTran.AfterBalance, &ledgerTran.CreatedAt)
		if errs != nil {
			com.PrintLog(fmt.Sprintf("(GETALLACCOUNT:0002) %s", err))
			return 501, fmt.Errorf("error on server side")
		}

		*ledger = append(*ledger, ledgerTran)
	}

	com.PrintLog(fmt.Sprintf("LIST : %v", ledger))

	com.PrintLog("GET MUTATION SUCCESS")

	return 200, nil
}
