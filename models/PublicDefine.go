package models

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	// "fmt"
	"math/big"
	"net/url"
	"strings"
	// "sync"

	"github.com/dgryski/dgoogauth"
	"github.com/ebitgo/ebitgo.com/models/logger"
)

const (
	POST_TYPE_FLAG             = "PostType"
	PT_CHECK_USERNAME          = "C0001"
	PT_CHECK_USERAUTH          = "C0002"
	PT_CHECK_GA_OPERA          = "C0003"
	PT_GET_WALLETS             = "C0004"
	PT_SEARCH_WALLETS          = "C0005"
	PT_GET_FRIENDS             = "C0006"
	PT_USER_REGISTRATION       = "R0001"
	PT_USER_MODIFY_PW          = "R0002"
	RT_USER_ADD_WALLET_ACCOUNT = "R0003"
	RT_UPDATE_WALLET_INFO      = "R0004"
	RT_USER_DEL_WALLET_ACCOUNT = "R0005"
	RT_USER_ADD_FRIEND         = "R0006"
	RT_USER_DEL_FRIEND         = "R0007"
	RT_UPDATE_WALLET_SSKEY     = "R0008"
	PT_USER_LOGIN              = "L0001"
	PT_USER_ACTIVE             = "L0002"

	POST_MARK_CHECK_USERNAME  = "checkUName"
	POST_MARK_REG_USERNAME    = "regUN"
	POST_MARK_REG_PASSWORD    = "regPW"
	POST_MARK_REG_VICODE      = "regViCode"
	POST_MARK_REG_VALUE1      = "regV1"
	POST_MARK_REG_VALUE2      = "regV2"
	POST_MARK_REG_VALUE3      = "regV3"
	POST_MARK_REG_OPERA       = "regOper"
	POST_MARK_USER_NAME       = "un"
	POST_MARK_PASS_WORD       = "pw"
	POST_MARK_NEW_PASS_WORD   = "npw"
	POST_MARK_VICODE          = "vic"
	POST_MARK_GACODE          = "gac"
	POST_MARK_GAKEY           = "gakey"
	POST_MARK_AUTHCODE        = "auth"
	POST_MARK_HEART           = "hbit"
	POST_MARK_FID             = "fid"
	POST_MARK_TID             = "tid"
	POST_MARK_MODIFY_GA       = "mgatype"
	POST_MARK_MODIFY_NICKNAME = "mnkname"
	POST_MARK_WALLET_NICKNAME = "nkname"
	POST_MARK_WALLET_ID       = "wid"
	POST_MARK_WALLET_PUBADDR  = "paddr"
	POST_MARK_WALLET_SKEY     = "skey"

	GA_PT_MODIFY_NEW    = "new"
	GA_PT_MODIFY_GET    = "get"
	GA_PT_MODIFY_DELETE = "delete"

	RT_M_SUCCESS     = "success"
	RT_M_USEREMAIL   = "login_user"
	RT_M_EXIST       = "IsExist"
	RT_M_USER_LVL    = "user_level"
	RT_M_USER_AUTH   = "user_auth"
	RT_M_AUTH_UPDATE = "update_auth"
	RT_M_GA_STATE    = "user_ga_used"
	RT_M_GA_KEY      = "user_ga_key"
	RT_M_GA_URI      = "user_ga_uri"
	RT_M_ERROR       = "Error"
	RT_M_RESULTS     = "results"

	RT_M_WALLET_NICKNAME   = "nickname"
	RT_M_WALLET_PUBLICADDR = "public_address"
	RT_M_WALLET_ID         = "wid"
	RT_M_WALLET_SECRETKEY  = "secret_key"
)

type IOperationInterface interface {
	GetResultData() *ReslutOperation
	QuaryExecute() error
	DecodeContext(vals url.Values)
}

func fromBase10(base10 string) *big.Int {
	i, ok := new(big.Int).SetString(base10, 10)
	if !ok {
		logmodels.Logger.Error(logmodels.SPrintError("PublicDefine : fromBase10", "panic bad number: %s", base10))
		panic("bad number: " + base10)
	}
	return i
}

func getRsaPrivateKey(ntime, uname, pw, auth string) *rsa.PrivateKey {
	tmp := ntime + ";" + uname + ";" + pw + ";" + auth

	prv, _ := rsa.GenerateKey(strings.NewReader(tmp), 64)
	return prv
}

func getSHA256Key(ntime, uname, pw, auth string) string {
	tmp := ntime + ";" + uname + ";" + pw + ";" + auth
	hash := sha256.New()
	hash.Write([]byte(tmp))
	ret := strings.Replace(base32.StdEncoding.EncodeToString(hash.Sum(nil)), "=", "", -1)
	// ret := hex.EncodeToString(hash.Sum(nil))
	return ret
}

func randomByte(length int) []byte {
	//rand Read
	k := make([]byte, length)
	if _, err := rand.Read(k); err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("PublicDefine : randomByte", "rand.Read() error : %v", err))
	}
	return k
}

func randomString(length int) string {
	//rand Read
	k := randomByte(length)
	return hex.EncodeToString(k)
	// ret := ""
	// for i := 0; i < length; i++ {
	// 	ret += fmt.Sprintf("%02x", k[i])
	// }
	// return ret
}

func randomBase32String() string {
	rd := randomByte(10)
	ret := base32.StdEncoding.EncodeToString(rd)
	return ret
}

func verifyGoogleAuth(key, code string) bool {
	var totp dgoogauth.OTPConfig
	totp.Secret = key
	totp.WindowSize = 2
	b, err := totp.Authenticate(code)
	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("PublicDefine : verifyGoogleAuth", "%v", err))
	}
	return b
}
