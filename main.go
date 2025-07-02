package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"workspace/github.com/Benjysparks/daily/internal/database"

	"github.com/distatus/battery"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Preferences struct {
	RawMessage []string `json:"RawMessage"`
}

var (
	weatherAPIKey   string
	newsAPIKey      string
	catAPIKey       string
	emailAddress    string
	emailPassword   string
	premTableAPIKey string
)

type apiConfig struct {
	db   *database.Queries
	cron *cron.Cron
}

var (
	jwtSecret   = []byte(os.Getenv("JWT_SECRET"))
	tokenExpiry = time.Hour * 24
)

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

func UnpackUserPreferences(row database.ShowUserPreferencesByEmailRow) ([]string, []string) {
    var prefs []string
    var extraData []string

    if row.Preferences.Valid {
        if err := json.Unmarshal(row.Preferences.RawMessage, &prefs); err != nil {
            prefs = []string{}
        }
    }

    if row.PreferenceVariables.Valid {
        if err := json.Unmarshal(row.PreferenceVariables.RawMessage, &extraData); err != nil {
            extraData = []string{}
        }
    }

    return prefs, extraData
}


func (cfg *apiConfig) emailerPrefs(email string) {
    user, err := cfg.db.ShowUserPreferencesByEmail(context.Background(), email)
    if err != nil {
        // handle error, maybe log and return early
        return
    }

    userPrefs, extraData := UnpackUserPreferences(user)

    var combinedHTML strings.Builder

    for i, moduleName := range userPrefs {
        fn, ok := moduleFunctions[moduleName]
        if !ok {
            continue
        }

        // Defensive: check if extraData exists for this index
        var arg string
        if i < len(extraData) {
            arg = extraData[i]
        } else {
            arg = "" // or some default
        }

        htmlPart := fn(arg)                 // pass the corresponding extraData
        combinedHTML.WriteString(htmlPart) // add it to the builder
    }

    emailBody := combinedHTML.String()

    emailer(emailBody)
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
	const port = "8000"

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
		db:   dbQueries,
		cron: c,
	}

	mux.Handle("/", http.FileServer(http.Dir(filepathRoot+"/html")))
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("GET /api/allusers", apiCfg.handlerShowAllUser)
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			http.ServeFile(w, r, filepathRoot+"/html/login.html")
		case http.MethodPost:
			apiCfg.handlerLogin(w, r)
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

	c.Start()
	select {}
}
