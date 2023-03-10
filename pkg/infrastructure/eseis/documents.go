package eseis

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func (e *EseisClient) GetDocument(uuid string) ([]byte, error) {
	path := fmt.Sprintf("/v1/sergic_documents?access_token=%s&uuid=%s", e.accessToken.accessToken, uuid)
	req, err := http.NewRequest("GET", e.buildURL(path), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create sergic_documents request: %w", err)
	}
	if err = e.setAuthentication(req); err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send sergic_documents request: %w", err)
	}
	defer resp.Body.Close()

	buffer := bytes.Buffer{}
	_, err = io.Copy(&buffer, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read sergic_documents response body: %w", err)
	}

	return buffer.Bytes(), nil
}
