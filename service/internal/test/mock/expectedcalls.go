package mock

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"testing"
)

type ExpectedCall interface {
	CallCounts() []int
	AllCalledExactlyOnce() bool
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
			var requestBodyBytes []byte
			_, err := request.Body.Read(requestBodyBytes)
			require.ErrorIs(t, err, io.EOF)
		} else {
			var actualRequestBody I
			require.NoError(t, json.NewDecoder(request.Body).Decode(&actualRequestBody))
			require.Equal(t, *e.ExpectedRequestBody, actualRequestBody)
		}
		responseBytes, err := json.Marshal(e.APIResponse)
		require.NoError(t, err)
		// can't see if e.APIResponse is nil because of generics, so
		// see if the marshal result was null instead
		if string(responseBytes) != "null" && string(responseBytes) != `""` {
			_, err = writer.Write(responseBytes)
			require.NoError(t, err)
		}
	}
}

func (e *ExpectedAPICall[_, _]) CallCounts() []int {
	return []int{e.callCount}
}

func (e *ExpectedAPICall[_, _]) AllCalledExactlyOnce() bool {
	return e.callCount == 1
}

func (e *ExpectedAPICall[_, _]) PathHandler(t *testing.T) (string, http.HandlerFunc) {
	return e.APIPath, e.HandlerFunction(t)
}

func (e *ExpectedAPICall[_, _]) Signature() string {
	return fmt.Sprintf("%s %s", e.Method, e.APIPath)
}

// ExpectedAPICallMulti is for cases where you expect multiple calls to the same APIPath with different bodies.
// Needed because you cannot register two handlers for the same path
// Bulk Proxy instance delete for example
type ExpectedAPICallMulti[IN, OUT any] struct {
	APIPath string
	Calls   []ExpectedAPICallData[IN, OUT]
}

type ExpectedAPICallData[IN, OUT any] struct {
	Method              string
	ExpectedRequestBody *IN
	APIResponse         OUT
	callCount           int
}

func (e *ExpectedAPICallMulti[I, O]) HandlerFunction(t *testing.T) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Read the actual body once and compare in slices.IndexFunc to see if
		// there is matching call in e
		var actualRequestBody I
		// saving the actual body bytes in case we need them as a string for an error message later onDo
		var actualRequestBodyBytes bytes.Buffer
		tee := io.TeeReader(request.Body, &actualRequestBodyBytes)
		requestBodyDecodeError := json.NewDecoder(tee).Decode(&actualRequestBody)
		if requestBodyDecodeError != nil && !errors.Is(requestBodyDecodeError, io.EOF) {
			// If there is an error other than EOF fail test now.
			// An EOF error is expected exactly when one of the expected calls has no request body
			// We will check for that in slices.IndexFunc
			require.NoError(t, requestBodyDecodeError)
		}

		callIndex := slices.IndexFunc(e.Calls, func(e ExpectedAPICallData[I, O]) bool {
			if e.Method != request.Method {
				return false
			}
			if e.ExpectedRequestBody == nil {
				return errors.Is(requestBodyDecodeError, io.EOF)
			}
			return reflect.DeepEqual(*e.ExpectedRequestBody, actualRequestBody)
		})
		require.GreaterOrEqual(t, callIndex, 0, "unexpected call to %s: method: %s, body %s", e.APIPath, request.Method, actualRequestBodyBytes.String())
		call := &e.Calls[callIndex]
		call.callCount += 1
		responseBytes, err := json.Marshal(call.APIResponse)
		require.NoError(t, err)
		// can't see if e.APIResponse is nil because of generics, so
		// see if the marshal result was null instead
		if string(responseBytes) != "null" && string(responseBytes) != `""` {
			_, err = writer.Write(responseBytes)
			require.NoError(t, err)
		}
	}
}

func (e *ExpectedAPICallMulti[_, _]) CallCounts() []int {
	var counts []int
	for _, call := range e.Calls {
		counts = append(counts, call.callCount)
	}
	return counts
}

func (e *ExpectedAPICallMulti[_, _]) AllCalledExactlyOnce() bool {
	for _, call := range e.Calls {
		if call.callCount != 1 {
			return false
		}
	}
	return true
}

func (e *ExpectedAPICallMulti[_, _]) PathHandler(t *testing.T) (string, http.HandlerFunc) {
	return e.APIPath, e.HandlerFunction(t)
}

func (e *ExpectedAPICallMulti[_, _]) Signature() string {
	var pieces []string
	for _, call := range e.Calls {
		pieces = append(pieces, call.Method)
	}
	pieces = append(pieces, e.APIPath)
	return strings.Join(pieces, " ")
}
