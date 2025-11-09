# Gin

Gin 是一个 Go 语言写的 Web 框架。

## 1 安装

<<< @/go/codes/gin/install.sh

## 2 HelloWorld

<<< @/go/codes/gin/hello_world.go

## 3 路由分组

假设现在有这些请求路径：`/goods/list`、`goods/add`、`goods/del`，在不使用路由分组的情况下，通常会这样写：

<<< @/go/codes/gin/router.go

这些路径都包含 `/goods`，可以通过路由分组进行管理：

<<< @/go/codes/gin/group.go

## 4 获取 URL 上的变量

假设现在需要通过 `/goods/[商品ID]` 访问对应商品的详情，通过如下方法可以识别 URL 中的 ID 变量：

<<< @/go/codes/gin/param.go

假设现在想通过 `/goods/[商品ID]/[动作]` 来对商品进行操作，可以这样做：

<<< @/go/codes/gin/params.go

```http
### 请求
GET http://localhost:8000/goods/123/delete

### 响应
{
  "action": "delete",
  "id": "123"
}
```

通过 `/goods/:id/*action` 也能实现，但是和 `/:action` 有些区别：

<<< @/go/codes/gin/params2.go

```http
### 请求
GET http://localhost:8000/goods/123/delete

### 响应
{
  "action": "/delete",
  "id": "123"
}


### 请求
GET http://localhost:8000/goods/123/delete/test

### 响应
{
  "action": "/delete/test",
  "id": "123"
}


### 请求
GET http://localhost:8000/goods/123

### 响应
{
  "action": "/",
  "id": "123"
}


### 请求
GET http://localhost:8000/goods/123/

### 响应
{
  "action": "/",
  "id": "123"
}
```

还可以直接绑定 Uri 到一个结构体：

<<< @/go/codes/gin/bind.go

```http
GET http://localhost:8000/goods/123/abc

# 响应
{
  "id": 123,
  "name": "abc"
}
```

## 5 获取表单信息

### 5.1 Get 请求

<<< @/go/codes/gin/get.go

测试：

```http
GET http://localhost:8000/hello

# 响应
{
  "framework": "Gin",
  "lang": ""
}
```

```http
GET http://localhost:8000/hello?lang=Java&framework=Spring

# 响应
{
  "framework": "Spring",
  "lang": "Java"
}
```

### 5.2 Post 请求

<<< @/go/codes/gin/post.go

测试：

```http
### 请求
POST http://localhost:8000/hello
Content-Type: application/x-www-form-urlencoded

lang = Go &
framework = Gin


# 响应
{
  "framework": "Gin",
  "lang": "Go"
}
```

## 6 Protobuf 渲染

定义 Protobuf 消息：

<<< @/go/codes/gin/msg.proto

生成 Go 包，在程序中导入使用：

<<< @/go/codes/gin/proto.go{12}

## 7 表单验证

<<< @/go/codes/gin/validate.go

验证不通过时，返回的信息是英文的，如下所示：

```json
{
    "error": "
		Key: 'SignUpInfo.Username' Error:Field validation for 'Username' failed on the 'min' tag\n
		Key: 'SignUpInfo.RePassword' Error:Field validation for 'RePassword' failed on the 'eqfield' tag\n
		Key: 'SignUpInfo.Email' Error:Field validation for 'Email' failed on the 'email' tag\n
		Key: 'SignUpInfo.Age' Error:Field validation for 'Age' failed on the 'lte' tag
	"
}
```

可将错误信息转成中文：

<<< @/go/codes/gin/translation.go

此时，错误信息将转为中文：

<<< @/go/codes/gin/error.json

## 8 中间件

<<< @/go/codes/gin/midware.go

中间件后续逻辑的执行终止，必须使用 `context.Abort()`，直接 `return` 无法阻止执行：

<<< @/go/codes/gin/abort.go

从添加中间件的 `Use` 函数源码上看：

<<< @/go/codes/gin/use.go

本质上是把中间件追加到了 `group.Handlers` 切片的后面。

而在 `GET/POST` 函数内部执行了一个 `combineHandlers` 函数：

<<< @/go/codes/gin/combineHandlers.go

这里将我们处理请求的函数和之前的 `Handlers` 拼在了一起。

而 `context.Next()` 的调用过程如下：

<<< @/go/codes/gin/next.go

如果在自定义中间件里直接 `return` 了，只代表当前中间件的逻辑结束了，`Handlers` 中后续的函数仍然会依次执行。

从 `context.Next()` 的源码逻辑可以看出，真正决定 `Handlers` 中的函数调用的是 `c.index`。而 `context.Abort()` 的作用正是修改这个变量：

<<< @/go/codes/gin/abortIndex.go

## 9 优雅退出

关闭程序的时候可能有请求还没有处理完，此时处理过程就会被迫中断。优雅退出其实就是在程序关闭时，不暴力关闭，而是要等待进程中的逻辑处理完成后，才关闭。

<<< @/go/codes/gin/quit.go
