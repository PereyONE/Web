package model

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Kasten struct {
	ID           string   `json:"_id"`
	Rev          string   `json:"_rev"`
	Kategorie    string   `json:"kategorie"`
	Subkategorie string   `json:"subkategorie"`
	Titel        string   `json:"titel"`
	Type         string   `json:"type"`
	Count        int      `json:"count"`
	Beschreibung string   `json:"beschreibung"`
	Sichtbarkeit bool     `json:"sichtbarkeit"`
	Ersteller    string   `json:"ersteller"`
	Lerner       []string `json:"lerner"`
	Done         bool     `json:"done"`
}

func GetKastenById(id string, username string) (Kasten Kasten, err error) {
	kasten, err := GetKasten(id)
	if kasten.Sichtbarkeit == false {
		if kasten.Ersteller != username {
			return kasten, errors.New("Kasten is not viewable")
		}
	}
	return kasten, err
}

func AddUserAsLerner(username string, kasten Kasten) (err error) {
	for i := 0; i < len(kasten.Lerner); i++ {
		if kasten.Lerner[i] == username {
			err = errors.New("Username schon drinne")
		} else {
			err = nil
		}
	}
	if err == nil {
		kasten.Lerner = append(kasten.Lerner, username)
		m, _ := kasten2Map(kasten)
		err = btDB.Set(kasten.ID, m)
	}
	return err
}

func UpdateKasten(kastenid string, username string, titel string, subkategorie string, beschreibung string, sichtbarkeit string) (kasten Kasten) {
	Kasten, _ := GetKastenById(kastenid, username)
	Kasten.Titel = titel
	Kasten.Subkategorie = subkategorie
	Kasten.Beschreibung = beschreibung
	if sichtbarkeit == "true" {
		Kasten.Sichtbarkeit = true
	} else {
		Kasten.Sichtbarkeit = false
	}

	if Kasten.Sichtbarkeit == false {
		Kasten.Lerner = nil
	}

	Naturwissenschaft := []string{"Biologie", "Chemie", "Elektrotechnik", "Informatik", "Mathematik", "Medizin", "Naturkunde", "Physik", "Sonstiges"}
	Sprachen := []string{"Chinesisch", "Deutsch", "Englisch", "Französisch", "Griechisch", "Italienisch", "Latein", "Russisch", "Sonstiges"}
	Gesellschaft := []string{"Ethik", "Geschichte", "Literatur", "Musik", "Politik", "Recht", "Soziales", "Sport", "Verkehrskunde", "Sonstiges"}
	Wirtschaft := []string{"BWL", "Finanzen", "Landwirtschaft", "Marketing", "VWL", "Sonstiges"}
	Geisteswissenschaften := []string{"Kriminologie", "Philosophie", "Psychologie", "Pädagogik", "Theologie", "Sonstiges"}

	if stringInSlice(Kasten.Subkategorie, Naturwissenschaft) == true {
		Kasten.Kategorie = "Naturwissenschaft"
	} else if stringInSlice(Kasten.Subkategorie, Sprachen) == true {
		Kasten.Kategorie = "Sprache"
	} else if stringInSlice(Kasten.Subkategorie, Gesellschaft) == true {
		Kasten.Kategorie = "Gesellschaft"
	} else if stringInSlice(Kasten.Subkategorie, Wirtschaft) == true {
		Kasten.Kategorie = "Wirtschaft"
	} else if stringInSlice(Kasten.Subkategorie, Geisteswissenschaften) == true {
		Kasten.Kategorie = "Geisteswissenschaften"
	}

	u, _ := kasten2Map(Kasten)
	_ = btDB.Set(Kasten.ID, u)

	kasten = Kasten
	return kasten
}

func GetKastenByTitel(titel string) Kasten {
	query := `
	{
		"selector": {
			"type": "kasten",
			"titel":"%s"
		 }
	 }`
	UserKasten, _ := btDB.QueryJSON(fmt.Sprintf(query, titel))
	kasten, _ := map2Kasten(UserKasten[0])
	return kasten
}

func GetAllKasten() ([]map[string]interface{}, error) {
	AllKasten, err := btDB.QueryJSON(`{
			"selector": {
			   "sichtbarkeit": true
			}
		}`)
	if err != nil {
		panic(err)
	} else {
		return AllKasten, nil
	}
}
func GetKasten(id string) (Kasten, error) {
	t, err := btDB.Get(id, nil)
	if err != nil {
		return Kasten{}, err
	}

	todo, err := map2Kasten(t)
	return todo, err
}

