package controllers

import (
    "fmt"
    "context"
    "net/http"
    "time"
    
    "geek-blog/config"
    "geek-blog/models"
    
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func GetPosts(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    collection := config.GetDB().Collection("posts")
    
    // 分页参数
    page := 1
    limit := 10
    if p := c.Query("page"); p != "" {
        page = parseInt(p, 1)
    }
    if l := c.Query("limit"); l != "" {
        limit = parseInt(l, 10)
    }
    
    skip := (page - 1) * limit
    
    // 查询选项
    findOptions := options.Find()
    findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
    findOptions.SetSkip(int64(skip))
    findOptions.SetLimit(int64(limit))
    
    // 查询条件
    filter := bson.M{}
    if tag := c.Query("tag"); tag != "" {
        filter["tags"] = tag
    }
    
    cursor, err := collection.Find(ctx, filter, findOptions)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer cursor.Close(ctx)
    
    var posts []models.Post
    if err := cursor.All(ctx, &posts); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // 获取总数
    total, _ := collection.CountDocuments(ctx, filter)
    
    c.JSON(http.StatusOK, gin.H{
        "posts": posts,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}

func GetPost(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    id := c.Param("id")
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }
    
    collection := config.GetDB().Collection("posts")
    
    var post models.Post
    err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
        return
    }
    
    // 增加浏览量
    collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$inc": bson.M{"view_count": 1}})
    
    c.JSON(http.StatusOK, post)
}

func CreatePost(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    var req models.CreatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    post := models.Post{
        Title:      req.Title,
        Content:    req.Content,
        Author:     req.Author,
        Tags:       req.Tags,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
        ViewCount:  0,
        CommentIDs: []primitive.ObjectID{},
    }
    
    collection := config.GetDB().Collection("posts")
    result, err := collection.InsertOne(ctx, post)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    post.ID = result.InsertedID.(primitive.ObjectID)
    c.JSON(http.StatusCreated, post)
}

func UpdatePost(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    id := c.Param("id")
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }
    
    var req models.UpdatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    update := bson.M{
        "$set": bson.M{
            "updated_at": time.Now(),
        },
    }
    
    if req.Title != "" {
        update["$set"].(bson.M)["title"] = req.Title
    }
    if req.Content != "" {
        update["$set"].(bson.M)["content"] = req.Content
    }
    if req.Tags != nil {
        update["$set"].(bson.M)["tags"] = req.Tags
    }
    
    collection := config.GetDB().Collection("posts")
    result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    if result.MatchedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

func DeletePost(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    id := c.Param("id")
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }
    
    collection := config.GetDB().Collection("posts")
    result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    if result.DeletedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
        return
    }
    
    // 删除相关评论
    commentCollection := config.GetDB().Collection("comments")
    commentCollection.DeleteMany(ctx, bson.M{"post_id": objectID})
    
    c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

func parseInt(s string, defaultValue int) int {
    var v int
    fmt.Sscanf(s, "%d", &v)
    if v == 0 {
        return defaultValue
    }
    return v
}