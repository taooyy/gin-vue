# Gin-Vue 项目代码审查报告

本文档旨在对当前的 Go (Gin) + Vue.js 全栈项目进行一次全面的代码审查，总结其架构、逻辑和潜在问题，为后续的开发和维护提供参考。

## 1. 整体架构

项目遵循了清晰的前后端分离设计，结构合理，易于理解。

### 1.1. 后端 (Go + Gin)

- **目录**: `server/`
- **架构**: 经典的三层架构 (Web/Handler -> Service -> Repository)。
  - `cmd/server/main.go`: 项目入口，负责初始化配置、数据库和路由。
  - `internal/router/router.go`: **依赖注入核心**。在这里，项目通过手动注入的方式，将 `Repository` (数据库操作) 注入到 `Service` (业务逻辑)，再将 `Service` 注入到 `Handler` (HTTP 请求处理)，实现了各层之间的解耦。
  - `internal/handler`: 处理 HTTP 请求的解析、验证和响应。
  - `internal/service`: 包含所有核心业务逻辑，不直接操作数据库。
  - `internal/repository`: 仅负责数据库的 CRUD 操作，是数据访问的唯一出口。

### 1.2. 前端 (Vue.js + Vite)

- **目录**: `web/`
- **架构**: 基于 Vue 3、Vite 和 Pinia 的现代化组件化架构。
  - `main.ts`: Vue 应用入口，负责初始化 Pinia、Vue Router 和 Element Plus。
  - `router/index.ts`: 定义所有前端路由和全局导航守卫 (`beforeEach`)。
  - `stores/`: Pinia 状态管理中心。`user.ts` 负责用户认证信息，`permission.ts` 负责菜单和权限。
  - `views/`: 页面级组件，代表一个完整的功能页面。
  - `components/`: 可复用的原子组件（如分页 `Pagination`）。
  - `api/`: 统一管理所有对后端 API 的请求，实现视图与网络请求的分离。
  - `composables/`: 可组合函数，用于抽离和复用逻辑（如 `useAuth.ts`）。

## 2. 认证与授权 (Auth Flow)

系统的认证授权流程设计稳健，职责分明。

1.  **用户登录**:
    - 前端在 `Login.vue` 页面通过 `useAuth` 调用 `loginApi`，将用户名和密码发送到后端 `POST /api/v1/system/login`。
    - 后端 `auth_service.go` 验证成功后，生成一个 JWT。
    - **JWT Payload**: `CustomClaims` (`pkg/jwt/jwt.go`) 是关键，它包含了 `UserID`, `Username`, `Role`, `OrgID`, `OrgType` 等核心信息，为后续所有权限判断提供依据。
    - 后端返回 JWT 和用户信息（包含 `OrgType`）。

2.  **状态持久化**:
    - 前端 `user.ts` (Pinia Store) 将 JWT 和用户信息存入 `sessionStorage`，以维持登录状态。

3.  **请求认证**:
    - 前端 `apiClient` (axios 实例) 会在每个请求的 Header 中自动附带 `Authorization: Bearer <jwt>`。
    - 后端 `middleware/auth.go` 中间件会拦截所有受保护的 API 请求，验证 JWT 的有效性。
    - 验证通过后，将 JWT 中的 `CustomClaims` 解码并存入 Gin 的请求上下文 `context` 中。

4.  **权限控制**:
    - **后端**: `middleware/authz.go` 等授权中间件从 `context` 中获取用户 Claims，根据 `Role` 或 `OrgType` 判断用户是否有权访问特定 API。
    - **前端**: `router/index.ts` 中的 `beforeEach` 导航守卫主要负责检查是否存在 token。更细粒度的权限控制体现在 `permission.ts` store 中，它根据用户的 `role` 动态生成可访问的侧边栏菜单，从而在 UI 上限制用户的操作范围。

## 3. 核心数据流 (以“账号管理”为例)

1.  **前端加载**: 用户访问“平台账号管理”页面 (`views/platform/AccountManagement.vue`)。
2.  **API 请求**: `onMounted`生命周期钩子触发 `getAccountList` 方法，该方法调用 `api/account.ts` 中的 `listAccountsApi`。
3.  **网络传输**: `listAccountsApi` 发起一个 `GET /api/v1/accounts` 请求。由于有 `Auth` 中间件保护，请求头中会携带 JWT。
4.  **后端处理**:
    - `router.go` 将请求路由到 `account_handler.go` 的 `ListAccounts` 方法。
    - `Handler` 从 `context` 中解析出当前登录用户的 `CustomClaims`。
    - `Handler` 调用 `account_service.go` 的 `ListAccounts` 方法，并传入 `creatorClaims`。
    - **核心逻辑**: `Service` 层根据 `creatorClaims.OrgType` 判断：
        - 如果是平台用户 (`OrgTypePlatform`)，则调用 `userRepo.ListByCreator` 并附带 `OrgType` 过滤器。
        - 如果是其他用户，调用 `userRepo.ListByCreator`，只按 `creatorID` 过滤。
    - `Repository` 层 (`user_repo_impl.go`) 执行最终的 SQL 查询并返回数据。
