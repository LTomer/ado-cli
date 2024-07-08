package core

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

func json(url *string, user string, pat string) {
	var objmap []map[string]interface{}
	// data, _ := getContent(url, user, pat)

	// if err := json.Unmarshal([]byte(*data), &objmap); err != nil {
	// 	//log.Fatal(err)
	// }

	fmt.Println(objmap[0]["href"]) // to parse out your value
}

func getContent(url *string, user string, pat string) (*string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", basicAuthHeader(user, pat))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	content := string(body)
	return &content, nil
}

func basicAuthHeader(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
