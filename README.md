# 极客博客系统 (Geek Blog)

一个使用 React + Gin + MongoDB 构建的现代化博客系统，支持 Markdown 和 LaTeX 渲染，具有极客风格的界面设计。

## 特性

- 📝 支持 Markdown 和 LaTeX 数学公式
- 🖼️ 图片上传功能
- 💬 评论系统（支持嵌套回复）
- 🎨 极客风格的终端界面设计
- 🔍 文章标签和分类
- 📱 响应式设计
- 🐳 Docker 容器化部署

## 技术栈

### 后端
- Go 1.21+
- Gin Web Framework
- MongoDB
- JWT 认证

### 前端
- React 18
- React Router v6
- React Markdown
- KaTeX (数学公式渲染)
- Highlight.js (代码高亮)

### 部署
- Docker & Docker Compose
- Nginx 反向代理

## 本地开发

### 前置要求
- Go 1.21+
- Node.js 18+
- MongoDB
- Docker (可选)

### 后端启动

```bash
cd backend
go mod download
go run main.go

cd frontend
npm install
npm start

###前端启动

```bash
cd frontend
npm install
npm start
