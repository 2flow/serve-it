package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func allowOriginFunc(r *http.Request, origin string) bool {
	if origin == "http://example.com" {
		return true
	}
	return false
}

func main() {
	var path string

	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		path = "./"
	}

	r := chi.NewRouter()
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc:  allowOriginFunc,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	//http.FileServer(http.Dir("./"))
	r.Handle("/*", handler())
	fmt.Println("Serve: " + path)
	http.ListenAndServe(":5001", r)
}

// SetContentType Setst the content type to the responsewriter by the filename or the content
func SetContentType(w http.ResponseWriter, file io.Reader, fileName string) {

	contentType := ""

	if strings.HasSuffix(fileName, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(fileName, ".js") {
		contentType = "application/javascript"
	} else if (strings.HasSuffix(fileName, ".ico")) || (strings.HasSuffix(fileName, ".cur")) {
		contentType = "image/x-icon"
	} else if strings.HasSuffix(fileName, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(fileName, ".html") {
		contentType = "text/html"
	} else if strings.HasSuffix(fileName, ".svg") {
		contentType = "image/svg+xml"
	}

	if contentType != "" {
		w.Header().Add("Content-Type", contentType)
	}
}

func handler() http.Handler {
	return http.HandlerFunc(serveDirectory)
}

// ServeFromAzureDirectory will us the filestorage to display the content
func serveDirectory(w http.ResponseWriter, r *http.Request) {
	//fetchDest := r.Header.Get("Sec-Fetch-Dest") // if index.html --> document otherwise script
	chiCtx := chi.RouteContext(r.Context())
	routPath := chiCtx.URLParams.Values[0]

	// // if the document is requested return the index.html
	// // this should work
	// if fetchDest == "document" {
	// 	routPath = "/index.html"
	// } else if fetchDest == "" {
	// 	parts := strings.Split(routPath, ".")
	// 	if len(parts) == 1 {
	// 		routPath = "/index.html"
	// 	}
	// }
	localPath := "./" + routPath

	reader, err := os.Open(localPath)

	if err != nil {
		io.WriteString(w, "NotFound")
		return
	}
	defer reader.Close()

	SetContentType(w, reader, routPath)
	w.WriteHeader(http.StatusOK)
	io.Copy(w, reader)
}
