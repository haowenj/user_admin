# 员工管理系统功能测试方案

## 1. 测试概述

### 1.1 测试目标
- 验证员工管理系统的核心功能正确性
- 确保API端点的稳定性和安全性
- 测试数据验证和错误处理机制
- 验证用户认证和授权功能

### 1.2 测试范围
- **用户管理模块**: 注册、登录、密码修改
- **员工管理模块**: 员工信息的增删改查
- **数据验证**: 输入验证、数据格式验证
- **安全性测试**: 密码安全、SQL注入防护
- **错误处理**: 异常场景处理

### 1.3 测试环境
- **操作系统**: macOS/Linux
- **Go版本**: 1.23.9
- **数据库**: MySQL 8.0+
- **测试工具**: Postman/curl, Go testing包

## 2. 测试数据设计

### 2.1 用户测试数据

#### 2.1.1 正常用户数据
```json
{
  "username": "testuser1",
  "password": "TestPass123!",
  "new_password": "NewPass456!"
}
```

#### 2.1.2 边界用户数据
```json
{
  "username": "a", // 最小长度
  "password": "A1!", // 最小长度
  "username": "very_long_username_that_exceeds_normal_limits", // 超长用户名
  "password": "VeryLongPasswordWithSpecialChars123!@#" // 复杂密码
}
```

#### 2.1.3 异常用户数据
```json
{
  "username": "", // 空用户名
  "password": "", // 空密码
  "username": "user with space", // 包含空格
  "password": "simple", // 弱密码
  "username": "<script>alert('xss')</script>", // XSS尝试
  "password": "123456" // 常见弱密码
}
```

### 2.2 员工测试数据

#### 2.2.1 正常员工数据
```json
{
  "name": "张三",
  "age": 28,
  "gender": "男",
  "department": "技术部",
  "position": "软件工程师",
  "hire_date": "2023-01-15"
}
```

#### 2.2.2 边界员工数据
```json
{
  "name": "A", // 最小长度
  "age": 18, // 最小年龄
  "name": "这是一个非常长的员工姓名测试长度限制", // 超长姓名
  "age": 65, // 最大年龄
  "hire_date": "2020-01-01", // 最早入职日期
  "hire_date": "2030-12-31" // 最晚入职日期
}
```

#### 2.2.3 异常员工数据
```json
{
  "name": "", // 空姓名
  "age": 17, // 未成年人
  "age": 66, // 超出年龄限制
  "gender": "未知", // 无效性别
  "department": "", // 空部门
  "position": "", // 空职位
  "hire_date": "2019-12-31", // 早于合理日期
  "hire_date": "invalid-date", // 无效日期格式
  "name": "<script>alert('xss')</script>" // XSS尝试
}
```

## 3. 功能测试用例

### 3.1 用户管理测试

#### 3.1.1 用户注册测试

**用例ID**: UR-001
**测试标题**: 正常用户注册
**测试步骤**:
1. 发送POST请求到 `/api/user/register`
2. 使用正常用户数据
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 200
- 响应包含用户ID
- 数据库中创建用户记录

**测试数据**:
```json
{
  "username": "newuser",
  "password": "NewPass123!"
}
```

**自动化脚本**:
```bash
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"username":"newuser","password":"NewPass123!"}'
```

---

**用例ID**: UR-002
**测试标题**: 重复用户名注册
**测试步骤**:
1. 先注册一个用户
2. 使用相同用户名再次注册
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 400
- 响应包含错误信息"用户名已存在"

**测试数据**:
```json
{
  "username": "existinguser",
  "password": "Pass123!"
}
```

---

**用例ID**: UR-003
**测试标题**: 空用户名注册
**测试步骤**:
1. 发送POST请求到 `/api/user/register`
2. 使用空用户名
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 400
- 响应包含错误信息"用户名不能为空"

**测试数据**:
```json
{
  "username": "",
  "password": "Pass123!"
}
```

---

**用例ID**: UR-004
**测试标题**: 弱密码注册
**测试步骤**:
1. 发送POST请求到 `/api/user/register`
2. 使用弱密码（如"123456"）
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 400
- 响应包含错误信息"密码强度不足"

