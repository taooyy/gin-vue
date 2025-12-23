# ⚙️ 校园食材供应链平台 - 后端 (Go + Gin)

欢迎来到校园食材供应链平台的后端项目！本项目是整个全栈应用的业务逻辑和数据服务层。

---

## 1. 项目概览 (Project Overview)

本项目采用 Go 语言和 Gin 框架构建，提供 RESTful API 服务，支持前端应用的数据交互和业务处理。核心功能包括用户认证、商品供应链管理、订单交易、全链路溯源以及财务结算等。

---

## 2. 技术栈 (Technology Stack)

*   **语言**: Go (最新稳定版)
*   **Web 框架**: Gin
*   **ORM**: GORM
*   **数据库**: MySQL 8.0
*   **配置管理**: Viper
*   **依赖管理**: Go Modules

---

## 3. 项目结构 (Project Structure)

```
server/
├── cmd/                    # 应用入口点
│   └── server/             # 主应用
│       └── main.go         # 应用程序启动文件
├── configs/                # 配置文件
│   └── config.yaml         # 环境变量和应用配置
├── internal/               # 私有应用代码 (高内聚、低耦合的核心业务逻辑)
│   ├── config/             # 配置加载模块
│   ├── handler/            # HTTP 请求处理器 (Controller 层)
│   ├── model/              # GORM 数据模型 (与数据库表结构对应)
│   │   ├── system.go       # 系统基础模块的模型定义
│   │   ├── scm.go          # 商品供应链模块的模型定义
│   │   ├── order.go        # 订单交易模块的模型定义
│   │   └── finance.go      # 财务结算模块的模型定义
│   ├── repository/         # 数据访问层 (DAO/Repository 层)
│   ├── router/             # Gin 路由配置
│   └── service/            # 业务逻辑层 (Service 层)
├── pkg/                    # 共享库和工具包
│   └── database/           # 数据库连接和初始化
│       └── mysql.go        # MySQL 数据库连接封装
├── go.mod                  # Go 模块定义文件
├── go.sum                  # Go 模块依赖校验文件
└── 数据库设计.md             # 详细的数据库设计文档
```

---

## 4. 设置与安装 (Setup and Installation)

1.  **进入后端项目目录**:
    ```bash
    cd server
    ```
2.  **安装 Go 依赖**:
    ```bash
    go mod tidy
    ```
    这将下载并安装所有项目所需的 Go 模块依赖。

---

## 5. 数据库设置 (Database Setup)

1.  **MySQL 数据库**: 确保您已安装并运行 MySQL 8.0 或更高版本。
2.  **创建数据库**:
    首先，在您的 MySQL 服务器上创建一个名为 `school_scm` 的数据库。
    ```sql
    CREATE DATABASE IF NOT EXISTS `school_scm` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
    ```
3.  **导入表结构**:
    详细的数据库表结构定义和初始化 SQL 脚本位于 `数据库设计.md` 文件中。您可以使用该文件末尾提供的 SQL 语句来创建所有表。
4.  **配置数据库连接**:
    打开 `configs/config.yaml` 文件，根据您的 MySQL 配置更新连接信息，特别是 `password` 字段。
    ```yaml
    mysql:
      host: "127.0.0.1"
      port: 3306
      user: "root"
      password: "your_mysql_password" # 请替换为您的MySQL密码
      dbname: "school_scm"
      config: "charset=utf8mb4&parseTime=True&loc=Local"
    ```

---

## 6. 可用脚本 (Available Scripts)

在项目目录下，你可以运行以下命令：

*   **`go run ./cmd/server/main.go`**:
    启动后端 API 服务。服务将根据 `configs/config.yaml` 中定义的端口运行 (默认为 `8080`)。

---

## 7. 编码规范 (Coding Standards)

本项目严格遵循 `../GEMINI.md` 中定义的全栈开发规范。请特别注意以下后端相关原则：

*   **架构与分层**: 严格遵守 Go Standard Project Layout (`cmd`, `internal`, `pkg`)。
*   **分层职责**: Controller 仅负责参数解析和响应；Service 包含所有业务逻辑；Repository 仅负责数据库操作。严禁跨层调用和职责混淆。
*   **依赖注入**: 推荐使用构造函数注入依赖，以提高代码的可测试性和解耦性。
*   **错误处理**: 必须显式处理 `error`。
*   **数据库操作**: 必须使用 GORM 进行数据库操作，严禁直接拼接 SQL 字符串。

---

感谢您的贡献！
