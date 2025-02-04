package schema

type MaxmindDBRegion struct {
	Continent         MaxmindDBRegionContinent `json:"continent" maxminddb:"continent"`
	Country           MaxmindDBRegionCountry   `json:"country" maxminddb:"country"`
	Location          MaxmindDBRegionLocation  `json:"location" maxminddb:"location"`
	RegisteredCountry MaxmindDBRegionCountry   `json:"registered_country" maxminddb:"registered_country"`
}

type MaxmindDBRegionContinent struct {
	Code      string               `json:"code" maxminddb:"code"`
	GeonameID int                  `json:"geoname_id" maxminddb:"geoname_id"`
	Names     MaxmindDBRegionNames `json:"names" maxminddb:"names"`
}

type MaxmindDBRegionNames struct {
	//De string `json:"de" maxminddb:"de"`
	En string `json:"en" maxminddb:"en"`
	//Es string `json:"es" maxminddb:"es"`
	//Fr string `json:"fr" maxminddb:"fr"`
	//Ja string `json:"ja" maxminddb:"ja"`
	//BR string `json:"pt-BR" maxminddb:"pt-BR"`
	//Ru string `json:"ru" maxminddb:"ru"`
	//Zh string `json:"zh-CN" maxminddb:"zh-CN"`
}

type MaxmindDBRegionCountry struct {
	GeonameID int                  `json:"geoname_id" maxminddb:"geoname_id"`
	IsoCode   string               `json:"iso_code" maxminddb:"iso_code"`
	Names     MaxmindDBRegionNames `json:"names" maxminddb:"names"`
}

type MaxmindDBRegionLocation struct {
	AccuracyRadius int     `json:"accuracy_radius" maxminddb:"accuracy_radius"`
	Latitude       float64 `json:"latitude" maxminddb:"latitude"`
	Longitude      float64 `json:"longitude" maxminddb:"longitude"`
	TimeZone       string  `json:"time_zone" maxminddb:"time_zone"`
}

type MaxmindDBAS struct {
	AutonomousSystemNumber       int    `json:"autonomous_system_number" maxminddb:"autonomous_system_number"`
	AutonomousSystemOrganization string `json:"autonomous_system_organization" maxminddb:"autonomous_system_organization"`
}
