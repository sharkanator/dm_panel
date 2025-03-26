
package api

import (
    "encoding/json"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

func GetMusicList(w http.ResponseWriter, r *http.Request) {
    files, err := os.ReadDir("web/static/audio")
    if err != nil {
        http.Error(w, "Could not read audio directory", http.StatusInternalServerError)
        return
    }

    var music []string
    for _, f := range files {
        if !f.IsDir() {
            name := f.Name()
            ext := strings.ToLower(filepath.Ext(name))
            if ext == ".mp3" || ext == ".ogg" || ext == ".wav" {
                music = append(music, name)
            }
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(music)
}
