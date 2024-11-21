package repository

import (
	"database/sql"
	"digibank/domain/entity"
	com "digibank/infrastructure/functions"
	"fmt"
	"net/http"
	"time"
)

func LoginUser(db *sql.DB, user *entity.UserBase, macAddress string, encryptedPass string) (rc int, err error) {

	defer com.PrintLog("====== LOGIN USER END ========")

	com.PrintLog("====== LOGIN USER START ========")

	com.PrintLog(fmt.Sprintf("USERBASE:\n%v\n", user))

	sqls := "SELECT id, username, password, user_type from user_base WHERE username = $1 AND status = 0"

	errors := db.QueryRow(sqls, user.Username).Scan(&user.ID, &user.Username, &user.Password, &user.UserType)

	if errors != nil {
		com.PrintLog(fmt.Sprintf("(LOGINUSER:9) : %s", err))
		return 401, fmt.Errorf("errors (LOGINUSER:9): wrong password or username")
	}

	if encryptedPass != user.Password {
		com.PrintLog(fmt.Sprintf("(LOGINUSER:10) : %s", err))
		return 401, fmt.Errorf("errors (LOGINUSER:10): wrong password or username")
	}

	if macAddress == "" {
		com.PrintLog(fmt.Sprintf("(LOGINUSER:10) : %s", err))
		return 401, fmt.Errorf("errors (LOGINUSER:10): Forbidden Login")
	}

	sqls = `
    		INSERT INTO user_mac_addresses (user_id, mac_address, created_at, updated_at, id)
    		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, (SELECT count(1) + 1 from user_mac_addresses))`
	res, err := db.Exec(sqls, user.ID, macAddress)
	if err != nil {
		com.PrintLog(fmt.Sprintf("(LOGINUSER:12) : %s", err))
		return 501, fmt.Errorf("server side error")
	}

	cnt, err := res.RowsAffected()

	if err != nil || cnt < 1 {
		com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
		if err != nil {
			com.PrintLog(fmt.Sprintf("(LOGINUSER:13): %s", err))
		}
		return 501, fmt.Errorf("error on server side")
	}

	user.Password = ""

	return 200, nil
}

func KeepLogin(db *sql.DB, username string) (userid int, rc int, err error) {
	defer com.PrintLog("====== KEEP LOGIN USER END ========")

	com.PrintLog("====== KEEP LOGIN USER START ========")

	var user int

	sqls := "SELECT id from user_base WHERE username = $1 AND status = 0"

	errors := db.QueryRow(sqls, username).Scan(&user)

	if errors != nil {
		if errors == sql.ErrNoRows {
			com.PrintLog(fmt.Sprintf("(KEEPLOGIN:0001) %s", errors))
			return 0, http.StatusForbidden, fmt.Errorf("forbidden login")
		} else {
			com.PrintLog(fmt.Sprintf("(KEEPLOGIN:0002) %s", errors))
			return 0, 501, fmt.Errorf("failed to login")
		}
	}

	com.PrintLog("KEEP LOGIN SUCCSS")

	return user, 200, nil
}

func RegisterUser(db *sql.DB, user *entity.UserBase, userInfo *entity.UserInfo) (rc int, err error) {

	defer com.PrintLog("====== REGISTER USER END ========")

	com.PrintLog("====== REGISTER USER START ========")

	tx, err := db.Begin()

	if err != nil {
		com.PrintLog(fmt.Sprintf("(REGISTERUSER:9) : %s", err))
		return 501, fmt.Errorf("errors (REGISTERUSER:9): ERROR DB BEGIN")
	}

	var cnt int
	var maxBaseId int
	var maxInfoId int

	layout := "2006-01-02"

	com.PrintLog(fmt.Sprintf("USERBASE:\n%v\n", user))
	com.PrintLog(fmt.Sprintf("USERINFO:\n%v\n", userInfo))

	time, err := time.Parse(layout, userInfo.BirthDate)

	sqls := `
			SELECT count(1)
			FROM   user_base
			WHERE  username = $1
	`

	errors := tx.QueryRow(sqls, user.Username).Scan(&cnt)

	if cnt != 0 {
		com.PrintLog(fmt.Sprintf("(REGISTERUSER:10) : %s", errors))
		return 409, fmt.Errorf("errors (REGISTERUSER:10): USER ALREADY EXIST")
	}

	sqls = `
			SELECT count(1) +1
			FROM   user_base
	`

	errors = tx.QueryRow(sqls).Scan(&maxBaseId)

	if errors != nil {
		com.PrintLog(fmt.Sprintf("(REGISTERUSER:11) : %s", errors))
		return 501, fmt.Errorf("errors (REGISTERUSER:11): FAILED TO REGISTER")
	}

	sqls = `
			SELECT count(1) + 1
			FROM   user_info
	`

	errors = tx.QueryRow(sqls).Scan(&maxInfoId)

	if errors != nil {
		com.PrintLog(fmt.Sprintf("(REGISTERUSER:12) : %s", errors))
		return 501, fmt.Errorf("errors (REGISTERUSER:12): FAILED TO REGISTER")
	}

	sqls = `
			INSERT INTO user_base (username, password, user_type, created_at, updated_at, id, status, pin)
			VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $4, 0, $5)
			RETURNING id`

	errors = tx.QueryRow(sqls, user.Username, user.Password, user.UserType, maxBaseId, user.Pin).Scan(&user.ID)

	if errors != nil {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("(REGISTERUSER:13) : %s", errors))
		return 501, fmt.Errorf("errors (REGISTERUSER:13): FAILED TO REGISTER")
	}

	com.PrintLog(fmt.Sprintf("USER BASE ID =   %d", user.ID))

	sqls = `
			INSERT INTO user_info (
			    user_id, name, gender, birth_date, address, 
			    occupation, job_place, email_address, phone_number, 
			    created_at, created_by, updated_at, updated_by, id
			) 
			VALUES (
			    $1, $2, $3, $4, $5, 
			    $6, $7, $8, $9, 
			    CURRENT_TIMESTAMP, $10, CURRENT_TIMESTAMP, $11, $12
			)
			RETURNING id`

	errors = tx.QueryRow(sqls, maxBaseId, userInfo.Name, userInfo.Gender, time, userInfo.Address, userInfo.Occupation, userInfo.JobPlace, userInfo.EmailAddress, userInfo.PhoneNumber, userInfo.CreatedBy, userInfo.UpdatedBy, maxInfoId).Scan(&userInfo.ID)

	if errors != nil {
		tx.Rollback()
		com.PrintLog(fmt.Sprintf("(REGISTERUSER:14) : %s", errors))
		return 501, fmt.Errorf("errors (REGISTERUSER:14): FAILED TO REGISTER")
	}

	tx.Commit()

	com.PrintLog(fmt.Sprintf("USER INFO ID =   %d", userInfo.ID))

	return 200, nil
}

