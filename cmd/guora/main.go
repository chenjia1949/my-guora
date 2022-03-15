package main

import (
	"my-guora/conf"
	"my-guora/internal/controller/rest"
	"my-guora/internal/controller/web"
	"my-guora/internal/middleware"
	"my-guora/internal/view"

	"github.com/gin-gonic/gin"
)

// SetupAPIRouter func
func SetupAPIRouter(r *gin.Engine) {

	r.Use(middleware.Logger())

	GroupAPI := r.Group("/api")
	{

		// Group: web
		GroupWeb := GroupAPI.Group("/web")
		{
			GroupWeb.POST("/security/sign", web.SecuritySign)
			GroupWeb.POST("/security/login", web.SecurityLogin)
			GroupWeb.POST("/security/logout", web.SecurityLogout)

			GroupWeb.POST("/question", middleware.Authorizer(), web.CreateQuestion)
			GroupWeb.POST("/answer", middleware.Authorizer(), web.CreateAnswer)
			GroupWeb.POST("/answer/:id/supporters", middleware.Authorizer(), web.CreateSupporter)
			GroupWeb.DELETE("/answer/:id/supporters", middleware.Authorizer(), web.DeleteSupporter)
			GroupWeb.POST("/comment", middleware.Authorizer(), web.CreateComment)
			GroupWeb.POST("/reply", middleware.Authorizer(), web.CreateReply)
			GroupWeb.POST("/file/avatar", middleware.Authorizer(), web.FileAvatarResolve)
		}

		// Group: rest
		GroupRest := GroupAPI.Group("/rest")
		GroupRest.Use(middleware.Authorizer()) //指定所有的rest接口都需要校验token
		{

			GroupRest.GET("/user/:id", rest.GetUser)
			GroupRest.GET("/users", rest.GetUsers)
			GroupRest.GET("/users/counts", rest.GetUsersCounts)
			GroupRest.PUT("/user/:id", rest.UpdateUser)
			GroupRest.DELETE("/user/:id", rest.DeleteUser)

			GroupRest.GET("/profile/:id", rest.GetProfile)
			GroupRest.GET("/profiles", rest.GetProfiles)
			GroupRest.GET("/profiles/counts", rest.GetProfilesCounts)
			GroupRest.PUT("/profile/:id", rest.UpdateProfile)
			GroupRest.DELETE("/profile/:id", rest.DeleteProfile)

			GroupRest.GET("/question/:id", rest.GetQuestion)
			GroupRest.GET("/questions", rest.GetQuestions)
			GroupRest.GET("/questions/counts", rest.GetQuestionsCounts)
			GroupRest.PUT("/question/:id", rest.UpdateQuestion)
			GroupRest.DELETE("/question/:id", rest.DeleteQuestion)

			GroupRest.GET("/answer/:id", rest.GetAnswer)
			GroupRest.GET("/answers", rest.GetAnswers)
			GroupRest.GET("/answers/counts", rest.GetAnswersCounts)
			GroupRest.PUT("/answer/:id", rest.UpdateAnswer)
			GroupRest.DELETE("/answer/:id", rest.DeleteAnswer)

			GroupRest.GET("/comment/:id", rest.GetComment)
			GroupRest.GET("/comments", rest.GetComments)
			GroupRest.GET("/comments/counts", rest.GetCommentsCounts)
			GroupRest.PUT("/comment/:id", rest.UpdateComment)
			GroupRest.DELETE("/comment/:id", rest.DeleteComment)

			GroupRest.GET("/reply/:id", rest.GetReply)
			GroupRest.GET("/replies", rest.GetReplies)
			GroupRest.GET("/replies/counts", rest.GetRepliesCounts)
			GroupRest.PUT("/reply/:id", rest.UpdateReply)
			GroupRest.DELETE("/reply/:id", rest.DeleteReply)
		}
	}
}

// SetupViewRouter func
func SetupViewRouter(r *gin.Engine) {

	// Default Group: view
	{
		r.Delims("\"/{{", "}}/\"")
		r.LoadHTMLGlob("../../web/*.html")
		r.Static("/static", "../../web/static")

		r.GET("/", middleware.Authorizer(), view.Index)
		r.GET("/profile", middleware.Authorizer(), view.Profile)
		r.GET("/question", middleware.Authorizer(), view.Question)
		r.GET("/answer", middleware.Authorizer(), view.Answer)
		r.GET("/admin", middleware.Authorizer(), middleware.Administrator(), view.Admin)
		r.GET("/login", view.Login)
		r.GET("/error", view.Error)
	}

}

func main() {
	/* var shouldInit = flag.Bool("init", false, "initialize all")
	flag.Parse()

	if *shouldInit {
		initAll(conf.Config())
	} */

	r := gin.Default()
	SetupAPIRouter(r)
	SetupViewRouter(r)

	r.Run(conf.Config().Address)
}
