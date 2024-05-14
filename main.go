package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Product struct {
	Data Data `json:"data"`
}
type Data struct {
	ID                   int         `json:"id,omitempty"`
	Name                 string      `json:"name,omitempty"`
	Description          string      `json:"description,omitempty"`
	Image                string      `json:"image,omitempty"`
	GameName             string      `json:"game_name,omitempty"`
	ExpiredDate          time.Time   `json:"expired_date"`
	ChainName            string      `json:"chain_name,omitempty"`
	ChainID              string      `json:"chain_id,omitempty"`
	TokenID              string      `json:"token_id"`
	Fee                  string      `json:"fee"`
	Price                string      `json:"price"`
	TsxPrice             float64     `json:"tsx_price"`
	CurrencyTokenAddress string      `json:"currency_token_address"`
	Params               string      `json:"params,omitempty"`
	ParamsJson           interface{} `json:"params_json"`
	ParamsTh             string      `json:"params_th,omitempty"`
	ParamsThJson         interface{} `json:"params_th_json"`
	ParamsEn             string      `json:"params_en,omitempty"`
	ParamsEnJson         interface{} `json:"params_en_json"`
	NftTokenAddress      string      `json:"nft_token_address"`
	Seller               string      `json:"seller"`
	UpdatedDate          string      `json:"updated_date"`
	CreatedDate          time.Time   `json:"created_date"`
}

func main() {
	e := echo.New()

	// allow cors
	e.Use(middleware.CORS())

	e.GET("/api/v1/products/:id", func(c echo.Context) error {
		itemId := c.Param("id")
		res, err := GetProductNFT(itemId)
		if err != nil {
			return c.JSON(404, map[string]interface{}{
				"data": nil,
			})
		}
		return c.JSON(200, res)
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func GetProductNFT(itemId string) (Product, error) {
	url := fmt.Sprintf("https://prod-mkp-api.astronize.com/mkp/item/nft/0x7d4622363695473062cc0068686d81964bb6e09f/%s", itemId)
	resp, err := http.Get(url)
	if err != nil {
		return Product{}, err
	}
	defer resp.Body.Close()

	var product Product
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		return Product{}, err
	}
	if err = json.Unmarshal([]byte(product.Data.Params), &product.Data.ParamsJson); err != nil {
		return Product{}, err
	}
	if err = json.Unmarshal([]byte(product.Data.ParamsTh), &product.Data.ParamsThJson); err != nil {
		return Product{}, err
	}
	if err = json.Unmarshal([]byte(product.Data.ParamsEn), &product.Data.ParamsEnJson); err != nil {
		return Product{}, err
	}

	price, _ := strconv.Atoi(product.Data.Price)
	product.Data.TsxPrice = float64(price) / 1e18

	return product, nil
}
