package httpx_test

import (
	"context"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eatmoreapple/httpx"
)

func TestRequestBuilder_Method(t *testing.T) {
	builder := httpx.New("http://localhost")
	builder.Method(http.MethodPost)

	req, err := builder.Build()
	require.NoError(t, err)
	assert.Equal(t, http.MethodPost, req.Method)
}

func TestRequestBuilder_Get(t *testing.T) {
	builder := httpx.New("http://localhost")
	builder.Get()

	req, err := builder.Build()
	require.NoError(t, err)
	assert.Equal(t, http.MethodGet, req.Method)
}

func TestRequestBuilder_Post(t *testing.T) {
	builder := httpx.New("http://localhost")
	builder.Post()

	req, err := builder.Build()
	require.NoError(t, err)
	assert.Equal(t, http.MethodPost, req.Method)
}

func TestRequestBuilder_SetHeader(t *testing.T) {
	builder := httpx.New("http://localhost")
	builder.SetHeader("Content-Type", "application/json")

	req, err := builder.Build()
	require.NoError(t, err)
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
}

func TestRequestBuilder_Form(t *testing.T) {
	builder := httpx.New("http://localhost")
	builder.Form(url.Values{"key": []string{"value"}})

	req, err := builder.Build()
	require.NoError(t, err)
	assert.Equal(t, "value", req.Form.Get("key"))
}

func TestRequestBuilder_Query(t *testing.T) {
	builder := httpx.New("http://localhost")
	builder.Query(map[string]string{"key": "value"})

	req, err := builder.Build()
	require.NoError(t, err)
	assert.Equal(t, "value", req.URL.Query().Get("key"))
}

func TestRequestBuilder_Json(t *testing.T) {
	builder := httpx.New("http://localhost")
	builder.Json(map[string]string{"key": "value"})

	req, err := builder.Build()
	require.NoError(t, err)
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(req.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"key":"value"}`, string(body))
}

func TestRequestBuilder_PostForm(t *testing.T) {
	builder := httpx.New("http://localhost")
	builder.PostForm(url.Values{"key": []string{"value"}})

	req, err := builder.Build()
	require.NoError(t, err)
	assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(req.Body)
	require.NoError(t, err)
	assert.Equal(t, "key=value", string(body))
}

func TestRequestBuilder_MultipartForm(t *testing.T) {
	builder := httpx.New("http://localhost")
	form := &multipart.Form{
		Value: map[string][]string{
			"key": {"value"},
		},
	}
	builder.MultipartForm(form)

	req, err := builder.Build()
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(req.Header.Get("Content-Type"), "multipart/form-data"))

	body, err := ioutil.ReadAll(req.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "key")
	assert.Contains(t, string(body), "value")
}

func TestRequestBuilder_BuildWithContext(t *testing.T) {
	builder := httpx.New("http://localhost")
	ctx := context.WithValue(context.Background(), "key", "value")
	req, err := builder.BuildWithContext(ctx)
	require.NoError(t, err)
	assert.Equal(t, "value", req.Context().Value("key"))
}
