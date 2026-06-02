package etl

import (
	"github.com/visagex/osrsdb-api/models"
	"github.com/visagex/osrsdb-api/wiki"
)

//fetch all the buckets using wiki package and normalize and structure

// return osrsItem slice

func Normalize() []models.OsrsItem {

	itemArray := []models.WikiItem{}
	bonusArray := []models.WikiBonus{}
	recipeArray := []models.WikiRecipe{}

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
		osrsItem := models.OsrsItem{
			Item_info:   v,
			Item_recipe: recipeMap[k],
			Item_bonus:  bonusMap[k],
		}

		osrsItems = append(osrsItems, osrsItem)
	}

	return osrsItems
}
