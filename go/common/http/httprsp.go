package http

//每一个HTTP响应都应该组合该结构
type BaseRsp struct {
	Code    int    `json:"code,omitempty"`    //业务响应码，可以是标准的HTTP ,也可以自定义
	Message string `json:"message,omitempty"` //处理消息。成功  ok   失败：详细原因
}
