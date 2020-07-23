package auth

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/casbin/casbin"
	"github.com/nydan/casbin/model"
)

const role = "role"

// Authorizer is middleware for authorization
func Authorizer(sesMngr *scs.SessionManager, e *casbin.Enforcer, users model.Users) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "role", r.Header.Get("role")))
			log.Println("context: ", r.Context())
			role := sesMngr.GetString(r.Context(), "role")

			log.Println("role: ", role)
			log.Println("url: ", r.URL.Path)
			log.Println("method: ", r.Method)

			if role == "" {
				role = "anonymous"
			}

			if role == "member" {
				uid := sesMngr.GetInt(r.Context(), "userID")
				if !users.Exist(uid) {
					writeError(http.StatusForbidden, "Forbidden", w, errors.New("user doesn't exist"))
					return
				}
			}

			if e.Enforce(role, r.URL.Path, r.Method) {
				next.ServeHTTP(w, r)
			} else {
				writeError(http.StatusForbidden, "Forbidden", w, errors.New("unauthorized"))
			}
		}
		return http.HandlerFunc(fn)
	}
}

func writeError(status int, message string, w http.ResponseWriter, err error) {
	log.Print("Error: ", err.Error())
	w.WriteHeader(status)
	w.Write([]byte(message))
}
