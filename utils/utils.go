// package utils contain useful helper functions
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrSubjectNotFound = errors.New("subject does not exist")
	ErrNoPermission    = errors.New("user has not done subject")
)

func checkSubjectExists(DB db.Env, subHash string) error {
	snap, err := DB.Restore("subjects", subHash)
	if snap == nil || !snap.Exists() {
		return ErrSubjectNotFound
	}
	return err
}

func checkSubjectRecords(DB db.Env, userHash, subHash string) error {
	col, err := DB.RestoreCollection("users/" + userHash + "/final_scores/" + subHash + "/records")
	if len(col) == 0 {
		return ErrNoPermission
	}
	return err
}

func CheckSubjectPermission(DB db.Env, userHash, subHash string) error {
	errSub, errRec := checkSubjectExists(DB, subHash), checkSubjectRecords(DB, userHash, subHash)
	if errSub != nil {
		return errSub
	} else if errRec != nil {
		return errRec
	}

	return nil
}

// MakeRequest runs a single request. This is used by test functions that run requests on the router
func MakeRequest(router *gin.Engine, method, endpoint string, body io.Reader, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, endpoint, body)

	for _, c := range cookies {
		req.AddCookie(c)
	}

	router.ServeHTTP(w, req)
	return w
}

// Encrypts using AES. keyString must be 128, 196 or 256 bits.
func AESEncrypt(stringToEncrypt string, keyString string) (encryptedString string, err error) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

// Encrypts using AES. keyString must be 128, 196 or 256 bits.
func AESDecrypt(encryptedString string, keyString string) (decryptedString string, err error) {
	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext), nil
}

func SHA256(body string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(body)))
}

func Bcrypt(body string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(body), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(pass), nil
}

func BcryptCompare(input, truth string) bool {
	return bcrypt.CompareHashAndPassword([]byte(truth), []byte(input)) == nil
}

func Min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
