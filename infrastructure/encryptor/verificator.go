package encryptor

import (
	"database/sql"
	"digibank/delivery/http/repository"
	com "digibank/infrastructure/functions"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"aidanwoods.dev/go-paseto"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/argon2"
)

func VerifyPassword(encryptedPass, password string) (isMatch bool, e error) {
	path, err := os.Getwd()
	if err != nil {
		com.PrintLog(fmt.Sprintf("(VERIFYPASS: 1000) %s", err))
		return false, err
	}
	currDir := fmt.Sprint(path, "/.env")
	er := godotenv.Load(currDir)

	if er != nil {
		com.PrintLog(fmt.Sprintf("(VERIFYPASS:1001) %s", er))
		return false, er
	}

	decodedPass, err := base64.StdEncoding.DecodeString(encryptedPass)

	if err != nil {
		com.PrintLog(fmt.Sprintf("(VERIFYPASS:1003) Error decoding base64: %s", err))
		return false, err
	}

	encrypted := argon2.IDKey([]byte(password), []byte(os.Getenv("salt")), 2, 64*1024, 8, 32)

	if string(decodedPass) != string(encrypted) {
		com.PrintLog("(VERIFYPASS:1002) PASSWORD IS NOT MATCH")
		return false, fmt.Errorf("PASSWORD IS NOT MATCH")
	}

	return true, nil
}

func VerifyField(encryptedToken string) (payloads map[string]interface{}, e error) {
	var payload map[string]interface{}
	key := GetStaticKey()
	path, err := os.Getwd()
	if err != nil {
		com.PrintLog(fmt.Sprintf("(VERIFY: 1000) %s", err))
		return nil, err
	}
	currDir := fmt.Sprint(path, "/.env")
	er := godotenv.Load(currDir)

	if er != nil {
		com.PrintLog(fmt.Sprintf("(VERIFY:1001) %s", er))
		return nil, er
	}

	parser := paseto.NewParser()

	token, er2 := parser.ParseV4Local(key, encryptedToken, nil)

	if er2 != nil {
		com.PrintLog(fmt.Sprintf("(VERIFY: 1003) %s", er2))
		return nil, er2
	}

	decrypted, err := token.GetString(os.Getenv("secretKey"))
	if err != nil {
		com.PrintLog(fmt.Sprintf("(VERIFY: 1004) %s\n", err))
		return nil, err
	}

	if errr := json.Unmarshal([]byte(decrypted), &payload); errr != nil {
		com.PrintLog(fmt.Sprintf("(VERIFY:1005) %s", errr))
		return nil, fmt.Errorf("error unmarshaling payload: %v", err)
	}

	return payload, nil
}

func Auth(encrypted map[string]interface{}, db *sql.DB, userId *int, usernm *string) (rc int, err error) {

	com.PrintLog(fmt.Sprintf("Token %v", encrypted["data"]))
	username := encrypted["data"].(map[string]interface{})["username"].(string)
	keepLogin, err := strconv.Atoi(encrypted["keepLogin"].(string))

	if err != nil {
		com.PrintLog(fmt.Sprintf("(Auth:0001): %s", err))
		return 501, fmt.Errorf("login failure")
	}

	if keepLogin != 1 {
		com.PrintLog(fmt.Sprintf("(Auth:0002): %s", err))
		return 409, fmt.Errorf("login failure")
	}

	com.PrintLog(fmt.Sprintf("USERNAME   :    %s", username))

	user, rcs, err := repository.KeepLogin(db, username)

	if err != nil {
		com.PrintLog(fmt.Sprintf("Auth:0003%s", err))
		return rcs, err
	}

	*userId = user
	*usernm = username

	return rcs, nil
}
