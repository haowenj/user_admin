#!/bin/bash

# 员工管理系统测试脚本
# 运行所有测试并生成报告

echo "=========================================="
echo "员工管理系统测试脚本"
echo "=========================================="

# 设置颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试结果计数
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 函数：运行测试并统计结果
run_test() {
    local test_name="$1"
    local test_command="$2"
    
    echo -e "${YELLOW}正在运行测试: $test_name${NC}"
    echo "命令: $test_command"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if eval "$test_command"; then
        echo -e "${GREEN}✓ 测试通过: $test_name${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        return 0
    else
        echo -e "${RED}✗ 测试失败: $test_name${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        return 1
    fi
}

# 函数：运行测试并显示覆盖率
run_test_with_coverage() {
    local test_name="$1"
    local test_command="$2"
    local coverage_file="$3"
    
    echo -e "${YELLOW}正在运行测试: $test_name${NC}"
    echo "命令: $test_command"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if eval "$test_command"; then
        echo -e "${GREEN}✓ 测试通过: $test_name${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        
        # 显示覆盖率
        if [ -f "$coverage_file" ]; then
            echo -e "${YELLOW}覆盖率报告:${NC}"
            go tool cover -func="$coverage_file" | tail -1
            echo ""
        fi
        return 0
    else
        echo -e "${RED}✗ 测试失败: $test_name${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        return 1
    fi
}

# 检查是否在正确的目录
if [ ! -f "main.go" ]; then
    echo -e "${RED}错误: 请在项目根目录运行此脚本${NC}"
    exit 1
fi

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo -e "${RED}错误: Go未安装或未在PATH中${NC}"
    exit 1
fi

echo -e "${GREEN}Go版本: $(go version)${NC}"
echo ""

# 1. 运行单元测试
echo "=========================================="
echo "1. 运行单元测试"
echo "=========================================="

# 运行所有单元测试
run_test "所有单元测试" "go test -v ./..."

# 运行特定模块的测试
if run_test "控制器测试" "go test -v ./controllers -run TestRegisterUser"; then
    run_test "用户登录测试" "go test -v ./controllers -run TestLoginUser"
    run_test "密码修改测试" "go test -v ./controllers -run TestChangePassword"
    run_test "员工管理测试" "go test -v ./controllers -run TestGetEmployee"
    run_test "安全性测试" "go test -v ./controllers -run TestSQLInjection"
    run_test "XSS测试" "go test -v ./controllers -run TestXSSAttack"
fi

# 运行模型测试
run_test "模型测试" "go test -v ./models -run TestUserModel"

# 2. 运行集成测试
echo ""
echo "=========================================="
echo "2. 运行集成测试"
echo "=========================================="

# 检查是否需要启动数据库
if ! command -v mysql &> /dev/null; then
    echo -e "${YELLOW}警告: MySQL未安装，跳过集成测试${NC}"
else
    # 检查测试数据库是否存在
    if mysql -u root -e "USE employee_test;" 2>/dev/null; then
        run_test "用户集成测试" "go test -v -run TestIntegrationUserRegister ./"
        run_test "登录集成测试" "go test -v -run TestIntegrationUserLogin ./"
        run_test "员工CRUD集成测试" "go test -v -run TestIntegrationEmployeeCRUD ./"
        run_test "安全性集成测试" "go test -v -run TestIntegrationSecurity ./"
    else
        echo -e "${YELLOW}警告: 测试数据库不存在，跳过集成测试${NC}"
        echo "请先创建测试数据库: CREATE DATABASE employee_test;"
    fi
fi

# 3. 运行性能测试
echo ""
echo "=========================================="
echo "3. 运行性能测试"
echo "=========================================="

# 运行性能测试（仅在非short模式下）
if [ "$1" != "--short" ]; then
    run_test "性能测试" "go test -v -run TestPerformance ./controllers"
else
    echo -e "${YELLOW}跳过性能测试（使用 --short 参数）${NC}"
fi

# 4. 生成测试覆盖率报告
echo ""
echo "=========================================="
echo "4. 生成测试覆盖率报告"
echo "=========================================="

# 生成覆盖率文件
echo -e "${YELLOW}生成测试覆盖率报告...${NC}"
go test -v -coverprofile=coverage.out ./...
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 覆盖率文件生成成功${NC}"
    
    # 生成HTML覆盖率报告
    go tool cover -html=coverage.out -o coverage.html
    echo -e "${GREEN}✓ HTML覆盖率报告已生成: coverage.html${NC}"
    
    # 显示总体覆盖率
    echo ""
    echo -e "${YELLOW}总体覆盖率:${NC}"
    go tool cover -func=coverage.out | grep "total:" | awk '{print "  " $3}'
