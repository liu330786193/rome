用户名：admin
密码： admin123.

# 本地运行-前端
### 拉取项目
git clone https://gitee.com/objs/mayfly-go.git

```
#### 进入前端项目文件夹
cd mayfly-go-web</code>

#### yarn与npm均可
#### 下载依赖等
yarn install || npm install

#### 本地运行
yarn run dev || npm run dev
```

# 本地运行-后端

```
# 拉取项目
git clone https://gitee.com/objs/mayfly-go.git

# 安装依赖
go mod tidy

# 进入server文件夹
cd server

# 本地运行项目
go run main.go

# 编译打包项目
# 打包amd版linux可执行文件
CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o mayfly-go-linux-amd64 main.go
# 打包arm版linux可执行文件
CGO_ENABLE=0 GOOS=linux GOARCH=arm64 go build -o mayfly-go-linux-arm64 main.go
# 打包windows版可执行文件
CGO_ENABLE=0 GOOS=windows GOARCH=amd64 go build -o mayfly-go-win main.go

```