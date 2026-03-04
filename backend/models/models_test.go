package models

import (
	"testing"
	"time"

	"employee-management/config"
)

func TestUserModel(t *testing.T) {
	// 初始化测试配置
	if err := config.InitConfig("config.yaml"); err != nil {
		t.Fatalf("加载配置失败: %v", err)
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
	if err := InitDB(); err != nil {
		t.Fatalf("连接数据库失败: %v", err)
	}

	// 自动迁移表
	Migrate()

	// 测试用户创建
	user := User{
		Username: "testuser",
		Password: "hashedpassword",
	}

	result := DB.Create(&user)
	if result.Error != nil {
		t.Errorf("创建用户失败: %v", result.Error)
	}

	// 验证用户ID是否正确生成
	if user.ID == 0 {
		t.Error("用户ID未正确生成")
	}

	// 测试用户查询
	var foundUser User
	result = DB.First(&foundUser, user.ID)
	if result.Error != nil {
		t.Errorf("查询用户失败: %v", result.Error)
	}

	if foundUser.Username != "testuser" {
		t.Errorf("期望用户名 'testuser'，实际 '%s'", foundUser.Username)
	}

	// 测试用户名唯一性约束
	duplicateUser := User{
		Username: "testuser", // 重复用户名
		Password: "anotherpassword",
	}
	result = DB.Create(&duplicateUser)
	if result.Error == nil {
		t.Error("应该因为用户名重复而创建失败")
	}

	// 清理测试数据
	DB.Delete(&user)
}

func TestEmployeeModel(t *testing.T) {
	// 初始化测试配置
	if err := config.InitConfig("config.yaml"); err != nil {
		t.Fatalf("加载配置失败: %v", err)
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
	if err := InitDB(); err != nil {
		t.Fatalf("连接数据库失败: %v", err)
	}

	// 自动迁移表
	Migrate()

	// 测试员工创建
	employee := Employee{
		Name:       "张三",
		Age:        28,
		Gender:     "男",
		Department: "技术部",
		Position:   "软件工程师",
		HireDate:   time.Now().Format("2006-01-02"),
	}

	result := DB.Create(&employee)
	if result.Error != nil {
		t.Errorf("创建员工失败: %v", result.Error)
	}

	// 验证员工ID是否正确生成
	if employee.ID == 0 {
		t.Error("员工ID未正确生成")
	}

	// 测试员工查询
	var foundEmployee Employee
	result = DB.First(&foundEmployee, employee.ID)
	if result.Error != nil {
		t.Errorf("查询员工失败: %v", result.Error)
	}

	if foundEmployee.Name != "张三" {
		t.Errorf("期望员工姓名 '张三'，实际 '%s'", foundEmployee.Name)
	}

	// 测试员工更新
	updatedData := Employee{
		Name: "张三（更新）",
		Age:  29,
	}
	result = DB.Model(&employee).Updates(updatedData)
	if result.Error != nil {
		t.Errorf("更新员工失败: %v", result.Error)
	}

	// 验证更新结果
	var updatedEmployee Employee
	DB.First(&updatedEmployee, employee.ID)
	if updatedEmployee.Name != "张三（更新）" {
		t.Errorf("期望更新后姓名 '张三（更新）'，实际 '%s'", updatedEmployee.Name)
	}

	// 测试员工软删除
	result = DB.Delete(&employee)
	if result.Error != nil {
		t.Errorf("删除员工失败: %v", result.Error)
	}

	// 验证软删除
	var deletedEmployee Employee
	result = DB.Unscoped().First(&deletedEmployee, employee.ID)
	if result.Error != nil {
		t.Errorf("查询已删除员工失败: %v", result.Error)
	}

	// 在正常查询中应该找不到已删除的员工
	result = DB.First(&deletedEmployee, employee.ID)
	if result.Error == nil {
		t.Error("已删除的员工不应该在正常查询中出现")
	}

	// 清理测试数据
	DB.Unscoped().Delete(&employee)
}

func TestEmployeeValidation(t *testing.T) {
	// 初始化测试配置
	if err := config.InitConfig("config.yaml"); err != nil {
		t.Fatalf("加载配置失败: %v", err)
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
	if err := InitDB(); err != nil {
		t.Fatalf("连接数据库失败: %v", err)
	}

	// 自动迁移表
	Migrate()

	tests := []struct {
		name     string
		employee Employee
		wantErr  bool
	}{
		{
			name: "正常员工数据",
			employee: Employee{
				Name:       "李四",
				Age:        25,
				Gender:     "女",
				Department: "市场部",
				Position:   "市场经理",
				HireDate:   time.Now().Format("2006-01-02"),
			},
			wantErr: false,
		},
		{
			name: "空姓名",
			employee: Employee{
				Name:       "",
				Age:        25,
				Department: "技术部",
			},
			wantErr: true,
		},
		{
			name: "年龄过小",
			employee: Employee{
				Name: "王五",
				Age:  17,
			},
			wantErr: true,
		},
		{
			name: "年龄过大",
			employee: Employee{
				Name: "赵六",
				Age:  66,
			},
			wantErr: true,
		},
		{
			name: "无效性别",
			employee: Employee{
				Name:   "钱七",
				Age:    30,
				Gender: "未知",
			},
			wantErr: true,
		},
		{
			name: "空部门",
			employee: Employee{
				Name: "孙八",
				Age:  28,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DB.Create(&tt.employee)

			if tt.wantErr {
				if result.Error == nil {
					t.Errorf("期望创建失败，但成功了")
				}
			} else {
				if result.Error != nil {
					t.Errorf("期望创建成功，但失败: %v", result.Error)
				}

				// 验证数据完整性
				if tt.employee.Name == "" {
					t.Error("员工姓名不应为空")
				}
				if tt.employee.Age < 18 || tt.employee.Age > 65 {
					t.Errorf("员工年龄应在18-65之间，实际: %d", tt.employee.Age)
				}
			}
		})
	}
}

func TestEmployeeTableName(t *testing.T) {
	employee := Employee{}
	tableName := employee.TableName()

	if tableName != "employees" {
		t.Errorf("期望表名 'employees'，实际 '%s'", tableName)
	}
}

func TestUserTableName(t *testing.T) {
	user := User{}
	tableName := user.TableName()

	if tableName != "users" {
		t.Errorf("期望表名 'users'，实际 '%s'", tableName)
	}
}
