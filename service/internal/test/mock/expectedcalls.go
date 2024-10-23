package mock

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

type ExpectedCall interface {
	CallCount() int
	PathHandler(t *testing.T) (string, http.HandlerFunc)
	Signature() string
}

type ExpectedAPICall[IN, OUT any] struct {
	Method              string
	APIPath             string
	ExpectedRequestBody *IN
	APIResponse         OUT
	callCount           int
}

func (e *ExpectedAPICall[I, _]) HandlerFunction(t *testing.T) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		e.callCount += 1
		require.Equal(t, e.Method, request.Method, "expected method %s for %s, got %s", e.Method, request.URL, request.Method)
		if e.ExpectedRequestBody == nil {
			var bytes []byte
			_, err := request.Body.Read(bytes)
			require.ErrorIs(t, err, io.EOF)
		} else {
			var actualRequestBody I
			require.NoError(t, json.NewDecoder(request.Body).Decode(&actualRequestBody))
			require.Equal(t, *e.ExpectedRequestBody, actualRequestBody)
		}
		responseBytes, err := json.Marshal(e.APIResponse)
		require.NoError(t, err)
		_, err = writer.Write(responseBytes)
		require.NoError(t, err)
	}
}

func (e *ExpectedAPICall[_, _]) CallCount() int {
	return e.callCount
}

func (e *ExpectedAPICall[_, _]) PathHandler(t *testing.T) (string, http.HandlerFunc) {
	return e.APIPath, e.HandlerFunction(t)
}

func (e *ExpectedAPICall[_, _]) Signature() string {
	return fmt.Sprintf("%s %s", e.Method, e.APIPath)
}
