package router

import (
	"net/http"
	"site/app/controller"

	"github.com/gorilla/mux"
)

func GetRouter() *mux.Router {

	r := mux.NewRouter()

	r.HandleFunc("/", controller.Indexx)
	r.HandleFunc("/kasten", controller.Kasten)
	r.HandleFunc("/register", controller.Register)

	r.HandleFunc("/LogIn", controller.LogIn)
	r.HandleFunc("/LogIn/kasten", controller.LogInKasten)
	r.HandleFunc("/meinekasten", controller.MeineKasten)
	r.HandleFunc("/profil", controller.Profil)
	r.HandleFunc("/erstellen", controller.Erstellen)
	r.HandleFunc("/bearbeiten", controller.Bearbeiten)
	r.HandleFunc("/anschauen", controller.Anschauen)
	r.HandleFunc("/lernen", controller.Lernen)
	r.HandleFunc("/registerIn", controller.AddUser)
	r.HandleFunc("/loginyo", controller.AuthenticateUser)
	r.HandleFunc("/updateprofile", controller.UpdateProfile)
	r.HandleFunc("/Logout", controller.Logout)
	r.HandleFunc("/SortKasten", controller.SortKasten)
	r.HandleFunc("/kastenerstellen", controller.KastenErstellen)
	r.HandleFunc("/karteerstellen", controller.KartenErstellen)
	r.HandleFunc("/karteupdaten", controller.KarteUpdaten)
	r.HandleFunc("/deletekarte", controller.DeleteKarte)
	r.HandleFunc("/deletekasten", controller.DeleteKasten)
	r.HandleFunc("/deleteuser", controller.DeleteUser)
	r.HandleFunc("/updatekasten", controller.UpdateKasten)
	r.HandleFunc("/upload", controller.Upload)

	r.HandleFunc("/test", controller.Test)

	r.PathPrefix("/css/").Handler(http.StripPrefix("/css", http.FileServer(http.Dir("./static/css"))))
	r.PathPrefix("/favicon/").Handler(http.StripPrefix("/favicon", http.FileServer(http.Dir("./static/favicon"))))
	r.PathPrefix("/font/").Handler(http.StripPrefix("/font", http.FileServer(http.Dir("./static/font"))))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js", http.FileServer(http.Dir("./static/js"))))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images", http.FileServer(http.Dir("./static/images"))))

	return r
}
