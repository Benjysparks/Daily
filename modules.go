package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"workspace/github.com/Benjysparks/daily/internal/database"
)

type ForecastStruct struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		WindchillC float64 `json:"windchill_c"`
		WindchillF float64 `json:"windchill_f"`
		HeatindexC float64 `json:"heatindex_c"`
		HeatindexF float64 `json:"heatindex_f"`
		DewpointC  float64 `json:"dewpoint_c"`
		DewpointF  float64 `json:"dewpoint_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Date      string `json:"date"`
			DateEpoch int    `json:"date_epoch"`
			Day       struct {
				MaxtempC          float64 `json:"maxtemp_c"`
				MaxtempF          float64 `json:"maxtemp_f"`
				MintempC          float64 `json:"mintemp_c"`
				MintempF          float64 `json:"mintemp_f"`
				AvgtempC          float64 `json:"avgtemp_c"`
				AvgtempF          float64 `json:"avgtemp_f"`
				MaxwindMph        float64 `json:"maxwind_mph"`
				MaxwindKph        float64 `json:"maxwind_kph"`
				TotalprecipMm     float64 `json:"totalprecip_mm"`
				TotalprecipIn     float64 `json:"totalprecip_in"`
				TotalsnowCm       float64 `json:"totalsnow_cm"`
				AvgvisKm          float64 `json:"avgvis_km"`
				AvgvisMiles       float64 `json:"avgvis_miles"`
				Avghumidity       int     `json:"avghumidity"`
				DailyWillItRain   int     `json:"daily_will_it_rain"`
				DailyChanceOfRain int     `json:"daily_chance_of_rain"`
				DailyWillItSnow   int     `json:"daily_will_it_snow"`
				DailyChanceOfSnow int     `json:"daily_chance_of_snow"`
				Condition         struct {
					Text string `json:"text"`
					Icon string `json:"icon"`
					Code int    `json:"code"`
				} `json:"condition"`
				Uv float64 `json:"uv"`
			} `json:"day"`
			Astro struct {
				Sunrise          string `json:"sunrise"`
				Sunset           string `json:"sunset"`
				Moonrise         string `json:"moonrise"`
				Moonset          string `json:"moonset"`
				MoonPhase        string `json:"moon_phase"`
				MoonIllumination int    `json:"moon_illumination"`
				IsMoonUp         int    `json:"is_moon_up"`
				IsSunUp          int    `json:"is_sun_up"`
			} `json:"astro"`
			Hour []struct {
				TimeEpoch int     `json:"time_epoch"`
				Time      string  `json:"time"`
				TempC     float64 `json:"temp_c"`
				TempF     float64 `json:"temp_f"`
				IsDay     int     `json:"is_day"`
				Condition struct {
					Text string `json:"text"`
					Icon string `json:"icon"`
					Code int    `json:"code"`
				} `json:"condition"`
				WindMph      float64 `json:"wind_mph"`
				WindKph      float64 `json:"wind_kph"`
				WindDegree   int     `json:"wind_degree"`
				WindDir      string  `json:"wind_dir"`
				PressureMb   float64 `json:"pressure_mb"`
				PressureIn   float64 `json:"pressure_in"`
				PrecipMm     float64 `json:"precip_mm"`
				PrecipIn     float64 `json:"precip_in"`
				SnowCm       float64 `json:"snow_cm"`
				Humidity     int     `json:"humidity"`
				Cloud        int     `json:"cloud"`
				FeelslikeC   float64 `json:"feelslike_c"`
				FeelslikeF   float64 `json:"feelslike_f"`
				WindchillC   float64 `json:"windchill_c"`
				WindchillF   float64 `json:"windchill_f"`
				HeatindexC   float64 `json:"heatindex_c"`
				HeatindexF   float64 `json:"heatindex_f"`
				DewpointC    float64 `json:"dewpoint_c"`
				DewpointF    float64 `json:"dewpoint_f"`
				WillItRain   int     `json:"will_it_rain"`
				ChanceOfRain int     `json:"chance_of_rain"`
				WillItSnow   int     `json:"will_it_snow"`
				ChanceOfSnow int     `json:"chance_of_snow"`
				VisKm        float64 `json:"vis_km"`
				VisMiles     float64 `json:"vis_miles"`
				GustMph      float64 `json:"gust_mph"`
				GustKph      float64 `json:"gust_kph"`
				Uv           float64 `json:"uv"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

