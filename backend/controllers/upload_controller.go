package controllers

import (
    "fmt"
    "net/http"
    "path/filepath"
    "time"
    
    "geek-blog/config"
    "github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
    file, err := c.FormFile("image")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
        return
    }
    
    // 检查文件类型
    ext := filepath.Ext(file.Filename)
    allowedExts := map[string]bool{
        ".jpg":  true,
        ".jpeg": true,
        ".png":  true,
        ".gif":  true,
        ".webp": true,
    }
    
    if !allowedExts[ext] {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
        return
    }
    
    // 生成唯一文件名
    filename := fmt.Sprintf("%d%s", time.Now().Unix(), ext)
    uploadPath := filepath.Join(config.GetConfig().UploadPath, filename)
    
    // 保存文件
    if err := c.SaveUploadedFile(file, uploadPath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        return
    }
    
    // 返回文件URL
    fileURL := fmt.Sprintf("/uploads/%s", filename)
    c.JSON(http.StatusOK, gin.H{
        "url":      fileURL,
        "filename": filename,
    })
}