package pennsieve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pennsieve/processor-post-metadata/service/util"
	"io"
	"net/http"
)

type Session struct {
	Token    string
	APIHost  string
	API2Host string
}

func NewSession(sessionToken, apiHost, api2Host string) *Session {
	return &Session{
		Token:    sessionToken,
		APIHost:  apiHost,
		API2Host: api2Host}
}

func (s *Session) newPennsieveRequest(method string, url string, structBody any) (*http.Request, error) {
	body, err := makeJSONBody(structBody)
	if err != nil {
		return nil, fmt.Errorf("error for %s %s request: %w",
			method, url, err)
	}
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating %s %s request: %w", method, url, err)
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.Token))
	return request, nil
}

func (s *Session) InvokePennsieve(method string, url string, structBody any) (*http.Response, error) {
	req, err := s.newPennsieveRequest(method, url, structBody)
	if err != nil {
		return nil, fmt.Errorf("error creating %s %s request: %w", method, url, err)
	}
	return util.Invoke(req)
}

func makeJSONBody(structBody any) (io.Reader, error) {
	if structBody == nil {
		return nil, nil
	}
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(structBody); err != nil {
		return nil, fmt.Errorf("error encoding body: %w", err)
	}
	return &buffer, nil
}
