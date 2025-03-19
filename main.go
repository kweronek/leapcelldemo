package main

import (
    "fmt"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type User struct {
    ID     string  `json:"id" gorm:"primaryKey"`
    Name   string  `json:"name"`
    Konten []Konto `json:"konten" gorm:"foreignKey:UserID"`
}

type Konto struct {
    ID      string  `json:"id" gorm:"primaryKey"`
    UserID  string  `json:"user_id" gorm:"index"`
    Balance float64 `json:"balance"`
}

var db *gorm.DB

func main() {
    dbPassword := os.Getenv("DB_PASSWORD")
    if dbPassword == "" {
        panic("Die Umgebungsvariable DB_PASSWORD ist nicht gesetzt!")
    }

    connStr := fmt.Sprintf(
        "host=ep-ancient-shadow-a5aui0rq-pooler.us-east-2.aws.neon.tech user=neondb_owner password=%s dbname=neondb sslmode=require",
        dbPassword,
    )

    var err error
    db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
    if err != nil {
        panic("Fehler beim Verbinden zur Datenbank: " + err.Error())
    }

    err = db.AutoMigrate(&User{}, &Konto{})
    if err != nil {
        panic("Fehler bei Migration: " + err.Error())
    }

    fmt.Println("✅ Verbunden mit DB und Migration erfolgreich!")

    r := gin.Default()

    // User-Routen
    r.POST("/user", createUser)
    r.GET("/user/:userID", getUserWithKonten)

    // Konto-Routen (ohne UserID in URL)
    r.GET("/konto/:kontoID", getKonto)
    r.POST("/konto", createKonto)
    r.PUT("/konto/:kontoID", updateKonto)
    r.DELETE("/konto/:kontoID", deleteKonto)

    // Konten sortiert nach Balance
    r.GET("/konten/sortiert", getKontenSortiert)

    r.Run(":8080")
}

func createUser(c *gin.Context) {
    var newUser User
    if err := c.ShouldBindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := db.Create(&newUser).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Erstellen"})
        return
    }

    c.JSON(http.StatusCreated, newUser)
}

func getUserWithKonten(c *gin.Context) {
    userID := c.Param("userID")
    var user User

    if err := db.Preload("Konten").First(&user, "id = ?", userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User nicht gefunden"})
        return
    }

    c.JSON(http.StatusOK, user)
}

func getKonto(c *gin.Context) {
    kontoID := c.Param("kontoID")
    var konto Konto

    if err := db.First(&konto, "id = ?", kontoID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Konto nicht gefunden"})
        return
    }
    c.JSON(http.StatusOK, konto)
}

func createKonto(c *gin.Context) {
    var newKonto Konto

    if err := c.ShouldBindJSON(&newKonto); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := db.Create(&newKonto).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Erstellen"})
        return
    }
    c.JSON(http.StatusCreated, newKonto)
}

func updateKonto(c *gin.Context) {
    kontoID := c.Param("kontoID")
    var konto Konto

    if err := db.First(&konto, "id = ?", kontoID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Konto nicht gefunden"})
        return
    }

    var updateData Konto
    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    konto.Balance = updateData.Balance
    db.Save(&konto)

    c.JSON(http.StatusOK, konto)
}

func deleteKonto(c *gin.Context) {
    kontoID := c.Param("kontoID")
    var konto Konto

    if err := db.First(&konto, "id = ?", kontoID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Konto nicht gefunden"})
        return
    }

    db.Delete(&konto)
    c.JSON(http.StatusOK, gin.H{"status": "Konto gelöscht"})
}

func getKontenSortiert(c *gin.Context) {
    var konten []Konto

    // Optional: Sortierreihenfolge via Query-Parameter (?order=asc|desc)
    order := c.DefaultQuery("order", "desc")
    if order != "asc" && order != "desc" {
        order = "desc"
    }

    if err := db.Order("balance " + order).Find(&konten).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Abrufen"})
        return
    }

    c.JSON(http.StatusOK, konten)
}
