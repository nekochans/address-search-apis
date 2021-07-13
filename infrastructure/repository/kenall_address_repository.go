package repository

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/nekochans/address-search-apis/domain"
	"github.com/nekochans/address-search-apis/infrastructure"
)

type KenallAddressRepository struct {
	HttpClient infrastructure.HttpClient
}

type Address struct {
	Jisx0402           string `json:"jisx0402"`
	OldCode            string `json:"old_code"`
	PostalCode         string `json:"postal_code"`
	PrefectureKana     string `json:"prefecture_kana"`
	CityKana           string `json:"city_kana"`
	TownKana           string `json:"town_kana"`
	TownKanaRaw        string `json:"town_kana_raw"`
	Prefecture         string `json:"prefecture"`
	City               string `json:"city"`
	Town               string `json:"town"`
	Koaza              string `json:"koaza"`
	KyotoStreet        string `json:"kyoto_street"`
	Building           string `json:"building"`
	Floor              string `json:"floor"`
	TownPartial        bool   `json:"town_partial"`
	TownAddressedKoaza bool   `json:"town_addressed_koaza"`
	TownChome          bool   `json:"town_chome"`
	TownMulti          bool   `json:"town_multi"`
	TownRaw            string `json:"town_raw"`
	Corporation        struct {
		Name       string `json:"name"`
		NameKana   string `json:"name_kana"`
		BlockLot   string `json:"block_lot"`
		PostOffice string `json:"post_office"`
		CodeType   int    `json:"code_type"`
	} `json:"corporation"`
}

type FindAddressesResponse struct {
	Version   string     `json:"version"`
	Addresses []*Address `json:"data"`
}

func (r *KenallAddressRepository) FindByPostalCode(postalCode string) (*domain.Address, error) {
	ctx := context.Background()
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.kenall.jp/v1/postalcode/"+postalCode, nil)
	req.Header.Set("Authorization", "Token "+os.Getenv("KENALL_API_KEY"))

	resp, err := r.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	var resBody FindAddressesResponse
	if err := json.NewDecoder(resp.Body).Decode(&resBody); err != nil {
		return nil, err
	}

	resAddress := resBody.Addresses[0]

	address := &domain.Address{
		PostalCode: postalCode,
		Prefecture: resAddress.Prefecture,
		Locality:   resAddress.City + resAddress.Town,
	}

	return address, nil
}
