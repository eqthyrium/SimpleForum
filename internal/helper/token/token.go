package token

import (
	"SimpleForum/internal/domain"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Token struct {
	UserId            int       `json:"userId"`
	UUID              string    `json:"uuid"`
	Role              string    `json:"role"`
	ExpireTimeActive  time.Time `json:"expireTimeActive"`
	ExpireTimeRefresh time.Time `json:"expireTimeRefresh"`
}

var mySecretKey string = "AddDeleteKey1618"

var MapUUID map[int]string = make(map[int]string)

func CreateSignedToken(userId int, role string) (string, error) {

	token := Token{
		UserId:            userId,
		UUID:              uuid.New().String(),
		Role:              role,
		ExpireTimeActive:  time.Now().Add(15 * time.Minute),
		ExpireTimeRefresh: time.Now().Add(24 * time.Hour),
	}

	MapUUID[token.UserId] = token.UUID

	tokenJson, err := json.Marshal(token)
	if err != nil {
		return "", fmt.Errorf("Token-CreateSignedToken, token marshalling failed: %v", err)
	}

	signature := createSignature(string(tokenJson), mySecretKey) // hashsum

	signatureToken := base64.URLEncoding.EncodeToString(tokenJson) + "." + signature

	return signatureToken, nil
}

func VerifyTokenHashSum(token string) (bool, error) {
	passedToken := strings.Split(token, ".")
	if len(passedToken) != 2 {
		return false, fmt.Errorf("Token-ValidateToken: %w", domain.ErrInvalidToken)
	}
	playLoad := passedToken[0]
	secondSignature := createSignature(playLoad, mySecretKey)
	if hmac.Equal([]byte(secondSignature), []byte(passedToken[1])) {
		return false, fmt.Errorf("Token-ValidateToken: %w", domain.ErrInvalidToken)
	} else {
		return true, nil
	}
}

func createSignature(tokenJson, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(tokenJson))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func ExtractDataFromToken(token string) (*Token, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Token-ExtractDataFromToken: %w")
	}
	playLoad, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("Token-ExtractDataFromToken: %w", err)
	}

	var tokenObject Token = Token{}
	err = json.Unmarshal(playLoad, &tokenObject)
	if err != nil {
		return nil, fmt.Errorf("Token-ExtractDataFromToken: %w", err)
	}

	return &tokenObject, nil
}

func VerifyTokenFields(token *Token) (bool, error) {

}
