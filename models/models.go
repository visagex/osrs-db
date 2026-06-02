package models

type WikiItem struct {
	PageName      string   `json:"page_name"`
	Item_id       []string `json:"item_id"`
	Item_name     string   `json:"item_name"`
	Examine       string   `json:"examine"`
	Tradeable     string   `json:"tradeable"`
	Weight        float64  `json:"weight"`
	HighAlchValue int      `json:"high_alchemy_value"`
}

type ItemBucketResponse struct {
	Items []WikiItem `json:"bucket"`
}

type RecipeBucketResponse struct {
	Recipes []WikiRecipe `json:"bucket"`
}

type BonusBucketResponse struct {
	Bonuses []WikiBonus `json:"bucket"`
}

type WikiBonus struct {
	PageName            string  `json:"page_name"`
	StabAttackBonus     int     `json:"stab_attack_bonus"`
	SlashAttackBonus    int     `json:"slash_attack_bonus"`
	CrushAttackBonus    int     `json:"crush_attack_bonus"`
	RangeAttackBonus    int     `json:"range_attack_bonus"`
	MagicAttackBonus    int     `json:"magic_attack_bonus"`
	StabDefenseBonus    int     `json:"stab_defence_bonus"`
	SlashDefenseBonus   int     `json:"slash_defence_bonus"`
	CrushDefenseBonus   int     `json:"crush_defence_bonus"`
	RangeDefenseBonus   int     `json:"range_defence_bonus"`
	MagicDefenseBonus   int     `json:"magic_defence_bonus"`
	StrengthBonus       int     `json:"strength_bonus"`
	RangedStrengthBonus int     `json:"ranged_strength_bonus"`
	PrayerBonus         int     `json:"prayer_bonus"`
	MagicDamageBonus    float64 `json:"magic_damage_bonus"`
	EquipmentSlot       string  `json:"equipment_slot"`
	AttackSpeed         int     `json:"weapon_attack_speed"`
	AttackRange         string  `json:"weapon_attack_range"`
	CombatStyle         string  `json:"combat_style"`
}

type WikiRecipe struct {
	PageName       string `json:"page_name"`
	ProductionJson string `json:"production_json"`
}

type ProductionJson struct {
	Ticks     int        `json:"ticks"`
	Materials []Material `json:"materials"`
	Skills    []Skill    `json:"skills"`
}

type Material struct {
	Quantity int    `json:"quantity"`
	Name     string `json:"name"`
}

type Skill struct {
	Experience int    `json:"experience"`
	Level      int    `json:"level"`
	Name       string `json:"name"`
	Boostable  string `json:"boostable"`
}

type OsrsItem struct {
	Item_info   WikiItem
	Item_recipe WikiRecipe
	Item_bonus  WikiBonus
}
