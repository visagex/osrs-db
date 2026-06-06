package etl

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
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
			Name:        k,
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

func initDb(db *sql.DB) {

	schema := `
    CREATE TABLE IF NOT EXISTS items (
        page_name           TEXT PRIMARY KEY,
        item_id             TEXT,
        item_name           TEXT,
        examine             TEXT,
        tradeable           TEXT,
        weight              DOUBLE,
        high_alchemy_value  INTEGER,
		quest				TEXT
    );

    CREATE TABLE IF NOT EXISTS bonuses (
        page_name               TEXT PRIMARY KEY,
        stab_attack_bonus       INTEGER,
        slash_attack_bonus      INTEGER,
        crush_attack_bonus      INTEGER,
        range_attack_bonus      INTEGER,
        magic_attack_bonus      INTEGER,
        stab_defense_bonus      INTEGER,
        slash_defense_bonus     INTEGER,
        crush_defense_bonus     INTEGER,
        range_defense_bonus     INTEGER,
        magic_defense_bonus     INTEGER,
        strength_bonus          INTEGER,
        ranged_strength_bonus   INTEGER,
        prayer_bonus            INTEGER,
        magic_damage_bonus      REAL,
        equipment_slot          TEXT,
        attack_speed            INTEGER,
        attack_range            TEXT,
        combat_style            TEXT,
        FOREIGN KEY (page_name) REFERENCES items(page_name)
    );

    CREATE TABLE IF NOT EXISTS recipes (
        page_name   TEXT PRIMARY KEY,
        ticks       TEXT,
        FOREIGN KEY (page_name) REFERENCES items(page_name)
    );

    CREATE TABLE IF NOT EXISTS materials (
        id          INTEGER PRIMARY KEY AUTOINCREMENT,
        page_name   TEXT,
        name        TEXT,
        quantity    REAL,
        FOREIGN KEY (page_name) REFERENCES recipes(page_name)
    );

    CREATE TABLE IF NOT EXISTS skills (
        id          INTEGER PRIMARY KEY AUTOINCREMENT,
        page_name   TEXT,
        name        TEXT,
        level       INTEGER,
        experience  TEXT,
        boostable   TEXT,
        FOREIGN KEY (page_name) REFERENCES recipes(page_name)
    );`

	_, err := db.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("created successfuly")
}

func insertItems(db *sql.DB, items []models.OsrsItem) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // no-op if tx.Commit() is called

	itemStmt, err := tx.Prepare(`
        INSERT OR IGNORE INTO items (page_name, item_id, item_name, examine, tradeable, weight, high_alchemy_value, quest)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		return err
	}
	defer itemStmt.Close()

	bonusStmt, err := tx.Prepare(`
        INSERT OR IGNORE INTO bonuses (page_name, stab_attack_bonus, slash_attack_bonus, crush_attack_bonus, range_attack_bonus, magic_attack_bonus, stab_defense_bonus, slash_defense_bonus, crush_defense_bonus, range_defense_bonus, magic_defense_bonus, strength_bonus, ranged_strength_bonus, prayer_bonus, magic_damage_bonus, equipment_slot, attack_speed, attack_range, combat_style)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		return err
	}
	defer bonusStmt.Close()

	recipeStmt, err := tx.Prepare(`
        INSERT OR IGNORE INTO recipes (page_name, ticks)
        VALUES (?, ?)
    `)
	if err != nil {
		return err
	}
	defer recipeStmt.Close()

	materialStmt, err := tx.Prepare(`
        INSERT OR IGNORE INTO materials (page_name, name, quantity)
        VALUES (?, ?, ?)
    `)
	if err != nil {
		return err
	}
	defer materialStmt.Close()

	skillStmt, err := tx.Prepare(`
        INSERT OR IGNORE INTO skills (page_name, name, level, experience, boostable)
        VALUES (?, ?, ?, ?, ?)
    `)
	if err != nil {
		return err
	}
	defer skillStmt.Close()

	for _, osrsItem := range items {
		itemID, _ := json.Marshal(osrsItem.Item_info.Item_id)

		_, err = itemStmt.Exec(
			osrsItem.Item_info.PageName,
			string(itemID),
			osrsItem.Item_info.Item_name,
			osrsItem.Item_info.Examine,
			osrsItem.Item_info.Tradeable,
			osrsItem.Item_info.Weight,
			osrsItem.Item_info.HighAlchValue,
			osrsItem.Item_info.Quest,
		)
		if err != nil {
			return err
		}

		_, err = bonusStmt.Exec(
			osrsItem.Item_bonus.PageName,
			osrsItem.Item_bonus.StabAttackBonus,
			osrsItem.Item_bonus.SlashAttackBonus,
			osrsItem.Item_bonus.CrushAttackBonus,
			osrsItem.Item_bonus.RangeAttackBonus,
			osrsItem.Item_bonus.MagicAttackBonus,
			osrsItem.Item_bonus.StabDefenseBonus,
			osrsItem.Item_bonus.SlashDefenseBonus,
			osrsItem.Item_bonus.CrushDefenseBonus,
			osrsItem.Item_bonus.RangeDefenseBonus,
			osrsItem.Item_bonus.MagicDefenseBonus,
			osrsItem.Item_bonus.StrengthBonus,
			osrsItem.Item_bonus.RangedStrengthBonus,
			osrsItem.Item_bonus.PrayerBonus,
			osrsItem.Item_bonus.MagicDamageBonus,
			osrsItem.Item_bonus.EquipmentSlot,
			osrsItem.Item_bonus.AttackSpeed,
			osrsItem.Item_bonus.AttackRange,
			osrsItem.Item_bonus.CombatStyle,
		)
		if err != nil {
			return err
		}

		_, err = recipeStmt.Exec(
			osrsItem.Item_info.PageName,
			osrsItem.Item_recipe.Ticks,
		)
		if err != nil {
			return err
		}

		for _, material := range osrsItem.Item_recipe.Materials {

			var quantityNum float64
			if material.Quantity != "" && material.Quantity != "NA" {
				quantityNum, err = strconv.ParseFloat(material.Quantity, 64)
				if err != nil {
					log.Fatal("Material", err)
				}
			}

			_, err = materialStmt.Exec(
				osrsItem.Item_info.PageName,
				material.Name,
				quantityNum,
			)
			if err != nil {
				return err
			}
		}

		for _, skill := range osrsItem.Item_recipe.Skills {
			var levelNum int

			if skill.Level != "" && skill.Level != "NA" && skill.Level != "?" {
				levelNum, err = strconv.Atoi(skill.Level)
				if err != nil {
					log.Fatal("Skill", err)
				}
			}

			_, err = skillStmt.Exec(
				osrsItem.Item_info.PageName,
				skill.Name,
				levelNum,
				skill.Experience,
				skill.Boostable,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func Run() {
	items := buildItems()

	db, err := sql.Open("sqlite3", "D:/Workspaces/osrsdb-api/osrs.db")
	if err != nil {
		log.Fatal("open error: ", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("ping error: ", err)
	}
	fmt.Println("db opened successfully")

	initDb(db)
	fmt.Println("init done")

	err = insertItems(db, items)
	if err != nil {
		log.Fatal("insert error: ", err)
	}

	fmt.Println("inserted ", len(items), " items")
}
