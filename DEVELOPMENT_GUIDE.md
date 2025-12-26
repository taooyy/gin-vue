# 食材供应链平台 - 项目架构与开发指南

本文档旨在为“食材供应链平台”项目提供一份清晰的架构解析、开发流程指南和未来迭代的架构性建议，以确保项目在多人协作和长期演进中保持高质量、高效率和高可维护性。

## 1. 项目愿景与核心概念

在开始之前，我们首先要明确项目的核心业务模型，所有技术决策都应围绕此模型展开。

- **核心业务**: 校园食材采购。
- **参与角色**: 平台、学校（站点）、供应商、食堂、商户。
- **核心层级**:
    1.  **平台**: 最高管理者，负责创建和管理**学校（站点）**。
    2.  **学校**: 平台下的租户，负责管理自己合作的**供应商**和**食堂**。
    3.  **供应商/食堂**: 学校下的二级租户，权限平级，各自拥有独立的后台功能。
- **关键特性**:
    - **功能下发**: 平台开发新功能后，可以“下发”给指定的学校。
    - **权限分配**: 学校管理员可以决定将哪些功能分配给自己的内部账号。
    - **数据溯源**: 全链路追踪食材从供应商到餐桌的每一个环节，是项目的核心价值之一。

## 2. 项目结构详解

项目采用前后端分离架构，目录结构清晰。

### 2.1. 后端 (`server/`)

后端采用 Go + Gin 框架，遵循经典的三层架构，实现了良好的职责分离。

-   `cmd/server/main.go`: **项目入口**。负责加载配置、初始化数据库连接、启动Gin服务器。
-   `internal/router/router.go`: **路由与依赖注入中心**。
    -   **路由**: 定义所有 API 端点（如 `/api/v1/accounts`）。
    -   **依赖注入**: 在此文件中，我们手动完成了 `Repository` -> `Service` -> `Handler` 的依赖链组装，确保了各层之间的解耦。
    -   **中间件**: 在这里为不同的路由组应用认证 (`AuthMiddleware`) 和授权 (`PlatformAdminAuth`) 中间件。
-   `internal/handler/`: **Web 层 (Handler)**。
    -   **职责**: 解析和校验 HTTP 请求参数，调用 `Service` 层处理业务，并将结果封装成 JSON 返回给前端。
    -   **原则**: **严禁**在此层编写任何业务逻辑。
-   `internal/service/`: **业务逻辑层 (Service)**。
    -   **职责**: 实现所有核心业务逻辑（如创建账号、处理订单等）。它可以调用一个或多个 `Repository` 来完成复杂操作。
    -   **原则**: **严禁**在此层直接编写 SQL 或操作数据库。
-   `internal/repository/`: **数据访问层 (Repository/DAO)**。
    -   **职责**: 封装对数据库的 **所有** CRUD 操作。`Service` 层通过调用 `Repository` 的方法来与数据库交互。
    -   **原则**: **严禁**在此层包含任何业务逻辑。只做纯粹的数据读写。
-   `internal/model/`: **数据模型层 (Model)**。
    -   **职责**: 定义 GORM 模型（对应数据库表结构）和业务中用到的常量（如 `OrgType`, `RoleKey`）。

### 2.2. 前端 (`web/`)

前端采用 Vue 3 + Vite + Element Plus + Pinia，是当前业界主流且高效的MVVM架构。

-   `main.ts`: **Vue 应用入口**。负责创建 Vue 实例，并全局注册 Vue Router、Pinia 和 Element Plus。
-   `router/index.ts`: **前端路由中心**。
    -   **路由定义**: 静态定义了所有页面路由，并利用嵌套路由清晰地划分了 `PlatformLayout`（平台布局）和 `TenantLayout`（租户布局）。
    -   **导航守卫**: `router.beforeEach` 是前端的“保安”。它在每次路由跳转前执行，检查用户是否已登录（`userStore.token`是否存在），如果未登录则强制跳转到登录页。
