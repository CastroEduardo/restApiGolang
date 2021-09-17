package controller

import (
	"encoding/json"
	"net/http"
	"os"
	"rest-api-golang/src/models/authinterfaces"

	u "rest-api-golang/src/utils"

	"github.com/dgrijalva/jwt-go"
)

var Login = func(w http.ResponseWriter, r *http.Request) {
	user := &authinterfaces.LoginUser{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	user.Password = ""
	tk := &authinterfaces.Token{IdSession: user.User}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_HASH")))

	resp := u.Message(true, "Successful")
	resp["data"] = user
	resp["token"] = tokenString
	u.Respond(w, resp)
}

var GetUser = func(w http.ResponseWriter, r *http.Request) {

	// profileService.Demo()
	// usersservice.Demo()
	// profile := &models.Profile{Name: "Juan", Age: 50, LastUpdated: time.Now(), Password: "passwrod"}
	// dbContext.InsertNew(*profile)
	//context := usersservice.GetList()

	// //user := dbContext.GetList()
	//fmt.Println(context)

	// for i, item := range context {
	// 	fmt.Println(" Aqui :"+item.Name, strconv.Itoa(i))
	// }

	// h := sha256.Sum256([]byte("demo"))
	// h.Write([]byte("demo"))
	// b := h.Sum(nil)
	// fmt.Println(h)

	// user1 := &authinterfaces.User{}
	// user2 := &authinterfaces.User{}
	// user1.Email = "cao.trung.thu@mail.com"
	// user2.Email = "cao.trung.thu@hot.com"

	// var users [2]*authinterfaces.User
	// users[0] = user1
	// users[1] = user2
	// resp := u.Message(true, "Successful")
	// resp["data"] = users
	// u.Respond(w, resp)

	//user := &models.LoginUser{}
	// err := json.NewDecoder(r.Body).Decode(user)
	// if err != nil {
	// 	u.Respond(w, u.Message(false, "Invalid request"))
	// 	return
	// }

	// user.Password = "25554"
	// user.User = "Castro23R554"
	// tk := &authinterfaces.Token{IdSession: user.User}
	// token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	// tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_HASH")))

	resp := u.Message(true, "Successful")
	//resp["data"] = "SD"
	// resp["token"] = "UUU"

	u.Respond(w, resp)
}
