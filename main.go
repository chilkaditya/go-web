package main
import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type MusicAlbum struct{
	Id string `json:"id`
	Title string `json:title`
	Artist string `json:artist`
	Movie_name string `json:movie_name`
	Language string `json:language`
	// Actor_actress string `json:actor_actress`
}

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
	rows, err := db.Query("SELECT id,title,artist,movie_name,language FROM MusicAlbums")
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	defer rows.Close()

	var albums []MusicAlbum
	for rows.Next() {
		var a MusicAlbum
		err := rows.Scan(&a.Id,&a.Title,&a.Artist,&a.Movie_name,&a.Language)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		albums = append(albums,a)
	}
	c.JSON(http.StatusOK,albums);
}

func createAlbum(c *gin.Context)  {
	var newAlbum MusicAlbum
	if err := c.BindJSON(&newAlbum); err != nil { // Its bind the input data to newAlbum in json format
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	stmt, err := db.Prepare("INSERT INTO MusicAlbums (id,title,artist,movie_name,language) VALUES (?,?,?,?,?)")
	// compiled the sql statement and later use Exec() to execute the the sql statement
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(newAlbum.Id,newAlbum.Title,newAlbum.Artist,newAlbum.Movie_name,newAlbum.Language) // herer we execuete the sql statement and put the '?' placeholder value
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not insert album"})
		return
	}
	c.JSON(http.StatusCreated,newAlbum)

}

func getAlbumbyID(c *gin.Context) {
	var a MusicAlbum
	id := c.Param("id") // getting the id from request url
	err := db.QueryRow("SELECT id,title,artist,movie_name,language FROM MusicAlbums WHERE id=?",id).
					Scan(&a.Id, &a.Title, &a.Artist, &a.Movie_name, &a.Language)
	if err != nil {
		c.JSON(http.StatusNotFound,gin.H{"message":"album not found"})
		return
	}
	
	c.JSON(http.StatusOK, a)

	
}	

func updateAlbum(c *gin.Context)  {
	var newAlbum MusicAlbum
	id := c.Param("id") // getting the id from request url
	if err := c.BindJSON(&newAlbum); err != nil { // parse the json request body into newAlbum
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	var existing MusicAlbum
	err := db.QueryRow("SELECT id,title,artist,movie_name,language FROM MusicAlbums WHERE id=?",id).Scan(&existing.Id,&existing.Title,&existing.Artist,&existing.Movie_name,&existing.Language)
	//current values of albums with the given id are scan into exisiting album
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
	if newAlbum.Movie_name != "" {
		existing.Movie_name = newAlbum.Movie_name
	}
	if newAlbum.Language != "" {
		existing.Language = newAlbum.Language
	}
	

	stmt, err := db.Prepare("UPDATE MusicAlbums SET title=?,artist=?,movie_name=?,language=? WHERE id=?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(existing.Title,existing.Artist,existing.Movie_name,existing.Language,id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	c.JSON(http.StatusOK,existing)
}

func deleteAlbum(c *gin.Context)  {
	id := c.Param("id")
	stmt, err := db.Prepare("DELETE FROM MusicAlbums WHERE id=?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	defer stmt.Close()
	rows, err := stmt.Exec(id) //rows is a result object containing metadata
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	rowsAffected, _ := rows.RowsAffected()// return how many rows are affected after deletion
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound,gin.H{"message":"Album not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Album deleted"})
}