-   `stores/`: **全局状态管理 (Pinia)**。
    -   `user.ts`: 存储当前用户的 Token、角色、组织类型等核心认证信息，并将其持久化到 `sessionStorage`。
    -   `permission.ts`: 负责根据用户角色动态生成侧边栏菜单，实现 UI 级别的权限控制。
-   `views/`: **页面级组件**。每个 `.vue` 文件代表一个完整的页面，如“账号管理”页。
-   `api/`: **API 服务层**。将所有对后端 API 的 HTTP 请求封装成独立的函数（如 `listAccountsApi`），让 `views` 组件只需调用函数，无需关心具体的 URL、请求方法和参数。
-   `layouts/`: **布局组件**。`PlatformLayout.vue` 和 `TenantLayout.vue` 分别定义了平台和租户（学校、供应商等）的通用页面框架（如顶栏、侧边栏），实现了多端UI的隔离。

## 3. 功能完善路线图与开发指南

本章节将替代宽泛的示例，转为一份针对本项目现状的、可执行的**功能完善路线图**。我们将首先分析各端现有菜单的疏漏和不合理之处，然后以一个具体的缺失功能为例，给出完整的全栈开发步骤。

### 3.1. 各端功能现状与完善建议

#### 3.1.1. 平台端 (Platform)

-   **现有菜单**: 平台总览、站点管理、权限管理、控制台、账号管理。
-   **现状分析**: 平台端的核心功能是“管钱”和“管人”，即对下级站点（学校）的管理和平台自身运维。当前菜单结构基本合理。
-   **缺失功能与建议**:
    1.  **[P0] 数据字典管理**: 后端已存在 `sys_dictionaries` 表，但前端缺失一个UI界面来管理这些全局配置（如商品单位“斤”、“箱”）。这是需要最优先补齐的基础功能。
    2.  **[P1] 运营广告管理**: 同上，后端有 `sys_banners` 表，但平台没有界面来管理首页的轮播图。

#### 3.1.2. 学校端 (School)

-   **现有菜单**: 供应链管理、食堂管理、订单管理、商品管理、溯源管理、结算管理、账号管理等。
-   **不合理的设计**:
    1.  **[P0] 菜单重复**: 学校后台同时存在“供应链管理 -> 供应商管理”和“供应商管理(新)”两个菜单，分别指向 `scm/SupplierManagement.vue` 和 `school/SupplierManagement.vue`。**这是冗余且混乱的**，应立即合并为一个，并统一组件。
    2.  **菜单命名**: “供应链管理 (scm)” 命名过于宽泛，建议直接将“供应商管理”提升为一级菜单，使结构更扁平、更清晰。
-   **缺失功能与建议**:
    1.  **[P1] 食堂与商户管理**: 在“食堂管理”菜单下，目前只有列表，缺少创建、编辑食堂和商户的功能。

#### 3.1.3. 供应商端 (Supplier)

-   **现有菜单**: 商品管理（上传、报价、修改）、配送管理（分拣、配送）、订单管理、账号管理、结算管理。
-   **现状分析**: 供应商端的功能较为完善，覆盖了核心的“管货”和“管钱”流程。
-   **缺失功能与建议**:
    1.  **[P1] 数据总览**: 缺少一个“工作台”或“数据总览”页面，用于展示关键指标，如“待处理订单数”、“待配送任务”、“本月销售额”等，以便供应商快速了解业务状况。

#### 3.1.4. 食堂端 (Canteen)

-   **现有菜单**: 商户管理、订单管理、账号管理。
-   **现状分析**: 目前的功能非常初级，仅限于查看。
-   **核心缺失功能**:
    1.  **[P0] 采购下单**: **这是食堂端最核心的功能，目前完全缺失。** 需要开发一个“商品市场”或“采购中心”页面，让食堂能够：
        -   浏览和搜索学校商品库中的商品。
        -   查看不同供应商对同一商品的报价。
        -   加购物车、下单、并管理自己的采购订单。

---

### 3.2. 开发指南示例：为平台端添加“数据字典管理”

