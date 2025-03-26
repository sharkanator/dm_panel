
# Utilise une image officielle de Go avec support SQLite
FROM golang:1.21

# Création du répertoire de travail
WORKDIR /app

# Copie du code source
COPY . .

# Téléchargement des dépendances
RUN go mod tidy

# Compilation
RUN go build -o dm-server ./cmd/server/main.go

# Port d’écoute de l’application
EXPOSE 8080

# Commande de démarrage
CMD ["./dm-server"]
