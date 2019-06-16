package main

import "encoding/json"

type ProductWrapper struct {
	Id            int                      `json:"id"       sql:"id, pk"`
	Name          string                   `json:"name"     sql:"name"`
	Options       []map[string]interface{} `json:"options"  pg:",json_number"`
	Children      []string                 `json:"-"        sql:"children, array"`
	CategoryId    string                   `json:"-"        sql:"category_id"`
	SubCategoryId string                   `json:"-"        sql:"sub_category_id"`
}

type Product struct {
	Id             string  `json:"Id"             sql:"id, pk"`
	FullName       string  `json:"Name"           sql:"full_name"`
	Model          string  `json:"Model"          sql:"-"`
	Group          string  `json:"Group"          sql:"-"`
	SubGroup       string  `json:"SubGroup"       sql:"-"`
	Size           string  `json:"Size"           sql:"size"`
	Colour         string  `json:"Colour"         sql:"colour"`
	CountInStock   float64 `json:"CountInStock"   sql:"count_in_stock"`
	PriceRetail    float64 `json:"PriceRetail"    sql:"price_retail"`
	PriceWholesale float64 `json:"PriceWholesale" sql:"price_wholesale"`
	HasImage       bool    `json:"-"              sql:"has_image"`
}

func (product Product) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"id": product.Id,
		"fullName": product.FullName,
		"size": product.Size,
		"color": product.Colour,
		"priceRetail": product.PriceRetail,
		"priceWholesale": product.PriceWholesale,
	}

	if product.HasImage {
		m["hasImage"] = true
	}

	if product.CountInStock > 5 {
		m["isCountInStockTooBig"] = true
		m["countInStock"] = 5
	} else {
		m["countInStock"] = product.CountInStock
	}

	return json.Marshal(m)
}
