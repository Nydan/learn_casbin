package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/casbin/casbin"
	"github.com/nydan/casbin/auth"
	"github.com/nydan/casbin/model"
)

func main() {
	// setup casbin auth rules
	authEnforce := casbin.NewEnforcer("./auth_model.conf", "./policy.csv")

	// setup in memory store engine for the session
	memStore := memstore.NewWithCleanupInterval(30 * time.Minute)

	// setup session manager
	sessionManager := scs.New()
	sessionManager.IdleTimeout = 30 * time.Minute
	sessionManager.Store = memStore

	users := createUsers()

	// setup handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/login", loginHandler(sessionManager, users))
	mux.HandleFunc("/member/current", currentMemberHandler(sessionManager))

	log.Print("Server started on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080",
		sessionManager.LoadAndSave(auth.Authorizer(sessionManager, authEnforce, users)(mux))))

}

func createUsers() model.Users {
	users := model.Users{}
	users = append(users, model.User{ID: 1, Name: "Admin", Role: "admin"})
	users = append(users, model.User{ID: 2, Name: "Member1", Role: "member"})
	users = append(users, model.User{ID: 3, Name: "Member2", Role: "member"})
	return users
}

func writeError(status int, message string, w http.ResponseWriter, err error) {
	log.Print("ERROR: ", err.Error())
	w.WriteHeader(status)
	w.Write([]byte(message))
}

func writeSuccess(message string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func loginHandler(ses *scs.SessionManager, users model.Users) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.PostFormValue("name")
		user, err := users.FindByName(name)
		if err != nil {
			writeError(http.StatusBadRequest, "WRONG_CREDENTIALS", w, err)
			return
		}
		// setup ession
		if err := ses.RenewToken(r.Context()); err != nil {
			writeError(http.StatusInternalServerError, "ERROR", w, err)
			return
		}
		ses.Put(r.Context(), "userID", user.ID)
		ses.Put(r.Context(), "role", user.Role)
		writeSuccess("SUCCESS", w)
	})
}

func currentMemberHandler(session *scs.SessionManager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid := session.GetInt(r.Context(), "userID")
		if uid == 0 {
			writeError(http.StatusInternalServerError, "ERROR", w, errors.New("member not found"))
			return
		}
		writeSuccess(fmt.Sprintf("User with ID: %d", uid), w)
	})
}
