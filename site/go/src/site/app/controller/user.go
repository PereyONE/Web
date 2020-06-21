package controller

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"io/ioutil"
	"net/http"
	"site/app/model"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store *sessions.CookieStore

func init() {
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key := make([]byte, 32)
	rand.Read(key)
	store = sessions.NewCookieStore(key)
}

// Register controller
func Register(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/logout.html", "template/register.html")
	t.ExecuteTemplate(w, "layout", nil)
}

// AddUser controller
func AddUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	mail := r.FormValue("mail")
	var Data int

	user := model.User{}
	user.Username = username
	user.Password = password
	user.Email = mail

	if model.UserExist(username) == true && model.EmailExist(mail) == true {
		Data = 1
		t, _ := template.ParseFiles("template/logout.html", "template/register.html")
		t.ExecuteTemplate(w, "layout", Data)
	} else if model.UserExist(username) == true && model.EmailExist(mail) != true {
		Data = 2
		t, _ := template.ParseFiles("template/logout.html", "template/register.html")
		t.ExecuteTemplate(w, "layout", Data)
	} else if model.UserExist(username) != true && model.EmailExist(mail) == true {
		Data = 3
		t, _ := template.ParseFiles("template/logout.html", "template/register.html")
		t.ExecuteTemplate(w, "layout", Data)
	} else {

		err := user.Add()
		if err != nil {
			data := struct {
				ErrorMsg string
				Type     bool
			}{
				ErrorMsg: "Username already exists!",
				Type:     true,
			}
			t, _ := template.ParseFiles("template/logout.html", "template/register.html")
			t.ExecuteTemplate(w, "layout", data)
		} else {
			t, _ := template.ParseFiles("template/logout.html", "template/success.html")
			t.ExecuteTemplate(w, "layout", user)
		}
	}
}

// Login controller
func Login(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "login.tmpl", nil)
}

// AuthenticateUser controller
func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var err error
	var user = model.User{}
	var data = struct {
		ErrorMsg string
		User     model.User
	}{
		ErrorMsg: "Username and/or password wrong!",
		User:     user,
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Authentication
	user, err = model.GetUserByUsername(username)
	if err == nil {
		// decode base64 String to []byte
		passwordDB, _ := base64.StdEncoding.DecodeString(user.Password)
		err = bcrypt.CompareHashAndPassword(passwordDB, []byte(password))
		if err == nil {
			session, _ := store.Get(r, "session")

			// Set user as authenticated
			session.Values["authenticated"] = true
			session.Values["username"] = username
			session.Save(r, w)
			http.Redirect(w, r, "/LogIn", http.StatusFound)
		} else {
			//falsches Passwort
			t, _ := template.ParseFiles("template/logout.html", "template/loginfail.html")
			t.ExecuteTemplate(w, "layout", data)
		}
	} else {
		//User nicht gefunden
		t, _ := template.ParseFiles("template/logout.html", "template/loginfail.html")
		t.ExecuteTemplate(w, "layout", data)
	}
}

