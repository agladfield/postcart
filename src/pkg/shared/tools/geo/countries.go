// Package geo wraps handling getting the iso2 country code from user input
package geo

import (
	"maps"
	"strings"
)

// copy map of names and iso3

func GetCountry(input string) string {
	normalizedCountry := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(input, "-", ""), " ", ""))

	iso2, exists := countryMap[normalizedCountry]
	if exists {
		return iso2
	}

	return "PC"
}

// Country mapping with normalized keys
var countryMap = mergeCountryMaps()

func mergeCountryMaps() map[string]string {
	newMap := make(map[string]string)
	maps.Copy(newMap, countryISO3s)
	maps.Copy(newMap, countryNamesToISO2)

	return newMap
}
