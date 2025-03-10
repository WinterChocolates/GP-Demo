package docs

import "github.com/swaggo/swag"

var SwaggerInfo = &swag.Spec{
	Version:     "1.0",
	Host:        "localhost:8080",
	BasePath:    "/api/v1",
	Schemes:     []string{"http", "https"},
	Title:       "招聘系统API文档",
	Description: "招聘系统后端API接口文档",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}