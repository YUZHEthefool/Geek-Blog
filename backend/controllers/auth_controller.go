package controllers

import (
    "net/http"
    "geek-blog/utils"
    
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // 这里简化处理，实际应该从数据库验证
    if req.Username != "admin" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
    
    // 验证密码（这里使用硬编码，实际应该从数据库获取hash）
    hashedPassword := "$2a$10$N9qo8uLOickgx2ZMRZoHK.6YuVGCTEH5.h5VvF5X8z5kANF.sX4Ey" // password: admin123
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
    
    // 生成 JWT
    token, err := utils.GenerateToken(req.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user": req.Username,
    })
}