// Logout controller
func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	session.Values["authenticated"] = false
	session.Values["username"] = ""
	session.Save(r, w)
	CountUser := model.GetUserCount()
	CountKasten := model.GetKastenCount()
	CountKarten := model.GetKartenCount()
	Count := struct {
		Users  int
		Karten int
		Kasten int
	}{
		Users:  CountUser,
		Karten: CountKarten,
		Kasten: CountKasten,
	}

	t, _ := template.ParseFiles("template/logout.html", "template/index.html")
	t.ExecuteTemplate(w, "layout", Count)
}
func LogIn(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)
	var User, err = model.GetUserByUsername(username)
	PublicCount := model.GetPublicKastenCount()
	PrivateCount := model.GetPrivateKastenCount(username)
	if err != nil {
		t, _ := template.ParseFiles("template/logout.html", "template/index.html")
		t.ExecuteTemplate(w, "layout", nil)
	} else {
		Data := struct {
			User         model.User
			PublicCount  int
			PrivateCount int
		}{
			User:         User,
			PublicCount:  PublicCount,
			PrivateCount: PrivateCount,
		}
		t, _ := template.ParseFiles("template/login.html", "template/indexlog.html")
		t.ExecuteTemplate(w, "layout", Data)
	}
}
func Upload(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)
	img, _, _ := r.FormFile("foto")
	defer img.Close()
	tempFile, _ := ioutil.TempFile("static/images", "upload-*.png")
	defer tempFile.Close()

	fileBytes, _ := ioutil.ReadAll(img)
	tempFile.Write(fileBytes)
	filename := tempFile.Name()
	model.UpdateFoto(username, filename)

	var User, _ = model.GetUserByUsername(username)
	PublicCount := model.GetPublicKastenCount()
	PrivateCount := model.GetPrivateKastenCount(username)
	KartenCount := model.GetKartenCountForUser(username)
	Data := struct {
		User         model.User
		PublicCount  int
		PrivateCount int
		KartenCount  int
		Fail         int
	}{
		User:         User,
		PublicCount:  PublicCount,
		PrivateCount: PrivateCount,
		KartenCount:  KartenCount,
		Fail:         0,
	}
	t, _ := template.ParseFiles("template/login.html", "template/profil.html")
	t.ExecuteTemplate(w, "layout", Data)

}

func Profil(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)
	var User, _ = model.GetUserByUsername(username)
	PublicCount := model.GetPublicKastenCount()
	PrivateCount := model.GetPrivateKastenCount(username)
	KartenCount := model.GetKartenCountForUser(username)
	Data := struct {
		User         model.User
		PublicCount  int
		PrivateCount int
		KartenCount  int
		Fail         int
	}{
		User:         User,
		PublicCount:  PublicCount,
		PrivateCount: PrivateCount,
		KartenCount:  KartenCount,
		Fail:         0,
	}
	t, _ := template.ParseFiles("template/login.html", "template/profil.html")
	t.ExecuteTemplate(w, "layout", Data)
}
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)
	oldpassword := r.FormValue("oldpassword")
	newpassword := r.FormValue("newpassword")
	User, _ := model.GetUserByUsername(username)
	PublicCount := model.GetPublicKastenCount()
	PrivateCount := model.GetPrivateKastenCount(username)
	KartenCount := model.GetKartenCountForUser(username)
	Data := struct {
		User         model.User
		Fail         int
		PublicCount  int
		PrivateCount int
		KartenCount  int
	}{
		User:         User,
		Fail:         0,
		PublicCount:  PublicCount,
		PrivateCount: PrivateCount,
		KartenCount:  KartenCount,
	}

	if model.CheckPassword(username, oldpassword) != true {
		Data.Fail = 1
		t, _ := template.ParseFiles("template/login.html", "template/profil.html")
		t.ExecuteTemplate(w, "layout", Data)
	} else {
		err := model.UpdatePassword(username, newpassword)
		User, _ := model.GetUserByUsername(username)
		if err != nil {
			Data.Fail = 2
			t, _ := template.ParseFiles("template/login.html", "template/profil.html")
			t.ExecuteTemplate(w, "layout", Data)

		} else {
			Data.User = User
			t, _ := template.ParseFiles("template/login.html", "template/profil.html")
			t.ExecuteTemplate(w, "layout", Data)
		}
	}

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)
	model.DeleteUser(username)

	CountUser := model.GetUserCount()
	CountKasten := model.GetKastenCount()
	CountKarten := model.GetKartenCount()
	Count := struct {
		Users  int
		Karten int
		Kasten int
	}{
		Users:  CountUser,
		Karten: CountKarten,
		Kasten: CountKasten,
	}
	t, _ := template.ParseFiles("template/logout.html", "template/index.html")
	t.ExecuteTemplate(w, "layout", Count)
}

// Auth is an authentication handler
func Auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")

		// Check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			h(w, r)
		}
	}
}
