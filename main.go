package main

import (
	"time"
	"fmt"
	"net/http"
	"github.com/joho/godotenv"
	"os"
	"encoding/json"
	"github.com/distatus/battery"
	"github.com/robfig/cron/v3"
	"log"
	mail "github.com/xhit/go-simple-mail/v2"
	_ "github.com/lib/pq"
	"database/sql"
	"workspace/github.com/Benjysparks/daily/internal/database"
	"context"
	"strings"
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

type Preferences struct {
		RawMessage []string `json:"RawMessage"`
	} 


var (
	weatherAPIKey 		string
	newsAPIKey   		string
	catAPIKey	  		string
	emailAddress 	 	string
	emailPassword 		string
	premTableAPIKey 	string
)

type apiConfig struct {
	db			   *database.Queries
	cron 		   *cron.Cron
}

var (
    jwtSecret   = []byte(os.Getenv("JWT_SECRET"))
    tokenExpiry = time.Hour * 24
)



var weatherHtmlBody string
var newsHTML		string 
var catHTML			string
var premTableHTML	string

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	weatherAPIKey = os.Getenv("APIKey")
	emailAddress = os.Getenv("ServerEmail")
	emailPassword = os.Getenv("ServerPassword")
	newsAPIKey = os.Getenv("NewsAPIKey")
	premTableAPIKey = os.Getenv("PremTableAPIKey")
}

var moduleFunctions = map[string]func() string {
  "news": 		getNews,    // returns HTML string
  "weather": 	getWeather,
  "sports": 	getPremTable,
  "cats":		getCatImage,
  // etc.
}


