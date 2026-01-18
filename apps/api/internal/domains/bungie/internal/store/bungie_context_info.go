package store

const bungieContextKey = "bungie"

type BungieContextInfo struct {
	Platform    string
	Gamertag    string
	Handler     string
	WeaponIndex int
	WeaponName  string
}
