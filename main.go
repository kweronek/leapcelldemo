package main

import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Konto struct {
    ID      string  `json:"id" gorm:"primaryKey"`
    UserID  string  `json:"user_id"`
    Balance float64 `json:"balance"`
}

var db *gorm.DB

func main() {
    /*
    connStr := "host=ep-ancient-shadow-a5aui0rq-pooler.us-east-2.aws.neon.tech user=neondb_owner password=npg_uHIeinM7x6Cp dbname=neondb sslmode=require"
    */
    // Lese das Datenbank-Passwort aus der Umgebungsvariable.
    dbPassword := os.Getenv("DB_PASSWORD")
    if dbPassword == "" {
        panic("Die Umgebungsvariable DB_PASSWORD ist nicht gesetzt!")
    }

    // Hier werden weitere Verbindungsparameter festgelegt.
    connStr := fmt.Sprintf(
        "host=ep-ancient-shadow-a5aui0rq-pooler.us-east-2.aws.neon.tech user=neondb_owner password=%s dbname=neondb sslmode=require",
        dbPassword
    )

    fmt.Println(connStr)
    
    var err error
    db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
    if err != nil {
        panic("Fehler beim Verbinden zur Datenbank: " + err.Error())
    }
    
    // Automatisches Migration der Tabelle, falls nicht vorhanden
    err = db.AutoMigrate(&Konto{})
    if err != nil {
        panic("Fehler bei Migration: " + err.Error())
    }

    fmt.Println("✅ Verbunden mit DB und Migration erfolgreich!")

    r := gin.Default()

    api := r.Group("/user/:userID/konto")
    {
        api.GET("/:kontoID", getKonto)
        api.POST("/", createKonto)
        api.PUT("/:kontoID", updateKonto)
        api.DELETE("/:kontoID", deleteKonto)
    }

    r.Run(":8080") // Startet Server
}

func getKonto(c *gin.Context) {
    userID := c.Param("userID")
    kontoID := c.Param("kontoID")
    var konto Konto

    if err := db.First(&konto, "id = ? AND user_id = ?", kontoID, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Konto nicht gefunden"})
        return
    }
    c.JSON(http.StatusOK, konto)
}

func createKonto(c *gin.Context) {
    userID := c.Param("userID")
    var newKonto Konto

    if err := c.ShouldBindJSON(&newKonto); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    newKonto.UserID = userID

    if err := db.Create(&newKonto).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Erstellen"})
        return
    }
    c.JSON(http.StatusCreated, newKonto)
}

func updateKonto(c *gin.Context) {
    userID := c.Param("userID")
    kontoID := c.Param("kontoID")
    var konto Konto

    if err := db.First(&konto, "id = ? AND user_id = ?", kontoID, userID).Error; err != nil {
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
    userID := c.Param("userID")
    kontoID := c.Param("kontoID")
    var konto Konto

    if err := db.First(&konto, "id = ? AND user_id = ?", kontoID, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Konto nicht gefunden"})
        return
    }

    db.Delete(&konto)
    c.JSON(http.StatusOK, gin.H{"status": "Konto gelöscht"})
}
