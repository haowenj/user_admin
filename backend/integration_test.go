package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"employee-management/config"
	"employee-management/models"
	"employee-management/routers"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func setupTestDB() {
	// 初始化测试配置
	if err := config.InitConfig("config.yaml"); err != nil {
		panic("failed to load test config")
	}

	// 修改配置为测试数据库
	config.AppConfig.Database = config.DatabaseConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "",
		DBName:   "employee_test",
		Charset:  "utf8mb4",
	}

	// 连接测试数据库
	if err := models.InitDB(); err != nil {
		panic("failed to connect test database: " + err.Error())
	}

	// 自动迁移表
	models.Migrate()
}

func teardownTestDB() {
	if models.DB != nil {
		models.DB.Exec("DELETE FROM users")
		models.DB.Exec("DELETE FROM employees")
		models.DB.Exec("DROP TABLE users")
		models.DB.Exec("DROP TABLE employees")
	}
}

func TestIntegrationUserRegister(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	router := setupTestRouter()

	// 测试正常用户注册
	t.Run("正常用户注册", func(t *testing.T) {
		userData := map[string]interface{}{
			"username": "integrationuser",
			"password": "IntPass123!",
		}

		jsonData, _ := json.Marshal(userData)
		req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d，实际 %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if msg, ok := response["message"].(string); ok && msg != "注册成功" {
			t.Errorf("期望消息 '注册成功'，实际 '%s'", msg)
		}

		// 验证数据库中的用户
		var user models.User
		if err := models.DB.Where("username = ?", "integrationuser").First(&user).Error; err != nil {
			t.Errorf("数据库中未找到用户: %v", err)
		}

		// 验证密码已加密
		if user.Password == "IntPass123!" {
			t.Error("密码未加密")
		}
	})

	// 测试重复注册
	t.Run("重复用户注册", func(t *testing.T) {
		userData := map[string]interface{}{
			"username": "integrationuser", // 重复用户名
			"password": "IntPass123!",
		}

		jsonData, _ := json.Marshal(userData)
		req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("期望状态码 %d，实际 %d", http.StatusBadRequest, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if msg, ok := response["error"].(string); ok && !strings.Contains(msg, "用户名已存在") {
			t.Errorf("期望错误消息包含 '用户名已存在'，实际 '%s'", msg)
		}
	})
}

func TestIntegrationUserLogin(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	// 先创建一个测试用户
	hashedPassword, _ := generateHash("LoginPass123!")
	user := models.User{
		Username: "loginuser",
		Password: hashedPassword,
	}
	db.Create(&user)

	router := setupTestRouter()

	// 测试正常登录
	t.Run("正常用户登录", func(t *testing.T) {
		loginData := map[string]interface{}{
			"username": "loginuser",
			"password": "LoginPass123!",
		}

		jsonData, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d，实际 %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if msg, ok := response["message"].(string); ok && msg != "登录成功" {
			t.Errorf("期望消息 '登录成功'，实际 '%s'", msg)
		}

		if _, ok := response["user_id"]; !ok {
			t.Error("响应中缺少 user_id")
		}
	})

	// 测试错误密码
	t.Run("错误密码登录", func(t *testing.T) {
		loginData := map[string]interface{}{
			"username": "loginuser",
			"password": "WrongPass123!",
		}

		jsonData, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("期望状态码 %d，实际 %d", http.StatusUnauthorized, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if msg, ok := response["error"].(string); ok && !strings.Contains(msg, "用户名或密码错误") {
			t.Errorf("期望错误消息包含 '用户名或密码错误'，实际 '%s'", msg)
		}
	})
}

func TestIntegrationEmployeeCRUD(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	router := setupTestRouter()

	// 测试创建员工
	t.Run("创建员工", func(t *testing.T) {
		employeeData := map[string]interface{}{
			"name":       "集成测试员工",
			"age":        30,
			"gender":     "男",
			"department": "技术部",
			"position":   "高级工程师",
		}

		jsonData, _ := json.Marshal(employeeData)
		req, _ := http.NewRequest("POST", "/api/employee", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d，实际 %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if name, ok := response["name"].(string); ok && name != "集成测试员工" {
			t.Errorf("期望员工姓名 '集成测试员工'，实际 '%s'", name)
		}

		// 验证数据库中的员工
		var employee models.Employee
		if err := db.Where("name = ?", "集成测试员工").First(&employee).Error; err != nil {
			t.Errorf("数据库中未找到员工: %v", err)
		}

		// 保存员工ID用于后续测试
		var employeeID uint
		if id, ok := response["id"].(float64); ok {
			employeeID = uint(id)
		}

		// 测试获取员工列表
		req, _ = http.NewRequest("GET", "/api/employee", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d，实际 %d", http.StatusOK, w.Code)
		}

		var employees []models.Employee
		json.Unmarshal(w.Body.Bytes(), &employees)

		if len(employees) == 0 {
			t.Error("员工列表为空")
		}

		// 测试获取单个员工
		req, _ = http.NewRequest("GET", "/api/employee/"+fmt.Sprintf("%d", employeeID), nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d，实际 %d", http.StatusOK, w.Code)
		}

		var singleEmployee models.Employee
		json.Unmarshal(w.Body.Bytes(), &singleEmployee)

		if singleEmployee.ID != employeeID {
			t.Errorf("期望员工ID %d，实际 %d", employeeID, singleEmployee.ID)
		}

		// 测试更新员工
		updateData := map[string]interface{}{
			"name":       "集成测试员工（更新）",
			"age":        31,
			"department": "研发部",
		}

		jsonData, _ = json.Marshal(updateData)
		req, _ = http.NewRequest("PUT", "/api/employee/"+fmt.Sprintf("%d", employeeID), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d，实际 %d", http.StatusOK, w.Code)
		}

		var updatedEmployee models.Employee
		db.First(&updatedEmployee, employeeID)

		if updatedEmployee.Name != "集成测试员工（更新）" {
			t.Errorf("期望更新后姓名 '集成测试员工（更新）'，实际 '%s'", updatedEmployee.Name)
		}

		// 测试删除员工
		req, _ = http.NewRequest("DELETE", "/api/employee/"+fmt.Sprintf("%d", employeeID), nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d，实际 %d", http.StatusOK, w.Code)
		}

		// 验证员工已被软删除
		var count int64
		db.Model(&models.Employee{}).Where("id = ?", employeeID).Count(&count)
		if count > 0 {
			t.Error("员工未被正确删除")
		}
	})
}

func TestIntegrationSecurity(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	router := setupTestRouter()

	// 测试SQL注入防护
	t.Run("SQL注入防护", func(t *testing.T) {
		userData := map[string]interface{}{
			"username": "admin'; DROP TABLE users;--",
			"password": "Pass123!",
		}

		jsonData, _ := json.Marshal(userData)
		req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 检查数据库是否仍然存在
		var count int64
		db.Model(&models.User{}).Count(&count)

		if count == 0 {
			t.Error("数据库被SQL注入破坏")
		}

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d，实际 %d", http.StatusOK, w.Code)
		}
	})

	// 测试XSS防护
	t.Run("XSS防护", func(t *testing.T) {
		userData := map[string]interface{}{
			"username": "<script>alert('XSS')</script>",
			"password": "Pass123!",
		}

		jsonData, _ := json.Marshal(userData)
		req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d，实际 %d", http.StatusOK, w.Code)
		}

		// 检查响应是否包含XSS脚本
		response := w.Body.String()
		if strings.Contains(response, "<script>") || strings.Contains(response, "onerror=") {
			t.Error("响应中包含XSS脚本")
		}
	})
}

func setupTestRouter() *gin.Engine {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建路由
	router := routers.SetupRouter()

	return router
}

// 添加一个辅助函数来生成密码哈希
func generateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
