# HTTP客户端测试工具
## 说明
- 本程序作为HTTP客户端，通过config.json配置，进行HTTP接口测试

## 配置说明
- **config.json**
```json
{
  "test_name": "readers_get",
  "url": "http://192.168.1.8:9000/",
  "basic": "xxxxx",
  "method": "POST",
  "post_content": "123",
  "total_count": 1000,
  "count_per_second": 100
}
```
- **说明：**

>"test_name":        测试名称

>"url":              http的url

>"basic":            basic认证的字符串

>"method":           http方法

>"post_content":     如果是post，put等方法，post的内容

>"total_count":      总请求次数

>"count_per_second": 每秒请求次数


  代码每0.1秒进行一次请求发送，每次发送count_per_second/10条数据
  
