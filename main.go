package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Message struct {
	Address Address
	ApiName string
}

type Address struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

const CEP = "81010300"
const URL_CDN = "https://cdn.apicep.com/file/apicep/" + CEP + ".json"
const URL_VIA = "http://viacep.com.br/ws/" + CEP + "/json/"

func main() {
	res1 := make(chan *http.Response)
	res2 := make(chan *http.Response)

	go func() {
		res := fetchRequest(URL_CDN)
		res1 <- res
	}()

	go func() {
		res := fetchRequest(URL_VIA)
		res2 <- res
	}()

	select {
	case res := <-res1:
		defer res.Body.Close()
		address := unmarshalResponse(res)

		responseObj := Message{
			Address: *address,
			ApiName: "CDN",
		}
		fmt.Printf("API Name: %s Address: %s\n", responseObj.ApiName, responseObj.Address)
	case res := <-res2:
		defer res.Body.Close()
		address := unmarshalResponse(res)

		responseObj := Message{
			Address: *address,
			ApiName: "VIA",
		}
		fmt.Printf("API Name: %s Address: %s\n", responseObj.ApiName, responseObj.Address)

	case <-time.After(time.Second * 50):
		fmt.Println("Timeout")
	}
}

func fetchRequest(url string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	return res
}

func unmarshalResponse(res *http.Response) *Address {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var address Address
	err = json.Unmarshal(body, &address)
	if err != nil {
		panic(err)
	}

	return &address
}
