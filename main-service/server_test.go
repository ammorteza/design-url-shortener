package main

import (
	"bytes"
	"encoding/json"
	"github.com/ammorteza/clean_architecture/entity"
	"github.com/ammorteza/clean_architecture/http/gin"
	gin2 "github.com/gin-gonic/gin"
	uuid "github.com/nu7hatch/gouuid"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUrl(t *testing.T){
	var input entity.Url
	t.Parallel()
	input.OriginalUrl = "http://google.com"
	userId, err := uuid.NewV4()
	if err != nil{
		t.Fatal(err)
	}
	input.User = userId.String()
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(input); err != nil{
		t.Fatal(err)
	}

	httpServer := initRoutes(gin.New()).GetDispatcher().(*gin2.Engine)
	request := httptest.NewRequest("POST", "/create_url", buf)
	response := httptest.NewRecorder()
	httpServer.ServeHTTP(response, request)
	if response.Code != http.StatusOK{
		t.Fatal("incorrect response")
	}
}

func TestRedirectUrl(t *testing.T){
	var url entity.Url
	httpServer := initRoutes(gin.New()).GetDispatcher().(*gin2.Engine)
	getRandomUrlRequest := httptest.NewRequest("POST", "/get_url_randomly", nil)
	getRandomUrlResponse := httptest.NewRecorder()
	httpServer.ServeHTTP(getRandomUrlResponse, getRandomUrlRequest)
	if getRandomUrlResponse.Code != http.StatusOK{
		t.Error("cannot get random url")
	}else{
		err := json.NewDecoder(getRandomUrlResponse.Body).Decode(&url)
		if err != nil || url.UniqueKey == ""{
			t.Fatal("cannot get random url")
		}else{
			request := httptest.NewRequest("GET", "/" + url.UniqueKey, nil)
			response := httptest.NewRecorder()

			httpServer.ServeHTTP(response, request)
			if response.Code != http.StatusSeeOther{
				t.Fatal("incorrect redirect")
			}
		}
	}
}