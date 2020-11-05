package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

type Coupon struct {
	Code string
}

type Coupons struct {
	// slice tipo um array infinito
	Coupon []Coupon
}

func (c Coupons) Check(code string) string {
	// _ = index
	// item = valor
	for _, item := range c.Coupon {
		if code == item.Code {
			return "valid"
		}
	}

	return "invalid"
}

type Result struct {
	Status string
}

var coupons Coupons
var codes []string

func main() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":9092", nil)
}

func init() {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	res, err := retryClient.Get("http://localhost:9093")

	if err == nil {
		defer res.Body.Close()
		data, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(data, &codes)

		for _, code := range codes {
			coupons.Coupon = append(coupons.Coupon, Coupon{
				Code: code,
			})
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	coupon := r.PostFormValue("coupon")
	valid := coupons.Check(coupon)

	result := Result{Status: valid}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))
}