下面，我们以上述分析中 **[P0] 优先级**的“数据字典管理”为例，展示一个完整功能的标准开发流程。

#### 步骤 1: 数据库与模型确认

-   **数据库**: `server/数据库设计.md` 中已定义 `sys_dictionaries` 表，包含 `dict_code`, `item_label`, `item_value` 等字段。
-   **GORM模型**: `server/internal/model/system.go` 中已定义 `SysDictionary` 模型。无需改动。

#### 步骤 2: 后端开发 (自底向上)

1.  **Repository** (`server/internal/repository/`):
    -   创建 `dictionary_repo.go` (接口) 和 `dictionary_repo_impl.go` (实现)。
    -   在接口中定义 `List(dictCode string)`, `Create(...)`, `Update(...)`, `Delete(...)` 等方法。
    -   在实现中，用 GORM 完成对 `sys_dictionaries` 表的增删改查。

2.  **Service** (`server/internal/service/`):
    -   创建 `dictionary_service.go`。
    -   实现对应的业务逻辑，例如创建前检查 `dict_code` 和 `item_value` 组合是否已存在。

3.  **Handler** (`server/internal/handler/`):
    -   创建 `dictionary_handler.go`。
    -   编写 `List`, `Create`, `Update`, `Delete` 等方法，负责解析 HTTP 请求，并调用 Service。

4.  **Router** (`server/internal/router/router.go`):
    -   在 `/api/v1` 路由组下，为平台管理员注册数据字典管理的路由。
        ```go
        // --- 依赖注入 ---
        dictRepo := repository.NewDictionaryRepository(database.DB)
        dictService := service.NewDictionaryService(dictRepo)
        dictHandler := handler.NewDictionaryHandler(dictService)

        // --- 路由注册 ---
        dictGroup := apiGroup.Group("/dictionaries")
        // 此路由只对平台管理员开放
        dictGroup.Use(middleware.AuthMiddleware(), middleware.PlatformAdminAuth())
        {
            dictGroup.GET("", dictHandler.List)
            dictGroup.POST("", dictHandler.Create)
            dictGroup.PUT("/:id", dictHandler.Update)
            dictGroup.DELETE("/:id", dictHandler.Delete)
        }
        ```

#### 步骤 3: 前端开发

1.  **API** (`web/src/api/`):
    -   创建 `dictionary.ts` 文件。
    -   封装对后端 `/api/v1/dictionaries` 端点的请求函数，如 `listDictionariesApi`, `createDictionaryApi` 等。

2.  **View** (`web/src/views/platform/`):
    -   创建 `DictionaryManagement.vue` 页面组件。
    -   使用 `el-table` 展示字典列表，提供“按 `dict_code` 筛选”功能。
    -   使用 `el-dialog` 和 `el-form` 实现新建和编辑字典项的弹窗。
    -   实现完整的增、删、改、查逻辑。

3.  **Router & Menu** (`web/src/router/index.ts`):
    -   在 `/platform` 的子路由中，添加数据字典页面的路由。建议放在“控制台”菜单下。
        ```typescript
        // in /platform children, under 'console'
        {
          path: 'dictionary', // 对应 /platform/console/dictionary
          component: () => import('@/views/platform/DictionaryManagement.vue'),
          meta: { title: '数据字典', icon: 'Collection' },
        },
        ```
    -   由于菜单是动态生成的，此步骤完成后，平台管理员登录即可在“控制台”下看到“数据字典”菜单项。

通过遵循以上流程，即可高效、规范地完成一个新功能的开发，并确保其完美融入现有项目架构。

## 4. 核心模块深度解析与架构建议

### 4.1. 权限系统：从“角色”到“功能下发”

**现状**:
目前系统是基于**角色（Role）**的访问控制。我们在路由 `meta` 中定义 `roles` 数组，前端 `permission.ts` 根据用户角色生成菜单，后端中间件也根据角色放行。这对于固定权限的场景是足够的。

**问题**:
无法满足“平台将功能下发给学校，学校再分配给自己的角色”这一动态需求。

