package main
import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type album struct{
	Id string `json:"id`
	Title string `json:title`
	Artist string `json:artist`
	Price float64 `json:price`
}

// var albums = []album{
// 	{Id:"1",Title:"Blue train",Artist:"John Coltrane",Price:45.90},
// 	{Id:"2",Title:"Muskurane",Artist:"Arijit sing",Price:200.89},
// 	{Id:"3",Title:"Dark knight",Artist:"Alan Walker",Price:90},
// }

func main()  {
	initDB()
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	// Serve the landing page
	router.GET("/",func (c *gin.Context)  {
		c.HTML(http.StatusOK,"index.html",nil)
	})
	
	router.GET("/albums",getAlbums)
	router.GET("/albums/:id",getAlbumbyID)
	router.POST("/albums",createAlbum)
	router.PATCH("/albums/:id", updateAlbum)
	router.DELETE("/albums/:id",deleteAlbum)
	

	router.Run("localhost:8000")
}

// func getAlbums(c *gin.Context)  {
// 	c.JSON(http.StatusOK,albums)
// }

func getAlbums(c *gin.Context)  {
	rows, err := db.Query("SELECT id,title,artist,price FROM albums")
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	defer rows.Close()

	var albums []album
	for rows.Next() {
		var a album
		err := rows.Scan(&a.Id,&a.Title,&a.Artist,&a.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		albums = append(albums,a)
	}
	c.JSON(http.StatusOK,albums);
}

func createAlbum(c *gin.Context)  {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	stmt, err := db.Prepare("INSERT INTO albums (id,title,artist,price) VALUES (?,?,?,?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(newAlbum.Id,newAlbum.Title,newAlbum.Artist,newAlbum.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not insert album"})
		return
	}
	c.JSON(http.StatusCreated,newAlbum)

}

func getAlbumbyID(c *gin.Context) {
	var a album
	id := c.Param("id")
	err := db.QueryRow("SELECT id,title,artist,price FROM albums WHERE id=?",id).
					Scan(&a.Id, &a.Title, &a.Artist, &a.Price)
	if err != nil {
		c.JSON(http.StatusNotFound,gin.H{"message":"album not found"})
		return
	}
	
	c.JSON(http.StatusOK, a)

	
}	

func updateAlbum(c *gin.Context)  {
	var newAlbum album
	id := c.Param("id")
	if err := c.BindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	var existing album
	err := db.QueryRow("SELECT id,title,artist,price FROM albums WHERE id=?",id).Scan(&existing.Id,&existing.Title,&existing.Artist,&existing.Price)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Album not found"})
		return
	}
	if newAlbum.Title != "" {
		existing.Title = newAlbum.Title
	}
	if newAlbum.Artist != "" {
		existing.Artist = newAlbum.Artist
	}
	if newAlbum.Price != 0 {
		existing.Price = newAlbum.Price
	}

	stmt, err := db.Prepare("UPDATE albums SET title=?,artist=?,price=? WHERE id=?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(existing.Title,existing.Artist,existing.Price,id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	c.JSON(http.StatusOK,existing)
}

func deleteAlbum(c *gin.Context)  {
	id := c.Param("id")
	stmt, err := db.Prepare("DELETE FROM albums WHERE id=?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound,gin.H{"message":"Album not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Album deleted"})
}