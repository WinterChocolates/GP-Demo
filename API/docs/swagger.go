package docs

// 定义API文档的基本信息
// @title 人力资源管理系统API
// @version 1.0
// @description 人力资源管理系统的后端API服务
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

// 定义通用的响应结构，解决swagger生成时的类型问题
type SwaggerResponse struct {
	Code      int         `json:"code" example:"200"`
	Message   string      `json:"message" example:"成功"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp" example:"1625123456"`
	RequestID string      `json:"request_id,omitempty" example:"req-123456"`
}

// 定义一些常用的响应模型，用于API文档

// 登录请求
type LoginRequest struct {
	Username string `json:"username" example:"admin" binding:"required"`
	Password string `json:"password" example:"password123" binding:"required"`
}

// 登录响应
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  struct {
		ID       uint   `json:"id" example:"1"`
		Username string `json:"username" example:"admin"`
		Usertype string `json:"usertype" example:"admin"`
	} `json:"user"`
}

// 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" example:"newuser" binding:"required"`
	Password string `json:"password" example:"password123" binding:"required"`
	Phone    string `json:"phone" example:"13800138000"`
	Email    string `json:"email" example:"user@example.com"`
}

// 通用分页响应
type PagedResponse struct {
	Total    int64       `json:"total" example:"100"`
	Page     int         `json:"page" example:"1"`
	PageSize int         `json:"page_size" example:"10"`
	Items    interface{} `json:"items"`
}
