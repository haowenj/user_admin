package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"employee-management/config"
	"employee-management/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var testUser models.User
var testEmployee models.Employee

func setupTestDB() {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

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

	// 创建测试用户
	hashedPassword, _ := generateHash("TestPass123!")
	testUser = models.User{
		Username: "testuser",
		Password: hashedPassword,
	}
	models.DB.Create(&testUser)

	// 创建测试员工
	testEmployee = models.Employee{
		Name:       "测试员工",
		Age:        25,
		Gender:     "男",
		Department: "技术部",
		Position:   "软件工程师",
		HireDate:   time.Now().Format("2006-01-02"),
	}
	models.DB.Create(&testEmployee)
}

func teardownTestDB() {
	// 清理测试数据
	if models.DB != nil {
		models.DB.Exec("DELETE FROM users")
		models.DB.Exec("DELETE FROM employees")
		models.DB.Exec("DROP TABLE users")
		models.DB.Exec("DROP TABLE employees")
	}
}

func generateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func TestRegisterUser(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	tests := []struct {
		name         string
		input        map[string]interface{}
		expectedCode int
		expectedMsg  string
	}{
		{
			name: "正常用户注册",
			input: map[string]interface{}{
				"username": "newuser",
				"password": "NewPass123!",
			},
			expectedCode: http.StatusOK,
			expectedMsg:  "注册成功",
		},
		// {
		// 	name: "重复用户名注册",
		// 	input: map[string]interface{}{
		// 		"username": "testuser",
		// 		"password": "Pass123!",
		// 	},
		// 	expectedCode: http.StatusBadRequest,
		// 	expectedMsg:  "用户名已存在",
		// },
		// {
		// 	name: "空用户名注册",
		// 	input: map[string]interface{}{
		// 		"username": "",
		// 		"password": "Pass123!",
		// 	},
		// 	expectedCode: http.StatusBadRequest,
		// 	expectedMsg:  "密码不能为空",
		// },
		// {
		// 	name: "短密码注册",
		// 	input: map[string]interface{}{
		// 		"username": "shortuser",
		// 		"password": "123",
		// 	},
		// 	expectedCode: http.StatusBadRequest,
		// 	expectedMsg:  "密码长度不能少于6位",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.POST("/api/user/register", Register)

			jsonData, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("期望状态码 %d，实际 %d", tt.expectedCode, w.Code)
			}

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if msg, ok := response["error"].(string); ok {
				if !strings.Contains(msg, tt.expectedMsg) {
					t.Errorf("期望错误消息包含 '%s'，实际 '%s'", tt.expectedMsg, msg)
				}
			} else if msg, ok := response["message"].(string); ok {
				if msg != tt.expectedMsg {
					t.Errorf("期望消息 '%s'，实际 '%s'", tt.expectedMsg, msg)
				}
			}
		})
	}
}

func TestLoginUser(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	tests := []struct {
		name         string
		input        map[string]interface{}
		expectedCode int
		expectedMsg  string
	}{
		{
			name: "正常用户登录",
			input: map[string]interface{}{
				"username": "testuser",
				"password": "TestPass123!",
			},
			expectedCode: http.StatusOK,
			expectedMsg:  "登录成功",
		},
		{
			name: "错误密码登录",
			input: map[string]interface{}{
				"username": "testuser",
				"password": "WrongPass123!",
			},
			expectedCode: http.StatusUnauthorized,
			expectedMsg:  "用户名或密码错误",
		},
		{
			name: "不存在用户登录",
			input: map[string]interface{}{
				"username": "nonexistentuser",
				"password": "AnyPass123!",
			},
			expectedCode: http.StatusUnauthorized,
			expectedMsg:  "用户名或密码错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.POST("/api/user/login", Login)

			jsonData, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("期望状态码 %d，实际 %d", tt.expectedCode, w.Code)
			}

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if msg, ok := response["error"].(string); ok {
				if !strings.Contains(msg, tt.expectedMsg) {
					t.Errorf("期望错误消息包含 '%s'，实际 '%s'", tt.expectedMsg, msg)
				}
			} else if msg, ok := response["message"].(string); ok {
				if msg != tt.expectedMsg {
					t.Errorf("期望消息 '%s'，实际 '%s'", tt.expectedMsg, msg)
				}
			}
		})
	}
}

