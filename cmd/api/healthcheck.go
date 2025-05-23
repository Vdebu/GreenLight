package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 返回服务器的状态 运行环境 版本
func (app *application) healthcheckHandler(c *gin.Context) {
	// 使用自定义类型封装Json数据
	env := envelop{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}
	// 自动初始化并输出json数据
	app.writeJson(c, http.StatusOK, env, nil)
}