5.  **前端渲染**: 数据沿链路返回，最终在 `AccountManagement.vue` 的表格中渲染出来。

## 4. 潜在问题与改进建议

通过对代码的审查，发现以下几点值得关注：

1.  **平台与站点的逻辑划分**:
    - **问题**: 最近的几个问题都围绕“平台”和“站点”的账号管理逻辑。根本原因在于，最初的列表查询逻辑 (`ListByCreator`) 只考虑了“谁创建”，而没有考虑“管理的对象是谁”。
    - **当前状态**: 在上一个问题中，这个逻辑已被修改。现在 `ListAccounts` 服务会区分平台管理员和其他管理员，平台管理员的列表会额外根据 `OrgType` 进行过滤。
    - **建议**: **这是一个正确的方向**。未来的开发中，应坚持“**平台查平台类型，租户查租户范围**”的原则。对于所有列表查询功能，都应检查当前用户的 `OrgType`，并应用对应的过滤规则，而不是简单地依赖 `CreatedBy`。

2.  **硬编码的角色/类型**:
    - **问题**: 在前端代码和后端代码中，曾出现过用数字 `1`, `2` 或字符串 `"platform_admin"` 来做逻辑判断的情况。
    - **当前状态**: <strong style="color: green;">已解决</strong>。在之前的重构中，大部分硬编码的数字和字符串已被 `OrgType` 枚举和角色常量所取代，显著提升了代码的可读性和可维护性。
    - **建议**: 保持这个好习惯。未来应将所有硬编码的标识符统一管理在 `constants` 或 `enum` 文件中。

3.  **前端路由守卫逻辑**:
    - **问题**: `router/index.ts` 中的 `beforeEach` 守卫在处理菜单生成时，曾依赖 `to.path` 来判断 `layoutType` (`platform` or `tenant`)，这在某些情况下可能不稳定。
    - **当前状态**: <strong style="color: green;">已解决</strong>。此问题已在修复“登录不跳转”的过程中被修正，现在的导航守卫依赖于 `userStore` 中更可靠的 `orgType` 状态来判断布局类型。
    - **建议**: 未来在编写导航守卫时，应优先从 `userStore` 中获取权威的用户状态信息，而不是从不确定的 `to` 或 `from` 路由对象中推断。

## 总结

项目整体架构清晰，遵循了现代全栈开发的最佳实践。近期遇到的问题主要集中在**业务逻辑层面对于多租户（平台 vs. 站点）隔离的实现细节**上，而非架构本身。

通过将 `OrgType` 集成到 JWT、服务层逻辑和数据库查询中，系统的健壮性已得到显著提升。建议团队在未来的功能开发中，继续强化和贯彻这一核心原则。

---

## 5. 最终设计原则：基于范围的数据权限

经过讨论，我们正式确立以下设计原则，作为未来所有功能开发，尤其是列表查询功能的权限设计基石。

> **核心原则：一个用户的“列表查看”权限，应由其“角色的管辖范围”（Scope）决定，而不是简单地由“所有权”（Ownership，即是否由其创建）决定。**

此原则意味着，一个用户能看到哪些数据，取决于他扮演的角色能管辖多大的范围，而不是他个人创建了哪些数据。

### 各角色范围定义与实现

1.  **平台管理员 (`platform_admin`)**
    -   **管辖范围**: 整个平台，即所有 `OrgType` 为 `Platform` 的组织及下属用户。
    -   **实现**: 当平台管理员请求一个列表时（如账号列表），查询逻辑**必须**筛选出所有隶属于“平台型”组织的数据。
    -   **SQL示例**: `... WHERE org_id IN (SELECT id FROM sys_organizations WHERE org_type = 0)`

2.  **学校管理员 (`school_admin`)**
    -   **管辖范围**: 自己所在学校内部，即与自己 `OrgID` 相同的组织及下属用户。
    -   **实现**: 当学校管理员请求一个列表时（如本校员工列表、本校供应商列表），查询逻辑**必须**用当前用户的 `OrgID` 作为过滤条件，筛选出目标数据。
    -   **SQL示例**: `... WHERE school_id = [currentUser.OrgID]` 或 `... WHERE org_id = [currentUser.OrgID]`

3.  **供应商管理员 (`supplier_admin`)**
    -   **管辖范围**: 自己所在供应商内部。
    -   **实现**: 逻辑同学校管理员，查询**必须**基于当前用户的 `OrgID` 进行数据隔离。

### “所有权”的使用场景

“所有权” (`CreatedBy`) 依然有用，但它的使用场景应该被严格限定。例如，在“编辑”或“删除”某个具体条目时，除了检查范围权限，还可以额外检查所有权，以实现更精细的控制（例如：“只有我自己创建的员工，我才能删除”）。但在**列表查询**层面，应优先使用范围原则。

通过严格遵守此原则，可以从根本上保证系统清晰的权限模型和租户间的数据安全。
