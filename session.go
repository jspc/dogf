package main

import (
	"net/http"

	"github.com/satori/go.uuid"
)

const (
	CookieName = "dogfucker"
)

type Cookie struct {
	Session string
}

func (rr Router) IsLoggedIn(r *http.Request) (loggedin bool, uuid string, err error) {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return
	}

	s, err := rr.decode(cookie)
	if err != nil {
		return
	}

	loggedin = true
	uuid = s

	return
}

func (rr Router) DoLogin(u, p string) (c *http.Cookie, err error) {
	err = rr.DB.AreValidCredentials(u, p)
	if err != nil {
		return
	}

	encoded, err := rr.encode(uuid.NewV4().String())
	if err != nil {
		return
	}

	c = &http.Cookie{
		Name:  CookieName,
		Value: encoded,
		Path:  "/",
	}

	return
}

func (rr Router) encode(uuid string) (string, error) {
	return rr.SecureCookie.Encode(CookieName, Cookie{uuid})
}

func (rr Router) decode(cookie *http.Cookie) (uuid string, err error) {
	c := &Cookie{}
	err = rr.SecureCookie.Decode(CookieName, cookie.Value, c)
	if err != nil {
		return
	}
	uuid = c.Session

	return
}