func TestChangePassword(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	tests := []struct {
		name         string
		input        map[string]interface{}
		expectedCode int
		expectedMsg  string
	}{
		{
			name: "正常密码修改",
			input: map[string]interface{}{
				"username":     "testuser",
				"old_password": "TestPass123!",
				"new_password": "NewPass456!",
			},
			expectedCode: http.StatusOK,
			expectedMsg:  "密码修改成功",
		},
		{
			name: "错误原密码修改",
			input: map[string]interface{}{
				"username":     "testuser",
				"old_password": "WrongOldPass!",
				"new_password": "NewPass456!",
			},
			expectedCode: http.StatusUnauthorized,
			expectedMsg:  "原密码错误",
		},
		{
			name: "弱新密码修改",
			input: map[string]interface{}{
				"username":     "testuser",
				"old_password": "TestPass123!",
				"new_password": "weak",
			},
			expectedCode: http.StatusOK,
			expectedMsg:  "密码修改成功",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.POST("/api/user/change-password", ChangePassword)

			jsonData, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/api/user/change-password", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("期望状态码 %d，实际 %d", tt.expectedCode, w.Code)
			}

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if msg, ok := response["error"].(string); ok {
				if !strings.Contains(msg, tt.expectedMsg) {
					t.Errorf("期望错误消息包含 '%s'，实际 '%s'", tt.expectedMsg, msg)
				}
			} else if msg, ok := response["message"].(string); ok {
				if msg != tt.expectedMsg {
					t.Errorf("期望消息 '%s'，实际 '%s'", tt.expectedMsg, msg)
				}
			}
		})
	}
}

func TestGetEmployees(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	router := gin.Default()
	router.GET("/api/employee", GetEmployees)

	req, _ := http.NewRequest("GET", "/api/employee", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 %d，实际 %d", http.StatusOK, w.Code)
	}

	var response []models.Employee
	json.Unmarshal(w.Body.Bytes(), &response)

	if len(response) == 0 {
		t.Error("期望返回员工列表，实际为空")
	}
}

func TestGetEmployee(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	router := gin.Default()
	router.GET("/api/employee/:id", GetEmployee)

	// 测试获取存在的员工
	req, _ := http.NewRequest("GET", "/api/employee/"+string(rune(testEmployee.ID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 %d，实际 %d", http.StatusOK, w.Code)
	}

	var response models.Employee
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.ID != testEmployee.ID {
		t.Errorf("期望员工ID %d，实际 %d", testEmployee.ID, response.ID)
	}

	// 测试获取不存在的员工
	req, _ = http.NewRequest("GET", "/api/employee/999", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("期望状态码 %d，实际 %d", http.StatusNotFound, w.Code)
	}
}

func TestCreateEmployee(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	tests := []struct {
		name         string
		input        map[string]interface{}
		expectedCode int
		expectedMsg  string
	}{
		{
			name: "正常员工创建",
			input: map[string]interface{}{
				"name":       "李四",
				"age":        30,
				"gender":     "女",
				"department": "市场部",
				"position":   "市场经理",
			},
			expectedCode: http.StatusOK,
			expectedMsg:  "",
		},
		{
			name: "创建员工时包含无效数据",
			input: map[string]interface{}{
				"name":       "",
				"age":        17,
				"gender":     "未知",
				"department": "",
				"position":   "",
			},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.POST("/api/employee", CreateEmployee)

			jsonData, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/api/employee", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("期望状态码 %d，实际 %d", tt.expectedCode, w.Code)
			}

			if tt.expectedMsg != "" {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)

				if msg, ok := response["error"].(string); ok {
					if !strings.Contains(msg, tt.expectedMsg) {
						t.Errorf("期望错误消息包含 '%s'，实际 '%s'", tt.expectedMsg, msg)
					}
				}
			}
		})
	}
}

func TestUpdateEmployee(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	tests := []struct {
		name         string
		input        map[string]interface{}
		expectedCode int
		expectedMsg  string
	}{
		{
			name: "正常员工更新",
			input: map[string]interface{}{
				"name":       "张三（更新）",
				"age":        29,
				"department": "研发部",
				"position":   "高级软件工程师",
			},
			expectedCode: http.StatusOK,
			expectedMsg:  "",
		},
		{
			name: "更新不存在员工",
			input: map[string]interface{}{
				"name": "新员工",
			},
			expectedCode: http.StatusNotFound,
			expectedMsg:  "员工不存在",
		},
		{
			name: "更新员工时包含无效数据",
			input: map[string]interface{}{
				"name":   "",
				"age":    67,
				"gender": "未知",
			},
			expectedCode: http.StatusOK,
			expectedMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.PUT("/api/employee/:id", UpdateEmployee)

			// 如果是测试更新不存在员工，使用999作为ID
			id := string(rune(testEmployee.ID))
			if tt.name == "更新不存在员工" {
				id = "999"
			}

			jsonData, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("PUT", "/api/employee/"+id, bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("期望状态码 %d，实际 %d", tt.expectedCode, w.Code)
			}

			if tt.expectedMsg != "" {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)

				if msg, ok := response["error"].(string); ok {
					if !strings.Contains(msg, tt.expectedMsg) {
						t.Errorf("期望错误消息包含 '%s'，实际 '%s'", tt.expectedMsg, msg)
					}
				}
			}
		})
	}
}

