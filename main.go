package main

import (
	"documentation/auth"
	"documentation/entity"
	"documentation/handler"
	"documentation/helper"
	"documentation/input"
	"documentation/repository"
	"documentation/service"
	webHandler "documentation/web/handler"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// godotenv not used for deploy on heroku
	// env := godotenv.Load()
	// if env != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	host := os.Getenv("DB_HOST")
	userHost := os.Getenv("DB_USER")
	userPass := os.Getenv("DB_PASSWORD")
	databaseName := os.Getenv("DB_DATABASE")
	databasePort := os.Getenv("DB_PORT")

	dsn := "host=" + host + " user=" + userHost + " password=" + userPass + " dbname=" + databaseName + " port=" + databasePort + " sslmode=require TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Feature{})
	db.AutoMigrate(&entity.FeatureDetail{})
	db.AutoMigrate(&entity.Prd{})

	fmt.Println("Database Connected")

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	authService := auth.NewService()
	userHandler := handler.NewUserHandler(userService, authService)

	prdRepository := repository.NewPrdRepository(db)
	prdService := service.NewPrdService(prdRepository)
	prdHandler := handler.NewPrdHandler(prdService)

	featureRepository := repository.NewFeatureRepository(db)
	featureService := service.NewFeatureService(featureRepository)
	featureHandler := handler.NewFeatureHandler(featureService)

	featureDetailRepository := repository.NewFeatureDetailRepository(db)
	featureDetailService := service.NewFeatureDetailService(featureDetailRepository)
	featureDetailHandler := handler.NewFeatureDetailHandler(featureDetailService)

	sessionWebHandler := webHandler.NewSessionHandler(userService, authService)
	userWebHandler := webHandler.NewUserHandler(userService)
	featureWebHandler := webHandler.NewFeatureHandler(featureService)
	featureDetailWebHandler := webHandler.NewFeatureDetailHandler(featureDetailService, featureService)
	prdWebHandler := webHandler.NewPRDHandler(prdService)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(CORSMiddleware())
	cookieStore := cookie.NewStore([]byte(auth.SECRET_KEY))
	router.Use(sessions.Sessions("documentations", cookieStore))
	router.LoadHTMLGlob("web/templates/**/*")
	router.Static("css", "./web/assets/css")
	router.Static("js", "./web/assets/js")
	router.Static("vuejs", "./web/assets/vue/vue")
	router.Static("other", "./web/assets/vue")
	router.Static("webfonts", "./web/assets/webfonts")
	router.Static("codemirror", "./web/assets/codemirror")
	router.Static("/icon", "./web/assets/icon")
	api := router.Group("/api/v1")

	api.POST("/login", userHandler.Login)
	api.POST("/register", userHandler.CreateUser)

	api.GET("/users", authMiddleware(authService, userService), userHandler.FetchUser)

	api.GET("/features", authMiddleware(authService, userService), featureHandler.GetFeatures)
	api.GET("/features/:id", authMiddleware(authService, userService), featureHandler.GetFeature)
	api.POST("/features", authMiddleware(authService, userService), featureHandler.CreateFeature)
	api.PUT("/features/:id", authMiddleware(authService, userService), featureHandler.UpdateFeature)
	api.DELETE("/features/:id", authMiddleware(authService, userService), featureHandler.DeleteFeature)
	api.GET("/features_ssr", authMiddleware(authService, userService), featureHandler.GetFeaturesSSR)

	api.GET("/feature_details/:id", authMiddleware(authService, userService), featureDetailHandler.GetFeatureDetails)
	api.GET("/feature_detail/:id", authMiddleware(authService, userService), featureDetailHandler.GetFeatureDetail)
	api.POST("/feature_detail", authMiddleware(authService, userService), featureDetailHandler.CreateFeatureDetail)
	api.PUT("/feature_detail/:id", authMiddleware(authService, userService), featureDetailHandler.UpdateFeatureDetail)
	api.DELETE("/feature_detail/:id", authMiddleware(authService, userService), featureDetailHandler.DeleteFeatureDetail)
	api.GET("/features_detail_ssr", authMiddleware(authService, userService), featureDetailHandler.GetFeaturesSSR)
	api.GET("/prd", authMiddleware(authService, userService), prdHandler.GetPrds)

	api.GET("/prd_ssr", authMiddleware(authService, userService), prdHandler.GetPRDSR)

	// WEB HANDLER
	router.GET("/", sessionWebHandler.New)
	router.GET("/dashboard", authAdminMiddleWare(), sessionWebHandler.Dashboard)
	router.POST("/sessions", sessionWebHandler.LoginAction)
	router.GET("/logout", sessionWebHandler.Logout)

	router.GET("/users", authAdminMiddleWare(), userWebHandler.Index)
	router.GET("/users/new", authAdminMiddleWare(), userWebHandler.New)
	router.GET("/users/update/:id", authAdminMiddleWare(), userWebHandler.Update)
	router.POST("/users", authAdminMiddleWare(), userWebHandler.Create)
	router.POST("/users/update_action/:id", authAdminMiddleWare(), userWebHandler.UpdateAction)
	router.GET("/users/delete/:id", authAdminMiddleWare(), userWebHandler.Delete)

	router.GET("/features", authAdminMiddleWare(), featureWebHandler.Index)
	router.GET("/features/new", authAdminMiddleWare(), featureWebHandler.New)
	router.GET("/features/update/:id", authAdminMiddleWare(), featureWebHandler.Update)
	router.GET("/features/delete/:id", authAdminMiddleWare(), featureWebHandler.Delete)
	router.POST("/features", authAdminMiddleWare(), featureWebHandler.Create)
	router.POST("/features/update_action/:id", authAdminMiddleWare(), featureWebHandler.UpdateAction)

	router.GET("/features_detail", authAdminMiddleWare(), featureDetailWebHandler.Index)
	router.GET("/features_detail/new", authAdminMiddleWare(), featureDetailWebHandler.New)
	router.GET("/features_detail/update/:id", authAdminMiddleWare(), featureDetailWebHandler.Update)
	router.GET("/features_detail/delete/:id", authAdminMiddleWare(), featureDetailWebHandler.Delete)
	router.POST("/features_detail", authAdminMiddleWare(), featureDetailWebHandler.Create)
	router.POST("/features_detail/update_action/:id", authAdminMiddleWare(), featureDetailWebHandler.UpdateAction)

	router.GET("/prds", authAdminMiddleWare(), prdWebHandler.Index)
	router.GET("/prds/new", authAdminMiddleWare(), prdWebHandler.New)
	router.GET("/prds/update/:id", authAdminMiddleWare(), prdWebHandler.Update)
	router.GET("/prds/delete/:id", authAdminMiddleWare(), prdWebHandler.Delete)
	router.POST("/prds", authAdminMiddleWare(), prdWebHandler.Create)
	router.POST("/prds/update_action/:id", authAdminMiddleWare(), prdWebHandler.UpdateAction)

	router.Run()
}

func authAdminMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		userIDSession := session.Get("userID")
		if userIDSession == nil {
			c.Redirect(http.StatusFound, "/")
			return
		}
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT,DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func authMiddleware(authService auth.Service, userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.ApiResponse("UnAuthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")

		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.ApiResponse(err.Error(), http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.ApiResponse("UnAuthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))
		var input input.InputIDUser
		input.ID = userID
		user, err := userService.UserServiceGetByID(input)

		if err != nil {
			response := helper.ApiResponse("UnAuthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)

	}

}
