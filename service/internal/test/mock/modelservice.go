package mock

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ModelService struct {
	Server        *httptest.Server
	ExpectedCalls []ExpectedCall
}

func NewModelService(t *testing.T, expectedCall ...ExpectedCall) *ModelService {
	mux := http.NewServeMux()
	for _, ph := range expectedCall {
		mux.HandleFunc(ph.PathHandler(t))
	}
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		require.Fail(t, "unexpected call to Pennsieve", "%s %s", request.Method, request.URL)
	})
	return &ModelService{
		Server:        httptest.NewServer(mux),
		ExpectedCalls: expectedCall,
	}
}

func (m *ModelService) AssertAllCalledExactlyOnce(t *testing.T) bool {
	for _, expectedCall := range m.ExpectedCalls {
		if !assert.Equal(t, 1, expectedCall.CallCount(), "%s was called %d times", expectedCall.Signature(), expectedCall.CallCount()) {
			return false
		}
	}
	return true
}

func (m *ModelService) Close() {
	m.Server.Close()
}

func (m *ModelService) URL() string {
	return m.Server.URL
}
