package routes

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"rest-api-golang/src/controller"
	"rest-api-golang/src/controller/authcontroller"
	"rest-api-golang/src/controller/settingscontroller"
	"rest-api-golang/src/dbContext/logssystemservice"
	"rest-api-golang/src/dbContext/sessionuserservice"
	"rest-api-golang/src/models/authinterfaces"
	u "rest-api-golang/src/utils"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var SetupServer = func(APP_PORT string) {
	var router = mux.NewRouter()

	// We use our custom CORS Middleware
	router.Use(CORS)

	//
	router.HandleFunc("/api/login", controller.Login).Methods("POST")
	router.HandleFunc("/api/getUser", controller.GetUser).Methods("GET")

	//auth Controller .
	router.HandleFunc("/api/auth/login", authcontroller.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/getDetailsLogin", authcontroller.GetDetailsLogin).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/logoutUser", authcontroller.LogoutUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/checkStatusToken", authcontroller.CheckStatusToken).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/updateProfileUser", authcontroller.UpdateProfileUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/updatePasswordUser", authcontroller.UpdatePasswordUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/sendConfigGeneralSystem", authcontroller.SendConfigSystem).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/sendListUsers", authcontroller.SendListUsers).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/sendListConfigSetting", authcontroller.SendListConfigSetting).Methods("POST", "OPTIONS")

	///

	router.HandleFunc("/api/auth/loginUser", authcontroller.LoginUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/auth/getUser", authcontroller.GetUser).Methods("GET", "OPTIONS")
	// router.HandleFunc("/api/auth/addCompany", authcontroller.AddCompany).Methods("GET", "OPTIONS")

	//settings Controller  with use CORS IN WEBAPI
	router.HandleFunc("/api/settings/sendListGeneral", settingscontroller.SendListGeneral).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/settings/login", settingscontroller.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/settings/loginUser", settingscontroller.LoginUser).Methods("GET")
	router.HandleFunc("/api/settings/getUser", settingscontroller.GetUser).Methods("GET")
	router.HandleFunc("/api/settings/addCompany", settingscontroller.AddCompany).Methods("GET")

	router.Use(JwtAuthentication)

	err := http.ListenAndServe(":"+APP_PORT, router)

	if err != nil {
		fmt.Print(err)
	} else {

	}
}

// CORS Middleware
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Next
		next.ServeHTTP(w, r)
		return
	})
}

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//if is METHOD GET SEND handler not need token Auth
		if r.Method == http.MethodGet {
			fmt.Println("Is Get Method")
			next.ServeHTTP(w, r)
			return
		}

		notAuth := []string{"/api/auth/login", "/api/login",
			"/api/auth/logoutUser", "/api/auth/checkStatusToken", "/api/settings/login"}
		requestPath := r.URL.Path //current request path
		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			response = u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			response = u.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		tk := &authinterfaces.Token{}

		//fmt.Println(tokenPart)
		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_HASH")), nil
		})

		//check session db
		sessionDb := sessionuserservice.FindToId(tk.IdSession)
		//fmt.Println(sessionDb.Active)
		if sessionDb.IdCompany == "" {
			response = u.Message(false, "Not session in db")
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			logssystemservice.Add(3, "Fallo session no activa..")
			return
		} else {
			if sessionDb.Active == false {
				response = u.Message(false, "Session disable en DB")
				w.WriteHeader(http.StatusUnauthorized)
				w.Header().Add("Content-Type", "application/json")
				u.Respond(w, response)
				logssystemservice.Add(3, "Fallo session desactivada en DB")
				return
			}

			now := time.Now()
			//if session expire
			if sessionDb.TokenExpire.Before(now) {

				response = u.Message(false, "Session expirada por fecha ")
				w.WriteHeader(http.StatusUnauthorized)
				w.Header().Add("Content-Type", "application/json")
				u.Respond(w, response)
				logssystemservice.Add(3, "Fallo session desactivada en DB")

				//disabled session Db
				sessionDb.Active = false
				sessionDb.DateLogout = time.Now()
				sessionuserservice.UpdateOne(sessionDb)
				return
			}

			//diff := now.Sub(sessionDb.TokenExpire)
			// convert diff to days
			//days := int(diff.Hours() / 24)

		}

		if err != nil { //Malformed token, returns with http code 403 as usual
			response = u.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			response = u.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk.IdSession)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}
