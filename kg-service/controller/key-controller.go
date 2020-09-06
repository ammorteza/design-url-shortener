package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

type keyController interface {
	GetUniqueKey(w http.ResponseWriter, r *http.Request)
}

func (c controller)GetUniqueKey(w http.ResponseWriter, r *http.Request){
	type uniqueKey struct {
		Key 		string			`json:"key"`
	}

	key, err := c.service.FindUniqueKey()
	if err != nil{
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var unique_key uniqueKey = uniqueKey{Key: key.Key}
	if err := json.NewEncoder(w).Encode(unique_key); err != nil{
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}