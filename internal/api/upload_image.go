
package api

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
)

func UploadImage(w http.ResponseWriter, r *http.Request) {
    r.ParseMultipartForm(10 << 20) // Max 10MB

    file, handler, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Error retrieving file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    dstPath := filepath.Join("web/static/images", filepath.Base(handler.Filename))
    dst, err := os.Create(dstPath)
    if err != nil {
        http.Error(w, "Error saving file", http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    _, err = io.Copy(dst, file)
    if err != nil {
        http.Error(w, "Error writing file", http.StatusInternalServerError)
        return
    }

    
    // --- Mise à jour DB ---
    db, err := sql.Open("sqlite3", "internal/database/dmpanel.db")
    if err != nil {
        fmt.Println("Erreur DB:", err)
        return
    }
    defer db.Close()

    _, _ = db.Exec("UPDATE images SET is_current = 0")

    _, err = db.Exec("INSERT OR REPLACE INTO images (filename, is_current) VALUES (?, 1)", handler.Filename)
    if err != nil {
        fmt.Println("Erreur insertion DB:", err)
    }

    fmt.Fprint(w, "ok")

}