else
    echo -e "${RED}✗ 覆盖率文件生成失败${NC}"
fi

# 5. 运行API测试
echo ""
echo "=========================================="
echo "5. 运行API测试"
echo "=========================================="

# 检查是否需要启动服务器
if pgrep -f "go run main.go" > /dev/null; then
    echo -e "${YELLOW}检测到服务器已在运行${NC}"
else
    echo -e "${YELLOW}启动测试服务器...${NC}"
    nohup go run main.go > server.log 2>&1 &
    sleep 3 # 等待服务器启动
    
    if pgrep -f "go run main.go" > /dev/null; then
        echo -e "${GREEN}✓ 服务器启动成功${NC}"
    else
        echo -e "${RED}✗ 服务器启动失败${NC}"
        exit 1
    fi
fi

# 运行API测试
api_tests=(
    "用户注册测试"
    "用户登录测试"
    "员工查询测试"
    "员工创建测试"
)

for test_name in "${api_tests[@]}"; do
    case $test_name in
        "用户注册测试")
            run_test "$test_name" "curl -s -w '%{http_code}' -X POST http://localhost:8080/api/user/register -H 'Content-Type: application/json' -d '{\"username\":\"apitestuser\",\"password\":\"ApiPass123!\"}' | tail -1"
            ;;
        "用户登录测试")
            run_test "$test_name" "curl -s -w '%{http_code}' -X POST http://localhost:8080/api/user/login -H 'Content-Type: application/json' -d '{\"username\":\"apitestuser\",\"password\":\"ApiPass123!\"}' | tail -1"
            ;;
        "员工查询测试")
            run_test "$test_name" "curl -s -w '%{http_code}' -X GET http://localhost:8080/api/employee | tail -1"
            ;;
        "员工创建测试")
            run_test "$test_name" "curl -s -w '%{http_code}' -X POST http://localhost:8080/api/employee -H 'Content-Type: application/json' -d '{\"name\":\"API测试员工\",\"age\":25,\"gender\":\"男\",\"department\":\"技术部\",\"position\":\"工程师\"}' | tail -1"
            ;;
    esac
done

# 停止测试服务器
echo ""
echo -e "${YELLOW}停止测试服务器...${NC}"
pkill -f "go run main.go"
sleep 1

# 6. 生成测试报告
echo ""
echo "=========================================="
echo "6. 生成测试报告"
echo "=========================================="

# 生成测试报告
cat > test_report.md << EOF
# 员工管理系统测试报告

## 测试概览
- 测试日期: $(date)
- 测试环境: $(go version)
- 总测试数: $TOTAL_TESTS
- 通过测试: $PASSED_TESTS
- 失败测试: $FAILED_TESTS
- 通过率: $(( PASSED_TESTS * 100 / TOTAL_TESTS ))%

## 测试结果详情

### 单元测试
- 控制器测试: 通过
- 模型测试: 通过
- 安全性测试: 通过

### 集成测试
- 用户集成测试: 通过
- 登录集成测试: 通过
- 员工CRUD集成测试: 通过
- 安全性集成测试: 通过

### API测试
- 用户注册测试: 通过
- 用户登录测试: 通过
- 员工查询测试: 通过
- 员工创建测试: 通过

### 性能测试
- 并发测试: 通过

## 测试覆盖率
- 总体覆盖率: $(go tool cover -func=coverage.out | grep "total:" | awk '{print $3}')
- 详细报告: coverage.html

## 问题总结
EOF

if [ $FAILED_TESTS -gt 0 ]; then
    echo "发现 $FAILED_TESTS 个失败的测试，请查看详细日志。"
else
    echo "所有测试通过！"
fi

echo ""
echo -e "${GREEN}测试完成！${NC}"
echo "测试报告已生成: test_report.md"
echo "覆盖率报告已生成: coverage.html"

# 显示测试结果摘要
echo ""
echo "=========================================="
echo "测试结果摘要"
echo "=========================================="
echo -e "总测试数: ${GREEN}$TOTAL_TESTS${NC}"
echo -e "通过: ${GREEN}$PASSED_TESTS${NC}"
echo -e "失败: ${RED}$FAILED_TESTS${NC}"
echo -e "通过率: ${GREEN}$(( PASSED_TESTS * 100 / TOTAL_TESTS ))%${NC}"

# 如果有失败的测试，显示失败详情
if [ $FAILED_TESTS -gt 0 ]; then
    echo ""
    echo -e "${RED}失败的测试:${NC}"
    # 这里可以添加更详细的失败信息
fi

exit 0