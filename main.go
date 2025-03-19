package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Konto struct {
    ID      string  `json:"id"`
    UserID  string  `json:"user_id"`
    Balance float64 `json:"balance"`
}

var konten = []Konto{
    {ID: "1", UserID: "7", Balance: 1000.0},
    {ID: "2", UserID: "7", Balance: 500.0},
}

func main() {
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

    for _, k := range konten {
        if k.ID == kontoID && k.UserID == userID {
            c.JSON(http.StatusOK, k)
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"error": "Konto nicht gefunden"})
}

func createKonto(c *gin.Context) {
    userID := c.Param("userID")
    var newKonto Konto

    if err := c.ShouldBindJSON(&newKonto); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    newKonto.UserID = userID
    konten = append(konten, newKonto)
    c.JSON(http.StatusCreated, newKonto)
}

func updateKonto(c *gin.Context) {
    userID := c.Param("userID")
    kontoID := c.Param("kontoID")
    var updateData Konto

    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    for i, k := range konten {
        if k.ID == kontoID && k.UserID == userID {
            konten[i].Balance = updateData.Balance
            c.JSON(http.StatusOK, konten[i])
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"error": "Konto nicht gefunden"})
}

func deleteKonto(c *gin.Context) {
    userID := c.Param("userID")
    kontoID := c.Param("kontoID")

    for i, k := range konten {
        if k.ID == kontoID && k.UserID == userID {
            konten = append(konten[:i], konten[i+1:]...)
            c.JSON(http.StatusOK, gin.H{"status": "Konto gelöscht"})
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"error": "Konto nicht gefunden"})
}
