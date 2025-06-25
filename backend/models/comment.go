package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    PostID    primitive.ObjectID `bson:"post_id" json:"post_id"`
    Author    string            `bson:"author" json:"author"`
    Email     string            `bson:"email" json:"email"`
    Content   string            `bson:"content" json:"content"`
    CreatedAt time.Time         `bson:"created_at" json:"created_at"`
    ParentID  *primitive.ObjectID `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
}

type CreateCommentRequest struct {
    PostID   string  `json:"post_id" binding:"required"`
    Author   string  `json:"author" binding:"required"`
    Email    string  `json:"email" binding:"required,email"`
    Content  string  `json:"content" binding:"required"`
    ParentID *string `json:"parent_id,omitempty"`
}