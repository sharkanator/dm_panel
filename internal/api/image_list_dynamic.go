
package api

import (
    "encoding/json"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

func GetImageList(w http.ResponseWriter, r *http.Request) {
    images := []string{}
    files, err := os.ReadDir("web/static/images")
    if err != nil {
        http.Error(w, "Erreur de lecture du dossier images", http.StatusInternalServerError)
        return
    }

    for _, f := range files {
        if !f.IsDir() {
            name := f.Name()
            ext := strings.ToLower(filepath.Ext(name))
            if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp" || ext == ".svg" || ext == ".bmp" {
                images = append(images, name)
            }
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(images)
}
