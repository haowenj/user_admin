package controllers

import (
	"fmt"
	"net/http"
	"time"

	"employee-management/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("注册请求 - 用户名: %s, 密码长度: %d\n", input.Username, len(input.Password))

	// 添加密码验证
	if input.Password == "" {
		fmt.Printf("错误: 密码为空\n")
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码不能为空"})
		return
	}

	if len(input.Password) < 6 {
		fmt.Printf("错误: 密码太短，长度: %d\n", len(input.Password))
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码长度不能少于6位"})
		return
	}

	var existingUser models.User
	if err := models.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		fmt.Printf("用户已存在: %s\n", input.Username)
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	fmt.Printf("开始加密密码...\n")
	fmt.Printf("原始密码: %s\n", input.Password)
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("密码加密失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	input.Password = string(hash)
	fmt.Printf("加密完成 - 哈希值: %s\n", input.Password)
	fmt.Printf("哈希长度: %d\n", len(input.Password))

	if err := models.DB.Create(&input).Error; err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}

	fmt.Printf("用户注册成功: %s\n", input.Username)
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("登录请求 - 用户名: %s, 密码: %s\n", input.Username, input.Password)

	var user models.User
	if err := models.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		fmt.Printf("用户查询错误: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	fmt.Printf("找到用户 - ID: %d, 用户名: %s\n", user.ID, user.Username)
	fmt.Printf("数据库密码: %s\n", user.Password)
	fmt.Printf("密码长度: %d\n", len(user.Password))

	// 检查密码是否为空
	if user.Password == "" {
		fmt.Printf("警告: 数据库中的密码为空！\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户数据异常"})
		return
	}

	// 检查密码格式
	if len(user.Password) < 20 {
		fmt.Printf("警告: 数据库中的密码格式异常！长度: %d\n", len(user.Password))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		fmt.Printf("密码校验失败: %v\n", err)
		fmt.Printf("尝试手动验证...\n")

		// 手动验证：用相同的密码重新哈希
		rehash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf("重新哈希失败: %v\n", err)
		} else {
			fmt.Printf("重新哈希结果: %s\n", rehash)
			fmt.Printf("哈希匹配: %v\n", user.Password == string(rehash))
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	fmt.Printf("密码校验成功！\n")
	c.JSON(http.StatusOK, gin.H{"message": "登录成功", "user_id": user.ID})
}

func ChangePassword(c *gin.Context) {
	var input struct {
		Username    string `json:"username" binding:"required"`
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := models.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "原密码错误"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	if err := models.DB.Model(&user).Update("password", string(hash)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码修改失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

func GetEmployees(c *gin.Context) {
	var employees []models.Employee
	if err := models.DB.Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取员工列表失败"})
		return
	}
	c.JSON(http.StatusOK, employees)
}

func GetEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee
	if err := models.DB.First(&employee, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "员工不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取员工信息失败"})
		}
		return
	}
	c.JSON(http.StatusOK, employee)
}

func CreateEmployee(c *gin.Context) {
	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 确保日期格式正确
	if employee.HireDate == "" {
		employee.HireDate = time.Now().Format("2006-01-02")
	}

	fmt.Printf("创建员工 - 姓名: %s, 部门: %s, 职位: %s, 入职日期: %s\n",
		employee.Name, employee.Department, employee.Position, employee.HireDate)

	if err := models.DB.Create(&employee).Error; err != nil {
		fmt.Printf("创建员工失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建员工失败"})
		return
	}

	c.JSON(http.StatusOK, employee)
}

func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee
	if err := models.DB.First(&employee, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "员工不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取员工信息失败"})
		}
		return
	}

	var input models.Employee
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 确保日期格式正确
	if input.HireDate == "" {
		input.HireDate = time.Now().Format("2006-01-02")
	}

	if err := models.DB.Model(&employee).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新员工信息失败"})
		return
	}

	c.JSON(http.StatusOK, employee)
}

func DeleteEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee
	if err := models.DB.First(&employee, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "员工不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取员工信息失败"})
		}
		return
	}

	if err := models.DB.Delete(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除员工失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
