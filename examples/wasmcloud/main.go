// 使用 TinyGo 编写的 WebAssembly 组件示例
// 构建: tinygo build -target=wasi -gc=leaking -o main.wasm main.go

package main

import (
	"encoding/json"
	"fmt"
	
	http "github.com/wasmcloud/wasmcloud/examples/go/http-hello-world/gen"
)

func init() {
	http.Exports.Handle = handleRequest
}

func handleRequest(request http.Request) http.Response {
	// 记录请求信息
	fmt.Printf("Received request: %s %s\n", request.Method, request.Path)
	
	// 构建响应
	response := map[string]interface{}{
		"message": "Hello from WebAssembly + Go!",
		"runtime": "wasmcloud",
		"language": "tinygo",
		"path":    request.Path,
		"method":  request.Method,
	}
	
	body, err := json.Marshal(response)
	if err != nil {
		return http.Response{
			StatusCode: 500,
			Headers: []http.Header{
				{Name: "Content-Type", Value: "text/plain"},
			},
			Body: []byte("Internal Server Error"),
		}
	}
	
	return http.Response{
		StatusCode: 200,
		Headers: []http.Header{
			{Name: "Content-Type", Value: "application/json"},
		},
		Body: body,
	}
}

func main() {}
