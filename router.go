package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	cookie "github.com/gorilla/securecookie"
)

type Router struct {
	SecureCookie *cookie.SecureCookie
	DB           database
	Renderer     Renderer
}

func NewRouter(hashKey, blockKey []byte, db string) (rr Router, err error) {
	rr.SecureCookie = cookie.New(hashKey, blockKey)
	rr.DB, err = NewDatabase(db)
	if err != nil {
		return
	}

	rr.Renderer, err = NewRenderer()
	if err != nil {
		return
	}

	return
}

func setHeaders(w http.ResponseWriter, origin string) (wDup http.ResponseWriter) {
	wDup = w
	wDup.Header().Set("Access-Control-Allow-Headers", "requested-with, Content-Type, origin, authorization, accept, client-security-token, cache-control, x-api-key")
	wDup.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	wDup.Header().Set("Access-Control-Allow-Origin", origin)
	wDup.Header().Set("Access-Control-Allow-Credentials", "true")
	wDup.Header().Set("Access-Control-Max-Age", "10")
	wDup.Header().Set("Cache-Control", "no-cache")

	return
}

func (rr Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/login":
		rr.login(w, r)
	default:
		potentialPath := path.Join("public", r.URL.Path)
		_, err := os.Stat(potentialPath)
		if os.IsNotExist(err) {
			rr.fourohfour(w, r.URL)
		} else {
			http.ServeFile(w, r, potentialPath)
		}
	}
}

func (rr Router) login(w http.ResponseWriter, r *http.Request) {
	loggedin, _, _ := rr.IsLoggedIn(r) // treat errors as not logged in

	if r.Method == "GET" {
		pageData := struct {
			LoggedIn bool
		}{
			LoggedIn: loggedin,
		}

		body, err := rr.Renderer.Render("login", pageData)
		if err != nil {
			rr.fivehundred(w, err)

			return
		}
		fmt.Fprintf(w, string(body))

		return
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			rr.fivehundred(w, err)
		}

		c, err := rr.DoLogin(r.FormValue("username"), r.FormValue("password"))
		if err != nil {
			log.Println(err)
			rr.fourohone(w)

			return
		}

		http.SetCookie(w, c)
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func (rr Router) fourohfour(w http.ResponseWriter, u *url.URL) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "%q not found\n", u)
}

func (rr Router) fourohone(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "login failed\n")
}

func (rr Router) fivehundred(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)

	log.Print(err)
	fmt.Fprintf(w, "an error happened\n")
}
