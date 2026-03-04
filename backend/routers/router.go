package routers

import (
	"employee-management/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	userGroup := r.Group("/api/user")
	{
		userGroup.POST("/register", controllers.Register)
		userGroup.POST("/login", controllers.Login)
		userGroup.POST("/change-password", controllers.ChangePassword)
	}

	employeeGroup := r.Group("/api/employee")
	{
		employeeGroup.GET("", controllers.GetEmployees)
		employeeGroup.GET("/:id", controllers.GetEmployee)
		employeeGroup.POST("", controllers.CreateEmployee)
		employeeGroup.PUT("/:id", controllers.UpdateEmployee)
		employeeGroup.DELETE("/:id", controllers.DeleteEmployee)
	}

	return r
}
