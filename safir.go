package ParsBale

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// V3 Client (API Key based)
type SafirClient struct {
	AccessKey string
	Client    *http.Client
}

func NewSafirClient(accessKey string) *SafirClient {
	return &SafirClient{
		AccessKey: accessKey,
		Client:    http.DefaultClient,
	}
}

func (s *SafirClient) SafirRequest(endpoint string, body interface{}) ([]byte, error) {
	apiUrl := fmt.Sprintf("%s/%s", SafirEndpoint, endpoint)
	jsonData, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	req.Header.Set("api-access-key", s.AccessKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, NetworkError{Err: err}
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	// بررسی وضعیت HTTP
	if resp.StatusCode >= 400 {
		var apiErr SafirError
		if json.Unmarshal(data, &apiErr) == nil && apiErr.Message != "" {
			return nil, apiErr
		}
		return nil, fmt.Errorf("safir request failed with status %d: %s", resp.StatusCode, string(data))
	}

	return data, nil
}

// --- Models V3 ---

type SafirMessageData struct {
	Message    *SafirMessage `json:"message,omitempty"`
	OTPMessage *SafirOTP     `json:"otp_message,omitempty"`
	IsSecure   bool          `json:"is_secure,omitempty"`
}

type SafirMessage struct {
	Text     string `json:"text,omitempty"`
	FileID   string `json:"file_id,omitempty"`
	CopyText string `json:"copy_text,omitempty"`
}

type SafirOTP struct {
	OTP string `json:"otp"`
}

type SafirSendRequest struct {
	RequestID   string           `json:"request_id,omitempty"`
	BotID       int64            `json:"bot_id"`
	PhoneNumber string           `json:"phone_number"`
	MessageData SafirMessageData `json:"message_data"`
}

type SafirSendResponse struct {
	MessageID string           `json:"message_id"`
	ErrorData []SafirErrorInfo `json:"error_data"`
}

type SafirErrorInfo struct {
	PhoneNumber string `json:"phone_number"`
	Code        int    `json:"code"`
	Description string `json:"description"`
}

// --- Methods V3 ---

func (s *SafirClient) SendMessage(botID int64, phone string, data SafirMessageData) (*SafirSendResponse, error) {
	req := SafirSendRequest{
		BotID:       botID,
		PhoneNumber: phone,
		MessageData: data,
	}

	raw, err := s.SafirRequest("send_message", req)
	if err != nil {
		return nil, err
	}
	var resp SafirSendResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, ParseError{Raw: raw, Err: err}
	}
	return &resp, nil
}

func (s *SafirClient) UploadFile(filePath string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	part, _ := writer.CreateFormFile("file", filePath)
	io.Copy(part, file)
	writer.Close()

	req, _ := http.NewRequest("POST", SafirEndpoint+"/upload_file", body)
	req.Header.Set("api-access-key", s.AccessKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	resBody, _ := io.ReadAll(resp.Body)

	// اصلاح شده: بررسی وضعیت HTTP
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("upload failed (%d): %s", resp.StatusCode, string(resBody))
	}

	var res struct {
		FileID string `json:"file_id"`
	}
	if err := json.Unmarshal(resBody, &res); err != nil {
		return "", err
	}
	return res.FileID, nil
}

// V2 Client (OAuth Token based for OTP Service)
type SafirOTPClient struct {
	ClientID     string
	ClientSecret string
	Client       *http.Client
	Token        string
}

func NewSafirOTPClient(clientID, clientSecret string) *SafirOTPClient {
	return &SafirOTPClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Client:       http.DefaultClient,
	}
}

func (s *SafirOTPClient) GetToken() error {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", s.ClientID)
	data.Set("client_secret", s.ClientSecret)
	data.Set("scope", "read")

	req, _ := http.NewRequest("POST", "https://safir.bale.ai/api/v2/auth/token", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.Client.Do(req)
	if err != nil {
		return NetworkError{Err: err}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var res struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return ParseError{Raw: body, Err: err}
	}
	s.Token = res.AccessToken
	return nil
}

func (s *SafirOTPClient) SendOTP(phone string, otp string) error {
	if s.Token == "" {
		if err := s.GetToken(); err != nil {
			return err
		}
	}

	body := map[string]interface{}{
		"phone": phone,
		"otp":   otp,
	}
	jsonData, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "https://safir.bale.ai/api/v2/send_otp", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return NetworkError{Err: err}
	}
	defer resp.Body.Close()

	resBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("otp error: %s", string(resBody))
	}
	return nil
}
