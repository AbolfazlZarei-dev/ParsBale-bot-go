package ParsBale

import (
	"encoding/json"
	"fmt"
)

type APIError struct {
	Code        int                 `json:"error_code"`
	Description string              `json:"description"`
	Parameters  *ResponseParameters `json:"parameters,omitempty"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("bale api error: [%d] %s", e.Code, e.Description)
}

type NetworkError struct {
	Err error
}

func (e NetworkError) Error() string {
	return fmt.Sprintf("network error: %v", e.Err)
}

type ParseError struct {
	Raw []byte
	Err error
}

func (e ParseError) Error() string {
	return fmt.Sprintf("parse error: %v (raw: %s)", e.Err, string(e.Raw))
}

type SafirError struct {
	Type    int    `json:"type"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e SafirError) Error() string {
	return fmt.Sprintf("safir error [%d]: %s", e.Code, e.Message)
}

func CheckAPIResponse(body []byte) error {
	var resp struct {
		Ok          bool                `json:"ok"`
		ErrorCode   int                 `json:"error_code,omitempty"`
		Description string              `json:"description,omitempty"`
		Parameters  *ResponseParameters `json:"parameters,omitempty"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return ParseError{Raw: body, Err: err}
	}
	if !resp.Ok {
		return APIError{
			Code:        resp.ErrorCode,
			Description: resp.Description,
			Parameters:  resp.Parameters,
		}
	}
	return nil
}
