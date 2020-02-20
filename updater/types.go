package updater

import "time"

type ProductWrapper struct {
	Id            int                      `json:"id"       sql:"id, pk"`
	Name          string                   `json:"name"     sql:"name"`
	Options       []map[string]interface{} `json:"options"  pg:",json_number"`
	Children      []string                 `json:"-"        sql:"children, array"`
	CategoryId    string                   `json:"-"        sql:"category_id"`
	SubCategoryId string                   `json:"-"        sql:"sub_category_id"`
	LastUpdate    time.Time                `json:"-"        sql:"last_update"`
}

type Product struct {
	Id              string             `json:"Id"                      sql:"id, pk"`
	FullName        string             `json:"Name"                    sql:"full_name"`
	Model           string             `json:"Model"                   sql:"-"`
	Group           string             `json:"Group"                   sql:"-"`
	SubGroup        string             `json:"SubGroup"                sql:"-"`
	Size            string             `json:"Size"                    sql:"size"`
	Colour          string             `json:"Colour"                  sql:"colour"`
	CountInStock    float64            `json:"CountInStock"            sql:"count_in_stock"`
	CountInStocks   map[string]float64 `json:"-"                       sql:"count_in_stocks"`
	HasImage        bool               `json:"-"                       sql:"has_image"`
	TechSpec        string             `json:"TechnicalSpecifications" sql:"tech_spec"`
	LastUpdate      time.Time          `json:"-"                       sql:"last_update"`
}