func UpdateUserInfo(db *sql.DB, userInfo *entity.UserInfo, username string) (rc int, err error) {

	defer com.PrintLog("================= UPDATE USER END ================")

	com.PrintLog("================= UPDATE USER START ================")

	com.PrintLog(fmt.Sprintf("USERINFO %v", userInfo))

	if userInfo.Name == "" || userInfo.PhoneNumber == "" {
		com.PrintLog("ERROR : Name and Phone numbe null")
		com.PrintLog(fmt.Sprintf("NAME :%s Phone Number %s", userInfo.Name, userInfo.PhoneNumber))
		return 409, fmt.Errorf("name or phone must be filled")
	}

	sqls := `UPDATE user_info
			SET
			    name = $1,
			    address = $2,
			    occupation = $3,
			    job_place = $4,
			    email_address = $5,
			    phone_number = $6,
			    updated_at = CURRENT_TIMESTAMP,
			    updated_by = $7
			WHERE user_id = $8;`

	res, err := db.Exec(sqls, userInfo.Name, userInfo.Address, userInfo.Occupation, userInfo.JobPlace, userInfo.EmailAddress, userInfo.PhoneNumber, username, userInfo.UserID)
	if err != nil {
		com.PrintLog(fmt.Sprintf("(UPDATEUSERINFO:0001) : %s", err))
		return 501, fmt.Errorf("server side error")
	}

	cnt, err := res.RowsAffected()

	if err != nil || cnt < 1 {
		com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
		if err != nil {
			com.PrintLog(fmt.Sprintf("(UPDATEUSERINFO:0002): %s", err))
		}
		return 501, fmt.Errorf("error on server side")
	}

	return 200, nil
}

func DeactiveUser(db *sql.DB, user *entity.UserBase, encryptedPass string) (rc int, err error) {

	defer com.PrintLog("============ DACTIVE USER END ===================")

	com.PrintLog("============ DACTIVE USER START ===================")

	var realPassword string

	com.PrintLog(fmt.Sprintf("USER : %v", user))

	sqls := "SELECT password from user_base WHERE username = $1 AND status = 0"

	errors := db.QueryRow(sqls, user.Username).Scan(&realPassword)

	if errors != nil {
		if errors == sql.ErrNoRows {
			com.PrintLog(fmt.Sprintf("(DACTIVEUSER:0001) %s", errors))
			return http.StatusForbidden, fmt.Errorf("forbidden login")
		} else {
			com.PrintLog(fmt.Sprintf("(DACTIVEUSER:0002) %s", errors))
			return 501, fmt.Errorf("failed to login")
		}
	}

	if realPassword != encryptedPass {
		com.PrintLog("Wrong password!")
		return http.StatusUnauthorized, fmt.Errorf("wrong password")
	}

	com.PrintLog("USER FOUND !!")

	sqls = `
			UPDATE user_base
			SET    status = 1
			WHERE  username = $1
	`

	res, err := db.Exec(sqls, user.Username)
	if err != nil {
		com.PrintLog(fmt.Sprintf("(DEACTIVEUSER:0003) : %s", err))
		return 501, fmt.Errorf("server side error")
	}

	cnt, err := res.RowsAffected()

	if err != nil || cnt < 1 {
		com.PrintLog(fmt.Sprintf("Row Affected %d", cnt))
		if err != nil {
			com.PrintLog(fmt.Sprintf("(DEACTIVEUSER:0003): %s", err))
		}
		return 501, fmt.Errorf("error on server side")
	}

	return 200, nil
}

func GetUserInfo(db *sql.DB, userInfo *entity.UserInfo) (rc int, err error) {
	defer com.PrintLog("============ GET USER END ===================")

	com.PrintLog("============ GET USER START ===================")

	sqls := `
	SELECT 
		name, 
		gender, 
		birth_date, 
		address, 
		occupation, 
		job_place, 
		email_address, 
		phone_number 
	FROM 
		user_info 
	WHERE 
		user_id = $1
`

	errors := db.QueryRow(sqls, userInfo.UserID).Scan(
		&userInfo.Name,
		&userInfo.Gender,
		&userInfo.BirthDate,
		&userInfo.Address,
		&userInfo.Occupation,
		&userInfo.JobPlace,
		&userInfo.EmailAddress,
		&userInfo.PhoneNumber,
	)

	if errors != nil {
		com.PrintLog(fmt.Sprintf("(GETUSERINFO:0001) : %s", err))
		return 501, fmt.Errorf("server side error")
	}

	return 200, nil
}
