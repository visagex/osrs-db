package wiki

//fetch buckets

//Bucket:Infobox bonuses
//Bucket:Infobox item
//Bucket:Recipe
//Bucket:Infobox npc
//Bucket:Infobox monster
//Bucket:Recommended equipment
//Bucket:Drop table sources

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/visagex/osrsdb-api/models"
)

// var itemArray []models.WikiItem
// var bonusArray []models.WikiBonus
// var recipeArray []models.WikiRecipe

const wikiURL = "https://oldschool.runescape.wiki/api.php"

var itemFields = []string{"item_id", "item_name", "examine", "tradeable", "weight", "value", "page_name", "high_alchemy_value", "quest"}
var bonusFields = []string{"page_name", "page_name_sub", "stab_attack_bonus", "slash_attack_bonus", "crush_attack_bonus", "range_attack_bonus", "magic_attack_bonus", "stab_defence_bonus", "slash_defence_bonus", "crush_defence_bonus",
	"range_defence_bonus", "magic_defence_bonus", "strength_bonus", "ranged_strength_bonus", "prayer_bonus", "magic_damage_bonus", "equipment_slot", "weapon_attack_speed", "weapon_attack_range", "combat_style"}
var recipeFields = []string{"page_name", "production_json"}

func buildFieldString(fields []string) string {
	res := ""
	quoted := make([]string, len(fields))
	for i, f := range fields {
		quoted[i] = "'" + f + "'"
	}

	res = strings.Join(quoted, ",")
	return res
}

//var itemMap = make(map[string]models.WikiItem)

//build a helper function that builds request easier, and returns *http.resquest

func buildRequest(limit int, offset int, rType string, fields []string) *http.Request {
	req, err := http.NewRequest("GET", wikiURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	f := buildFieldString(fields)

	var query string = fmt.Sprintf("bucket('%s').select(%s).limit(%d).offset(%d).run()", rType, f, limit, offset)

	params := req.URL.Query()
	params.Add("action", "bucket")
	params.Add("query", query)
	params.Add("format", "json")
	req.URL.RawQuery = params.Encode()
	req.Header.Set("User-Agent", "osrs-item-db/1.0 (kolbydeaguiar@gmail.com)")

	return req
}

// func sendRequest(request *http.Request) *http.Response

func fetchItems(itemArray *[]models.WikiItem, wg *sync.WaitGroup) {
	var limit int = 5000
	var offset int = 0
	client := &http.Client{}

	for {
		req := buildRequest(limit, offset, "infobox_item", itemFields)
		//fmt.Println("Requesting:", req.URL.String())

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println("Status:", resp.StatusCode)
		// fmt.Println("Raw:", string(body))

		var result models.ItemBucketResponse
		if err := json.Unmarshal(body, &result); err != nil {
			log.Fatal(err)
		}

		*itemArray = append(*itemArray, result.Items...)

		resp.Body.Close()

		if len(result.Items) < limit {
			wg.Done()
			break
		}

		offset += limit
	}
}

func fetchBonuses(bonusArray *[]models.WikiBonus, wg *sync.WaitGroup) {
	var limit int = 5000
	var offset int = 0
	client := &http.Client{}

	for {
		req := buildRequest(limit, offset, "infobox_bonuses", bonusFields)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		//fmt.Println("Status:", resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		//fmt.Println("Raw:", string(body))

		var result models.BonusBucketResponse
		if err := json.Unmarshal(body, &result); err != nil {
			log.Fatal(err)
		}

		*bonusArray = append(*bonusArray, result.Bonuses...)

		resp.Body.Close()

		if len(result.Bonuses) < limit {
			wg.Done()
			break
		}

		offset += limit

	}

}

func fetchRecipes(recipeArray *[]models.WikiRecipe, wg *sync.WaitGroup) {
	var limit int = 5000
	var offset int = 0
	client := &http.Client{}

	for {
		req := buildRequest(limit, offset, "recipe", recipeFields)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		//fmt.Println("Status:", resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var result models.RecipeBucketResponse
		if err := json.Unmarshal(body, &result); err != nil {
			log.Fatal(err)
		}

		*recipeArray = append(*recipeArray, result.Recipes...)

		resp.Body.Close()

		if len(result.Recipes) < limit {
			wg.Done()
			break
		}

		offset += limit
	}

}

func FetchAll(itemArray *[]models.WikiItem, bonusArray *[]models.WikiBonus, recipeArray *[]models.WikiRecipe) {
	var wg sync.WaitGroup
	wg.Add(1)
	go fetchItems(itemArray, &wg)
	wg.Add(1)
	go fetchBonuses(bonusArray, &wg)
	wg.Add(1)
	go fetchRecipes(recipeArray, &wg)
	wg.Wait()
	fmt.Println("fetched all buckets")
}
