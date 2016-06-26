package commons_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/vil-coyote-acme/go-concurrency/commons"
	"io"
	"net/http"
	"testing"
)

func Test_unmarshallOrder_should_unmarshal_without_error(t *testing.T) {
	// given
	expectedOrder := commons.Order{Id: 1, Quantity: 5, Type: commons.Beer, CallBackUrl: "http://callback.com/money"}
	order := new(commons.Order)
	body, _ := json.Marshal(expectedOrder)
	var req http.Request
	req.Body = nopCloser{bytes.NewBuffer(body)}
	req.ContentLength = int64(len(body))
	// when
	err := commons.UnmarshalOrderFromHttp(&req, order)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrder, *order)
}

func Test_unmarshallOrder_should_unmarshal_with_error(t *testing.T) {
	// given
	order := new(commons.Order)
	var req http.Request
	req.Body = nopCloser{bytes.NewBuffer(make([]byte, 0))}
	req.ContentLength = int64(0)
	// when
	err := commons.UnmarshalOrderFromHttp(&req, order)
	assert.NotNil(t, err)
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error {
	return nil
}
