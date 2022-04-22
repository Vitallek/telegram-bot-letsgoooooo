package misc

func GetFullCityName(city string, region string, country string, sep string) string {
	if country == "Белоруссия" {
		country = "Беларусь"
	}
	return city + sep + region + sep + country
}
