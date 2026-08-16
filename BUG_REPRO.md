# 修复前故障复现（Docker）

## 项目与标准命令

Recipe Planner Service 是一个本地运行的菜谱与周计划服务，支持录入菜谱、生成用餐计划及查询采购清单。在仓库根目录可执行以下标准命令：

```bash
go build ./...
go run ./cmd/server
go test ./...
```

## 环境构建与编译

已实际执行以下命令，`linux/amd64` 和 `linux/arm64` 的镜像构建与容器内编译均成功：

```bash
docker build --platform linux/amd64 --file benzhi.Dockerfile --tag recipe-planner-bug003-base:amd64 .
docker run --rm --platform linux/amd64 recipe-planner-bug003-base:amd64 go build ./...
docker build --platform linux/arm64 --file benzhi.Dockerfile --tag recipe-planner-bug003-base:arm64 .
docker run --rm --platform linux/arm64 recipe-planner-bug003-base:arm64 go build ./...
```

## 故障触发步骤

从仓库根目录执行以下命令：

```bash
docker run --rm --platform linux/arm64 recipe-planner-bug003-base:arm64 go test ./...
```

该命令已连续执行 20 次，均稳定失败。

## 实际错误输出

```text
?   	github.com/1260124186-cc/recipe-planner-service/cmd/server	[no test files]
?   	github.com/1260124186-cc/recipe-planner-service/internal/domain	[no test files]
--- FAIL: TestGeneratePlanStopsWhenRequestIsCanceledAfterRecipesAreRead (0.00s)
    planner_context_test.go:60: GeneratePlan error = <nil>, want context cancellation
--- FAIL: TestMemoryStoreRejectsCanceledPlanSave (0.00s)
    planner_context_test.go:78: SavePlan error = <nil>, want context cancellation
FAIL
FAIL	github.com/1260124186-cc/recipe-planner-service/internal/service	0.012s
?   	github.com/1260124186-cc/recipe-planner-service/internal/store	[no test files]
ok  	github.com/1260124186-cc/recipe-planner-service/internal/transport	0.008s
FAIL
退出状态：1
```

## 期望行为

当用户在生成周计划的过程中取消请求时，操作应返回取消结果，且不应创建新的计划。随后查询计划时，不应看到该次已取消操作对应的计划。
