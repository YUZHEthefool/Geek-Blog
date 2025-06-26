# 构建阶段
FROM golang:1.21-alpine AS builder
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod ./
COPY go.sum ./ 
 # 明确复制 go.sum 是个好习惯

# 设置 Go 模块代理为国内源
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 下载 go.mod 中定义的依赖
RUN go mod download

# 复制所有源代码
COPY . .

# 整理依赖，确保 go.mod 和 go.sum 是最新的
RUN go mod tidy

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 运行阶段
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/main .
RUN mkdir -p /app/uploads
EXPOSE 8080
CMD ["./main"]
