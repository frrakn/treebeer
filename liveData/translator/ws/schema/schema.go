package schema

import "encoding/json"

type Game struct {
	TeamStats   map[string]map[string]json.RawMessage
	PlayerStats map[string]map[string]json.RawMessage
	GameID      string
	MatchID     string
	//	GeneratedName string
	//	Realm         string
	GameComplete bool
	T            int32
}

/* type Team struct {
	//	TeamID           int32
	BaronsKilled     int32
	DragonsKilled    int32
	FirstBlood       bool
	TowersKilled     int32
	InhibitorsKilled int32
	MatchVictory     int32
	MatchDefeat      int32
	Color            string
}

type Player struct {
	//	ParticipantID   int32
	//	TeamID          int32
	//	SummonerName    string
	ChampionName string
	ChampionID   int32
	//	SkinIndex       int32
	//	ProfileIconID   int32
	SummonerSpell1  int32
	SummonerSpell2  int32
	Kills           int32
	Deaths          int32
	Assists         int32
	DoubleKills     int32
	TripleKills     int32
	QuadraKills     int32
	PentaKills      int32
	Level           int32
	MaxHealth       int32
	MaxPower        int32
	AttackDamage    int32
	AbilityPower    int32
	Armor           int32
	MagicResist     int32
	AttackSpeed     int32
	CcReduction     int32
	MovementSpeed   int32
	SpellVamp       int32
	Lifesteal       int32
	ArmorPen        int32
	MagicPen        int32
	ArmorPenPercent int32
	MagicPenPercent int32
	HealthRegen     int32
	PowerRegen      int32
	WardsPlaced     int32
	WardsKilled     int32
	Runes           []*Rune
	Masteries       []*Mastery
	Items           []int32
	Skills          map[string]int32
	Mk              int32 // Minion Kills
	Cg              int32 // Current Gold
	Tg              int32 // Total Gold
	Xp              int32 // Experience
	X               int32 // X-Coordinate
	Y               int32 // Y-Coordinate
	H               int32 // Health
	P               int32 // Power (Mana, Energy, etc.)
	Td              int32 // Total Damage
	Pd              int32 // Physical Damage
	Md              int32 // Magic Damage
	Trd             int32 // True Damage
	Tdc             int32 // Total Damage to Champions
	Pdc             int32 // Physical Damage to Champions
	Mdc             int32 // Magic Damage to Champions
	Trdc            int32 // True Damage to Champions
	PlayerID        string
}

type Rune struct {
	RuneID int32
	Count  int32
}

type Mastery struct {
	MasteryID int32
	Rank      int32
} */
