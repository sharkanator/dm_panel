
package api

import (
    "encoding/json"
    "net/http"
        "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

var dbPath = "./internal/db/dmpanel.db"

type ImagePayload struct {
    Filename   string `json:"filename"`
    FogEnabled bool   `json:"fog"`
}

type FogPayload struct {
    FogState string `json:"fog_state"`
}

func SetCurrentImage(w http.ResponseWriter, r *http.Request) {
    var payload ImagePayload
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Données invalides", http.StatusBadRequest)
        return
    }

    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        http.Error(w, "Erreur base de données", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    db.Exec("UPDATE images SET is_current = 0")
    res, err := db.Exec("UPDATE images SET is_current = 1, fog_enabled = ? WHERE filename = ?", payload.FogEnabled, payload.Filename)
    if err != nil {
        http.Error(w, "Erreur mise à jour", http.StatusInternalServerError)
        return
    }

    affected, _ := res.RowsAffected()
    if affected == 0 {
        db.Exec("INSERT INTO images (filename, is_current, fog_enabled) VALUES (?, 1, ?)", payload.Filename, payload.FogEnabled)
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Image mise à jour"))
}

func GetCurrentImage(w http.ResponseWriter, r *http.Request) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        http.Error(w, "Erreur base de données", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    var filename string
    var fogEnabled bool
    var fogState string

    err = db.QueryRow("SELECT filename, fog_enabled, fog_state FROM images WHERE is_current = 1").Scan(&filename, &fogEnabled, &fogState)
    if err != nil {
        http.Error(w, "Aucune image trouvée", http.StatusNotFound)
        return
    }

    response := map[string]interface{}{
        "filename":    filename,
        "fog_enabled": fogEnabled,
        "fog_state":   fogState,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// Nouvelle route pour mettre à jour l'état du fog
func UpdateFogState(w http.ResponseWriter, r *http.Request) {
    var payload FogPayload
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Données invalides", http.StatusBadRequest)
        return
    }

    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        http.Error(w, "Erreur base de données", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    _, err = db.Exec("UPDATE images SET fog_state = ? WHERE is_current = 1", payload.FogState)
    if err != nil {
        http.Error(w, "Erreur mise à jour fog", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Fog mis à jour"))
}
