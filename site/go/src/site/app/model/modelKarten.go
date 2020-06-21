package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
)

type Karten struct {
	ID           string `json:"_id"`
	Rev          string `json:"_rev"`
	Titel        string `json:"titel"`
	Frage        string `json:"frage"`
	Antwort      string `json:"antwort"`
	Ersteller    string `json:"ersteller"`
	Type         string `json:"type"`
	KastenID     string `json:"kastenid"`
	KartenNummer int    `json:"kartennummer"`
}

func GetRandomKarte(kastenid string) (karte Karten, err error) {
	karten, _ := GetKartenByKasten(kastenid)
	if len(karten) <= 1 {
		if len(karten) < 1 {
			err = errors.New("Keine Karten Digger")
		} else {
			karte, err = Map2Karte(karten[0])
		}
	} else {
		ran := rand.Intn(len(karten))
		karte, err = Map2Karte(karten[ran])
	}
	return karte, err
}

func GetKartenByKasten(kastenid string) ([]map[string]interface{}, error) {
	query := `{
		"selector": {
			"type":"karte",
			"kastenid": "%s"
			}
			 }`
	AllKasten, err := btDB.QueryJSON(fmt.Sprintf(query, kastenid))
	if err != nil {
		panic(err)
	} else {
		return AllKasten, nil
	}
}

func KarteErstellen(username string, kastenid string, titel string, frage string, antwort string) (karte Karten) {
	//Karte erstellen
	karte.Ersteller = username
	karte.KastenID = kastenid
	karte.Titel = titel
	karte.Frage = frage
	karte.Antwort = antwort
	karte.Type = "karte"
	karte.KartenNummer = 1
	yo, _ := GetKartenByKasten(kastenid)
	for i := 0; i < len(yo); i++ {
		karte.KartenNummer = karte.KartenNummer + 1
	}
	karte.Add()

	//Update des Kastencounts
	Kasten, _ := GetKastenById(kastenid, username)
	Kasten.Count = Kasten.Count + 1
	m, _ := kasten2Map(Kasten)
	_ = btDB.Set(Kasten.ID, m)

	return karte
}
func KarteUpdaten(kartenid string, titel string, frage string, antwort string) (karte Karten) {

	karte, _ = GetKarteByIdOnly(kartenid)
	karte.Titel = titel
	karte.Frage = frage
	karte.Antwort = antwort

	m := karten2Map(karte)
	_ = btDB.Set(karte.ID, m)

	return karte
}
func DeleteKarte(kartenid string, kastenid string, username string) {
	//Update Kasten Count
	Kasten, _ := GetKastenById(kastenid, username)
	Kasten.Count = Kasten.Count - 1
	m, _ := kasten2Map(Kasten)
	_ = btDB.Set(Kasten.ID, m)

	//Update kartennummern
	Karten, _ := GetKartenByKasten(kastenid)
	Karte, _ := GetKarteByIdOnly(kartenid)
	if len(Karten) > 1 {
		for i := Karte.KartenNummer; i < len(Karten); i++ {
			a, _ := Map2Karte(Karten[i])
			a.KartenNummer = a.KartenNummer - 1

			m := karten2Map(a)
			_ = btDB.Set(a.ID, m)
		}
	}

	//Delete Karte
	btDB.Delete(kartenid)
}

func (t Karten) Add() error {
	karten := karten2Map(t)

	delete(karten, "_id")
	delete(karten, "_rev")

	_, _, err := btDB.Save(karten, nil)
	if err != nil {
		fmt.Printf("[Add] error : %s", err)
	}
	return err
}

func Remove(slice []map[string]interface{}, s int) []map[string]interface{} {

	return append(slice[:s], slice[s+1:]...)

}
func karten2Map(t Karten) map[string]interface{} {
	var doc map[string]interface{}
	tJSON, _ := json.Marshal(t)
	json.Unmarshal(tJSON, &doc)

	return doc
}

func GetKarten(id string) (Karten, error) {
	t, err := btDB.Get(id, nil)
	if err != nil {
		return Karten{}, err
	}

	todo, err := Map2Karte(t)
	return todo, err
}
func GetKarteById(allKarten []map[string]interface{}, id string) (karten Karten, err error) {
	var aktiv Karten
	for i := 0; i < len(allKarten); i++ {
		if allKarten[i]["_id"] == id {
			aktiv, _ = Map2Karte(allKarten[i])
			return aktiv, err
		}
	}
	return aktiv, err
}
func GetKarteByIdOnly(id string) (karte Karten, err error) {
	query := `{
		"selector": {
			"type":"karte",
			"_id": "%s"
			}
			 }`
	kartenMap, err := btDB.QueryJSON(fmt.Sprintf(query, id))
	if err != nil || len(kartenMap) == 0 {
		err = errors.New("Karte nicht gefunden")
	} else {
		karte, _ = Map2Karte(kartenMap[0])
	}
	return karte, err
}

// Convert from map[string]interface{} to Todo struct as required by golang-couchdb methods
func Map2Karte(karten map[string]interface{}) (t Karten, err error) {
	uJSON, err := json.Marshal(karten)
	json.Unmarshal(uJSON, &t)

	return t, err
}
