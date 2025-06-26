# --- 构建阶段 ---
FROM node:18-alpine AS builder

WORKDIR /app

# 为了更好地利用缓存，先只复制 package.json 和 package-lock.json
COPY package*.json ./

# 设置国内镜像源并安装依赖
# 这一步只依赖 package-lock.json，只要它不变，这一层就会被缓存
RUN npm config set registry https://registry.npmmirror.com && npm install

# 现在，复制项目的所有其他文件
# (包括 src, public 等)
COPY . .

# 在所有文件都就位后，运行构建命令
RUN npm run build


# --- 运行阶段 ---
FROM nginx:alpine

# 从构建阶段复制出编译好的静态文件
COPY --from=builder /app/build /usr/share/nginx/html

# 复制你的 Nginx 配置文件
# (确保 nginx-frontend.conf 在 ./frontend 目录下)
COPY nginx-frontend.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
