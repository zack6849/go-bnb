package domain

type City struct {
	DatabaseDrivenList
	Location Location
	//iso_3166-2, eg: US-NY (essentially a state)
	SubDivision string `gorm:"column:subdivision"`
	FullName    string `gorm:"column:full_name"`
	Timezone    string `gorm:"column:timezone"`
	Hash        string `gorm:"column:hash"`
}
