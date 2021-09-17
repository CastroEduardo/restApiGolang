package settingscontroller

import (
	"encoding/json"
	"net/http"
	"os"
	"reflect"
	"rest-api-golang/src/dbContext/companyservice"
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
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

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

	checkUser := usersservice.CheckUserPasswordForEmail(user.User, user.Password)

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
			expireDate = now

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
			logUser := "Inicio Session .."
			logsuserservice.Add(1, checkUser.ID, checkUser.IdCompany, logUser)

			msg = "Usuario Logeado "

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

var LoginUser = func(w http.ResponseWriter, r *http.Request) {

	// tokenHeader := r.Header.Get("Authorization") //Grab the token from the header
	// idSession := sessionuserservice.GetClaimForToken(tokenHeader)
	// fmt.Println(idSession.User.NickName)

	// list := logsuserservice.GetList()

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
	resp["data"] = true
	resp["token"] = ""
	u.Respond(w, resp)
}

type Bird struct {
	Id   string
	Name string
	Id2  string
}

func Filter(arr interface{}, cond func(interface{}) bool) interface{} {
	contentType := reflect.TypeOf(arr)
	contentValue := reflect.ValueOf(arr)

	newContent := reflect.MakeSlice(contentType, 0, 0)
	for i := 0; i < contentValue.Len(); i++ {
		if content := contentValue.Index(i); cond(content.Interface()) {
			newContent = reflect.Append(newContent, content)
		}
	}
	return newContent.Interface()
}

func isNoBarAndLessThanTenChar(a models.GenericList, value1 string) bool {
	return !strings.HasPrefix(a.Value1, value1)
}

var GetUser = func(w http.ResponseWriter, r *http.Request) {

	// listProvinces := genericlistservice.GetListForIdCompanyAndIdentity("", "provinces")
	// list := genericlistservice.GetListForIdCompanyAndIdentity("", "municipies")

	// for _, s := range listProvinces {
	// 	fmt.Printf("%s ==> %s %s \n", s.IdKey, s.Name, s.Value1)

	// 	for _, m := range list {
	// 		if m.Value1 == s.IdKey {
	// 			fmt.Println("************ " + m.Name)
	// 		}
	// 	}
	// 	// result := Choose(list, isNoBarAndLessThanTenChar)
	// 	// fmt.Println(result) // [foo_super]

	// }

	// var bird []Bird

	//  `[{"id":"1","name":"Name","id2":"4"},
	// {"id":"1","Name":"Name","Id2":"4"}]`

	// json.Unmarshal([]byte(myJsonString), &bird)

	// for _, s := range bird {
	// 	fmt.Printf("%s ==> %s %s \n", s.Id, s.Name, s.Id2)

	// 	newGeneric := models.GenericList{
	// 		IdKey:     s.Id,
	// 		Name:      s.Name,
	// 		Date:      time.Now(),
	// 		Identity:  "municipies",
	// 		Status:    1,
	// 		Value1:    s.Id2,
	// 		IdCompany: "",
	// 		Note:      "Municipie " + s.Name + " " + " IdProvince is Value1 => " + s.Id2,
	// 	}

	// 	genericlistservice.Add(newGeneric)

	//
	// fmt.Println(result)

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

	// newUser := models.IUser{
	// 	DateAdd:   time.Now(),
	// 	IdCompany: "55555555",
	// 	IdRol:     "144444",
	// 	Image:     "imagen",
	// 	LastLogin: time.Now(),
	// 	LastName:  "apellido",
	// 	Name:      "NOmbre",
	// 	Password:  utils.Encript([]byte("231154")),
	// 	Status:    1,
	// 	NickName:  "usuario1",
	// }

	// result2 := usersservice.Add(newUser)
	// fmt.Println(result2)

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
	resp["data"] = true
	//resp["token"] = tokenString

	u.Respond(w, resp)
}

var AddCompany = func(w http.ResponseWriter, r *http.Request) {

	//AddCompany
	newCompany := authinterfaces.Company{
		Address:   "Direccion /c larga #21 demo",
		DateAdd:   time.Now(),
		Image:     "logodemo.png",
		NameLong:  "Nombre EMpresa largo",
		NameShort: "Nombre Corto",
		Other:     "Otras Configuraciones",
		Phone:     "809-521-2144 / 20-52222",
		Rnc:       "004-251111-2",
		Slogan:    "Slogan de Empresa..",
		Status:    1,
	}
	idCompany := companyservice.Add(newCompany)

	//create rol User
	newRolUser := authinterfaces.RolUser{
		IdCompany: idCompany,
		Date:      time.Now(),
		Name:      "Administradores",
		Note:      "Todos los Privilegios",
		Status:    1,
	}
	idRolUser := rolesuserservice.Add(newRolUser)

	//add PrivilegeRol
	newPrivilege := authinterfaces.PrivilegeRolUser{
		IdCompany: idCompany,
		IdRol:     idRolUser,
		WebAccess: true,
		Config:    1,
		TypeUser:  1,
	}
	privilegeroluserservice.Add(newPrivilege)

	//add new User
	newUser := authinterfaces.User{
		IdCompany: idCompany,
		DateAdd:   time.Now(),
		City:      "Santo Domingo",
		Gender:    "0",
		Contact:   "809-545-5444",
		IdRol:     idRolUser,
		Image:     "user.png",
		LastLogin: time.Now(),
		LastName:  "Apellido del Usuario",
		Name:      "Nombre del Usuario",
		NickName:  strings.ToLower("castro2354"),
		Password:  utils.Encript([]byte("231154")),
		ForcePass: true,
		Public:    0,
		Status:    1,
		Email:     "castro@gmail.com",
		Note:      "Alguna Nota Larga Para el Usuario --> Para describir algo",
	}
	usersservice.Add(newUser)

	//add  logs Systems
	logssystemservice.Add(1, "Agrego Nueva Empresa..: "+idCompany)

	resp := u.Message(true, "Successful")
	resp["data"] = "{}"
	//resp["token"] = tokenString

	u.Respond(w, resp)
}

//Send Config General System
var SendListGeneral = func(w http.ResponseWriter, r *http.Request) {

	requestData := &models.RequestListGeneral{}
	err := json.NewDecoder(r.Body).Decode(requestData)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	listSend := usersservice.GetListFromIdCompany("asdsad")

	resp := u.Message(false, "ok")
	resp["data"] = listSend

	// resp["provincies"] = listProvinces
	// resp["municipies"] = listMunicipies
	u.Respond(w, resp)
}
