package controllers

import (
    "context"
    "net/http"
    "time"
    
    "geek-blog/config"
    "geek-blog/models"
    
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func GetComments(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    postID := c.Query("post_id")
    if postID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
        return
    }
    
    objectID, err := primitive.ObjectIDFromHex(postID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post_id"})
        return
    }
    
    collection := config.GetDB().Collection("comments")
    cursor, err := collection.Find(ctx, bson.M{"post_id": objectID})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer cursor.Close(ctx)
    
    var comments []models.Comment
    if err := cursor.All(ctx, &comments); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, comments)
}

func CreateComment(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    var req models.CreateCommentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    postID, err := primitive.ObjectIDFromHex(req.PostID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post_id"})
        return
    }
    
    comment := models.Comment{
        PostID:    postID,
        Author:    req.Author,
        Email:     req.Email,
        Content:   req.Content,
        CreatedAt: time.Now(),
    }
    
    if req.ParentID != nil {
        parentID, err := primitive.ObjectIDFromHex(*req.ParentID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parent_id"})
            return
        }
        comment.ParentID = &parentID
    }
    
    collection := config.GetDB().Collection("comments")
    result, err := collection.InsertOne(ctx, comment)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    comment.ID = result.InsertedID.(primitive.ObjectID)
    
    // 更新文章的评论ID列表
    postCollection := config.GetDB().Collection("posts")
    postCollection.UpdateOne(ctx, bson.M{"_id": postID}, bson.M{
        "$push": bson.M{"comment_ids": comment.ID},
    })
    
    c.JSON(http.StatusCreated, comment)
}

func DeleteComment(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    id := c.Param("id")
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }
    
    collection := config.GetDB().Collection("comments")
    
    // 获取评论信息
    var comment models.Comment
    err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&comment)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
        return
    }
    
    // 删除评论
    result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    if result.DeletedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
        return
    }
    
    // 从文章中移除评论ID
    postCollection := config.GetDB().Collection("posts")
    postCollection.UpdateOne(ctx, bson.M{"_id": comment.PostID}, bson.M{
        "$pull": bson.M{"comment_ids": objectID},
    })
    
    c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}