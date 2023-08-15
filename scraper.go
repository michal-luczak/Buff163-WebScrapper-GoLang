package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Response struct {
	Code string `json:"code"`
	Data struct {
		Item []struct {
			Id int32 `json:"id"`
			BuyMaxPrice string `json:"buy_max_price"`
			SellMinPrice string `json:"sell_min_price"`
		} `json:"items"`
		TotalPage int `json:"total_page"`
	} `json:"data"`
	Msg interface{} `json:"msg"`
}

func getDataFromPage(page int, client *http.Client) int {
	items := fmt.Sprintf("https://buff.163.com/api/market/goods?game=csgo&page_num=%d&min_price=50&max_price=200&page_size=80&sort_by=price.asc", page)
	request, err := http.NewRequest("GET", items, nil)

	request.Header.Add("Accept-Language", "pl-PL")
	request.Header.Add("Cookie", "Device-Id=ZNs4dwiDdc9xeI6We42D; Locale-Supported=en; game=csgo; session=1-LW4479wzdu6dayTIY8F5MB9EHjg_6jV9fsyWQmDyOf652032876009; csrf_token=IjgzMGI5MmNlZTQ1NDQ4NTBiODUyNWNlNDgwNTcxZjZjMzgzMzQ3ZmEi.F64m1g.P7UuaH3r8JB2BWI4c24Q-GAIBfw")
		
	if err != nil {
		fmt.Println("Error while creating request:", err)
		return 1
	}

	// Wykonanie żądania
	resp, err := client.Do(request)


	if err != nil {
		fmt.Println("Error sending request:", err)
		return 1
	}

	defer resp.Body.Close()

	// Dekodowanie odpowiedzi JSON
	var response Response

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		fmt.Println("Error decoding response:", err)
		return 1
	}

	TotalPage := response.Data.TotalPage

	for _, item := range response.Data.Item {
		BuyMaxPriceFloat, err := strconv.ParseFloat(item.BuyMaxPrice, 32)
		SellMinPriceFloat, err := strconv.ParseFloat(item.SellMinPrice, 32)

		if err != nil {
			fmt.Println("Error converting string to float:", err)
			return 1
		}

		ratio := (SellMinPriceFloat * 0.975 - BuyMaxPriceFloat) / BuyMaxPriceFloat

		if ratio > 0.2 {
			fmt.Printf("Highest sell offer: %.2f\n", SellMinPriceFloat)
			fmt.Printf("Highest buy order: %.2f\n", BuyMaxPriceFloat)
			fmt.Printf("Profit percentage: %.2f%%\n", ratio * 100)
			fmt.Printf("Id: %d\n\n", item.Id)
		}
	}

	return TotalPage
}

func main() {
	
	client := &http.Client{}

	i := 1
	TotalPage := getDataFromPage(i, client)
	for i <= TotalPage {
		getDataFromPage(i, client)
		fmt.Println(i)
		fmt.Println(TotalPage)
		i++
	}
}