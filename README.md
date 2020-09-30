# Panda
Panda 是一个由 Golang 实现的代理服务器

# TODO

服务端兼容常见代理软件的非加密模式和部分对称加密

# HTTP GET

```
GET http://www.huang314.cn/img/blog.png HTTP/1.1\r\nHost: www.huang314.cn\r\nUser-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:80.0) Gecko/20100101 Firefox/80.0\r\nAccept: image/webp,*/*\r\nAccept-Language: zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2\r\nAccept-Encoding: gzip, deflate\r\nConnection: keep-alive\r\nReferer: http://www.huang314.cn/\r\nCookie: JSESSIONID=AFA081E95160BE57AD4A441D9D9D374B\r\n\r\n
```

# 注意

1. go 的切片只有 append 才能拓展

2. append([], []...)是要截取前面的切片，否则会将后面的切片会连在，前面的切片的未使用区间的后面

# 参考

[socks5 协议详解](https://jiajunhuang.com/articles/2019_06_06-socks5.md.html)

[渗透基础——使用Go语言开发socks代理工具](https://3gstudent.github.io/3gstudent.github.io/%E6%B8%97%E9%80%8F%E5%9F%BA%E7%A1%80-%E4%BD%BF%E7%94%A8Go%E8%AF%AD%E8%A8%80%E5%BC%80%E5%8F%91socks%E4%BB%A3%E7%90%86%E5%B7%A5%E5%85%B7/)

[socks5proxy](https://github.com/shikanon/socks5proxy)