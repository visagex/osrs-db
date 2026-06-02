package etl

import (
	"encoding/json"
	"fmt"

	"github.com/visagex/osrsdb-api/models"
	"github.com/visagex/osrsdb-api/wiki"
)

//fetch all the buckets using wiki package and build into osrsItem struct slice

// return osrsItem slice

func buildItems() []models.OsrsItem {

	itemArray := []models.WikiItem{}
	bonusArray := []models.WikiBonus{}
	recipeArray := []models.WikiRecipe{}

	fmt.Println("fetching...")
	wiki.FetchAll(&itemArray, &bonusArray, &recipeArray)

	itemMap := make(map[string]models.WikiItem)
	bonusMap := make(map[string]models.WikiBonus)
	recipeMap := make(map[string]models.WikiRecipe)

	for _, item := range itemArray {
		itemMap[item.PageName] = item
	}

	for _, bonus := range bonusArray {
		bonusMap[bonus.PageName] = bonus
	}

	for _, recipe := range recipeArray {
		recipeMap[recipe.PageName] = recipe
	}

	osrsItems := []models.OsrsItem{}

	for k, v := range itemMap {

		prodJson := models.ProductionJson{}
		if recipeMap[k].ProductionJson != "" {
			json.Unmarshal([]byte(recipeMap[k].ProductionJson), &prodJson)
		}

		osrsItem := models.OsrsItem{
			Item_info:   v,
			Item_recipe: prodJson,
			Item_bonus:  bonusMap[k],
		}

		osrsItems = append(osrsItems, osrsItem)
	}

	// for _, item := range osrsItems {
	// 	fmt.Println(item)
	// }

	return osrsItems
}
