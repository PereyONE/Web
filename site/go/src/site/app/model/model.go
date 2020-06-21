package model

type (
	// Body vom login request
	Login struct {
		Username string `json:"User"`
		Password string `json:"Passwort"`
	}

	// Body vom registrierungs request
	Registrierung struct {
		Name       string `json:"Name"`
		EMail      string `json:"EMail"`
		Passwort   string `json:"Passwort"`
		Akzeptiert bool   `json:"Akzeptiert"`
	}

	// UpdateProfil body of profil put
	UpdateProfil struct {
		EMail    string `json:"EMail"`
		Passwort string `json:"Passwort"`
		Neu      string `json:"Neu"`
	}

	// LernenErgebnis body of lern post
	LernenErgebnis struct {
		Index    int  `json:"Index"`
		Ergebnis bool `json:"Ergebnis"`
	}

	// Karteikasten body of edit post/put
	Karteikasten struct {
		Kategorie      string `json:"Kategorie"`
		Unterkategorie string `json:"Unterkategorie"`
		Titel          string `json:"Titel"`
		Beschreibung   string `json:"Beschreibung"`
		Public         bool   `json:"Public"`
	}

	// Karteikarte body of edit2 post/put
	Karteikarte struct {
		Titel   string `json:"Titel"`
		Frage   string `json:"Frage"`
		Antwort string `json:"Antwort"`
	}
)
