package router

import (
	"cuboid-challenge/app/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Setup returns the app router.
func Setup() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	bag := r.Group("/bags")
	{
		bag.GET("", controller.ListBags)
		bag.GET("/:bagID", controller.GetBag)
		bag.POST("", controller.CreateBag)
		bag.DELETE("/:bagID", controller.DeleteBag)
	}

	cuboid := r.Group("/cuboids")
	{
		cuboid.GET("", controller.ListCuboids)
		cuboid.GET("/:cuboidId", controller.GetCuboid)
		cuboid.POST("", controller.CreateCuboid)
		cuboid.PUT("/:cuboidId", controller.UpdateCuboid)
		cuboid.DELETE("/:cuboidId", controller.DeleteCuboid)
	}

	return r
}
