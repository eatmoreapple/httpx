package httpx_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eatmoreapple/httpx"
)

func TestRequestBuilder_Do(t *testing.T) {
	// Step 1: Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, test!"))
	}))
	defer server.Close()

	// Step 2: Use the RequestBuilder to create a request
	builder := httpx.New(server.URL)
	builder.Get()

	// Step 3: Send the request and get the response
	resp, err := builder.Do()
	require.NoError(t, err)

	// Step 4: Verify the response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "Hello, test!", string(body))
}

func TestRequestBuilder_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Post test!"))
	}))
	defer server.Close()

	builder := httpx.New(server.URL)
	builder.Post()

	resp, err := builder.Do()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "Post test!", string(body))
}

func TestRequestBuilder_SetHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("SetHeader test!"))
	}))
	defer server.Close()

	builder := httpx.New(server.URL)
	builder.SetHeader("Content-Type", "application/json")

	resp, err := builder.Do()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "SetHeader test!", string(body))
}

func TestRequestBuilder_Body(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Equal(t, "Body test!", string(body))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Body test!"))
	}))
	defer server.Close()

	builder := httpx.New(server.URL)
	builder.Body(io.NopCloser(strings.NewReader("Body test!")))

	resp, err := builder.Do()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "Body test!", string(body))
}

func TestRequestBuilder_PostForm(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Equal(t, "foo=bar", string(body))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PostForm test!"))
	}))
	defer server.Close()

	builder := httpx.New(server.URL)
	builder.PostForm(map[string][]string{"foo": {"bar"}})

	resp, err := builder.Do()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "PostForm test!", string(body))
}

func TestRequestBuilder_Json(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Equal(t, `{"foo":"bar"}`, string(body))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Json test!"))
	}))
	defer server.Close()

	builder := httpx.New(server.URL)
	builder.Json(map[string]string{"foo": "bar"})

	resp, err := builder.Do()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "Json test!", string(body))
}
