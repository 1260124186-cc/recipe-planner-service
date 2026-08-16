# 修复前故障复现（Docker）

## 项目与标准命令

Recipe Planner Service 是一个本地运行的菜谱、食材库存与用餐计划服务，提供计划生成、查询、烹饪和采购清单等 HTTP 接口。

在仓库根目录可执行的标准命令：

```bash
go build ./...
go run ./cmd/server
go test ./...
```

## 环境构建与编译

已实际执行以下 `linux/arm64` Docker 镜像构建与容器内编译命令，均成功：

```bash
docker build --platform linux/arm64 -f benzhi.Dockerfile -t recipe-planner-service-bug-004-base:verify .
docker run --rm --platform linux/arm64 --entrypoint go recipe-planner-service-bug-004-base:verify build ./...
```

已实际执行以下 `linux/amd64` Docker 镜像构建与容器内编译命令，均成功：

```bash
docker build --platform linux/amd64 -f benzhi.Dockerfile -t recipe-planner-service-bug-004-base:amd64 .
docker run --rm --platform linux/amd64 --entrypoint go recipe-planner-service-bug-004-base:amd64 build ./...
```

## 故障触发步骤

在仓库根目录先完成上述 `linux/arm64` 镜像构建，然后执行：

```bash
docker run --rm --platform linux/arm64 --entrypoint go recipe-planner-service-bug-004-base:verify test ./...
```

该命令稳定失败；连续执行 20 次均可复现。

## 实际错误输出

```text
?   	github.com/1260124186-cc/recipe-planner-service/cmd/server	[no test files]
?   	github.com/1260124186-cc/recipe-planner-service/internal/domain	[no test files]
ok  	github.com/1260124186-cc/recipe-planner-service/internal/service	0.019s
--- FAIL: TestMissingPlanReadDoesNotBlockLaterWrites (0.15s)
    memory_defer_test.go:25: a failed plan operation left the memory store blocked
--- FAIL: TestMissingPlanCompletionDoesNotBlockLaterReads (0.15s)
    memory_defer_test.go:41: a failed plan operation left the memory store blocked
FAIL
FAIL	github.com/1260124186-cc/recipe-planner-service/internal/store	0.168s
ok  	github.com/1260124186-cc/recipe-planner-service/internal/transport	0.004s
FAIL
```

命令以非零状态退出。

## 期望行为

当不存在的计划请求返回错误后，后续有效的计划创建、查询或烹饪请求仍应及时完成，服务不应长期无响应，也不需要依靠重启恢复。