type NewsArticles struct {
	TotalArticles int `json:"totalArticles"`
	Articles      []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Content     string `json:"content"`
		URL         string `json:"url"`
		Image       string `json:"image"`
		PublishedAt string `json:"publishedAt"`
		Source      struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"source"`
	} `json:"articles"`
}

type CatStruct []struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type PremTable struct {
	Filters struct {
		Season string `json:"season"`
	} `json:"filters"`
	Area struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
		Flag string `json:"flag"`
	} `json:"area"`
	Competition struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Code   string `json:"code"`
		Type   string `json:"type"`
		Emblem string `json:"emblem"`
	} `json:"competition"`
	Season struct {
		ID              int         `json:"id"`
		StartDate       string      `json:"startDate"`
		EndDate         string      `json:"endDate"`
		CurrentMatchday int         `json:"currentMatchday"`
		Winner          interface{} `json:"winner"`
	} `json:"season"`
	Standings []struct {
		Stage string      `json:"stage"`
		Type  string      `json:"type"`
		Group interface{} `json:"group"`
		Table []struct {
			Position int `json:"position"`
			Team     struct {
				ID        int    `json:"id"`
				Name      string `json:"name"`
				ShortName string `json:"shortName"`
				Tla       string `json:"tla"`
				Crest     string `json:"crest"`
			} `json:"team"`
			PlayedGames    int         `json:"playedGames"`
			Form           interface{} `json:"form"`
			Won            int         `json:"won"`
			Draw           int         `json:"draw"`
			Lost           int         `json:"lost"`
			Points         int         `json:"points"`
			GoalsFor       int         `json:"goalsFor"`
			GoalsAgainst   int         `json:"goalsAgainst"`
			GoalDifference int         `json:"goalDifference"`
		} `json:"table"`
	} `json:"standings"`
}

var weatherHtmlBody string
var newsHTML string
var catHTML string
var premTableHTML string

var moduleFunctions = map[string]func(string) string{
	"news":    getNews, // returns HTML string
	"weather": getWeather,
	"sports":  getleagueTable,
	"cats":    getCatImage,
	// etc.
}

type Option struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type Module struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	NeedsInput  bool     `json:"needsInput"`
	InputType   string   `json:"inputType"` // "text" or "select"
	InputLabel  string   `json:"inputLabel"`
	Options     []Option `json:"options"` // only used if InputType == "select"
}

var Modules = []Module{
	{
		ID:          "news",
		Title:       "News",
		Description: "Top 10 news stories.",
		Image:       "/images/NewsIcon.jpg",
		NeedsInput:  true,
		InputType:   "select",
		InputLabel:  "Select your category",
		Options: []Option{
			{Value: "general", Label: "General"},
			{Value: "world", Label: "World"},
			{Value: "nation", Label: "National"},
			{Value: "business", Label: "Business"},
			{Value: "technology", Label: "Technology"},
			{Value: "entertainment", Label: "Entertainment"},
			{Value: "sports", Label: "Sports"},
			{Value: "science", Label: "Science"},
			{Value: "health", Label: "Health"},
		},
	},
	{
		ID:          "sports",
		Title:       "League Table",
		Description: "Choose your football league.",
		Image:       "/images/league.jpg",
		NeedsInput:  true,
		InputType:   "select",
		InputLabel:  "Select your league",
		Options: []Option{
			{Value: "PL", Label: "Premier League"},
			{Value: "ELC", Label: "Championship"},
			{Value: "EL1", Label: "League One"},
			{Value: "EL2", Label: "League Two"},
			{Value: "PD", Label: "La Liga"},
			{Value: "SA", Label: "Serie A"},
			{Value: "BL1", Label: "Bundesliga"},
			{Value: "FL1", Label: "Ligue 1"},
			{Value: "DED", Label: "Eredivisie"},
			{Value: "PPL", Label: "Primeira Liga"},
			{Value: "BSA", Label: "Brasileirão"},
			{Value: "ARGPD", Label: "Argentine Primera"},
			{Value: "MLS", Label: "MLS"},
		},
	},
	{
		ID:          "weather",
		Title:       "Weather",
		Description: "Hourly forecast for your city.",
		Image:       "/images/Weather.jpg",
		NeedsInput:  true,
		InputType:   "text",
		InputLabel:  "Enter Postcode",
	},
	{
		ID:          "cats",
		Title:       "Cats",
		Description: "Random daily cat image or gif.",
		Image:       "/images/Cats.jpg",
	},
}

func (cfg *apiConfig) handleGetModulesMeta(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Modules)
}

func (cfg *apiConfig) handlerGetModuleExtraData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, `{"error": "Missing Authorization header"}`, http.StatusUnauthorized)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, `{"error": "Invalid Authorization header format"}`, http.StatusUnauthorized)
		return
	}

	tokenString := parts[1]

	prefs, err := cfg.db.GetPreferencesByToken(r.Context(), tokenString)
	if err != nil {
		http.Error(w, `{"error": "Token not found in database"}`, http.StatusUnauthorized)
	}

	userPrefs := database.UserPreference{
		Preferences:         prefs.Preferences,
		PreferenceVariables: prefs.PreferenceVariables,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userPrefs)
}

func getleagueTable(league string) string {
	url := fmt.Sprintf("https://api.football-data.org/v4/competitions/%v/standings", league)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Request creation failed:", err)
		return ""
	}

	req.Header.Add("X-Auth-Token", premTableAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Bad status code:", resp.StatusCode)
		return ""
	}

	var premTable PremTable
	err = json.NewDecoder(resp.Body).Decode(&premTable)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(premTable.Competition.Name)

	// Start building HTML with league name
	var premTableHTML strings.Builder

	premTableHTML.WriteString(fmt.Sprintf(`
<div style="max-width: 100%%; overflow-x: auto; border: 1px solid #ccc; border-radius: 8px; padding: 10px; background-color: #ffffff; font-family: Arial, sans-serif; font-size: 13px; color: #333;">
  <h1 style="font-size: 18px; text-align: center;">%v</h1>
  <table cellpadding="6" cellspacing="0" border="0" style="border-collapse: collapse; min-width: 600px; width: 100%%;">
    <thead>
      <tr style="background-color: #f0f0f0;">
        <th align="left">#</th>
        <th align="left">Team</th>
        <th align="left">GP</th>
        <th align="left">Pts</th>
        <th align="left">W</th>
        <th align="left">L</th>
        <th align="left">GF</th>
        <th align="left">GA</th>
        <th align="left">GD</th>
      </tr>
    </thead>
    <tbody>
`, premTable.Competition.Name))

	// Iterate over standings and add rows
	for _, standing := range premTable.Standings {
		if standing.Type != "TOTAL" {
			continue // usually only want TOTAL standings
		}
		for _, team := range standing.Table {
			premTableHTML.WriteString(fmt.Sprintf(`
      <tr>
        <td>%d</td>
		<td><img src="%s" alt="Crest" style="width: 25px; height: 25px; vertical-align: middle;"></td>
        <td>%s</td>
        <td>%d</td>
        <td>%d</td>
        <td>%d</td>
        <td>%d</td>
        <td>%d</td>
        <td>%d</td>
        <td>%d</td>
      </tr>`,
				team.Position,
				team.Team.Crest,
				team.Team.Name,
				team.PlayedGames,
				team.Points,
				team.Won,
				team.Lost,
				team.GoalsFor,
				team.GoalsAgainst,
				team.GoalDifference,
			))
		}
	}

	premTableHTML.WriteString(`
    </tbody>
  </table>
</div>
`)

	return premTableHTML.String()
}

func getCatImage(_ string) string {
	url := fmt.Sprintf("https://api.thecatapi.com/v1/images/search?size=med&format=json&api_key=%v", catAPIKey)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(err)
		return ""
	}

	var catImage CatStruct
	err = json.NewDecoder(resp.Body).Decode(&catImage)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	catHTML = fmt.Sprintf(`<h3 style="color: white">Here's a cat!</h3>
    <img src="%v" alt="Image" style="border: 1px solid #ccc;
    border-radius: 8px; width: 100%%; max-width: 400px; height: auto;
    display: block; margin: 12px auto;">`, catImage[0].URL)

	return catHTML
}