func getPremTable() string {
	url := "https://api.football-data.org/v4/competitions/PL/standings"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Request creation failed:", err)
		fmt.Println(err)
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
		fmt.Println(err)
		return ""
	}

	var premTable PremTable
	err = json.NewDecoder(resp.Body).Decode(&premTable)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	premTableHTML = `<div style="max-width: 100%; overflow-x: auto; border: 1px solid #ccc; border-radius: 8px; padding: 10px; background-color: #ffffff; font-family: Arial, sans-serif; font-size: 13px; color: #333;">
  	<h1 style="font-size: 18px; text-align: center;">Premier League Table</h1>
	<table cellpadding="6" cellspacing="0" border="0" style="border-collapse: collapse; min-width: 600px; width: 100%;">
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
 `

	for _, standing := range premTable.Standings {
		if standing.Type != "TOTAL" {
			continue
		}
		for _, team := range standing.Table {
			teamRow := fmt.Sprintf(`<tr>
		<td>%d</td>
		<td><img src="%s" alt="" width="16" height="16" style="vertical-align: middle; margin-right: 4px; border-radius: 2px;"> %s</td>
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
				team.Team.ShortName,
				team.PlayedGames,
				team.Points,
				team.Won,
				team.Lost,
				team.GoalsFor,
				team.GoalsAgainst,
				team.GoalDifference,
			)
			premTableHTML += teamRow
		}
	}

	premTableHTML += `</tbody></table></div>`
	fmt.Println("working")


	return premTableHTML
}


func getCatImage() string {
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

func getWeather() string {
	
	url := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%v&q=tn27%%209fj&days=1&aqi=no&alerts=no", weatherAPIKey)

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

	weatherHtmlBody = `<h3 style="color: white">Daily Weather</h3><div style="overflow-x: auto; white-space: nowrap; padding: 6px;">`

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

func makeNewsCard(title, publishDate, imageURL, description, url string) string {
    return fmt.Sprintf(`
    <a href="%s" target="_blank" style="text-decoration: none; color: inherit;">
      <div style="border: 1px solid #ddd; border-radius: 8px; padding: 16px; font-family: Arial, sans-serif; max-width: 350px; display: flex; flex-direction: column; cursor: pointer; margin-bottom: 16px;">
        
        <h2 style="margin: 0 0 12px 0; font-size: 18px; color: #222; line-height: 1.2;">
          %s
        </h2>
        
        <small style="color: #777; margin-bottom: 8px; font-size: 12px;">
          Published: %s
        </small>
        
        <div style="display: flex; gap: 12px; align-items: flex-start;">
          <img src="%s" alt="Article Image" style="width: 120px; height: auto; border-radius: 4px; flex-shrink: 0; object-fit: cover;">
          
          <p style="margin: 0; color: #444; font-size: 14px; line-height: 1.4; flex-grow: 1;">
            %s
          </p>
        </div>
        
      </div>
    </a>
    `, url, title, publishDate, imageURL, description)
}

func getNews() string {
	url := fmt.Sprintf("https://gnews.io/api/v4/search?q=example&lang=en&country=us&max=10&apikey=%v", newsAPIKey)

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


	newsHTML = `<h3 style="color: white">Top news</h3><div style="max-width: 375px; height: 600px; overflow-y: auto; padding: 8px; border: 1px solid #ccc; margin: 0 auto;">`

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

func ordinalSuffix(day int) string {
	if day >= 11 && day <= 13 {
		return fmt.Sprintf("%dth", day)
	}
	switch day % 10 {
	case 1:
		return fmt.Sprintf("%dst", day)
	case 2:
		return fmt.Sprintf("%dnd", day)
	case 3:
		return fmt.Sprintf("%drd", day)
	default:
		return fmt.Sprintf("%dth", day)
	}
}

func dateStringFunc() string {
	now := time.Now()
	weekday := now.Weekday().String()                 
	dayWithSuffix := ordinalSuffix(now.Day())        
	month := now.Month().String()                   
	year := now.Year()                               

	return fmt.Sprintf("%s %s of %s %d", weekday, dayWithSuffix, month, year)
}

func emailer(modulePreferences string) {
	
	server := mail.NewSMTPClient()
	server.Host = "smtp.gmail.com"
	server.Port = 587
	server.Username = emailAddress
	server.Password = emailPassword
	server.Encryption = mail.EncryptionSTARTTLS
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}

	email := mail.NewMSG()
	email.SetFrom("Daily Bot <your_email@gmail.com>").
		AddTo("benmatthews92@hotmail.co.uk").
		SetSubject("Your Daily Read!")

	
	body := fmt.Sprintf(`
		<body bgcolor="#f0f0f0" style="margin:0; padding:0;">
		<h1 style="color: yellow">Good Morning!</h1>
		<h2 style="color: white">Todays date is %v</h2>
		%v	
	`, dateStringFunc(), modulePreferences)
	email.SetBody(mail.TextHTML, body)

	if email.Error != nil {
		fmt.Println("Email setup error:", email.Error)
		return
	}

	err = email.Send(smtpClient)
	if err != nil {
		fmt.Println("Error sending email:", err)
	} else {
		fmt.Println("✅ Email sent successfully!")
	}
}

func UnpackUserPreferences(row database.ShowUserPreferencesByEmailRow) ([]string) {

    if !row.Preferences.Valid {
        return []string{}
    }

    var prefs []string
    err := json.Unmarshal(row.Preferences.RawMessage, &prefs)
    if err != nil {
        return []string{}
    }
	return prefs
	
}


func (cfg *apiConfig) emailerPrefs(email string) {

	user, _ := cfg.db.ShowUserPreferencesByEmail(context.Background(), email)

	userPrefs := UnpackUserPreferences(user)

	var combinedHTML strings.Builder

	for _, moduleName := range userPrefs {
		fmt.Println(moduleName)
		fn, ok := moduleFunctions[moduleName]
		if !ok {
			continue
		}
		htmlPart := fn()        // call the function to get HTML
		combinedHTML.WriteString(htmlPart) // add it to the builder
	}

	emailBody := combinedHTML.String()

	emailer(emailBody)
}

func scheduleDailyTask(hour, minute int) {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
		if next.Before(now) {
			next = next.Add(24 * time.Hour)
		}
		duration := next.Sub(now)
		fmt.Printf("Waiting %v until next run...\n", duration)
		time.Sleep(duration)


		fmt.Println("Running scheduled task at", time.Now())


		time.Sleep(24 * time.Hour)
	}
}

func BatteryCheck() {
    batteries, err := battery.GetAll()
    if err != nil {
        log.Fatalf("Failed to get battery info: %v", err)
    }

    for i, bat := range batteries {
        percent := (bat.Current / bat.Full) * 100
        fmt.Printf("Battery #%d: %.2f%% (State: %s)\n", i, percent, bat.State)
    }
}

func (cfg *apiConfig) scheduleUserEmail(email string, hour, minute int32) error {
	fmt.Println("function called")
	// Validate hour and minute
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return fmt.Errorf("invalid hour or minute")
	}

	// Format cron spec: "30 8 * * *"
	spec := fmt.Sprintf("%d %d * * *", minute, hour)
	fmt.Printf("%v: %v", email, spec)

	// Add function to cron
	_, err := cfg.cron.AddFunc(spec, func() {
	cfg.emailerPrefs(email)
	})
		return err

}


func main() {

	

	

	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Print("Cound not open connection to database")
	}
	dbQueries := database.New(db)
	c := cron.New()

	apiCfg := apiConfig{
		db:				dbQueries,
		cron: 			c,
	}
	
	// apiCfg.emailerPrefs("Benjy@hotmail.com")

	mux.Handle("/", http.FileServer(http.Dir(filepathRoot + "/html")))
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("GET /api/allusers", apiCfg.handlerShowAllUser)
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        http.ServeFile(w, r, filepathRoot+"/html/login/login.html")
    case http.MethodPost:
        apiCfg.handlerLogin(w, r) // or jwtMiddleware(apiCfg.handlerLogin)(w, r) if needed
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
	})
	mux.HandleFunc("POST /api/preferences", apiCfg.handlerUpdatePreferences)
	mux.HandleFunc("GET /api/showpreferences", apiCfg.handlerShowUserPreferences)
	mux.HandleFunc("GET /api/userinfo", apiCfg.handlerSendInfoToFront)
	mux.HandleFunc("GET /api/clearuser", apiCfg.handlerClearUsers)



	go func() {
		log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
		log.Fatal(srv.ListenAndServe())
	}()

	// c.AddFunc("20 18 * * *", createEmail(variable))
	// c.AddFunc("@hourly", BatteryCheck)

	c.Start()
	select {}
}
