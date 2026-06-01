package etl

import (
	"github.com/visagex/osrsdb-api/models"
	"github.com/visagex/osrsdb-api/wiki"
)

//fetch all the buckets using wiki package and normalize and structure

// return osrs-db item
func Normalize() {
	itemArray := []models.WikiItem{}
	bonusArray := []models.WikiBonus{}
	recipeArray := []models.WikiRecipe{}

	wiki.FetchAll(&itemArray, &bonusArray, &recipeArray)
}
