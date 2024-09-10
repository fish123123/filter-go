package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// GET 请求处理
	router.GET("/getparam", func(c *gin.Context) {
		requestInfo := printRequest(c.Request)
		param := c.Query("param")
		c.JSON(200, gin.H{
			"message":     "发送的数据: " + param,
			"requestInfo": requestInfo,
		})
	})

	// POST 请求处理
	router.POST("/postparam", func(c *gin.Context) {
		requestInfo := printRequest(c.Request)
		var requestBody struct {
			Param string `json:"param"`
		}
		if err := c.ShouldBindJSON(&requestBody); err == nil {
			c.JSON(200, gin.H{
				"message":     "发送的数据: " + requestBody.Param,
				"requestInfo": requestInfo,
			})
		} else {
			c.JSON(400, gin.H{"error": err.Error()})
		}
	})

	// 敏感字符处理接口
	router.GET("/sensitive", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sensitive.html", gin.H{
			"message": "请求中存在敏感字符，请检查后重新提交",
		})
	})

	// 加载模版文件
	router.LoadHTMLGlob("templates/*")

	// 渲染首页
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.Run(":8080")
}

// 打印请求信息
func printRequest(r *http.Request) string {
	var b bytes.Buffer
	b.WriteString(r.Method + " " + r.URL.String() + " " + r.Proto + "\n")
	host := r.Header.Get("Host")
	if host != "" {
		b.WriteString("Host: " + host + "\n")
	}

	for k, v := range r.Header {
		if k != "Host" {
			b.WriteString(k + ": " + v[0] + "\n")
		}
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	b.WriteString("\n")
	b.WriteString(string(bodyBytes))

	var requestBody interface{}
	json.Unmarshal(bodyBytes, &requestBody)

	return b.String()
}
