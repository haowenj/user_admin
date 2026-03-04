# 员工管理系统测试代码说明

## 概述

本测试代码基于 functional_test_plan.md 生成，包含了完整的单元测试、集成测试和API测试代码。

## 测试文件结构

### 1. 单元测试文件

#### `controllers/controller_test.go`
- **功能**: 测试控制器层的所有功能
- **测试内容**:
  - 用户注册测试（正常、重复、空用户名、短密码）
  - 用户登录测试（正常、错误密码、不存在用户）
  - 密码修改测试（正常、错误原密码、弱新密码）
  - 员工查询测试（获取列表、获取单个、不存在员工）
  - 员工创建测试（正常、无效数据）
  - 员工更新测试（正常、不存在员工、无效数据）
  - 员工删除测试（正常、不存在员工）
  - 安全性测试（SQL注入、XSS攻击）
  - 性能测试（并发注册、并发查询）

#### `models/models_test.go`
- **功能**: 测试数据模型的正确性
- **测试内容**:
  - 用户模型测试（创建、查询、唯一性约束）
  - 员工模型测试（创建、查询、更新、软删除）
  - 数据验证测试（姓名、年龄、性别、部门等验证）
  - 表名测试

### 2. 集成测试文件

#### `integration_test.go`
- **功能**: 测试整个系统的集成
- **测试内容**:
  - 用户注册集成测试
  - 用户登录集成测试
  - 员工CRUD集成测试
  - 安全性集成测试

### 3. 测试工具文件

#### `test_server.go`
- **功能**: 独立的测试服务器
- **用途**: 为API测试提供HTTP服务

#### `run_tests.sh`
- **功能**: 自动化测试脚本
- **功能**:
  - 运行所有单元测试
  - 运行集成测试
  - 运行性能测试
  - 生成测试覆盖率报告
  - 运行API测试
  - 生成测试报告

## 测试运行方法

### 1. 运行所有测试

```bash
# 赋予执行权限
chmod +x run_tests.sh

# 运行所有测试
./run_tests.sh

# 跳过性能测试
./run_tests.sh --short
```

### 2. 运行特定测试

```bash
# 运行单元测试
go test -v ./...

# 运行控制器测试
go test -v ./controllers

# 运行模型测试
go test -v ./models

# 运行特定测试函数
go test -v ./controllers -run TestRegisterUser
go test -v ./models -run TestEmployeeModel
```

### 3. 生成测试覆盖率报告

```bash
# 生成覆盖率文件
go test -v -coverprofile=coverage.out ./...

# 生成HTML覆盖率报告
go tool cover -html=coverage.out -o coverage.html

# 查看覆盖率详情
go tool cover -func=coverage.out
```

### 4. 运行API测试

```bash
# 启动测试服务器
go run test_server.go &

# 等待服务器启动
sleep 3

# 运行API测试
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"TestPass123!"}'

# 停止服务器
pkill -f "go run test_server.go"
```

## 测试数据

### 测试数据库配置

测试使用独立的数据库 `employee_test`，连接配置：
```
用户名: test_user
密码: test_password
主机: localhost
端口: 3306
数据库名: employee_test
```

### 测试数据准备

测试脚本会自动创建测试数据，包括：
- 测试用户：`testuser`
- 测试员工：`测试员工`

## 测试覆盖范围

### 1. 功能测试
- ✅ 用户管理（注册、登录、密码修改）
- ✅ 员工管理（增删改查）
- ✅ 数据验证
- ✅ 错误处理

### 2. 安全性测试
- ✅ SQL注入防护
- ✅ XSS攻击防护
- ✅ 密码安全验证

### 3. 性能测试
- ✅ 并发用户注册
- ✅ 并发员工查询
- ✅ 响应时间测试

### 4. 集成测试
- ✅ 完整业务流程
- ✅ 数据库操作
- ✅ API端点测试

## 测试报告

### 1. 测试结果报告

测试完成后会生成 `test_report.md`，包含：
- 测试概览
- 测试结果统计
- 通过率
- 问题总结

### 2. 覆盖率报告

覆盖率报告 `coverage.html` 包含：
- 总体覆盖率
- 各文件覆盖率
- 函数覆盖率

## 测试最佳实践

### 1. 测试原则
- 每个测试函数应该独立运行
- 测试数据应该在测试后清理
- 使用有意义的测试名称
- 测试应该快速执行

### 2. 测试维护
- 随着代码变更更新测试
- 保持测试代码的清洁
- 定期审查测试覆盖率
- 添加新的测试用例

### 3. 持续集成
- 将测试集成到CI/CD流程
- 在代码提交时自动运行测试
- 设置覆盖率阈值
- 定期生成测试报告

## 故障排除

### 1. 数据库连接问题

如果测试失败，检查：
- MySQL是否运行
- 测试数据库是否存在
- 数据库连接参数是否正确

创建测试数据库：
```sql
CREATE DATABASE employee_test;
```

### 2. 端口冲突

如果8080端口被占用，修改配置文件中的端口设置。

### 3. 依赖问题

确保所有依赖都已安装：
```bash
go mod tidy
```

## 扩展测试

### 1. 添加新的测试用例

在相应的测试文件中添加新的测试函数：
```go
func TestNewFeature(t *testing.T) {
    // 测试逻辑
}
```

### 2. 添加性能测试

在 `controller_test.go` 中添加性能测试：
```go
func TestPerformanceNewFeature(t *testing.T) {
    // 性能测试逻辑
}
```

### 3. 添加集成测试

在 `integration_test.go` 中添加新的集成测试：
```go
func TestNewIntegration(t *testing.T) {
    // 集成测试逻辑
}
```

## 总结

本测试代码提供了完整的测试覆盖，确保员工管理系统的质量和稳定性。通过运行这些测试，可以：

1. 验证系统功能的正确性
2. 确保系统的安全性
3. 测试系统的性能
4. 验证系统的集成能力
5. 生成详细的测试报告

建议定期运行这些测试，并将测试集成到开发流程中，以确保系统的质量和稳定性。