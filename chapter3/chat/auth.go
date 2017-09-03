package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("auth"); err == http.ErrNoCookie || cookie.Value == "" {
		// not authenticated yet
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// something wrong
		panic(err.Error())
	} else {
		// successful, call next wrapped handler
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// Path style : /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	if len(segs) == 4 {
		action := segs[2]
		provider := segs[3]
		switch action {
		case "login":
			provider, err := gomniauth.Provider(provider)
			if err != nil {
				log.Fatalln("Failed to get auth provider:", provider, "-", err)
			}
			loginURL, err := provider.GetBeginAuthURL(nil, nil)
			if err != nil {
				log.Fatalln("Error occurs in calling GetBeginAuthURL:", provider, "-", err)
			}
			w.Header().Set("Location", loginURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
		case "callback":
			provider, err := gomniauth.Provider(provider)
			if err != nil {
				log.Fatalln("Failed to get provider:", provider, "-", err)
			}
			creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
			if err != nil {
				log.Fatalln("Could not finish authentication:", provider, "-", err)
			}
			user, err := provider.GetUser(creds)
			if err != nil {
				log.Fatalln("Failed to get user data:", provider, "-", err)
			}
			authCookieValue := objx.New(map[string]interface{}{
				"name":       user.Name(),
				"avatar_url": user.AvatarURL(),
				"email":      user.Email(),
			}).MustBase64()
			http.SetCookie(w, &http.Cookie{
				Name:  "auth",
				Value: authCookieValue,
				Path:  "/"})
			w.Header()["Location"] = []string{"/chat"}
			w.WriteHeader(http.StatusTemporaryRedirect)
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Not supported action: %s", action)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not supported request: %s", r.URL.Path)
	}
}
