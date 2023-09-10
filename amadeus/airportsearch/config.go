package airportsearch

type Meta struct {
	Count int   `json:"count"`
	Links Links `json:"links"`
}

type Links struct {
	Self string `json:"self"`
}

type Location struct {
	Type           string    `json:"type"`
	SubType        string    `json:"subType"`
	Name           string    `json:"name"`
	DetailedName   string    `json:"detailedName"`
	ID             string    `json:"id"`
	Self           SelfLink  `json:"self"`
	TimeZoneOffset string    `json:"timeZoneOffset"`
	IataCode       string    `json:"iataCode"`
	GeoCode        GeoCode   `json:"geoCode"`
	Address        Address   `json:"address"`
	Analytics      Analytics `json:"analytics"`
}

type SelfLink struct {
	Href    string   `json:"href"`
	Methods []string `json:"methods"`
}

type GeoCode struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Address struct {
	CityName    string `json:"cityName"`
	CityCode    string `json:"cityCode"`
	CountryName string `json:"countryName"`
	CountryCode string `json:"countryCode"`
	RegionCode  string `json:"regionCode"`
}

type Analytics struct {
	Travelers Travelers `json:"travelers"`
}

type Travelers struct {
	Score int `json:"score"`
}
