package ParsBale

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

const (
	APIEndpoint   = "https://tapi.bale.ai/bot%s/%s"
	FileEndpoint  = "https://tapi.bale.ai/file/bot%s/%s"
	SafirEndpoint = "https://safir.bale.ai/api/v3"
)

type FileUpload struct {
	Name    string
	Content io.Reader
}

type HTTPClient struct {
	Client  *http.Client
	Token   string
	BaseURL string
}

func NewHTTPClient(token string) *HTTPClient {
	return &HTTPClient{
		Token:   token,
		BaseURL: APIEndpoint,
		Client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (h *HTTPClient) Request(method string, params map[string]string, body interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf(h.BaseURL, h.Token, method)

	var reqBody io.Reader
	contentType := "application/x-www-form-urlencoded"

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
		contentType = "application/json"
	} else if len(params) > 0 {
		form := make(url.Values)
		for k, v := range params {
			form.Set(k, v)
		}
		reqBody = bytes.NewBufferString(form.Encode())
	}

	req, err := http.NewRequest("POST", endpoint, reqBody)
	if err != nil {
		return nil, NetworkError{Err: err}
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, NetworkError{Err: err}
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NetworkError{Err: err}
	}

	if err := CheckAPIResponse(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (h *HTTPClient) RequestMultipart(method string, params map[string]string, files map[string]FileUpload) ([]byte, error) {
	endpoint := fmt.Sprintf(h.BaseURL, h.Token, method)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	for field, file := range files {
		part, err := writer.CreateFormFile(field, file.Name)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, file.Content)
		if err != nil {
			return nil, err
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, NetworkError{Err: err}
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, NetworkError{Err: err}
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NetworkError{Err: err}
	}

	if err := CheckAPIResponse(data); err != nil {
		return nil, err
	}

	return data, nil
}
