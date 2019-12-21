package mystatus

type WttrInResponse struct {
	CurrentCondition []CurrentCondition `json:"current_condition"`
	Request          []Request          `json:"request"`
	Weather          []Weather          `json:"weather"`
}
type WeatherDesc struct {
	Value string `json:"value"`
}
type WeatherIconURL struct {
	Value string `json:"value"`
}
type CurrentCondition struct {
	FeelsLikeC      string           `json:"FeelsLikeC"`
	FeelsLikeF      string           `json:"FeelsLikeF"`
	Cloudcover      string           `json:"cloudcover"`
	Humidity        string           `json:"humidity"`
	ObservationTime string           `json:"observation_time"`
	PrecipMM        string           `json:"precipMM"`
	Pressure        string           `json:"pressure"`
	TempC           string           `json:"temp_C"`
	TempF           string           `json:"temp_F"`
	UvIndex         int              `json:"uvIndex"`
	Visibility      string           `json:"visibility"`
	WeatherCode     string           `json:"weatherCode"`
	WeatherDesc     []WeatherDesc    `json:"weatherDesc"`
	WeatherIconURL  []WeatherIconURL `json:"weatherIconUrl"`
	Winddir16Point  string           `json:"winddir16Point"`
	WinddirDegree   string           `json:"winddirDegree"`
	WindspeedKmph   string           `json:"windspeedKmph"`
	WindspeedMiles  string           `json:"windspeedMiles"`
}
type Request struct {
	Query string `json:"query"`
	Type  string `json:"type"`
}
type Astronomy struct {
	MoonIllumination string `json:"moon_illumination"`
	MoonPhase        string `json:"moon_phase"`
	Moonrise         string `json:"moonrise"`
	Moonset          string `json:"moonset"`
	Sunrise          string `json:"sunrise"`
	Sunset           string `json:"sunset"`
}
type Hourly struct {
	DewPointC        string           `json:"DewPointC"`
	DewPointF        string           `json:"DewPointF"`
	FeelsLikeC       string           `json:"FeelsLikeC"`
	FeelsLikeF       string           `json:"FeelsLikeF"`
	HeatIndexC       string           `json:"HeatIndexC"`
	HeatIndexF       string           `json:"HeatIndexF"`
	WindChillC       string           `json:"WindChillC"`
	WindChillF       string           `json:"WindChillF"`
	WindGustKmph     string           `json:"WindGustKmph"`
	WindGustMiles    string           `json:"WindGustMiles"`
	Chanceoffog      string           `json:"chanceoffog"`
	Chanceoffrost    string           `json:"chanceoffrost"`
	Chanceofhightemp string           `json:"chanceofhightemp"`
	Chanceofovercast string           `json:"chanceofovercast"`
	Chanceofrain     string           `json:"chanceofrain"`
	Chanceofremdry   string           `json:"chanceofremdry"`
	Chanceofsnow     string           `json:"chanceofsnow"`
	Chanceofsunshine string           `json:"chanceofsunshine"`
	Chanceofthunder  string           `json:"chanceofthunder"`
	Chanceofwindy    string           `json:"chanceofwindy"`
	Cloudcover       string           `json:"cloudcover"`
	Humidity         string           `json:"humidity"`
	PrecipMM         string           `json:"precipMM"`
	Pressure         string           `json:"pressure"`
	TempC            string           `json:"tempC"`
	TempF            string           `json:"tempF"`
	Time             string           `json:"time"`
	UvIndex          string           `json:"uvIndex"`
	Visibility       string           `json:"visibility"`
	WeatherCode      string           `json:"weatherCode"`
	WeatherDesc      []WeatherDesc    `json:"weatherDesc"`
	WeatherIconURL   []WeatherIconURL `json:"weatherIconUrl"`
	Winddir16Point   string           `json:"winddir16Point"`
	WinddirDegree    string           `json:"winddirDegree"`
	WindspeedKmph    string           `json:"windspeedKmph"`
	WindspeedMiles   string           `json:"windspeedMiles"`
}
type Weather struct {
	Astronomy   []Astronomy `json:"astronomy"`
	Date        string      `json:"date"`
	Hourly      []Hourly    `json:"hourly"`
	MaxtempC    string      `json:"maxtempC"`
	MaxtempF    string      `json:"maxtempF"`
	MintempC    string      `json:"mintempC"`
	MintempF    string      `json:"mintempF"`
	SunHour     string      `json:"sunHour"`
	TotalSnowCm string      `json:"totalSnow_cm"`
	UvIndex     string      `json:"uvIndex"`
}
