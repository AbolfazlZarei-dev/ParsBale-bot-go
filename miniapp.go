package ParsBale

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type WebAppInitData struct {
	QueryID  string      `json:"query_id"`
	User     *WebAppUser `json:"user"`
	AuthDate int64       `json:"auth_date"`
	Hash     string      `json:"hash"`
}

type WebAppUser struct {
	ID              int64  `json:"id"`
	FirstName       string `json:"first_name"`
	Username        string `json:"username,omitempty"`
	AllowsWriteToPM bool   `json:"allows_write_to_pm,omitempty"`
}

func ValidateWebAppData(initData string, botToken string) (*WebAppInitData, error) {
	vals, err := url.ParseQuery(initData)
	if err != nil {
		return nil, err
	}

	hash := vals.Get("hash")
	if hash == "" {
		return nil, fmt.Errorf("hash is missing")
	}
	vals.Del("hash")

	var keys []string
	for k := range vals {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var dataCheckArr []string
	for _, k := range keys {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%s", k, vals.Get(k)))
	}
	dataCheckString := strings.Join(dataCheckArr, "\n")

	secretKeyHmac := hmac.New(sha256.New, []byte("WebAppData"))
	secretKeyHmac.Write([]byte(botToken))
	secretKey := secretKeyHmac.Sum(nil)

	dataHmac := hmac.New(sha256.New, secretKey)
	dataHmac.Write([]byte(dataCheckString))
	resultHash := hex.EncodeToString(dataHmac.Sum(nil))

	if resultHash != hash {
		return nil, fmt.Errorf("invalid hash")
	}

	var waid WebAppInitData
	if u := vals.Get("user"); u != "" {
		if err := json.Unmarshal([]byte(u), &waid.User); err != nil {
			return nil, err
		}
	}
	waid.QueryID = vals.Get("query_id")
	fmt.Sscanf(vals.Get("auth_date"), "%d", &waid.AuthDate)
	waid.Hash = hash

	return &waid, nil
}
