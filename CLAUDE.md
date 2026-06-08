# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

TodoList（todolist）— 一个现代化、功能丰富的待办任务管理应用，支持双存储模式：离线模式（IndexedDB）和在线模式（Go 后端 + PostgreSQL）。中文名"轻量待办任务"。

## 技术栈

- **前端 (web/)**：React 18 + TypeScript + Vite + Tailwind CSS + shadcn/ui（Radix 原语）
- **状态管理**：Zustand（store）+ React Context（Auth、Project、Task、Sidebar）
- **数据请求**：@tanstack/react-query（在线模式）
- **后端 (server/)**：Go + Gin + PostgreSQL + sqlc（类型安全 SQL 代码生成）+ golang-migrate
- **桌面端 (desktop/)**：Wails v2（Go + Web 前端封装）
- **包管理**：pnpm（前端）、Go modules（后端）
- **测试**：Vitest（前端，支持 fast-check 属性测试）、Go test（后端）
- **Markdown 编辑器**：Milkdown（基于 ProseMirror），用于任务描述
- **拖拽排序**：@hello-pangea/dnd

## 常用命令

### 前端（在 `web/` 目录下执行）
```bash
pnpm dev              # 启动开发服务器（端口 8080）
pnpm build            # 生产构建
pnpm build:dev        # 开发构建
pnpm lint             # 运行 ESLint
pnpm test             # 运行测试（vitest run）
pnpm test:watch       # 监听模式运行测试
pnpm test:ui          # Vitest 可视化界面
```

### 后端（在 `server/` 目录下执行）
```bash
make dev              # 热重载启动（需要 air）
make run              # 直接运行
make build            # 构建二进制 → bin/todolist-server
make test             # 运行 Go 测试
make test-coverage    # 带覆盖率报告的测试
make lint             # golangci-lint 检查
make fmt              # go fmt + goimports 格式化
make sqlc             # 从 db/queries/ 生成 sqlc 代码
make migrate-up       # 执行数据库迁移
make migrate-down     # 回滚迁移
```

### 桌面端（Wails）
```bash
cd web && pnpm wails:dev          # Wails 开发模式
cd web && pnpm wails:build        # Wails 构建
cd web && pnpm wails:build:prod   # Wails 生产构建
```

### Docker（根目录下执行）
```bash
make compose-up        # 启动 docker-compose（docker/docker-compose.yml）
make compose-down      # 停止
make compose-logs      # 查看日志
make docker-build      # 构建 Docker 镜像
make docker-all        # 构建 + 测试 + 推送
```

### 根目录快捷命令
```bash
make dev       # cd web && pnpm dev
make build     # cd web && pnpm build
make test      # cd web && pnpm test
make dev-all   # 同时启动前端和后端（scripts/dev.sh）
```

## 架构

### 双存储模式（核心模式）

应用支持两种存储后端，运行时通过 `VITE_STORAGE_MODE` 环境变量或 localStorage 键 `todo_storage_mode` 选择：

- **离线模式** (`offline`)：`IndexedDBAdapter` — 所有数据存储在本地 IndexedDB
- **在线模式** (`online`，默认)：`OnlineStorageAdapter` — 数据通过 Go 后端 API 存储

两者都实现 `StorageAdapter` 接口（`web/src/storage/types.ts`）。工厂函数在 `web/src/storage/index.ts` 中根据配置创建单例。所有数据访问必须通过 `storage/operations.ts`，不要直接访问适配器。

切换模式需要刷新页面（单例模式）。配置优先级：localStorage > 环境变量 > 默认值（`online`）。

### 前端数据流

```
页面 (src/pages/) → 上下文 (src/contexts/) → 存储操作 (src/storage/operations.ts) → StorageAdapter
                                                                                   ↕
                                                                         在线：apiClient (src/lib/apiClient.ts)
                                                                         离线：IndexedDB
```

