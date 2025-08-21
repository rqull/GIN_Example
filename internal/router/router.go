package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/rqull/GIN_Example/internal/controller"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	ctl := controller.NewBioskopController(db)

	api := r.Group("/api/v1")
	{
		api.POST("/bioskop", ctl.Create)
		api.GET("/bioskop", ctl.GetAll)
		api.GET("/bioskop/:id", ctl.GetByID)
		api.PUT("/bioskop/:id", ctl.Update)
		api.DELETE("/bioskop/:id", ctl.Delete)
	}

	return r
}
