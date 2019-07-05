package auth

// RSAKey is the encoding params for the auth
import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"strconv"

	"github.com/vincentserpoul/mangosteam"
)

// RSAKey represent the key for encoding password
type RSAKey struct {
	PublicKeyModulus  string `json:"publickey_mod"`
	PublicKeyExponent string `json:"publickey_exp"`
	Timestamp         string `json:"timestamp"`
}

// GetRSAKeyURI is the URI used to get the RSA Key in prod
const GetRSAKeyURI = "/login/getrsakey"

// GetRSAKey queries steam to get the key needed to encode the password
func GetRSAKey(username string) (*RSAKey, error) {
	// Est-ce testable voir solution login_test.go
	resp, err := http.PostForm(
		mangosteam.BaseSteamWebURL+GetRSAKeyURI,
		url.Values{"username": {username}},
	)
	if err != nil {
		return nil, fmt.Errorf("mangosteam GetRSAKey(%s): %v", username, err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	key, err := extractRSAKeyFromJSON(body)
	if err != nil {
		return nil, fmt.Errorf("mangosteam GetRSAKey(%s): %v", username, err)
	}

	return key, nil
}

func extractRSAKeyFromJSON(JSONBytes []byte) (*RSAKey, error) {
	key := RSAKey{}
	if err := json.Unmarshal(JSONBytes, &key); err != nil {
		return nil, err
	}
	if key.PublicKeyExponent == "" ||
		key.PublicKeyModulus == "" ||
		key.Timestamp == "" {
		return nil, fmt.Errorf(
			"websteam extractRSAKeyFromJSON: incomplete RSAKey unmarshalled: %+v", key)
	}

	return &key, nil
}

// EncryptPassword will RSA encode the password with the key provided
func EncryptPassword(password string, rsaKey *RSAKey) (string, error) {

	if password == "" {
		return "", fmt.Errorf("websteam EncryptPassword(password, rsaKey) " +
			"requires a non empty password as an argument")
	}
	if rsaKey == nil {
		return "", fmt.Errorf("websteam EncryptPassword(password, rsaKey) " +
			"requires a non empty rsaKey as an argument")
	}

	// convert the hex string to int
	pubkeyExpInt, err := strconv.ParseInt(rsaKey.PublicKeyExponent, 16, 64)

	if err != nil {
		return "", fmt.Errorf("mangosteam EncryptPassword(): %v", err)
	}

	// convert the hex string to a big int (we can't use ParseInt)
	pubkeyModInt := new(big.Int)
	pubkeyModInt.SetString(rsaKey.PublicKeyModulus, 16)

	var realKey rsa.PublicKey

	realKey.E = int(pubkeyExpInt)
	realKey.N = pubkeyModInt

	encryptedPass, err := rsa.EncryptPKCS1v15(rand.Reader, &realKey, []byte(password))

	if err != nil {
		return "", fmt.Errorf("mangosteam EncryptPassword(): %v", err)
	}

	return base64.StdEncoding.EncodeToString(encryptedPass), nil
}
