package authcontroller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	genericlistservice "rest-api-golang/src/dbContext/genericListService"
	"rest-api-golang/src/dbContext/logssystemservice"
	"rest-api-golang/src/dbContext/logsuserservice"
	"rest-api-golang/src/dbContext/privilegeroluserservice"
	"rest-api-golang/src/dbContext/rolesuserservice"
	"rest-api-golang/src/dbContext/sessionuserservice"
	"rest-api-golang/src/dbContext/usersservice"
	"rest-api-golang/src/models"
	"rest-api-golang/src/models/authinterfaces"
	"rest-api-golang/src/utils"
	u "rest-api-golang/src/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
)

///#######Funtional

// Login ... This function is just to handle the incoming request
var Login = func(w http.ResponseWriter, r *http.Request) {

	var tokenString string
	var expireDate time.Time
	var msg string

	user := &authinterfaces.LoginUser{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	//check UserFirts
	checkUser := usersservice.CheckUserPasswordForUser(user.User, user.Password)

	if checkUser.NickName == "" {
		//validate user and password Email
		checkUser = usersservice.CheckUserPasswordForEmail(user.User, user.Password)
	}

	// fmt.Printf("User %s \n", checkUser.Name)
	if checkUser.NickName != "" {
		if checkUser.Status == 1 {

			//generate session User
			newSession := authinterfaces.SessionUser{
				Token:     "",
				Active:    true,
				DateAdd:   time.Now(),
				IdCompany: checkUser.IdCompany,
				IdUser:    checkUser.ID,
				Remember:  user.Remember,
			}

			getIdSession := sessionuserservice.Add(newSession)

			//generate token
			user.Password = ""
			tk := &authinterfaces.Token{IdSession: getIdSession}
			token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
			generateToken, _ := token.SignedString([]byte(os.Getenv("TOKEN_HASH")))
			tokenString = generateToken

			//define expire token
			now := time.Now()

			if user.Remember {
				expireDate = now.AddDate(0, 0, 1)
			} else {
				expireDate = now.Add(2 * time.Minute)
			}

			//update IdSession
			sessionUser := sessionuserservice.FindToId(getIdSession)
			sessionUser.Token = tokenString
			sessionUser.TokenExpire = expireDate
			sessionuserservice.UpdateOne(sessionUser)

			//update LasLoginUser
			userFound := usersservice.FindToId(checkUser.ID)
			userFound.LastLogin = time.Now()
			usersservice.UpdateLastLogin(userFound)

			logUser := "Inicio Session .."
			logsuserservice.Add(1, checkUser.ID, checkUser.IdCompany, logUser)

			msg = "Usuario Logeado"

		} else {

			logUser := "Intento Loggin User Desactivado .."
			logsuserservice.Add(1, checkUser.ID, checkUser.IdCompany, logUser)
			msg = "Usuario Desactivado.."
		}

	} else {
		logSystem := "Intento Login Fallido usuario : " + user.User + "  "
		logssystemservice.Add(3, logSystem)
		msg = "Usuario Invalido"
	}

	//fmt.Println(checkUser)
	resp := u.Message(true, msg)
	resp["token"] = tokenString
	resp["expire"] = expireDate
	resp["msg"] = msg
	u.Respond(w, resp)
}

//send Details Login User
var GetDetailsLogin = func(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("getDetailsLogin")
	tokenHeader := r.Header.Get("Authorization") //Grab the token from the header
	claimSession := sessionuserservice.GetClaimForToken(tokenHeader)
	//fmt.Println(claimSession.User)

	//list := logsuserservice.GetList()
	//fmt.Println(checkUser)
	resp := u.Message(true, "ok")
	resp["user"] = claimSession.User
	resp["company"] = claimSession.Company
	resp["privilegesUser"] = claimSession.PrivilegeRolUser

	u.Respond(w, resp)
}

//Logout Session userlogin
var LogoutUser = func(w http.ResponseWriter, r *http.Request) {
	//get user on request
	tokenHeader := r.Header.Get("Authorization") //Grab the token from the header
	//fmt.Println(tokenHeader)
	if tokenHeader != "Bearer" {
		//if exist Token
		claim := sessionuserservice.GetClaimForToken(tokenHeader)
		sessionuserservice.LogoutSessionToIdUser(claim.User.ID)

		logUser := "Salio del Sistema."
		logsuserservice.Add(1, claim.User.ID, claim.Company.ID, logUser)
	}

	resp := u.Message(true, "ok")
	u.Respond(w, resp)
}

//checkSatus Token
var CheckStatusToken = func(w http.ResponseWriter, r *http.Request) {

	tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

	idSession := sessionuserservice.GetIdSessionToToken(tokenHeader) //sessionuserservice.LogoutSessionToIdUser(modelRequest.IdUser)
	session := sessionuserservice.FindToId(idSession)
	resp := u.Message(false, "ok")
	resp["result"] = session.Active
	u.Respond(w, resp)
}

//Update Profile UserSystem
var UpdateProfileUser = func(w http.ResponseWriter, r *http.Request) {

	//get user on request
	tokenHeader := r.Header.Get("Authorization") //Grab the token from the header
	claim := sessionuserservice.GetClaimForToken(tokenHeader)

	user := &authinterfaces.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		fmt.Println("fallo")
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	//find user to Update
	userUpdate := usersservice.FindToId(user.ID)

	userUpdate.Name = user.Name
	userUpdate.LastName = user.LastName
	userUpdate.Email = user.Email
	userUpdate.City = user.City
	userUpdate.Gender = user.Gender
	userUpdate.Contact = user.Contact
	userUpdate.Note = user.Note
	//update fields
	result := usersservice.UpdateOne1(userUpdate)

	if result {
		logSystem := "IdUser:" + claim.User.ID + "Actualizo Perfil de usuarioId: " + user.NickName + "  "
		logssystemservice.Add(3, logSystem)
	}

	//fmt.Println(result)

	// Do something with the Person struct...
	//fmt.Fprintf(w, "Person: %+v", user)

	resp := u.Message(false, "ok")
	resp["result"] = result
	u.Respond(w, resp)
}

//Update Password User
var UpdatePasswordUser = func(w http.ResponseWriter, r *http.Request) {

	//get user on request
	tokenHeader := r.Header.Get("Authorization") //Grab the token from the header
	claim := sessionuserservice.GetClaimForToken(tokenHeader)

	requestData := &authinterfaces.RequestUpdatePassUser{}
	err := json.NewDecoder(r.Body).Decode(requestData)
	if err != nil {
		fmt.Println("fallo")
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	//find user to Update Pass
	userUpdate := usersservice.FindToId(requestData.ID)

	var result bool = false
	//if Is Old Pass
	oldPass := utils.Descrypt(userUpdate.Password, []byte(requestData.OldPassword))
	if oldPass {
		newPass := utils.Encript([]byte(requestData.NewPassword))
		//update Password User
		userUpdate.Password = newPass
		result = usersservice.UpdateOne1(userUpdate)
		//remove session User Update
		sessionuserservice.LogoutSessionToIdUser(userUpdate.ID)
	}

	if result {
		logSystem := "IdUser:" + claim.User.ID + "Actualizo Password De usuario: " + userUpdate.ID + "  "
		logssystemservice.Add(3, logSystem)
	}
	resp := u.Message(false, "ok")
	resp["result"] = result
	u.Respond(w, resp)
}

//Send Config General System
var SendConfigSystem = func(w http.ResponseWriter, r *http.Request) {

	//get user on request
	tokenHeader := r.Header.Get("Authorization") //Grab the token from the header
	claim := sessionuserservice.GetClaimForToken(tokenHeader)

	listRoles := rolesuserservice.GetListForIdCompany(claim.Company.ID)
	listPrivilegesUser := privilegeroluserservice.GetListForIdCompany(claim.Company.ID)
	listProvinces := genericlistservice.GetListForIdCompanyAndIdentity("", "provinces")
	listMunicipies := genericlistservice.GetListForIdCompanyAndIdentity("", "municipies")

	// list1, _ := json.Marshal(listRoles)
	// list2, _ := json.Marshal(listPrivilegesUser)
	//jsonString := `{"name":` + string(list2) + `,"name2":` + string(list1) + `}`

	resp := u.Message(false, "ok")
	resp["roles"] = listRoles
	resp["privileges"] = listPrivilegesUser
	resp["provincies"] = listProvinces
	resp["municipies"] = listMunicipies
	u.Respond(w, resp)
}

//Send listUsers of System
var SendListUsers = func(w http.ResponseWriter, r *http.Request) {

	//get user on request
	tokenHeader := r.Header.Get("Authorization") //Grab the token from the header
	claim := sessionuserservice.GetClaimForToken(tokenHeader)

	listUsers := usersservice.GetListFromIdCompany(claim.Company.ID)

	resp := u.Message(false, "ok")
	resp["data"] = listUsers

	u.Respond(w, resp)
}

//Send list Config Settings
var SendListConfigSetting = func(w http.ResponseWriter, r *http.Request) {

	//get user on request
	tokenHeader := r.Header.Get("Authorization") //Grab the token from the header
	claim := sessionuserservice.GetClaimForToken(tokenHeader)

	params := &models.RequestGeneric{}
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	var idCompany = ""

	if params.Params2 != "" {
		idCompany = claim.Company.ID
	}

	var listSend []byte

	//fmt.Println(claim.Company.ID)
	switch params.Params1 {
	case "roles":
		listRolesUsers := rolesuserservice.GetListForIdCompany(claim.Company.ID)

		//convert list Roles to Generic List to Send
		var newList = []models.GenericList{}
		for _, item := range listRolesUsers {
			newItem := models.GenericList{
				ID:     "",
				Name:   item.Name,
				IdKey:  item.ID,
				Value1: item.ID,
				Status: item.Status,
			}
			newList = append(newList, newItem)
		}

		getList, _ := json.Marshal(newList)
		listSend = getList
	case "":
		break
	default:
		listGeneric := genericlistservice.GetListForIdCompanyAndIdentity(idCompany, params.Params1)
		getList, _ := json.Marshal(listGeneric)
		listSend = getList

	}

	resp := u.Message(false, "ok")
	resp["data"] = string(listSend)

	u.Respond(w, resp)
}

///###########Funtional

// LoginUser ... This function is just to handle the incoming request
var LoginUser = func(w http.ResponseWriter, r *http.Request) {
	fmt.Println("is Here")
	tokenHeader := r.Header.Get("Authorization") //Grab the token from the header
	idSession := sessionuserservice.GetClaimForToken(tokenHeader)
	fmt.Println("IdSession: " + idSession.User.NickName)
	list := logsuserservice.GetList()
	// emp := list
	// e, err := json.Marshal(emp)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(string(e))

	// err := json.Unmarshal([]byte(list), &models.LogUser)
	// if err != nil {

	// }

	// // jsonInfo, _ := json.Marshal(list)

	// fmt.Println(string(list))

	// newSession := authinterfaces.SessionUser{
	// 	IdUser:      "5e7a3d93a248a6e1c5d6698b",
	// 	Active:      true,
	// 	DateAdd:     time.Now(),
	// 	IdCompany:   "asdasd",
	// 	Token:       "asdsadsadsadsda",
	// 	TokenExpire: time.Now(),
	// 	Remember:    true,
	// }

	// sessionuserservice.Add(newSession)

	//fmt.Println(result)
	// user.Password = ""
	// tk := &authinterfaces.Token{UserId: user.Email}
	// token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	// tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_HASH")))

	resp := u.Message(true, "Successful")
	resp["data"] = list
	resp["token"] = ""
	u.Respond(w, resp)
}

// GetUser ... This function is just to handle the incoming request
var GetUser = func(w http.ResponseWriter, r *http.Request) {

	// newCompany := authinterfaces.Company{
	// 	Address:   "Direccion ",
	// 	DateAdd:   time.Now(),
	// 	Image:     "logo.png",
	// 	NameLong:  "Nombre Largo Empresa",
	// 	NameShort: "Nombre Corto",
	// 	Other:     "Otros Datos",
	// 	Phone:     "809-561-2512 / 809-245-5444",
	// 	Rnc:       "001-0215211-0",
	// 	Slogan:    "Slogan Company",
	// 	Status:    1,
	// }

	// result := companyservice.Add(newCompany)

	// fmt.Println(result)

	// for i := 0; i <= 5000; i++ {
	// 	newLogSystem := models.LogSystem{
	// 		Log:    "Update :" + strconv.Itoa(i),
	// 		Level:  1,
	// 		Status: 1,
	// 		Date:   time.Now(),
	// 	}

	// 	logssystemservice.Add(newLogSystem)
	// 	//time.Sleep(5 * time.Second)
	// }

	//list := logssystemservice.GetList()

	// list := usersservice.GetList()

	// for i, item := range list {
	// 	fmt.Println(item.Log, strconv.Itoa(i))
	// }
	// result := usersservice.FindToId("5e795d655d554045401496e6")
	// result.NickName = "ADMIN23"
	// fmt.Println(usersservice.UpdateOne(result))
	// fmt.Println(result.ID)

	newUser := models.IUserLogin{
		DateAdd:   time.Now(),
		IdCompany: "55555555",
		IdRol:     "144444",
		Image:     "imagen",
		LastLogin: time.Now(),
		LastName:  "apellido",
		Name:      "NOmbre",
		Password:  utils.Encript([]byte("231154")),
		Status:    1,
		NickName:  "usuario1",
	}

	result2 := usersservice.Add(newUser)
	fmt.Println(result2)

	// profileService.Demo()
	// usersservice.Demo()

	// profile := &models.Profile{Name: "Juan", Age: 50, LastUpdated: time.Now(), Password: "passwrod"}
	// dbContext.InsertNew(*profile)
	// context := usersservice.GetList()
	// // //user := dbContext.GetList()
	// //fmt.Println(context)

	// for i, item := range context {
	// 	fmt.Println(item.Name, strconv.Itoa(i))
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

	user := &authinterfaces.LoginUser{}
	// err := json.NewDecoder(r.Body).Decode(user)
	// if err != nil {
	// 	u.Respond(w, u.Message(false, "Invalid request"))
	// 	return
	// }

	user.Password = "25554"
	user.User = "Castro2354"
	//tk := &authinterfaces.Token{UserId: user.Email}
	//token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	//tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_HASH")))

	resp := u.Message(true, "Successful")
	resp["data"] = "{}"
	//resp["token"] = tokenString

	u.Respond(w, resp)
}
