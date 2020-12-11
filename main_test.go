package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/iotex-antenna-go/v2/jwt"
	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
	"github.com/stretchr/testify/require"

	"github.com/iotexproject/phoenix-gem/config"
	"github.com/iotexproject/phoenix-gem/handler"
	"github.com/iotexproject/phoenix-gem/log"
	"github.com/iotexproject/phoenix-gem/server"
)

const (
	defaultContentType = "application/json"
)

func fakeS3Server() *httptest.Server {
	// fake s3
	backend := s3mem.New()
	faker := gofakes3.New(backend)
	return httptest.NewServer(faker.Server())
}

func Test_ServerWithS3Storage(t *testing.T) {
	var jwtToken string
	var urlPath string
	var err error

	r := require.New(t)
	testFile, err := ioutil.TempFile(os.TempDir(), "test-config")
	path := testFile.Name()
	r.NoError(err)
	testFile.Close()
	defer func() {
		r.NoError(os.Remove(path))
	}()

	s3Server := fakeS3Server()
	defer s3Server.Close()

	cfg, err := config.New("config.yaml")
	r.NoError(err)
	cfg.Server.DBPath = path
	r.NoError(log.InitLoggers(cfg.Log, cfg.SubLogs))

	ts := server.New(cfg)
	go func() {
		ts.Start()
	}()
	defer func() {
		ts.Shutdown(context.Background())
	}()
	time.Sleep(time.Second)
	Addr := "http://127.0.0.1" + ts.Addr
	registerData := bytes.NewReader([]byte(`{ "name": "s3", "region":"www", "endpoint":"xxx", "key":"yyy", "token":"zzz"}`))

	t.Run("with no authorized", func(t *testing.T) {
		//register
		urlPath = "/register"
		res, body, err := testRequest("POST", Addr+urlPath, "", "", registerData)
		t.Logf("status=> %v body => %s", res.StatusCode, body)
		r.NoError(err)
		r.Equal(res.StatusCode, http.StatusUnauthorized)
	})
	t.Run("with error authorized", func(t *testing.T) {
		//register
		urlPath = "/register"
		jwtToken = "xxxxxx"
		res, body, err := testRequest("POST", Addr+urlPath, "", jwtToken, registerData)
		t.Logf("status=> %v body => %s", res.StatusCode, body)
		r.NoError(err)
		r.Equal(res.StatusCode, http.StatusUnauthorized)
	})
	t.Run("with authorized", func(t *testing.T) {
		key, _ := crypto.HexStringToPrivateKey("bc145bb9f00d55a3571e22660ef5fd1bfa596e272b80add2919735b82c273004")
		issue := time.Now().Unix()
		expire := time.Now().Add(time.Hour * 240).Unix()
		subject := "s3"
		jwtToken, err = jwt.SignJWT(issue, expire, subject, jwt.CREATE, key)
		r.NoError(err)

		//register
		urlPath = "/register"
		registerData = bytes.NewReader([]byte(`{ "name": "s3", "region":"www", "endpoint":"` + s3Server.URL + `", "key":"yyy", "token":"zzz"}`))
		res, body, err := testRequest("POST", Addr+urlPath, "", jwtToken, registerData)
		t.Logf("status=> %v body => %s", res.StatusCode, body)
		r.NoError(err)
		r.Equal(res.StatusCode, http.StatusOK)
		r.Contains(body, "successful")

		//createBucket
		urlPath = "/pods"
		res, body, err = testRequest("POST", Addr+urlPath, "", jwtToken, bytes.NewReader([]byte(`{ "name": "test10"}`)))
		t.Logf("status=> %v body => %s", res.StatusCode, body)
		r.NoError(err)
		r.Equal(res.StatusCode, http.StatusOK)
		r.Contains(body, "successful")

		//deleteBucket
		urlPath = "/pods/test10"
		res, body, err = testRequest("DELETE", Addr+urlPath, "", jwtToken, nil)
		t.Logf("status=> %v body => %s", res.StatusCode, body)
		r.NoError(err)
		r.Equal(res.StatusCode, http.StatusForbidden)
		r.Contains(body, handler.ErrorPermissionDenied.Error())

		delToken, err := jwt.SignJWT(issue, expire, subject, jwt.DELETE, key)
		r.NoError(err)
		res, body, err = testRequest("DELETE", Addr+urlPath, "", delToken, nil)
		t.Logf("status=> %v body => %s", res.StatusCode, body)
		r.NoError(err)
		r.Equal(res.StatusCode, http.StatusOK)
		r.Contains(body, "successful")

		//after register bucket
		urlPath = "/pods"
		r.NoError(err)
		res, body, err = testRequest("POST", Addr+urlPath, "", jwtToken, bytes.NewReader([]byte(`{ "name": "test"}`)))
		t.Logf("status=> %v body => %s", res.StatusCode, body)
		r.NoError(err)
		r.Equal(res.StatusCode, http.StatusOK)
		r.Contains(body, "successful")

		//createObject
		urlPath = "/pea/test/foobar.txt"
		jwtToken, err = jwt.SignJWT(issue, expire, subject, jwt.UPDATE, key)
		r.NoError(err)
		res, body, err = testRequest("POST", Addr+urlPath, "", jwtToken, bytes.NewReader([]byte(`foobar`)))
		t.Logf("status=> %v body => %s", res.StatusCode, body)
		r.NoError(err)
		r.Equal(res.StatusCode, http.StatusOK)
		r.Contains(body, "successful")

		//getObject
		urlPath = "/pea/test/foobar.txt"
		jwtToken, err = jwt.SignJWT(issue, expire, subject, jwt.READ, key)
		r.NoError(err)
		res, body, err = testRequest("GET", Addr+urlPath, "", jwtToken, nil)
		t.Logf("status=> %v body => %s", res.StatusCode, body)
		r.NoError(err)
		r.Equal(res.StatusCode, http.StatusOK)
		r.Contains(body, "foobar")

		//getObjects
		urlPath = "/pea/test"
		jwtToken, err = jwt.SignJWT(issue, expire, subject, jwt.READ, key)
		r.NoError(err)
		res, body, err = testRequest("GET", Addr+urlPath, "", jwtToken, nil)
		t.Logf("status=> %v body => %s", res.StatusCode, body)
		r.NoError(err)
		r.Equal(res.StatusCode, http.StatusOK)
		r.Contains(body, "foobar.txt")

		//deleteObjects
		urlPath = "/pea/test/foobar.txt"
		jwtToken, err = jwt.SignJWT(issue, expire, subject, jwt.DELETE, key)
		r.NoError(err)
		res, body, err = testRequest("DELETE", Addr+urlPath, "", jwtToken, nil)
		t.Logf("status=> %v body => %s", res.StatusCode, body)
		r.NoError(err)
		r.Equal(res.StatusCode, http.StatusOK)
		r.Contains(body, "foobar.txt")

		//unregister
		urlPath = "/register/s3"
		jwtToken, err = jwt.SignJWT(issue, expire, subject, jwt.DELETE, key)
		r.NoError(err)
		res, body, err = testRequest("DELETE", Addr+urlPath, "", jwtToken, nil)
		t.Logf("status=> %v body => %s", res.StatusCode, body)
		r.NoError(err)
		r.Equal(res.StatusCode, http.StatusOK)
		r.Contains(body, "successful")

		// createBucket fail test, after unregister
		urlPath = "/pods"
		res, body, err = testRequest("POST", Addr+urlPath, "", jwtToken, bytes.NewReader([]byte(`{ "name": "test10"}`)))
		t.Logf("status=> %v body => %s", res.StatusCode, body)
		r.NoError(err)
		r.Equal(res.StatusCode, http.StatusNoContent)
	})
}

func testRequest(
	method, path string,
	contentType string,
	jwtToken string,
	body io.Reader) (*http.Response, string, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, "", err
	}
	if jwtToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
	}
	if contentType == "" {
		contentType = defaultContentType
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	return resp, string(respBody), nil
}
