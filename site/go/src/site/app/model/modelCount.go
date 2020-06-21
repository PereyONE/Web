package model

import (
	"encoding/json"
	"fmt"

	couchdb "github.com/leesper/couchdb-golang"
)

type Count struct {
	Id     string `json: "_id"`
	Rev    string `json: "_rev"`
	Users  string `json: "users"`
	Kasten string `json: "kasten"`
	Karten string `json: "karten"`
	Done   bool   `json:"done"`
	couchdb.Document
}

var btDB *couchdb.Database

func init() {
	var err error
	btDB, err = couchdb.NewDatabase("http://localhost:5984/dream")
	if err != nil {
		panic(err)
	}
}
func GetPrivateKastenCount(username string) int {
	query := `
	{
		"selector": {
			 "type": "kasten",
			 "ersteller": "%s"
		}
	}`
	u, _ := btDB.QueryJSON(fmt.Sprintf(query, username))
	user := len(u)
	return user
}

func GetPublicKastenCount() int {
	query := `
	{
		"selector": {
			 "type": "kasten",
			 "sichtbarkeit": true
		}
	}`
	u, _ := btDB.QueryJSON(fmt.Sprintf(query))
	user := len(u)
	return user
}

func GetKartenCountForUser(username string) int {
	Kasten, _, _, _ := GetKastenByUser(username)
	count := 0
	for i := 0; i < len(Kasten); i++ {
		a, _ := map2Kasten(Kasten[i])
		b, _ := GetKartenByKasten(a.ID)
		count = count + len(b)
	}
	return count
}

func GetUserCount() int {
	query := `
	{
		"selector": {
			 "type": "User"
		}
	}`
	u, _ := btDB.QueryJSON(fmt.Sprintf(query))
	user := len(u)
	return user
}
func GetKastenCount() int {
	query := `
	{
		"selector": {
			 "type": "kasten"
		}
	}`
	u, _ := btDB.QueryJSON(fmt.Sprintf(query))
	user := len(u)
	return user
}
func GetKartenCount() int {
	query := `
	{
		"selector": {
			 "type": "karte"
		}
	}`
	u, _ := btDB.QueryJSON(fmt.Sprintf(query))
	user := len(u)
	return user
}

func (t Count) countToggleStatus() error {
	if t.Done {
		t.Done = false
	} else {
		t.Done = true
	}
	err := btDB.Set(t.Id, count2Map(t))

	if err != nil {
		fmt.Printf("[ToggleStatus] error: %s", err)
	}
	return err
}

func count2Map(t Count) map[string]interface{} {
	var doc map[string]interface{}
	tJSON, _ := json.Marshal(t)
	json.Unmarshal(tJSON, &doc)

	return doc
}
func map2Count(count map[string]interface{}) (u Count, err error) {
	uJSON, err := json.Marshal(count)
	json.Unmarshal(uJSON, &u)

	return u, err
}