**测试数据**:
```json
{
  "username": "weakpassuser",
  "password": "123456"
}
```

#### 3.1.2 用户登录测试

**用例ID**: UL-001
**测试标题**: 正常用户登录
**测试步骤**:
1. 发送POST请求到 `/api/user/login`
2. 使用正确的用户名和密码
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 200
- 响应包含用户ID
- 密码正确验证

**测试数据**:
```json
{
  "username": "testuser",
  "password": "TestPass123!"
}
```

---

**用例ID**: UL-002
**测试标题**: 错误密码登录
**测试步骤**:
1. 发送POST请求到 `/api/user/login`
2. 使用正确的用户名但错误的密码
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 400
- 响应包含错误信息"用户名或密码错误"

**测试数据**:
```json
{
  "username": "testuser",
  "password": "WrongPass123!"
}
```

---

**用例ID**: UL-003
**测试标题**: 不存在用户登录
**测试步骤**:
1. 发送POST请求到 `/api/user/login`
2. 使用不存在的用户名
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 400
- 响应包含错误信息"用户名或密码错误"

**测试数据**:
```json
{
  "username": "nonexistentuser",
  "password": "AnyPass123!"
}
```

#### 3.1.3 密码修改测试

**用例ID**: UC-001
**测试标题**: 正常密码修改
**测试步骤**:
1. 用户登录获取用户ID
2. 发送POST请求到 `/api/user/change-password`
3. 使用正确的原密码和新密码
4. 验证响应状态码和内容

**预期结果**:
- 状态码: 200
- 密码成功修改
- 新密码可以用于登录

**测试数据**:
```json
{
  "old_password": "OldPass123!",
  "new_password": "NewPass456!"
}
```

---

**用例ID**: UC-002
**测试标题**: 错误原密码修改
**测试步骤**:
1. 发送POST请求到 `/api/user/change-password`
2. 使用错误的原密码
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 400
- 响应包含错误信息"原密码错误"

**测试数据**:
```json
{
  "old_password": "WrongOldPass!",
  "new_password": "NewPass456!"
}
```

---

**用例ID**: UC-003
**测试标题**: 弱新密码修改
**测试步骤**:
1. 发送POST请求到 `/api/user/change-password`
2. 使用弱密码作为新密码
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 400
- 响应包含错误信息"新密码强度不足"

**测试数据**:
```json
{
  "old_password": "CorrectOldPass!",
  "new_password": "weak"
}
```

### 3.2 员工管理测试

#### 3.2.1 员工查询测试

**用例ID**: EQ-001
**测试标题**: 获取员工列表
**测试步骤**:
1. 发送GET请求到 `/api/employee`
2. 验证响应状态码和内容

**预期结果**:
- 状态码: 200
- 响应包含员工列表
- 列表格式正确

---

**用例ID**: EQ-002
**测试标题**: 获取单个员工信息
**测试步骤**:
1. 发送GET请求到 `/api/employee/{id}`
2. 使用有效的员工ID
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 200
- 响应包含员工详细信息
- 数据格式正确

---

**用例ID**: EQ-003
**测试标题**: 获取不存在员工信息
**测试步骤**:
1. 发送GET请求到 `/api/employee/{id}`
2. 使用不存在的员工ID
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 404
- 响应包含错误信息"员工不存在"

#### 3.2.2 员工创建测试

**用例ID**: EC-001
**测试标题**: 正常员工创建
**测试步骤**:
1. 发送POST请求到 `/api/employee`
2. 使用正常员工数据
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 200
- 响应包含创建的员工信息
- 数据库中创建员工记录
- 入职日期自动设置为当前日期

**测试数据**:
```json
{
  "name": "李四",
  "age": 30,
  "gender": "女",
  "department": "市场部",
  "position": "市场经理"
}
```

---

**用例ID**: EC-002
**测试标题**: 创建员工时指定入职日期
**测试步骤**:
1. 发送POST请求到 `/api/employee`
2. 使用包含入职日期的员工数据
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 200
- 入职日期与指定日期一致