func getWeather(location string) string {

	location = strings.ReplaceAll(location, " ", "")
	url := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%v&q=%v&days=1&aqi=no&alerts=no", weatherAPIKey, location)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(err)
		return ""
	}

	var forecast ForecastStruct
	err = json.NewDecoder(resp.Body).Decode(&forecast)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	weatherHtmlBody = fmt.Sprintf(`<h3 style="color: white">Todays weather for %v</h3><div style="overflow-x: auto; white-space: nowrap; padding: 6px;">`, forecast.Location.Name)

	for _, hour := range forecast.Forecast.Forecastday[0].Hour[:24] {
		card := fmt.Sprintf(`
			<div style="display: inline-block; vertical-align: top; min-width: 112px; max-width: 135px; margin-right: 6px; border: 1px solid #ddd; border-radius: 6px; padding: 9px; font-family: Arial, sans-serif; text-align: center;">
				<h2 style="margin: 0 0 9px 0; font-size: 12px; color: #333;">%s</h2>
				<img src="https:%s" alt="Weather Icon" style="width: 45px; height: 45px; object-fit: contain; margin-bottom: 6px;">
				<p style="margin: 3px 0; font-size: 10.5px; color: #555;">%s</p>
				<p style="margin: 3px 0; font-size: 10.5px; color: #555;">%.1f°C</p>
				<p style="margin: 3px 0; font-size: 10.5px; color: #555;">Humidity: %d%%</p>
			</div>
		`, hour.Time[len(hour.Time)-5:], hour.Condition.Icon, hour.Condition.Text, hour.TempC, hour.Humidity)

		weatherHtmlBody += card
	}

	weatherHtmlBody += `</div>`

	return weatherHtmlBody
}

func getNews(category string) string {
	url := fmt.Sprintf("https://gnews.io/api/v4/top-headlines?category=%v&lang=en&country=uk&max=10&apikey=%v", category, newsAPIKey)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Bad status News code:", resp.Status)
		return ""
	}

	var newsArticles NewsArticles
	err = json.NewDecoder(resp.Body).Decode(&newsArticles)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	newsHTML = fmt.Sprintf(`<h3 style="color: white">Top %v news</h3><div style="max-width: 375px; height: 600px; overflow-y: auto; padding: 8px; border: 1px solid #ccc; margin: 0 auto;">`, category)

	for _, article := range newsArticles.Articles {
		t, err := time.Parse(time.RFC3339, article.PublishedAt)
		var publishDate string
		if err != nil {
			fmt.Println("Error parsing date:", err)
			publishDate = article.PublishedAt // fallback to raw string
		} else {
			publishDate = t.Format("02/01/2006 15:04")
		}
		newsHTML += makeNewsCard(article.Title, publishDate, article.Image, article.Description, article.URL)
	}

	newsHTML += `</div>`
	return newsHTML
}
