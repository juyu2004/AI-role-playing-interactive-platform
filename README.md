AI 角色扮演互动平台

概述
- 基于 AI 大模型的语音角色扮演聊天 Web 应用。
- 前端: React + Vite + TypeScript；后端: Go（主）/ C++（stub）；API: OpenAPI；开发容器: docker-compose。

目录结构
```
frontend/        # Web 前端 (Vite + React + TS)
backend-go/      # Go 后端 (主服务)
backend-cpp/     # C++ 后端 stub（可选扩展）
api/             # OpenAPI 规范与共享 schema
infra/           # DevOps、Dockerfile 等
```

快速开始（开发）
1) 前端
```
cd frontend
npm install
npm run dev
```

2) Go 后端
```
cd backend-go
go run ./cmd/server
```

3) Docker（可选）
```
docker-compose up --build
```

里程碑
- M1: 需求/设计、架构与 UI 草图
- M2: 前后端基础框架与首个角色聊天闭环
- M3: 功能完善、性能优化与测试
- M4: 部署上线

许可证
MIT
✨ 本项目旨在利用先进的AI大模型技术，打造一个平台，让用户能够与他们喜爱的虚拟角色（如小说人物、历史人物、影视角色等）进行实时、自然的语音对话，获得独特的互动体验。✨
