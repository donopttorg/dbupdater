package updater

import (
	"net/http"
	"os"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"github.com/juju/errors"
)

func SendPost(url string, params map[string]interface{}) ([]byte, error) {
	jsonStr, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Trace(err)
	}

	req, err := http.NewRequest("POST", os.Getenv("DB_API_SERVER") + url,
		bytes.NewBuffer(jsonStr))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}


func SendPostWithUrlParams(apiURL string, params map[string]string) ([]byte, int, error) {
	u, err := url.Parse(os.Getenv("DB_API_SERVER") + apiURL)
	if err != nil {
		return nil, -1, errors.Trace(err)
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("POST", 	u.String(), bytes.NewBuffer([]byte{}))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, -1, errors.Trace(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return body, resp.StatusCode, nil
}


func SendProtectedPostWithUrlParams(apiURL string, params map[string]string) ([]byte, int, error) {
	params["JToken"] = os.Getenv("JToken")
	return SendPostWithUrlParams(apiURL, params)
}
