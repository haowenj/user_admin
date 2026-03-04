# 员工管理系统测试代码

## 文件说明

### 1. `controllers/controller_test.go`
控制器层单元测试文件，包含：
- 用户管理测试（注册、登录、密码修改）
- 员工管理测试（增删改查）
- 安全性测试（SQL注入、XSS）
- 性能测试（并发测试）

### 2. `models/models_test.go`
数据模型测试文件，包含：
- 用户模型测试
- 员工模型测试
- 数据验证测试

### 3. `integration_test.go`
集成测试文件，包含：
- 完整业务流程测试
- API端点测试
- 安全性集成测试

### 4. `test_server.go`
独立测试服务器，用于API测试

### 5. `run_tests.sh`
自动化测试脚本，包含：
- 单元测试运行
- 集成测试运行
- 性能测试运行
- 覆盖率报告生成
- API测试运行
- 测试报告生成

### 6. `README_TEST.md`
测试代码详细说明文档

## 使用方法

### 运行所有测试
```bash
chmod +x run_tests.sh
./run_tests.sh
```

### 运行特定测试
```bash
# 单元测试
go test -v ./...

# 控制器测试
go test -v ./controllers

# 模型测试
go test -v ./models

# 特定测试函数
go test -v ./controllers -run TestRegisterUser
```

### 生成覆盖率报告
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 测试覆盖范围

- ✅ 用户管理功能
- ✅ 员工管理功能
- ✅ 数据验证
- ✅ 安全性测试
- ✅ 性能测试
- ✅ 集成测试

## 测试数据库

测试使用 `employee_test` 数据库，连接参数：
- 用户名: test_user
- 密码: test_password
- 主机: localhost
- 端口: 3306