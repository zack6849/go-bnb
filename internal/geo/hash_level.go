package geo

type HashLevel uint

// https://en.wikipedia.org/wiki/Geohash#Digits_and_precision_in_km

const (
	_                     = iota
	HashLevelContinent    // +- 1,600 mi
	HashLevelCountry      // +- 390 mi
	HashLevelState        // +- 48 mi
	HashLevelCity         // +- 12 mi
	HashLevelNeighborhood // +- 1.5 mi
	HashLevelBlock        // +- 0.38 mi
	HashLevelBuilding     // +- 0.047 mi
	HashLevelHouse        // +- 0.0012 mi
)
