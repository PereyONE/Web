package controller

import (
	"errors"
	"html/template"
	"net/http"
	"site/app/model"
)

var tmpl *template.Template

//LogOut
func Indexx(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/logout.html", "template/index.html")
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
	t.ExecuteTemplate(w, "layout", Count)
}

func Kasten(w http.ResponseWriter, r *http.Request) {
	Naturwissenschaften, _ := model.GetKastenByKat("Naturwissenschaft")
	Sprache, _ := model.GetKastenByKat("Sprache")
	Gesellschaft, _ := model.GetKastenByKat("Gesellschaft")
	Wirtschaft, _ := model.GetKastenByKat("Wirtschaft")
	Geisteswissenschaften, _ := model.GetKastenByKat("Geisteswissenschaften")
	Data := struct {
		Naturwissenschaften   []map[string]interface{}
		Sprache               []map[string]interface{}
		Gesellschaft          []map[string]interface{}
		Wirtschaft            []map[string]interface{}
		Geisteswissenschaften []map[string]interface{}
		PublicCount           int
		PrivateCount          int
	}{
		Naturwissenschaften:   Naturwissenschaften,
		Sprache:               Sprache,
		Gesellschaft:          Gesellschaft,
		Wirtschaft:            Wirtschaft,
		Geisteswissenschaften: Geisteswissenschaften,
	}
	t, _ := template.ParseFiles("template/logout.html", "template/karteikasten.html")
	t.ExecuteTemplate(w, "layout", Data)

}

func LogInKasten(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)
	var User, _ = model.GetUserByUsername(username)
	Naturwissenschaften, _ := model.GetKastenByKat("Naturwissenschaft")
	Sprache, _ := model.GetKastenByKat("Sprache")
	Gesellschaft, _ := model.GetKastenByKat("Gesellschaft")
	Wirtschaft, _ := model.GetKastenByKat("Wirtschaft")
	Geisteswissenschaften, _ := model.GetKastenByKat("Geisteswissenschaften")
	PublicCount := model.GetPublicKastenCount()
	PrivateCount := model.GetPrivateKastenCount(username)
	Data := struct {
		User                  model.User
		Naturwissenschaften   []map[string]interface{}
		Sprache               []map[string]interface{}
		Gesellschaft          []map[string]interface{}
		Wirtschaft            []map[string]interface{}
		Geisteswissenschaften []map[string]interface{}
		PublicCount           int
		PrivateCount          int
	}{
		User:                  User,
		Naturwissenschaften:   Naturwissenschaften,
		Sprache:               Sprache,
		Gesellschaft:          Gesellschaft,
		Wirtschaft:            Wirtschaft,
		Geisteswissenschaften: Geisteswissenschaften,
		PublicCount:           PublicCount,
		PrivateCount:          PrivateCount,
	}
	t, _ := template.ParseFiles("template/login.html", "template/karteikastenlog.html")
	t.ExecuteTemplate(w, "layout", Data)
}

func SortKasten(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)
	kategorie := r.FormValue("sort")
	var User, _ = model.GetUserByUsername(username)
	kasten, _ := model.GetKastenByKat(kategorie)
	PublicCount := model.GetPublicKastenCount()
	PrivateCount := model.GetPrivateKastenCount(username)
	Data := struct {
		User         model.User
		Kasten       []map[string]interface{}
		Kategorie    string
		PublicCount  int
		PrivateCount int
	}{
		User:         User,
		Kasten:       kasten,
		Kategorie:    kategorie,
		PublicCount:  PublicCount,
		PrivateCount: PrivateCount,
	}
	t, _ := template.ParseFiles("template/login.html", "template/sortedkasten.html")
	t.ExecuteTemplate(w, "layout", Data)

}

func MeineKasten(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)
	var User, _ = model.GetUserByUsername(username)
	UserKasten, _, LernKasten, _ := model.GetKastenByUser(username)
	PublicCount := model.GetPublicKastenCount()
	PrivateCount := model.GetPrivateKastenCount(username)
	Data := struct {
		User         model.User
		UserKasten   []map[string]interface{}
		LernKasten   []map[string]interface{}
		PublicCount  int
		PrivateCount int
	}{
		User:         User,
		UserKasten:   UserKasten,
		LernKasten:   LernKasten,
		PublicCount:  PublicCount,
		PrivateCount: PrivateCount,
	}
	t, _ := template.ParseFiles("template/login.html", "template/meinekasten.html")
	t.ExecuteTemplate(w, "layout", Data)
}

