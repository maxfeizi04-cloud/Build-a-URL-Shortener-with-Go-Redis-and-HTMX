package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"url_shortener/utils"
)

//go:embed templates/*.html
var templatesFS embed.FS

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseFS(templatesFS, "templates/*.html"))
}

func main() {
	dbClient := utils.NewRedisClient()
	ctx := context.Background()

	if err := dbClient.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	myHost := os.Getenv("MY_HOST")
	if myHost == "" {
		myHost = "http://localhost:8080"
	}

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Template execution error:", err)
		}
	})

	http.HandleFunc("POST /shorten", func(w http.ResponseWriter, r *http.Request) {
		url := r.FormValue("url")
		if url == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}
		log.Println("Payload:", url)

		shortURL := utils.GetShortCode()
		fullShortURL := fmt.Sprintf("%s/r/%s", myHost, shortURL)
		log.Printf("Generated short URL: %s", fullShortURL)

		if err := utils.SetKey(ctx, dbClient, shortURL, url, 0); err != nil {
			http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
			log.Println("Redis error:", err)
			return
		}

		fmt.Fprintf(w, `<p class="mt-4 text-green-600">Shortened URL: <a
			href="/r/%s" class="underline">%s</a></p>`, shortURL, fullShortURL)
	})

	http.HandleFunc("GET /r/{code}", func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		if code == "" {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		longURL, err := utils.GetLongURL(ctx, dbClient, code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Redirect(w, r, longURL, http.StatusPermanentRedirect)
	})

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
