package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

/*
 sync  : 8.01 sec
 async : 3.00 sec
*/
var (
	pizza = flag.String("pizza", "", "Pizza type")
	store = flag.String("store", "", "Pizza store name")
	price = flag.String("price", "", "Pizza price")
)

const (
	OrderService   = "http://localhost:8081"
	PaymentService = "http://localhost:8082"
	StoreService   = "http://localhost:8083"
)

type PizzaOrderRequest struct {
	Pizza, Store, Price string
}

func main() {
	flag.Parse()
	order := PizzaOrderRequest{
		Pizza: *pizza,
		Store: *store,
		Price: *price,
	}
	body, _ := json.Marshal(&order)
	start := time.Now()
	orderChan := make(chan *http.Response, 1)
	paymentChan := make(chan *http.Response, 1)
	storeChan := make(chan *http.Response, 1)
	go func() {
		err := SendPostRequestAsync(OrderService, body, orderChan)
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		err := SendPostRequestAsync(PaymentService, body, paymentChan)
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		err := SendPostRequestAsync(StoreService, body, storeChan)
		if err != nil {
			panic(err)
		}
	}()
	orderServiceResp := <-orderChan
	defer orderServiceResp.Body.Close()
	bs, _ := ioutil.ReadAll(orderServiceResp.Body)
	fmt.Println(string(bs))
	paymentServiceResp := <-paymentChan
	defer paymentServiceResp.Body.Close()
	bs, _ = ioutil.ReadAll(paymentServiceResp.Body)
	fmt.Println(string(bs))
	storeServiceResp := <-storeChan
	defer orderServiceResp.Body.Close()
	bs, _ = ioutil.ReadAll(storeServiceResp.Body)
	fmt.Println(string(bs))
	end := time.Now()
	fmt.Printf("Order processed after %v secondes \n", end.Sub(start).Seconds())
}
func SendPostRequest(url string, body []byte) *http.Response {
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	return resp
}
func SendPostRequestAsync(url string, body []byte, rc chan *http.Response) error {
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err == nil {
		rc <- resp
	}
	return err
}
