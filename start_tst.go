package main

import (
	"encoding/xml"
	"io/ioutil"
	l "log"
	"net/http"
	"net/url"

	"golang.org/x/net/html/charset"
)

func httpRequest(req string) (*http.Response, error) {
	res, err := http.PostForm(req, url.Values{"email": {itest.User}, "pass": {itest.Pass}})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func xmlDecoder(res *http.Response) *xml.Decoder {
	dec := xml.NewDecoder(res.Body)
	dec.CharsetReader = charset.NewReaderLabel
	return dec
}

func startTest() {
	l.Println("Start test")
	res, err := httpRequest(entry.Request)
	if err != nil {
		l.Println("Error http request. Error= ", err)
	}
	// defer res.Body.Close()
	// fmt.Println(res.Body)
	// decoder := xmlDecoder(res)
	// var testinit structs.TestInitiation
	// if err := decoder.Decode(&testinit); err != nil {
	// 	l.Println("Error decode http response. Error: ", err)
	// }
	// fmt.Println(testinit)
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		l.Println("Не смог прочитать тело ответа. Ошибка=", err)
	}
	l.Println(string(body))
}
