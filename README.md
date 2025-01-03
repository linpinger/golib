# golib
some tool for golang

- 一些自己常用的函数放到这个库里，可能随时改的妈都不认识^_^


## Log

- 2025-01-03: FML.go 修复 `DescDelBlankPage` 一直以来的一个bug，会删除最后正常的一章
- 2025-01-02: FML.go 修改 `DescDelBlankPage` 方法: if 第一章字节数 小于 6000  全清 else 倒序清空
- 2024-10-25: http.go 添加 环境变量代理 `Proxy: http.ProxyFromEnvironment`
- 2022-04-19: http.go 中添加 `io.Copy(ioutil.Discard, response.Body)` 以重用连接
- 2024-10-10: FML.go 修改 `GetAllBlankPages` 参数: bool 改为 int（空白章节最大字节）
- 2024-08-30: FML.go 添加 `DescDelBlankPage` 方法: 倒序清空空白章节（貌似有bug）
- 2024-08-28: FML.go 修改 排序 方法
- 2023-08-07: novelsite.go 修改起点无效的旧url为新的桌面版url
- 2022-04-19: http.go 中添加 `io.Copy(ioutil.Discard, response.Body)` 以重用连接
- 2021-12-16: 将 foxbook-golang 中的库移过来