func Erstellen(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)
	var User, _ = model.GetUserByUsername(username)
	PublicCount := model.GetPublicKastenCount()
	PrivateCount := model.GetPrivateKastenCount(username)
	Data := struct {
		User         model.User
		PublicCount  int
		PrivateCount int
		Update       string
	}{
		User:         User,
		PublicCount:  PublicCount,
		PrivateCount: PrivateCount,
	}
	t, _ := template.ParseFiles("template/login.html", "template/erstellen.html")
	t.ExecuteTemplate(w, "layout", Data)
}

func KastenErstellen(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)
	var User, _ = model.GetUserByUsername(username)
	titel := r.FormValue("titel")
	subkategorie := r.FormValue("subkategorie")
	beschreibung := r.FormValue("beschreibung")
	sichtbarkeit := r.FormValue("sichtbarkeit")
	_ = model.KastenErstellen(username, titel, subkategorie, beschreibung, sichtbarkeit)
	PublicCount := model.GetPublicKastenCount()
	PrivateCount := model.GetPrivateKastenCount(username)
	kasten := model.GetKastenByTitel(titel)
	Data := struct {
		User         model.User
		Kasten       model.Kasten
		PublicCount  int
		PrivateCount int
		Erste        model.Karten
		Karten       []map[string]interface{}
		Aktiv        model.Karten
	}{
		User:         User,
		Kasten:       kasten,
		PublicCount:  PublicCount,
		PrivateCount: PrivateCount,
	}
	t, _ := template.ParseFiles("template/login.html", "template/erstellen2.html")
	t.ExecuteTemplate(w, "layout", Data)
}

func UpdateKasten(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)
	kastenid := r.FormValue("kastenid")
	updaten := r.FormValue("updaten")
	titel := r.FormValue("titel")
	subkategorie := r.FormValue("subkategorie")
	beschreibung := r.FormValue("beschreibung")
	sichtbarkeit := r.FormValue("sichtbarkeit")
	var User, _ = model.GetUserByUsername(username)
	Kasten, _ := model.GetKastenById(kastenid, username)
	PublicCount := model.GetPublicKastenCount()
	PrivateCount := model.GetPrivateKastenCount(username)
	if updaten != "" {
		Data := struct {
			User         model.User
			Kasten       model.Kasten
			PublicCount  int
			PrivateCount int
			Updaten      string
		}{
			Kasten:       Kasten,
			User:         User,
			PublicCount:  PublicCount,
			PrivateCount: PrivateCount,
			Updaten:      updaten,
		}
		t, _ := template.ParseFiles("template/login.html", "template/erstellen1.html")
		t.ExecuteTemplate(w, "layout", Data)
	} else {
		Kasten = model.UpdateKasten(kastenid, username, titel, subkategorie, beschreibung, sichtbarkeit)
		Karten, _ := model.GetKartenByKasten(Kasten.ID)
		Data := struct {
			User         model.User
			Kasten       model.Kasten
			PublicCount  int
			PrivateCount int
			Updaten      string
			Karten       []map[string]interface{}
			Erste        model.Karten
			Aktiv        model.Karten
		}{
			Kasten:       Kasten,
			User:         User,
			PublicCount:  PublicCount,
			PrivateCount: PrivateCount,
			Updaten:      updaten,
			Karten:       Karten,
		}
		ErsteKarte, _ := model.Map2Karte(Karten[0])
		Karten = model.Remove(Karten, 0)
		Data.Erste = ErsteKarte
		Data.Karten = Karten
		Data.Aktiv = ErsteKarte
		t, _ := template.ParseFiles("template/login.html", "template/erstellen2.html")
		t.ExecuteTemplate(w, "layout", Data)
	}
}

