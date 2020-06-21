package model

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User data structure
type User struct {
	ID              string `json:"_id"`
	Rev             string `json:"_rev"`
	Type            string `json:"type"`
	Done            bool   `json:"done"`
	Username        string `json:"name"`
	Password        string `json:"password"`
	Email           string `json:"email"`
	Eigenkasten     []int  `json:"eigenkasten"`
	ErstellteKarten int    `json:"erstellteKarten"`
	ErstellteKasten int    `json:"erstellteKasten"`
	ErstelltAm      string `json:"erstelltAm"`
	Bild            string `json:"bild"`
}

func UpdateFoto(username string, adresse string) {
	User, _ := GetUserByUsername(username)

	adr := strings.SplitAfterN(adresse, "\\", -1)
	fmt.Println(adr)
	c := adr[1] + adr[2]
	User.Bild = c
	u, _ := user2Map(User)
	_ = btDB.Set(User.ID, u)
}

// Add User
func (user User) Add() (err error) {
	// Check wether username already exists
	userInDB, err := GetUserByUsername(user.Username)
	if err == nil && userInDB.Username == user.Username {
		return errors.New("username exists already")
	}

	// Hash password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	b64HashedPwd := base64.StdEncoding.EncodeToString(hashedPwd)

	time := time.Now().String()
	time3 := strings.SplitAfterN(time, " ", 2)

	user.Password = b64HashedPwd
	user.Type = "User"
	user.ErstelltAm = time3[0]

	// Convert Todo struct to map[string]interface as required by Save() method
	u, err := user2Map(user)

	// Delete _id and _rev from map, otherwise DB access will be denied (unauthorized)
	delete(u, "_id")
	delete(u, "_rev")

	// Add todo to DB
	_, _, err = btDB.Save(u, nil)

	if err != nil {
		fmt.Printf("[Add] error: %s", err)
	}

	return err
}

// GetUserByUsername retrieve User by username
func GetUserByUsername(username string) (user User, err error) {
	if username == "" {
		return User{}, errors.New("no username provided")
	}

	query := `
	{
		"selector": {
			 "type": "User",
			 "name": "%s"
		}
	}`
	u, err := btDB.QueryJSON(fmt.Sprintf(query, username))
	if err != nil || len(u) < 1 {
		return User{}, err
	}

	user, err = map2User(u[0])
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func UserExist(username string) bool {
	query := `
	{
		"selector": {
			 "type": "User",
			 "name": "%s"
		}
	}`
	u, _ := btDB.QueryJSON(fmt.Sprintf(query, username))
	if len(u) < 1 {
		return false
	}
	return true
}

func EmailExist(email string) bool {
	query := `
	{
		"selector": {
			 "type": "User",
			 "email": "%s"
		}
	}`
	u, _ := btDB.QueryJSON(fmt.Sprintf(query, email))
	if len(u) < 1 {
		return false
	}
	return true
}

func CheckPassword(username, password string) bool {
	user, err := GetUserByUsername(username)
	passwordDB, _ := base64.StdEncoding.DecodeString(user.Password)
	err = bcrypt.CompareHashAndPassword(passwordDB, []byte(password))
	if err == nil {
		return true
	} else {
		//falsches Passwort
		return false
	}
}

func UpdatePassword(username, password string) (err error) {
	userInDB, err := GetUserByUsername(username)
	passwordDB, _ := base64.StdEncoding.DecodeString(userInDB.Password)
	err = bcrypt.CompareHashAndPassword(passwordDB, []byte(password))
	if err != nil {
		// Hash password
		hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		b64HashedPwd := base64.StdEncoding.EncodeToString(hashedPwd)

		userInDB.Password = b64HashedPwd
		Docid := userInDB.ID
		// Convert Todo struct to map[string]interface as required by Save() method
		u, err := user2Map(userInDB)
		// Add todo to DB
		err = btDB.Set(Docid, u)

		if err != nil {
			fmt.Printf("[Add] error: %s", err)
		}

		return err
	} else {
		return errors.New("same pw")
	}
}

func DeleteUser(username string) {
	User, _ := GetUserByUsername(username)
	Kasten, _, _, _ := GetKastenByUser(username)
	for i := 0; i < len(Kasten); i++ {
		a, _ := map2Kasten(Kasten[i])
		DeleteKasten(a.ID)
	}

	_ = btDB.Delete(User.ID)
}

// ---------------------------------------------------------------------------
// Internal helper functions
// ---------------------------------------------------------------------------

// Convert from User struct to map[string]interface{} as required by golang-couchdb methods
func user2Map(u User) (user map[string]interface{}, err error) {
	uJSON, err := json.Marshal(u)
	json.Unmarshal(uJSON, &user)

	return user, err
}

// Convert from map[string]interface{} to User struct as required by golang-couchdb methods
func map2User(user map[string]interface{}) (u User, err error) {
	uJSON, err := json.Marshal(user)
	json.Unmarshal(uJSON, &u)

	return u, err
}
