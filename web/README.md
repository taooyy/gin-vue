# 🌐 校园食材供应链平台 - 前端 (Vue 3 + Element Plus)

欢迎来到校园食材供应链平台的前端项目！本项目是整个全栈应用的用户界面部分，旨在提供一个高效、直观的操作体验。

---

## 1. 技术栈 (Technology Stack)

*   **框架**: Vue 3 (Composition API, Script Setup)
*   **语言**: TypeScript 5.x
*   **UI 组件库**: Element Plus
*   **状态管理**: Pinia
*   **构建工具**: Vite
*   **路由**: Vue Router 4.x
*   **样式预处理器**: SCSS

---

## 2. 项目结构 (Project Structure)

```
web/
├── public/                 # 静态资源 (例如 favicon, robots.txt)
├── src/                    # 源代码目录
│   ├── assets/             # 静态资源 (图片、字体等)
│   ├── composables/        # 可复用组合式函数 (逻辑抽离)
│   ├── layouts/            # 页面布局组件 (例如 PlatformLayout, TenantLayout)
│   │   └── components/     # 布局相关的子组件 (例如 Header, Sidebar)
│   ├── router/             # Vue Router 配置
│   ├── stores/             # Pinia 状态管理模块
│   ├── styles/             # 全局样式和主题配置
│   │   ├── element/        # Element Plus 变量覆盖 (主题定制)
│   │   └── main.scss       # 全局样式入口
│   ├── types/              # TypeScript 类型定义
│   │   ├── api/            # 后端 API 数据类型 (待完善)
│   │   └── config.ts       # 全局配置类型 (例如 ROLES, MenuItem)
│   ├── views/              # 页面组件
│   │   ├── auth/           # 认证相关页面 (例如 Login)
│   │   ├── platform/       # 平台管理页面
│   │   └── workspace/      # 工作区页面 (学校/供应商/食堂)
│   ├── App.vue             # 根组件
│   ├── main.ts             # 应用入口文件
│   └── vite-env.d.ts       # Vite 环境声明文件
├── .env                    # 环境变量配置
├── package.json            # 项目依赖和脚本
├── tsconfig.json           # TypeScript 配置
└── vite.config.ts          # Vite 配置
```

---

## 3. 设置与安装 (Setup and Installation)

1.  **进入前端项目目录**:
    ```bash
    cd web
    ```
2.  **安装依赖**:
    ```bash
    npm install
    ```

---

## 4. 可用脚本 (Available Scripts)

在项目目录下，你可以运行以下命令：

*   **`npm run dev`**:
    在开发模式下运行应用。在浏览器中打开 `http://localhost:5173` (或其他Vite分配的端口) 查看。热重载功能将被启用。
*   **`npm run build`**:
    为生产环境构建应用。构建产物将输出到 `dist/` 目录。
*   **`npm run preview`**:
    在本地预览生产构建的应用。这不会启动一个完整的服务器，只是用 Vite 提供的简单服务器来预览 `dist/` 目录的内容。

---

## 5. 编码规范 (Coding Standards)

本项目严格遵循 `../GEMINI.md` 中定义的全栈开发规范。请特别注意以下前端相关原则：

*   **UI 组件使用**: 优先使用 Element Plus。所有自定义样式和主题配置应在 `styles/element/` 或 `assets/scss` 中进行。
*   **代码组织**: 复杂的业务逻辑应提取为独立的 Composable (`use*.ts`)。`.vue` 文件应专注于视图渲染。
*   **类型安全**: 必须使用 TypeScript 接口或类型定义。API 响应数据的类型定义应放在 `types/api/` 目录下。
*   **状态管理**: 使用 Pinia 管理全局状态，局部状态保留在组件内部 `ref` 中。
*   **换行符**: 统一使用 `LF` (Linux 风格)。

---

感谢您的贡献！