func DeleteKasten(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)
	var User, _ = model.GetUserByUsername(username)
	kastenid := r.FormValue("kastenid")
	model.DeleteKasten(kastenid)
	UserKasten, _, LernKasten, _ := model.GetKastenByUser(username)
	PublicCount := model.GetPublicKastenCount()
	PrivateCount := model.GetPrivateKastenCount(username)
	Data := struct {
		User         model.User
		UserKasten   []map[string]interface{}
		LernKasten   []map[string]interface{}
		PublicCount  int
		PrivateCount int
	}{
		User:         User,
		UserKasten:   UserKasten,
		LernKasten:   LernKasten,
		PublicCount:  PublicCount,
		PrivateCount: PrivateCount,
	}
	t, _ := template.ParseFiles("template/login.html", "template/meinekasten.html")
	t.ExecuteTemplate(w, "layout", Data)
}

func KartenErstellen(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)
	var User, _ = model.GetUserByUsername(username)
	titel := r.FormValue("titel")
	frage := r.FormValue("frage")
	antwort := r.FormValue("antwort")
	kastenid := r.FormValue("kastenid")
	karte := model.KarteErstellen(username, kastenid, titel, frage, antwort)
	Karten, _ := model.GetKartenByKasten(kastenid)
	kasten, _ := model.GetKastenById(kastenid, username)
	PublicCount := model.GetPublicKastenCount()
	PrivateCount := model.GetPrivateKastenCount(username)

	Data := struct {
		User         model.User
		Karten       []map[string]interface{}
		Kasten       model.Kasten
		Aktiv        model.Karten
		Erste        model.Karten
		PublicCount  int
		PrivateCount int
	}{
		User:         User,
		Kasten:       kasten,
		PublicCount:  PublicCount,
		PrivateCount: PrivateCount,
	}
	ErsteKarte, _ := model.Map2Karte(Karten[0])
	Karten = model.Remove(Karten, 0)
	Data.Erste = ErsteKarte
	Data.Karten = Karten
	Data.Aktiv = karte
	t, _ := template.ParseFiles("template/login.html", "template/erstellen2.html")
	t.ExecuteTemplate(w, "layout", Data)

}

func KarteUpdaten(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)
	var User, _ = model.GetUserByUsername(username)
	titel := r.FormValue("titel")
	frage := r.FormValue("frage")
	antwort := r.FormValue("antwort")
	kastenid := r.FormValue("kastenid")
	kartenid := r.FormValue("kartenid")
	karte := model.KarteUpdaten(kartenid, titel, frage, antwort)
	Karten, _ := model.GetKartenByKasten(kastenid)
	kasten, _ := model.GetKastenById(kastenid, username)
	PublicCount := model.GetPublicKastenCount()
	PrivateCount := model.GetPrivateKastenCount(username)

	Data := struct {
		User         model.User
		Karten       []map[string]interface{}
		Kasten       model.Kasten
		Aktiv        model.Karten
		Erste        model.Karten
		PublicCount  int
		PrivateCount int
	}{
		User:         User,
		Kasten:       kasten,
		PublicCount:  PublicCount,
		PrivateCount: PrivateCount,
	}
	ErsteKarte, _ := model.Map2Karte(Karten[0])
	Karten = model.Remove(Karten, 0)
	Data.Erste = ErsteKarte
	Data.Karten = Karten
	Data.Aktiv = karte
	t, _ := template.ParseFiles("template/login.html", "template/erstellen2.html")
	t.ExecuteTemplate(w, "layout", Data)
}

func DeleteKarte(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)
	var User, _ = model.GetUserByUsername(username)
	kastenid := r.FormValue("kastenid")
	kartenid := r.FormValue("kartenid")
	model.DeleteKarte(kartenid, kastenid, username)
	Karten, _ := model.GetKartenByKasten(kastenid)
	kasten, _ := model.GetKastenById(kastenid, username)
	PublicCount := model.GetPublicKastenCount()
	PrivateCount := model.GetPrivateKastenCount(username)

	Data := struct {
		User         model.User
		Karten       []map[string]interface{}
		Kasten       model.Kasten
		Aktiv        model.Karten
		Erste        model.Karten
		PublicCount  int
		PrivateCount int
	}{
		User:         User,
		Kasten:       kasten,
		PublicCount:  PublicCount,
		PrivateCount: PrivateCount,
	}
	ErsteKarte, _ := model.Map2Karte(Karten[0])
	Karten = model.Remove(Karten, 0)
	Data.Erste = ErsteKarte
	Data.Karten = Karten
	Data.Aktiv = ErsteKarte
	t, _ := template.ParseFiles("template/login.html", "template/erstellen2.html")
	t.ExecuteTemplate(w, "layout", Data)
}

