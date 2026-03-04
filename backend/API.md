# 后端 API 接口文档

## 基础信息

- **框架**: Go + Gin + Gorm
- **数据库**: SQLite
- **Base URL**: `http://localhost:8080`

---

## 用户接口

### 1. 用户注册

**POST** `/api/user/register`

**请求体**:
```json
{
  "username": "string",
  "password": "string"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名（唯一） |
| password | string | 是 | 密码（至少6位） |

**响应**:
- 成功: `200 OK`
```json
{
  "message": "注册成功"
}
```
- 失败: `400 Bad Request`
```json
{
  "error": "错误信息"
}
```

---

### 2. 用户登录

**POST** `/api/user/login`

**请求体**:
```json
{
  "username": "string",
  "password": "string"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**响应**:
- 成功: `200 OK`
```json
{
  "message": "登录成功",
  "user_id": 1
}
```
- 失败: `401 Unauthorized`
```json
{
  "error": "用户名或密码错误"
}
```

---

### 3. 修改密码

**POST** `/api/user/change-password`

**请求头**: 无需认证

**请求体**:
```json
{
  "username": "string",
  "old_password": "string",
  "new_password": "string"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| old_password | string | 是 | 原密码 |
| new_password | string | 是 | 新密码 |

**响应**:
- 成功: `200 OK`
```json
{
  "message": "密码修改成功"
}
```
- 失败: `401 Unauthorized`
```json
{
  "error": "原密码错误"
}
```

---

## 员工接口

### 4. 获取所有员工

**GET** `/api/employee`

**响应**:
- 成功: `200 OK`
```json
[
  {
    "id": 1,
    "name": "张三",
    "age": 30,
    "gender": "男",
    "department": "技术部",
    "position": "工程师",
    "hire_date": "2024-01-15",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  }
]
```

---

### 5. 获取单个员工

**GET** `/api/employee/:id`

**路径参数**:
| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 员工ID |

**响应**:
- 成功: `200 OK`
```json
{
  "id": 1,
  "name": "张三",
  "age": 30,
  "gender": "男",
  "department": "技术部",
  "position": "工程师",
  "hire_date": "2024-01-15",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```
- 失败: `404 Not Found`
```json
{
  "error": "员工不存在"
}
```

---

### 6. 创建员工

**POST** `/api/employee`

**请求体**:
```json
{
  "name": "张三",
  "age": 30,
  "gender": "男",
  "department": "技术部",
  "position": "工程师",
  "hire_date": "2024-01-15"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 员工姓名 |
| age | int | 否 | 年龄 |
| gender | string | 否 | 性别 |
| department | string | 否 | 部门 |
| position | string | 否 | 职位 |
| hire_date | string | 否 | 入职日期（格式: YYYY-MM-DD） |

**响应**:
- 成功: `200 OK`
```json
{
  "id": 1,
  "name": "张三",
  "age": 30,
  "gender": "男",
  "department": "技术部",
  "position": "工程师",
  "hire_date": "2024-01-15",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```
- 失败: `400 Bad Request`
```json
{
  "error": "错误信息"
}
```

---

### 7. 更新员工

**PUT** `/api/employee/:id`

**路径参数**:
| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 员工ID |

**请求体**:
```json
{
  "name": "张三",
  "age": 31,
  "department": "技术部",
  "position": "高级工程师"
}
```

**响应**:
- 成功: `200 OK`
```json
{
  "id": 1,
  "name": "张三",
  "age": 31,
  "gender": "男",
  "department": "技术部",
  "position": "高级工程师",
  "hire_date": "2024-01-15",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-16T10:00:00Z"
}
```
- 失败: `404 Not Found`
```json
{
  "error": "没有查询到员工数据"
}
```

---

### 8. 删除员工

**DELETE** `/api/employee/:id`

**路径参数**:
| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 员工ID |

**响应**:
- 成功: `200 OK`
```json
{
  "message": "删除成功"
}
```
- 失败: `404 Not Found`
```json
{
  "error": "员工不存在"
}
```

---

## 错误码汇总

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 认证失败/未授权 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |
