# zero-api-middleware

| 名称        | 作用                                 |
| ----------- | ------------------------------------ |
| auth        | 基本认证，摘要认证                   |
| bodylimit   | 限制请求体大小                       |
| casbin      | 访问控制(未实现)                     |
| cors        | 跨域控制                             |
| csrf        | 跨站请求伪造防御                     |
| jwt         | jwt 验证                             |
| limiter     | 限流，全局                           |
| logger      | 请求日志                             |
| must-param  | 必要参数检查                         |
| newrelic    | 监控                                 |
| nonce       | 随机参数 nonce 重复检查              |
| opentracing | 追踪                                 |
| sign        | 签名验证                             |
| throttle    | 限流，默认指定每一个 ip 的每一个请求 |
| timestamp   | 时间戳检查，与当前时间不得相差太多   |
