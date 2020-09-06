package controller

import (
	"encoding/json"
	"errors"
	"github.com/ammorteza/clean_architecture/entity"
	"log"
	"net/http"
	"time"
)

type urlController interface {
	CreateUrl(w http.ResponseWriter, r *http.Request)
	RedirectUrl(w http.ResponseWriter, r *http.Request)
	GetUrlRandomly(w http.ResponseWriter, r *http.Request)
}

func (c controller)GetUrlRandomly(w http.ResponseWriter, r *http.Request){
	url, err := c.service.GetUrlRandomly()
	if err != nil{
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if url.UniqueKey == ""{
		log.Println("there isn't any urls")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(&url); err != nil{
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (c controller)RedirectUrl(w http.ResponseWriter, r *http.Request)  {
	uniqueKey := r.URL.Query().Get("key")
	var original_url entity.Url
	original_url.UniqueKey = uniqueKey
	orgUrl, err := c.service.FetchUrl(original_url)
	if err != nil{
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, orgUrl.OriginalUrl, http.StatusSeeOther)
}

func (c controller)CreateUrl(w http.ResponseWriter, r *http.Request){
	var input entity.Url

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil{
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	input.CreatedAt = time.Now()
	input.ExpiredAt = (time.Now()).AddDate(1,0,0)
	key, err := getUniqueKey()
	if err != nil{
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	input.UniqueKey = key
	if err := c.service.CreateUrl(&input); err != nil{
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	type shortenUrl struct {
		Url 		string 			`json:"url"`
	}
	if err := json.NewEncoder(w).Encode(shortenUrl{Url: "http://172.28.1.20/" + input.UniqueKey}); err != nil{
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func getUniqueKey() (key string, err error){
	type uniqueKey struct {
		Key 		string			`json:"key"`
	}
	res, err := http.Post("http://172.28.2.20/unique_key", "application/json", nil)
	if err != nil{
		return "", err
	}

	defer func() {
		if res != nil{
			err = res.Body.Close()
		}
	}()
	if res.StatusCode != http.StatusOK{
		return "", errors.New("error in the unique key generation")
	}
	var unique_key uniqueKey
	if err = json.NewDecoder(res.Body).Decode(&unique_key); err != nil{
		return "", err
	}

	return unique_key.Key, nil
}

func (c controller)Redirect(w http.ResponseWriter, r *http.Request){

}

