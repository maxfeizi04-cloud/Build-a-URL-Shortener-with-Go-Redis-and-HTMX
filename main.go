package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"url_shortener/utils"
)

//go:embed templates/*.html
var templatesFS embed.FS

var ctx = context.Background()

func main() {
	// We create the DB connection here and use it in the handlers
	// 我们在此处创建数据库连接，并在处理器中使用它。
	dbClient := utils.NewRedisClient()
	if dbClient == nil {
		fmt.Println("Failed to connect to Redis")
		return
	}

	// we use the http.HandleFunc to define our routes
	// and their corresponding handler functions
	// 我们使用 HandleFunc 函数来定义我们的路由及其对应的处理函数
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, nil)
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
	})

	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		// Shorten the provided URL, store it and return it to our UI
		// 缩短提供的 URL,存储它并返回我们的 UI
		// Get the URL to shorten from the request
		// 从请求中获取要缩短的URL
		url := r.FormValue("url")

		// Close the boby when done
		fmt.Println("Payload", url)

		// Shorten the URL
		// 缩短 URL
		shortURL := utils.GetShortCode()
		fullShortURL := fmt.Sprintf("MY_HOST/r/%s", shortURL)

		// Generated short URL
		// 生成短链接
		// 输出到控制台 log to console
		fmt.Printf("Generated short URL: %s\n", fullShortURL)

		// @TODO: Store {shortcode: url} in Redis
		// 缩短的 URL 存入到 Redis中
		utils.SetKey(&ctx, dbClient, shortURL, url, 0)

		// @TODO return the shortened URL in the UI
		// 短 URL 返回到 用户 UI
		fmt.Fprintf(w, `<p class="mt-4 text-green-600">Shortened URL: <a 
			href="/r/%s" class="underline">%s</a></p>`, shortURL, fullShortURL)
	})

	http.HandleFunc("/r/{code}", func(w http.ResponseWriter, r *http.Request) {
		key := r.FormValue("code")
		if key == "" {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		longURL, err := utils.GetLongURL(&ctx, dbClient, key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Redirect(w, r, longURL, http.StatusPermanentRedirect)
	})

	// Start the server on port 8080
	// 在端口 8080 上启动服务
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