// GetAllTodosForUser retrieves Todos of provided username from DB
func GetKastenByUser(username string) (UserKasten []map[string]interface{}, err1 error, LernKasten []map[string]interface{}, err2 error) {
	query1 := `
	{
		"selector": {
			"type": "kasten",
			"ersteller": "%s"
		}
	 }`
	query2 := `
	{
		"selector": {
			"type": "kasten",
			"lerner": {
			   "$elemMatch": {
				  "$eq": "%s"
			   }
			}
		 }
	 }`
	UserKasten, err1 = btDB.QueryJSON(fmt.Sprintf(query1, username))
	if err1 != nil {
		err1 = errors.New("Da is war inne Friddn mit den User Karten")
	}
	LernKasten, err2 = btDB.QueryJSON(fmt.Sprintf(query2, username))
	if err2 != nil {
		err2 = errors.New("Da is war inne Friddn mit den Lern Karten")
	}
	return UserKasten, err1, LernKasten, err2
}
func GetKastenByKat(kategorie string) ([]map[string]interface{}, error) {
	query := `
	{
		"selector": {
			"sichtbarkeit": true,
			"kategorie": "%s"
		}
	 }`
	allKasten, err := btDB.QueryJSON(fmt.Sprintf(query, kategorie))

	if err != nil || len(allKasten) == 0 {
		allKasten = nil
		return allKasten, err
	}
	return allKasten, nil
}
func KastenErstellen(username string, titel string, kategorie string, beschreibung string, sichtbarkeit string) (Kasten Kasten) {
	Kasten.Subkategorie = kategorie
	Kasten.Titel = titel
	Kasten.Type = "kasten"
	Kasten.Count = 0
	Kasten.Beschreibung = beschreibung
	if sichtbarkeit == "true" {
		Kasten.Sichtbarkeit = true
	} else {
		Kasten.Sichtbarkeit = false
	}
	Kasten.Ersteller = username
	var lerner []string
	Kasten.Lerner = lerner
	Kasten.Done = false

	Naturwissenschaft := []string{"Biologie", "Chemie", "Elektrotechnik", "Informatik", "Mathematik", "Medizin", "Naturkunde", "Physik", "Sonstiges"}
	Sprachen := []string{"Chinesisch", "Deutsch", "Englisch", "Französisch", "Griechisch", "Italienisch", "Latein", "Russisch", "Sonstiges"}
	Gesellschaft := []string{"Ethik", "Geschichte", "Literatur", "Musik", "Politik", "Recht", "Soziales", "Sport", "Verkehrskunde", "Sonstiges"}
	Wirtschaft := []string{"BWL", "Finanzen", "Landwirtschaft", "Marketing", "VWL", "Sonstiges"}
	Geisteswissenschaften := []string{"Kriminologie", "Philosophie", "Psychologie", "Pädagogik", "Theologie", "Sonstiges"}

	if stringInSlice(Kasten.Subkategorie, Naturwissenschaft) == true {
		Kasten.Kategorie = "Naturwissenschaft"
	} else if stringInSlice(Kasten.Subkategorie, Sprachen) == true {
		Kasten.Kategorie = "Sprache"
	} else if stringInSlice(Kasten.Subkategorie, Gesellschaft) == true {
		Kasten.Kategorie = "Gesellschaft"
	} else if stringInSlice(Kasten.Subkategorie, Wirtschaft) == true {
		Kasten.Kategorie = "Wirtschaft"
	} else if stringInSlice(Kasten.Subkategorie, Geisteswissenschaften) == true {
		Kasten.Kategorie = "Geisteswissenschaften"
	}
	Kasten.Add()
	return Kasten
}

//Helperfunction für KastenErstellen
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (t Kasten) Add() error {
	kasten, _ := kasten2Map(t)

	delete(kasten, "_id")
	delete(kasten, "_rev")

	_, _, err := btDB.Save(kasten, nil)
	if err != nil {
		fmt.Printf("[Add] error : %s", err)
	}
	return err
}

func (t Kasten) ToggleStatus() error {
	if t.Done {
		t.Done = false
	} else {
		t.Done = true
	}
	u, _ := kasten2Map(t)
	err := btDB.Set(t.ID, u)

	if err != nil {
		fmt.Printf("[ToggleStatus] error: %s", err)
	}
	return err
}
func DeleteKasten(kastenid string) {
	karten, _ := GetKartenByKasten(kastenid)
	for i := 0; i < len(karten); i++ {
		karte, _ := Map2Karte(karten[i])
		_ = btDB.Delete(karte.ID)
	}
	_ = btDB.Delete(kastenid)
}

func kasten2Map(u Kasten) (kasten map[string]interface{}, err error) {
	uJSON, err := json.Marshal(u)
	json.Unmarshal(uJSON, &kasten)

	return kasten, err
}

// Convert from map[string]interface{} to Todo struct as required by golang-couchdb methods
func map2Kasten(kasten map[string]interface{}) (t Kasten, err error) {
	uJSON, err := json.Marshal(kasten)
	json.Unmarshal(uJSON, &t)

	return t, err
}
