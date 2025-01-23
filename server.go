package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/adilanuarbek/form-service/controllers"
	"github.com/adilanuarbek/form-service/models"
	"github.com/adilanuarbek/form-service/templates"
	"github.com/adilanuarbek/form-service/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

type config struct {
	PSQL       models.PostgresConfig
	SessionKey string
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		fmt.Printf("error loading env config: %v\n", err)
	}
	err = run(cfg)
	if err != nil {
		fmt.Printf("error running the server: %v\n", err)
	}
}

func run(cfg config) error {
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		return err
	}
	defer db.Close()

	store := sessions.NewCookieStore([]byte(cfg.SessionKey))
	store.Options.HttpOnly = true
	store.Options.SameSite = http.SameSiteLaxMode

	// Services
	userService := &models.UserService{
		DB: db,
	}
	formService := &models.FormService{
		DB: db,
	}

	// Controllers
	handlersC := controllers.Handlers{
		FormService: formService,
	}
	usersC := controllers.Users{
		UserService: userService,
		Session:     store,
	}
	umw := controllers.UserMiddleware{
		UserService: userService,
		Session:     store,
	}

	r := chi.NewRouter()
	r.Use(umw.SetUser)
	// Home and contact pages
	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "home.html", "tailwind.html"))))
	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "contact.html", "tailwind.html"))))

	// Sign up and Sign in
	usersC.Templates.Signup = views.Must(views.ParseFS(templates.FS, "signup.html", "tailwind.html"))
	r.Get("/signup", usersC.SignUp)
	r.Post("/signup", usersC.SignUpHandler)

	usersC.Templates.Signin = views.Must(views.ParseFS(templates.FS, "signin.html", "tailwind.html"))
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.SignInHandler)
	r.Post("/signout", usersC.SignOutHandler)

	// Dashboard
	r.Get("/dashboard", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "dashboard.html", "tailwind.html"))))
	handlersC.Templates.Profile = views.Must(views.ParseFS(templates.FS, "profile.html", "tailwind.html"))
	r.Get("/profile", handlersC.ProfileHandler)
	r.Post("/create-form", handlersC.PostDashboardHandler)
	handlersC.Templates.ViewForm = views.Must(views.ParseFS(templates.FS, "viewform.html", "tailwind.html"))
	r.Get("/form/{formSTR}", handlersC.ViewFormData)

	fmt.Println("Starting the server on 8080...")
	return http.ListenAndServe(":8080", r)
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	}

	if cfg.PSQL.Host == "" && cfg.PSQL.Port == "" {
		return cfg, fmt.Errorf("no PSQL config provided")
	}

	cfg.SessionKey = os.Getenv("SESSION_KEY")

	return cfg, nil
}
