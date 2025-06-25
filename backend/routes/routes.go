package routes

import (
    "geek-blog/controllers"
    "geek-blog/middleware"
    "github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
    // 中间件
    r.Use(middleware.CORS())
    
    // 静态文件
    r.Static("/uploads", "./uploads")
    
    // API 路由
    api := r.Group("/api")
    {
        // 文章路由
        api.GET("/posts", controllers.GetPosts)
        api.GET("/posts/:id", controllers.GetPost)
        api.POST("/posts", controllers.CreatePost)
        api.PUT("/posts/:id", controllers.UpdatePost)
        api.DELETE("/posts/:id", controllers.DeletePost)
        
        // 评论路由
        api.GET("/comments", controllers.GetComments)
        api.POST("/comments", controllers.CreateComment)
        api.DELETE("/comments/:id", controllers.DeleteComment)
        
        // 上传路由
        api.POST("/upload", controllers.UploadImage)
    }
}