func Bearbeiten(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/login.html", "template/erstellen2.html")
	t.ExecuteTemplate(w, "layout", nil)
}
func Anschauen(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)
	id := r.FormValue("kastenid")
	kartenid := r.FormValue("kartenid")
	bearbeiten := r.FormValue("bearbeiten")
	var User, _ = model.GetUserByUsername(username)
	Kasten, err := model.GetKastenById(id, username)
	PublicCount := model.GetPublicKastenCount()
	PrivateCount := model.GetPrivateKastenCount(username)
	if err != nil {
		Data := struct {
			Error        error
			User         model.User
			PublicCount  int
			PrivateCount int
		}{
			Error:        err,
			User:         User,
			PublicCount:  PublicCount,
			PrivateCount: PrivateCount,
		}
		t, _ := template.ParseFiles("template/login.html", "template/anschauen.html")
		t.ExecuteTemplate(w, "layout", Data)
	} else {
		Data := struct {
			User         model.User
			Erste        model.Karten
			Aktiv        model.Karten
			Kasten       model.Kasten
			Error        error
			Karten       []map[string]interface{}
			PublicCount  int
			PrivateCount int
		}{
			User:         User,
			Kasten:       Kasten,
			PublicCount:  PublicCount,
			PrivateCount: PrivateCount,
		}
		var Erste model.Karten
		var Aktiv model.Karten
		Karten, _ := model.GetKartenByKasten(id)
		if len(Karten) == 0 {
			err = errors.New("Es sind noch keine Karten für diesen Kasten verfügbar lege doch eine an")
		} else {
			ErsteKarte, _ := model.Map2Karte(Karten[0])
			Karten = model.Remove(Karten, 0)
			Erste = ErsteKarte
			if kartenid == "" {
				Aktiv = Erste
			} else if kartenid == "neue" {
				var neueKarte model.Karten
				neueKarte.ID = "newkarte"
				Aktiv = neueKarte
			} else {
				Aktiv, err = model.GetKarteById(Karten, kartenid)
			}
		}
		Data.Erste = Erste
		Data.Karten = Karten
		Data.Error = err
		Data.Aktiv = Aktiv

		if bearbeiten != "" {
			t, _ := template.ParseFiles("template/login.html", "template/erstellen2.html")
			t.ExecuteTemplate(w, "layout", Data)
		} else {
			t, _ := template.ParseFiles("template/login.html", "template/anschauen.html")
			t.ExecuteTemplate(w, "layout", Data)
		}
	}
}

func Lernen(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)
	user, _ := model.GetUserByUsername(username)
	id := r.FormValue("kastenid")
	kartenid := r.FormValue("kartenid")
	var karte model.Karten
	kasten, _ := model.GetKastenById(id, username)
	PublicCount := model.GetPublicKastenCount()
	PrivateCount := model.GetPrivateKastenCount(username)
	var err error
	Data := struct {
		User         model.User
		Kasten       model.Kasten
		Karte        model.Karten
		Error        error
		PublicCount  int
		PrivateCount int
	}{
		User:         user,
		Kasten:       kasten,
		Error:        err,
		PublicCount:  PublicCount,
		PrivateCount: PrivateCount,
	}
	if kartenid == "" {
		if kasten.Ersteller != username {
			model.AddUserAsLerner(username, kasten)
		}
		karte, err = model.GetRandomKarte(id)
		Data.Karte = karte
		Data.Error = err
		t, _ := template.ParseFiles("template/login.html", "template/lernen1.html")
		t.ExecuteTemplate(w, "layout", Data)
	} else {
		karte, err = model.GetKarteByIdOnly(kartenid)
		Data.Karte = karte
		Data.Error = err
		t, _ := template.ParseFiles("template/login.html", "template/lernen2.html")
		t.ExecuteTemplate(w, "layout", Data)
	}
}

func Test(w http.ResponseWriter, r *http.Request) {
	//t, _ := template.ParseFiles("template/test.html")
	//t.Execute(w, )

}