func TestDeleteEmployee(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	router := gin.Default()
	router.DELETE("/api/employee/:id", DeleteEmployee)

	// 测试删除存在的员工
	req, _ := http.NewRequest("DELETE", "/api/employee/"+string(rune(testEmployee.ID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 %d，实际 %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if msg, ok := response["message"].(string); ok && msg != "删除成功" {
		t.Errorf("期望消息 '删除成功'，实际 '%s'", msg)
	}

	// 测试删除不存在的员工
	req, _ = http.NewRequest("DELETE", "/api/employee/999", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("期望状态码 %d，实际 %d", http.StatusNotFound, w.Code)
	}
}

// SQL注入测试
func TestSQLInjection(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	router := gin.Default()
	router.POST("/api/user/register", Register)

	tests := []struct {
		name         string
		input        map[string]interface{}
		expectedCode int
	}{
		{
			name: "用户名SQL注入",
			input: map[string]interface{}{
				"username": "admin';--",
				"password": "Pass123!",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "用户名SQL注入2",
			input: map[string]interface{}{
				"username": "Robert'); DROP TABLE users;--",
				"password": "Pass123!",
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 检查数据库是否仍然存在
			var count int64
			models.DB.Model(&models.User{}).Count(&count)

			if count == 0 {
				t.Error("数据库被SQL注入破坏")
			}

			if w.Code != tt.expectedCode {
				t.Errorf("期望状态码 %d，实际 %d", tt.expectedCode, w.Code)
			}
		})
	}
}

// XSS攻击测试
func TestXSSAttack(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	router := gin.Default()
	router.POST("/api/user/register", Register)

	tests := []struct {
		name         string
		input        map[string]interface{}
		expectedCode int
	}{
		{
			name: "用户名XSS攻击",
			input: map[string]interface{}{
				"username": "<script>alert('XSS')</script>",
				"password": "Pass123!",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "员工信息XSS攻击",
			input: map[string]interface{}{
				"name": "<img src='x' onerror='alert(\"XSS\")'>",
				"age":  30,
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("期望状态码 %d，实际 %d", tt.expectedCode, w.Code)
			}

			// 检查响应是否包含XSS脚本
			response := w.Body.String()
			if strings.Contains(response, "<script>") || strings.Contains(response, "onerror=") {
				t.Error("响应中包含XSS脚本")
			}
		})
	}
}

// 性能测试
func TestPerformance(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	router := gin.Default()
	router.POST("/api/user/register", Register)

	// 测试并发注册
	t.Run("并发用户注册", func(t *testing.T) {
		if testing.Short() {
			t.Skip("跳过性能测试")
		}

		concurrency := 10
		requests := 100

		done := make(chan bool, concurrency)
		errors := make(chan error, concurrency)

		for i := 0; i < concurrency; i++ {
			go func(id int) {
				defer func() { done <- true }()

				for j := 0; j < requests/concurrency; j++ {
					input := map[string]interface{}{
						"username": fmt.Sprintf("user%d_%d", id, j),
						"password": "Pass123!",
					}

					jsonData, _ := json.Marshal(input)
					req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(jsonData))
					req.Header.Set("Content-Type", "application/json")

					w := httptest.NewRecorder()
					start := time.Now()
					router.ServeHTTP(w, req)
					duration := time.Since(start)

					if w.Code != http.StatusOK {
						errors <- fmt.Errorf("请求 %d_%d 失败: 状态码 %d", id, j, w.Code)
					}

					if duration > time.Second*5 {
						errors <- fmt.Errorf("请求 %d_%d 超时: %v", id, j, duration)
					}
				}
			}(i)
		}

		// 等待所有goroutine完成
		for i := 0; i < concurrency; i++ {
			select {
			case <-done:
			case err := <-errors:
				t.Error(err)
			case <-time.After(time.Second * 10):
				t.Error("测试超时")
			}
		}
	})
}
