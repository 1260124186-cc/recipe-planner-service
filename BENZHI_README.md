# Recipe Planner Service Docker 交付说明

Recipe Planner Service 用于在本地管理菜谱、食材库存和每日用餐计划，并根据未完成的计划计算采购缺口。服务不连接数据库或第三方在线服务，容器只提供固定的 Go 编译与运行环境，代码始终在容器内从源码编译。

## 本地标准命令

在仓库根目录执行：

```bash
go build ./...
go run ./cmd/server
go test ./...
```

服务默认监听容器和宿主机的 `8080` 端口。启动后可执行：

```bash
curl --fail http://127.0.0.1:8080/health
curl --fail http://127.0.0.1:8080/recipes
```

前一条命令应返回 HTTP 200 和 `{"status":"ok"}`；后一条命令会返回当前菜谱 JSON 数组。

## 版本固定方式

`go.mod` 固定 Go 语言版本为 `1.26.2`。实际 Dockerfile 为 `benzhi.Dockerfile`，其基础镜像固定为 `golang:1.26.2-bookworm`，并设置 `GOTOOLCHAIN=local`，因此容器不会自行下载或切换 Go 工具链。Dockerfile 先执行 `go mod download`，复制源码后执行 `go build ./...`，不会复制宿主机编译产物。

## 双架构标准验收

从仓库根目录分别执行：

```bash
./build_benzhi_docker.sh linux/amd64
./build_benzhi_docker.sh linux/arm64
```

脚本对指定平台依次完成镜像构建、在容器内执行 `go build ./...`、启动 HTTP 服务，并从宿主机请求 `http://127.0.0.1:18080/health`。每一步均以退出码 `0` 为通过条件；健康检查还必须得到 HTTP 200。

## 手工 Docker 操作

以下命令可逐步完成同样的验收：

```bash
docker build --platform linux/amd64 -f benzhi.Dockerfile -t recipe-planner-service:benzhi .
docker run --rm --platform linux/amd64 --entrypoint go recipe-planner-service:benzhi build ./...
docker run -d --rm --name recipe-planner-manual -p 18080:8080 recipe-planner-service:benzhi
curl --fail http://127.0.0.1:18080/health
docker rm -f recipe-planner-manual
```

将四条命令中的 `linux/amd64` 一致替换为 `linux/arm64`，即可进行另一平台的验收。手工运行容器后，健康检查返回 HTTP 200 且容器日志没有启动错误，即表示运行通过；`go build ./...` 和镜像构建都必须以退出码 `0` 结束。
