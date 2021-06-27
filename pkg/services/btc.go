package services

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type BTCService struct {

}

func NewBTCService() *BTCService {
	return &BTCService{}
}

func (b *BTCService) GetPriceBTCInUAH() (float64, error) {
	type Response struct {
		Date string `json:"date"`
		BTC struct{
			UAH float64 `json:"uah"`
		} `json:"btc"`
	}

	resp, err := http.Get("https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/latest/currencies/btc.json")
	if err != nil {
		return .0, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Println(err)
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return .0, err
	}

	var response Response

	if err = json.Unmarshal(body, &response); err != nil {
		return .0, err
	}

	return response.BTC.UAH, nil
}