**测试数据**:
```json
{
  "name": "王五",
  "age": 25,
  "gender": "男",
  "department": "技术部",
  "position": "前端工程师",
  "hire_date": "2023-06-01"
}
```

---

**用例ID**: EC-003
**测试标题**: 创建员工时包含无效数据
**测试步骤**:
1. 发送POST请求到 `/api/employee`
2. 使用包含无效数据的员工数据
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 400
- 响应包含具体的错误信息

**测试数据**:
```json
{
  "name": "",
  "age": 17,
  "gender": "未知",
  "department": "",
  "position": ""
}
```

#### 3.2.3 员工更新测试

**用例ID**: EU-001
**测试标题**: 正常员工更新
**测试步骤**:
1. 发送PUT请求到 `/api/employee/{id}`
2. 使用有效的员工ID和更新数据
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 200
- 响应包含更新后的员工信息
- 数据库中员工信息已更新

**测试数据**:
```json
{
  "name": "张三（更新）",
  "age": 29,
  "department": "研发部",
  "position": "高级软件工程师"
}
```

---

**用例ID**: EU-002
**测试标题**: 更新不存在员工
**测试步骤**:
1. 发送PUT请求到 `/api/employee/{id}`
2. 使用不存在的员工ID
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 404
- 响应包含错误信息"员工不存在"

---

**用例ID**: EU-003
**测试标题**: 更新员工时包含无效数据
**测试步骤**:
1. 发送PUT请求到 `/api/employee/{id}`
2. 使用包含无效数据的更新数据
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 400
- 响应包含具体的错误信息

**测试数据**:
```json
{
  "name": "",
  "age": 67,
  "gender": "未知"
}
```

#### 3.2.4 员工删除测试

**用例ID**: ED-001
**测试标题**: 正常员工删除
**测试步骤**:
1. 发送DELETE请求到 `/api/employee/{id}`
2. 使用有效的员工ID
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 200
- 响应包含删除成功信息
- 数据库中员工记录被软删除

---

**用例ID**: ED-002
**测试标题**: 删除不存在员工
**测试步骤**:
1. 发送DELETE请求到 `/api/employee/{id}`
2. 使用不存在的员工ID
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 404
- 响应包含错误信息"员工不存在"

## 4. 安全性测试

### 4.1 SQL注入测试

**用例ID**: SI-001
**测试标题**: 用户名SQL注入
**测试步骤**:
1. 发送POST请求到 `/api/user/register`
2. 使用包含SQL注入的用户名
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 400
- 系统应正确处理注入，不执行恶意SQL
- 数据库不应被破坏

**测试数据**:
```json
{
  "username": "admin';--",
  "password": "Pass123!"
}
```

---

**用例ID**: SI-002
**测试标题**: 员工信息SQL注入
**测试步骤**:
1. 发送POST请求到 `/api/employee`
2. 使用包含SQL注入的员工数据
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 400
- 系统应正确处理注入，不执行恶意SQL

**测试数据**:
```json
{
  "name": "Robert'); DROP TABLE employees;--",
  "age": 30,
  "gender": "男",
  "department": "技术部",
  "position": "工程师"
}
```

### 4.2 XSS攻击测试

**用例ID**: XSS-001
**测试标题**: 用户名XSS攻击
**测试步骤**:
1. 发送POST请求到 `/api/user/register`
2. 使用包含XSS脚本的用户名
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 400
- 系统应正确处理XSS，不执行恶意脚本
- 响应应进行HTML转义

**测试数据**:
```json
{
  "username": "<script>alert('XSS')</script>",
  "password": "Pass123!"
}
```

---

**用例ID**: XSS-002
**测试标题**: 员工信息XSS攻击
**测试步骤**:
1. 发送POST请求到 `/api/employee`
2. 使用包含XSS脚本的员工数据
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 400
- 系统应正确处理XSS，不执行恶意脚本

**测试数据**:
{
  "name": "<img src='x' onerror='alert(\"XSS\")'>",
  "age": 30,
  "gender": "男",
  "department": "技术部",
  "position": "工程师"
}

### 4.3 密码安全测试

