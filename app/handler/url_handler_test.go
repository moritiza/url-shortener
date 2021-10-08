package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/moritiza/url-shortener/app/dto"
	"github.com/moritiza/url-shortener/app/helper"
	"github.com/moritiza/url-shortener/app/repository"
	"github.com/moritiza/url-shortener/app/service"
	"github.com/moritiza/url-shortener/config"
	"github.com/moritiza/url-shortener/test"
)

var (
	shortUrl string

	testUrlRepository repository.UrlRepository
	testUrlService    service.UrlService
	testUrlHandler    UrlHandler
)

func TestCreateShortUrl(t *testing.T) {
	var (
		jsonString = []byte(`{"original_url":"http://google.com"}`)
		response   helper.Response
		url        dto.CreateShortUrl
	)

	req, err := http.NewRequest(http.MethodPost, "/api/create-url", bytes.NewBuffer(jsonString))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	cfg := test.Prepare()
	defer config.DisconnectDatabase(cfg)

	testUrlRepository = repository.NewUrlRepository(cfg.Database)
	testUrlService = service.NewUrlService(*cfg.Logger, testUrlRepository)
	testUrlHandler = NewUrlHandler(*cfg.Logger, *cfg.Validator, testUrlService)

	testUrlHandler.CreateShortUrl(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status created; got %v", res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	}

	err = json.Unmarshal(b, &response)
	if err != nil {
		t.Fatalf("Response is invalid. Error: %v", err)
	}

	if response.Data == nil {
		t.Fatalf("Response data is empty")
	}

	url.OriginalUrl = response.Data.(map[string]interface{})["original_url"].(string)
	url.ShortUrl = response.Data.(map[string]interface{})["short_url"].(string)

	if url.OriginalUrl != "http://google.com" || url.ShortUrl == "" {
		t.Fatalf("Unexpected Response")
	}

	shortUrl = url.ShortUrl
}

func TestRedirect(t *testing.T) {
	if shortUrl == "" {
		TestCreateShortUrl(t)
	}

	shortUrlParts := strings.Split(shortUrl, "/")
	urlName := string(shortUrlParts[len(shortUrlParts)-1])

	req, err := http.NewRequest(http.MethodGet, "/"+urlName, nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	// Fake gorilla/mux vars
	vars := map[string]string{
		"url_name": urlName,
	}

	req = mux.SetURLVars(req, vars)
	rec := httptest.NewRecorder()
	cfg := test.Prepare()
	defer config.DisconnectDatabase(cfg)

	testUrlRepository = repository.NewUrlRepository(cfg.Database)
	testUrlService = service.NewUrlService(*cfg.Logger, testUrlRepository)
	testUrlHandler = NewUrlHandler(*cfg.Logger, *cfg.Validator, testUrlService)

	testUrlHandler.Redirect(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusSeeOther {
		t.Fatalf("Expected status see other (303); got %v", res.StatusCode)
	}
}