- **Pages**：路由级组件（Index、Auth、Settings、Pomodoro、SearchResults 等）
- **Contexts**：AuthContext、ProjectContext、TaskContext/TaskProvider、SidebarContext — 向组件提供领域状态
- **Zustand Stores**：`taskStore`、`projectStore`、`userProfileStore` — 轻量状态切片
- **Hooks**：业务逻辑隔离在 hooks 中（`useTaskOperation`、`useTaskFilter`、`usePomodoroTimer` 等）
- **Services**：跨领域关注点（`dataTransferService`、`deadlineService`、`projectMemberService`）

### 后端分层架构

```
handler（HTTP 处理）→ service（业务逻辑）→ repository（数据访问，调用 sqlc）→ PostgreSQL
```

- **handler/**：请求解析、参数校验、响应构造
- **service/**：业务逻辑、事务编排、权限校验
- **repository/**：调用 sqlc 生成的代码，不在此层写原始 SQL
- **model/**：领域实体、DTO、请求/响应结构体
- **db/queries/**：原始 SQL 文件 — sqlc 从此生成 Go 代码
- **db/migrations/**：版本化 SQL 迁移（golang-migrate 格式）

核心设计：后端提供聚合查询接口（如 `/api/v1/overview`）避免前端 N+1 问题。前端不应循环请求每个项目的数据。

### 组件组织

- `src/components/ui/` — shadcn/ui 基础组件（基于 Radix），与业务逻辑解耦
- `src/components/tasks/` — 任务相关业务组件（TaskList、TaskDetail、TaskItem、编辑器、筛选器等）
- `src/components/projects/` — 项目对话框（编辑、分享、加入）
- `src/components/sidebar/` — 侧边栏导航
- `src/components/attachments/` — 文件附件处理
- `src/components/checkin/` — 每日签到功能
- `src/components/search/` — 搜索界面
- `src/components/settings/` — 设置面板

路径别名：`@/` 映射到 `web/src/`。

## 编码规范

- **语言**：代码和注释使用英文；Git 提交信息使用英文，遵循 Conventional Commits 格式
- **TypeScript**：导出 API 使用显式类型；局部可推断变量可省略类型。禁止 `any` 和类型断言逃逸，除非有明确理由且局部封装
- **命名**：使用全词而非缩写（如 `useTaskOperation`、`formatDateText`）。函数用动词，变量用名词
- **组件**：受控组件与单向数据流。复杂 UI 抽取为组件，可复用逻辑抽取为 hooks
- **加载状态**：使用骨架屏（Skeleton），严禁使用"加载中..."等文字提示
- **注释**：只解释"为什么"，不解释"做了什么"。删除过时注释
- **Git 提交**：Conventional Commits 格式。多段落信息使用多个 `-m` 参数，不使用 `\n` 转义换行
- **DRY**：重复 UI → 抽取组件；重复逻辑 → 抽取 hook 或工具函数
- **早返回**：使用守卫式写法，避免深层嵌套
- **列表渲染**：使用稳定 key，避免不必要的 re-render（合理使用 memo/useMemo/useCallback）

## Supabase 迁移

在 `supabase/migrations/` 中创建迁移文件时，遵循命名规范：`YYYYMMDDHHmmss_short_description.sql`（UTC 时间戳）。新建表必须启用 RLS。每个操作（select/insert/update/delete）每个角色（anon/authenticated）创建独立策略，不要合并策略。SQL 全部小写。对破坏性操作添加充分注释。

## 环境变量

前端（在 `web/` 目录下）：
- `VITE_STORAGE_MODE` — `online`（默认）或 `offline`
- `VITE_SUPABASE_URL`、`VITE_SUPABASE_ANON_KEY` — 使用 Supabase 时配置

后端（在 `server/` 目录下）：
- 完整列表见 `server/.env.example`（数据库、JWT、CORS、日志配置）
- 配置优先级：环境变量 > `.env` 文件 > 代码默认值
