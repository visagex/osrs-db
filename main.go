package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const wikiURL = "https://oldschool.runescape.wiki/api.php"

type WikiItem struct {
	Item_id   []string `json:"item_id"`
	Item_name string   `json:"item_name"`
	Examine   string   `json:"examine"`
	Tradeable string   `json:"tradeable"`
	Weight    float64  `json:"weight"`
}

type BucketResponse struct {
	Items []WikiItem `json:"bucket"`
}

func main() {
	var limit int = 5000
	var offset int = 0

	for {
		req, err := http.NewRequest("GET", wikiURL, nil)
		if err != nil {
			log.Fatal(err)
		}

		var query string = fmt.Sprintf("bucket('infobox_item').select('item_id','item_name','examine','tradeable','weight','value').limit(%d).offset(%d).run()", limit, offset)

		params := req.URL.Query()
		params.Add("action", "bucket")
		params.Add("query", query)
		params.Add("format", "json")
		req.URL.RawQuery = params.Encode()

		req.Header.Set("User-Agent", "osrs-item-db/1.0 (kolbydeaguiar@gmail.com)")
		//fmt.Println("Requesting:", req.URL.String())

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
			resp.Body.Close()
		}

		//fmt.Println("Status:", resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		//fmt.Println("Raw:", string(body))

		var result BucketResponse
		if err := json.Unmarshal(body, &result); err != nil {
			log.Fatal(err)
		}

		resp.Body.Close()

		for _, item := range result.Items {
			if len(item.Item_id) < 1 {
				continue
			}

			fmt.Printf("ID: %s  Name: %s  Tradeable: %s  Weight: %.3f\n",
				item.Item_id[0],
				item.Item_name,
				item.Tradeable,
				item.Weight,
			)
		}

		if len(result.Items) < limit {
			break
		}

		fmt.Println(offset)
		offset += limit

	}

}
