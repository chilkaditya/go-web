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

var albums = []album{
	{Id:"1",Title:"Blue train",Artist:"John Coltrane",Price:45.90},
	{Id:"2",Title:"Muskurane",Artist:"Arijit sing",Price:200.89},
	{Id:"3",Title:"Dark knight",Artist:"Alan Walker",Price:90},
}

func main()  {
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

func getAlbums(c *gin.Context)  {
	c.JSON(http.StatusOK,albums)
}

func createAlbum(c *gin.Context)  {
	var newAlbum album;
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}
	albums = append(albums,newAlbum)
	c.JSON(http.StatusCreated,newAlbum)
}

func getAlbumbyID(c *gin.Context) {
	id := c.Param("id")
	for _,a := range(albums) {
		if a.Id == id {
			c.JSON(http.StatusOK,a)
			return
		}
	}
	c.JSON(http.StatusNotFound,gin.H{"message":"album not found"})
}	

func updateAlbum(c *gin.Context)  {
	var newAlbum album
	id := c.Param("id")
	if err := c.BindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	for i,a := range(albums) {
		if a.Id == id {
			// Only update fields if they're non-zero (basic check)
			if newAlbum.Title != "" {
				albums[i].Title = newAlbum.Title
			}
			if newAlbum.Artist != "" {
				albums[i].Artist = newAlbum.Artist
			}
			if newAlbum.Price != 0 {
				albums[i].Price = newAlbum.Price
			}
			c.JSON(http.StatusOK, albums[i])
			return
		}
	}
	c.JSON(http.StatusNotFound,gin.H{"message":"album not found"})

}

func deleteAlbum(c *gin.Context)  {
	id := c.Param("id")
	for i,a := range(albums) {
		if a.Id == id {
			albums = append(albums[:i], albums[i+1:]...)
			c.JSON(http.StatusOK,gin.H{"message": "Album deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound,gin.H{"message":"album not found"})
}