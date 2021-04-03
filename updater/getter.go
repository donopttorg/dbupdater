package updater

import (
	"encoding/json"
	"github.com/juju/errors"
	"math"
	"strings"
)


func getAllProductsFromServer() ([]*ProductWrapper, []*Product, error) {
	// downloading data from server
	rawData, err := func() ([]*Product, error) {
		raw, _, err := SendProtectedPostWithUrlParams("/GoodsChangeStokForSale", map[string]string{
			"LastVersionStock": "0",
			"LastVersionPrice": "0",
			"LastVersionGoods": "0",
		})
		if err != nil {
			return nil, errors.Trace(err)
		}

		parsed := make([]*Product, 0)
		err = json.Unmarshal(raw, &parsed)
		if err != nil {
			return nil, errors.Trace(err)
		}

		for i := 0; i < len(parsed); {
			val := parsed[i]

			if len(strings.Replace(val.Model, " ", "", -1)) == 0 {
				//log.Println(i, len(parsed))
				parsed = parsed[:i+copy(parsed[i:], parsed[i+1:])]
			} else {
				i++
			}
		}

		// rounding count in stocks
		for _, item := range parsed {
			for key, countInStock := range item.CountInStocks {
				item.CountInStocks[key] = math.Ceil(countInStock*100) / 100
			}
		}

		return parsed, nil
	}()
	if err != nil {
		return nil, nil, errors.Trace(err)
	}

	//processing all data
	rawRes := func(rawData []*Product) map[string]*ProductWrapper {
		parseModel := func(pr *Product) (modelName string, options []string) {
			if strings.Contains(pr.Model, "_") {
				options = strings.Split(pr.Model, "_")
				modelName = options[0]
				options = options[1:]
				options = append(options, pr.Colour)

				if len(pr.Size) != 0 {
					options = append(options, pr.Size)
				}
				return
			} else {
				return pr.Model, []string{pr.Colour, pr.Size}
			}
		}

		createPath := func(pw *ProductWrapper, pr *Product, options []string) {
			root := &pw.Options

			for i, optName := range options {
				found := false

				for _, opt := range *root {

					if opt["name"].(string) == optName {
						if i != (len(options) - 1) {
							found = true
							root = opt["options"].(*[]map[string]interface{})
							break
						} else {
							return
						}
					}

				}

				if !found {
					if i == (len(options) - 1) {
						*root = append(*root, map[string]interface{}{
							"name": optName,
							"id": pr.Id,
						})
					} else {
						newRoot := &[]map[string]interface{}{}

						*root = append(*root, map[string]interface{}{
							"name": optName,
							"options": newRoot,
						})
						root = newRoot
					}
				}
			}
		}

		removeSpaces := func(st string) string {
			st = strings.TrimLeft(st, " ")
			st = strings.TrimRight(st, " ")
			return st
		}

		resData := make(map[string]*ProductWrapper, 0)
		for _, val := range rawData {
			modelName, options := parseModel(val)
			//log.Println(val.Model, modelName, options)

			if resData[modelName] == nil {
				resData[modelName] = &ProductWrapper{
					Name: modelName,
					Children: make([]string, 0),
					Options: make([]map[string]interface{}, 0, 0),
				}
			}

			pw := resData[modelName]
			pw.Children = append(pw.Children, val.Id)
			createPath(pw, val, options)

			if len(pw.CategoryId) == 0 {
				pw.CategoryId = removeSpaces(val.Group)
				pw.SubCategoryId = removeSpaces(val.SubGroup)
			}
		}

		return resData
	}(rawData)

	res := make([]*ProductWrapper, len(rawRes))
	i := 0
	for _, v := range rawRes {
		res[i] = v
		i++
	}

	return res, rawData, nil
}

func getCountInStocksFromServer() ([]*Product, error) {
	products, err := func() ([]*Product, error) {
		raw, _, err := SendProtectedPostWithUrlParams("/GoodsChangeStokForSale", map[string]string{
			"LastVersionStock": "0",
			"LastVersionPrice": "0",
			"LastVersionGoods": "0",
		})
		if err != nil {
			return nil, errors.Trace(err)
		}

		parsed := make([]*Product, 0)
		err = json.Unmarshal(raw, &parsed)
		if err != nil {
			return nil, errors.Trace(err)
		}

		for i := 0; i < len(parsed); {
			val := parsed[i]

			if len(strings.Replace(val.Model, " ", "", -1)) == 0 {
				//log.Println(i, len(parsed))
				parsed = parsed[:i+copy(parsed[i:], parsed[i+1:])]
			} else {
				i++
			}
		}

		return parsed, nil
	}()
	if err != nil {
		return  nil, errors.Trace(err)
	}

	rawCountInStocks, _, err := SendProtectedPostWithUrlParams("/BalancesForCustomerTwoWarehouses", map[string]string{})
	if err != nil {
		return nil, errors.Trace(err)
	}

	countInStocksResp := make(map[string]interface{})
	err = json.Unmarshal(rawCountInStocks, &countInStocksResp)
	if err != nil {
		return nil, errors.Trace(err)
	}

	countInStocksData := countInStocksResp["listOfGoodsForDisclosures"].([]interface{})
	for _, product := range products {
		for _, data := range countInStocksData {
			parsedData := data.(map[string]interface{})
			if parsedData["Id"].(string) == product.Id {
				parsedCountInStocks := make(map[string]float64)
				for key, val := range parsedData["Quantities"].(map[string]interface{}) {
					parsedCountInStocks[key] = val.(float64)
				}
				product.CountInStocks = parsedCountInStocks
			}
		}
	}

	return products, nil
}