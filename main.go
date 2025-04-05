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

	// Serve the landing page
	router.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Welcome</title>
			</head>
			<body>
				<h1>Welcome to the Landing Page</h1>
				<p>This is a simple HTML landing page.</p>
			</body>
			</html>
		`))
	})
	router.GET("/albums",getAlbums)
	router.POST("/albums",postAlbum)
	router.GET("/albums/:id",getAlbumbyID)

	router.Run("localhost:8000")
}

func getAlbums(c *gin.Context)  {
	c.IndentedJSON(http.StatusOK,albums)
}

func postAlbum(c *gin.Context)  {
	var newAlbum album;
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}
	albums = append(albums,newAlbum)
	c.IndentedJSON(http.StatusCreated,newAlbum)
}

func getAlbumbyID(c *gin.Context) {
	id := c.Param("id")

	for _,a := range(albums) {
		if a.Id == id {
			c.IndentedJSON(http.StatusOK,a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound,gin.H{"message":"album not found"})
}	