**架构升级建议**:
引入**功能（Feature）**作为权限分配的最小单元。

1.  **新建数据库表**:
    -   `sys_features`: 功能表。由平台定义，如 `{ id: 1, key: 'supplier_management', name: '供应商管理' }`。
    -   `sys_role_features`: **角色-功能关联表**。`{ role_id, feature_id }`，用于定义一个“角色模板”包含哪些功能。
    -   `sys_school_features`: **学校-功能关联表**。`{ school_id, feature_id }`，这张表是“功能下发”的核心。平台通过操作这张表，来决定给哪个学校开放哪些功能。

2.  **新流程**:
    -   **平台**: 在“站点管理”页面，可以为每个学校勾选需要开通的功能。操作的结果就是向 `sys_school_features` 表中增删数据。
    -   **学校**: 学校管理员在“角色管理”页面，只能看到自己学校被授权的功能列表。他创建或修改角色时，只能从这个列表中为角色分配权限。操作的结果是修改 `sys_role_features` 表。
    -   **登录与鉴权**: 用户登录后，后端需要根据其 `RoleID` 和 `SchoolID`，计算出他最终拥有的所有 `feature_key` 列表，并放入 JWT。前端和后端的权限检查都改为**检查用户是否拥有某个 `feature_key`**，而不是检查他属于哪个 `role_key`。

### 4.2. 溯源模块：保证数据的“不可篡改”

**现状**:
`ord_order_items` 和 `ord_item_traces` 表的设计体现了“快照”和“扁平化”思想，这是非常正确的方向。

-   **交易快照**: 下单时，将商品的关键信息（品名、价格、规格）复制一份存入 `ord_order_items`，确保了历史订单的不可篡改性。
-   **溯源快照**: `ord_item_traces` 将供应商、商户、证书等信息“拍扁”存入 JSON 字段，使得扫码查询时无需联表，性能极高。

**架构建议**:
在现有基础上，进一步强化**溯源时间轴 (Time Line)**。

1.  **事件驱动的溯源日志**:
    -   将“下单”、“分拣”、“出库”、“配送”、“签收”等每一个关键动作，都视为一个“溯源事件”。
    -   可以新建一个 `trace_events` 表：`{ id, trace_id, event_type, operator_name, location, timestamp, extra_data_json }`。
    -   例如，“分拣”时，`supplier_staff` 在 App 上点击“开始分拣”，就向这张表里插入一条记录。
    -   最终展示溯源信息时，按时间顺序拉取一个 `trace_id` 下的所有事件，就能形成一个完整、动态、可信的时间轴。

2.  **数据采集**:
    -   溯源的价值取决于数据的真实性和丰富度。建议在供应商和食堂端的相关页面，通过表单、文件上传（如现场照片、质检单）、甚至 GPS 定位等方式，尽可能多地采集与环节相关的数据，并存入 `trace_events` 的 `extra_data_json` 字段中。

## 5. 开发规范与约定

1.  **代码风格**: 遵循项目现有代码风格。Go 使用 `gofmt`，前端使用 Prettier。
2.  **常量管理**: **严禁**在代码中使用硬编码的“魔术数字”或字符串（如 `org_type = 1`, `role = "admin"`）。所有这类值都应在 `internal/model/` 下的常量文件中统一定义。
3.  **权限第一**: 开发任何新功能，**第一步**就是思考其权限模型：谁能访问？谁能操作？并将其体现在路由、中间件和UI的控制中。
4.  **遵循分层**: 严格遵守 `Handler` -> `Service` -> `Repository` 的调用顺序，禁止跨层调用。
5.  **提交规范**: Git Commit Message 遵循 Angular 规范 (`feat:`, `fix:`, `refactor:`, `docs:`)，便于追溯和生成更新日志。

---

我已根据您的要求，结合代码审查结果和业务逻辑，撰写了这份详细的指南。请您审阅。在我们就这些原则和方向达成一致后，便可以着手解决您之前提到的具体问题了。
