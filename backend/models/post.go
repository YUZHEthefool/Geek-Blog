package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
    ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    Title       string              `bson:"title" json:"title"`
    Content     string              `bson:"content" json:"content"`
    Author      string              `bson:"author" json:"author"`
    Tags        []string            `bson:"tags" json:"tags"`
    CreatedAt   time.Time           `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time           `bson:"updated_at" json:"updated_at"`
    ViewCount   int                 `bson:"view_count" json:"view_count"`
    CommentIDs  []primitive.ObjectID `bson:"comment_ids" json:"comment_ids"`
}

type CreatePostRequest struct {
    Title   string   `json:"title" binding:"required"`
    Content string   `json:"content" binding:"required"`
    Author  string   `json:"author" binding:"required"`
    Tags    []string `json:"tags"`
}

type UpdatePostRequest struct {
    Title   string   `json:"title"`
    Content string   `json:"content"`
    Tags    []string `json:"tags"`
}