**用例ID**: PS-001
**测试标题**: 弱密码检测
**测试步骤**:
1. 发送POST请求到 `/api/user/register`
2. 使用常见弱密码
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 400
- 系统应拒绝弱密码

**测试数据**:
```json
{
  "username": "weakuser",
  "password": "password"
}
```

---

**用例ID**: PS-002
**测试标题**: 密码复杂度要求
**测试步骤**:
1. 发送POST请求到 `/api/user/register`
2. 使用不符合复杂度要求的密码
3. 验证响应状态码和内容

**预期结果**:
- 状态码: 400
- 系统应拒绝不符合复杂度要求的密码

**测试数据**:
```json
{
  "username": "complexuser",
  "password": "simple"
}
```

## 5. 性能测试

### 5.1 并发测试

**用例ID**: PT-001
**测试标题**: 用户注册并发测试
**测试步骤**:
1. 使用并发工具（如JMeter）模拟多个用户同时注册
2. 发送多个并发注册请求
3. 监控系统响应时间和资源使用

**预期结果**:
- 系统应正确处理并发请求
- 响应时间在可接受范围内
- 数据库操作正确，无数据竞争

**并发参数**:
- 并发用户数: 10
- 请求总数: 100
- 超时时间: 30秒

---

**用例ID**: PT-002
**测试标题**: 员工查询并发测试
**测试步骤**:
1. 使用并发工具模拟多个用户同时查询员工信息
2. 发送多个并发查询请求
3. 监控系统响应时间和资源使用

**预期结果**:
- 系统应正确处理并发查询
- 响应时间在可接受范围内
- 数据库查询正确，无数据竞争

**并发参数**:
- 并发用户数: 20
- 请求总数: 200
- 超时时间: 30秒

### 5.2 负载测试

**用例ID**: LT-001
**测试标题**: 系统负载测试
**测试步骤**:
1. 逐步增加系统负载
2. 监控系统性能指标
3. 观察系统在不同负载下的表现

**预期结果**:
- 系统应能承受预期负载
- 响应时间随负载增加而合理增长
- 系统不应崩溃或出现严重错误

**负载参数**:
- 初始负载: 10并发用户
- 最大负载: 100并发用户
- 负载增长步长: 10用户
- 每个负载级别持续时间: 5分钟

## 6. 测试执行计划

### 6.1 测试环境准备

#### 6.1.1 环境配置
```yaml
测试环境:
  数据库: MySQL 8.0
  数据库名: employee_test
  端口: 3306
  用户名: test_user
  密码: test_password
  
应用配置:
  服务器端口: 8080
  数据库连接: test_user:test_password@tcp(localhost:3306)/employee_test?charset=utf8mb4&parseTime=True&loc=Local
  JWT密钥: test_jwt_secret_key
  运行模式: debug
```

#### 6.1.2 测试数据准备
```sql
-- 创建测试数据库
CREATE DATABASE employee_test;
USE employee_test;

-- 创建用户表
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 创建员工表
CREATE TABLE employees (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    age INT,
    gender VARCHAR(10),
    department VARCHAR(100),
    position VARCHAR(100),
    hire_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

### 6.2 测试执行步骤

#### 6.2.1 单元测试执行
```bash
# 进入项目目录
cd /Users/wenjuhao/code/cms/backend

# 运行所有测试
go test -v ./...

# 运行特定测试
go test -v -run TestUserRegister
go test -v -run TestEmployeeCreate

# 生成测试覆盖率报告
go test -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

#### 6.2.2 集成测试执行
```bash
# 启动应用服务器
go run main.go &

# 等待服务器启动
sleep 3

# 执行API测试
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"TestPass123!"}'

# 停止服务器
pkill -f "go run main.go"
```

