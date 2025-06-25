package config

import (
    "context"
    "log"
    "os"
    "time"
    
    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
    Port        string
    MongoURI    string
    DBName      string
    JWTSecret   string
    UploadPath  string
}

var (
    cfg *Config
    DB  *mongo.Database
)

func InitConfig() {
    godotenv.Load()
    
    cfg = &Config{
        Port:       getEnv("PORT", "8080"),
        MongoURI:   getEnv("MONGO_URI", "mongodb://localhost:27017"),
        DBName:     getEnv("DB_NAME", "geekblog"),
        JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),
        UploadPath: getEnv("UPLOAD_PATH", "./uploads"),
    }
    
    // 创建上传目录
    os.MkdirAll(cfg.UploadPath, os.ModePerm)
}

func ConnectDB() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }
    
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal("Failed to ping MongoDB:", err)
    }
    
    DB = client.Database(cfg.DBName)
    log.Println("Connected to MongoDB")
}

func GetConfig() *Config {
    return cfg
}

func GetDB() *mongo.Database {
    return DB
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}