#### 6.2.3 自动化测试脚本
```bash
#!/bin/bash

# 自动化测试脚本
echo "开始执行自动化测试..."

# 测试用户注册
echo "测试用户注册..."
response=$(curl -s -w "%{http_code}" -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"username":"autotestuser","password":"AutoPass123!"}')

if [ "$response" -eq 200 ]; then
    echo "✓ 用户注册测试通过"
else
    echo "✗ 用户注册测试失败"
fi

# 测试员工创建
echo "测试员工创建..."
response=$(curl -s -w "%{http_code}" -X POST http://localhost:8080/api/employee \
  -H "Content-Type: application/json" \
  -d '{"name":"测试员工","age":25,"gender":"男","department":"技术部","position":"工程师"}')

if [ "$response" -eq 200 ]; then
    echo "✓ 员工创建测试通过"
else
    echo "✗ 员工创建测试失败"
fi

echo "自动化测试完成"
```

### 6.3 测试结果记录

#### 6.3.1 测试结果模板
```markdown
# 测试执行报告

## 测试概览
- 测试日期: 2024-01-15
- 测试环境: macOS + MySQL 8.0
- 测试人员: 测试工程师
- 测试版本: v1.0.0

## 测试结果统计
- 总用例数: 25
- 通过用例数: 23
- 失败用例数: 2
- 通过率: 92%

## 失败用例详情
### 用例ID: UR-004
- 测试标题: 弱密码注册
- 实际结果: 系统接受了弱密码"123456"
- 预期结果: 系统应拒绝弱密码
- 严重程度: 高
- 建议修复: 加强密码复杂度验证

### 用例ID: XSS-002
- 测试标题: 员工信息XSS攻击
- 实际结果: 系统接受了包含XSS脚本的数据
- 预期结果: 系统应拒绝XSS攻击
- 严重程度: 高
- 建议修复: 加强输入验证和输出转义
```

## 7. 测试工具推荐

### 7.1 单元测试工具
- **Go testing**: 内置测试框架
- ** testify**: 提供丰富的断言和测试工具
- **gomock**: 用于模拟依赖项

### 7.2 API测试工具
- **Postman**: 功能强大的API测试工具
- **curl**: 命令行HTTP客户端
- **httpie**: 更友好的命令行HTTP客户端

### 7.3 性能测试工具
- **JMeter**: 开源性能测试工具
- **wrk**: 高性能HTTP基准测试工具
- **ab**: Apache HTTP服务器基准测试工具

### 7.4 持续集成工具
- **GitHub Actions**: 自动化测试和部署
- **GitLab CI**: 持续集成平台
- **Jenkins**: 开源CI/CD服务器

## 8. 测试最佳实践

### 8.1 测试设计原则
1. **测试覆盖**: 确保所有核心功能都有测试覆盖
2. **边界测试**: 测试正常、边界和异常情况
3. **数据驱动**: 使用多种测试数据验证系统行为
4. **自动化**: 尽可能实现测试自动化
5. **独立性**: 测试用例之间相互独立

### 8.2 测试执行建议
1. **定期执行**: 每次代码提交后运行测试
2. **持续集成**: 将测试集成到CI/CD流程
3. **测试报告**: 生成详细的测试报告
4. **问题跟踪**: 跟踪和管理测试发现的bug
5. **定期回顾**: 定期回顾和改进测试策略

### 8.3 测试维护
1. **测试更新**: 随着代码变更更新测试用例
2. **测试重构**: 重构和优化测试代码
3. **测试文档**: 维护测试文档和测试计划
4. **测试培训**: 对团队成员进行测试培训
5. **测试文化建设**: 培养良好的测试文化

## 9. 总结

本测试方案为员工管理系统提供了全面的功能测试覆盖，包括：

- **用户管理功能**: 注册、登录、密码修改
- **员工管理功能**: 员工信息的增删改查
- **数据验证**: 输入验证、数据格式验证
- **安全性测试**: SQL注入、XSS攻击、密码安全
- **性能测试**: 并发测试、负载测试

通过执行这些测试用例，可以确保系统的功能正确性、安全性和稳定性。建议：

1. **实现自动化测试**: 将测试脚本集成到CI/CD流程
2. **定期执行测试**: 每次代码变更后运行测试
3. **持续改进**: 根据测试结果不断改进系统
4. **扩展测试范围**: 根据业务发展扩展测试覆盖

这个测试方案可以作为系统测试的基础，随着系统的演进不断扩展和完善。