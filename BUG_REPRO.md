# 修复前故障复现（Docker）

## 项目与标准命令
Recipe Planner Service 是本地运行的菜谱与周计划服务，可按标签生成每日用餐计划并提供采购清单。请在仓库根目录执行：

```bash
go build ./...
go run ./cmd/server
go test ./...
```

## 环境构建与编译
已实际执行并通过以下 linux/arm64 命令：

```bash
docker build --platform linux/arm64 -f benzhi.Dockerfile -t recipe-planner-service:benzhi .
docker run --rm --platform linux/arm64 --entrypoint go recipe-planner-service:benzhi build ./...
```

已实际执行并通过以下 linux/amd64 命令：

```bash
docker build --platform linux/amd64 -f benzhi.Dockerfile -t recipe-planner-service:benzhi .
docker run --rm --platform linux/amd64 --entrypoint go recipe-planner-service:benzhi build ./...
```

两个平台的镜像构建和容器内编译均成功。

## 故障触发步骤
在仓库根目录执行：

```bash
go test ./...
```

## 实际错误输出
```text
?   	github.com/1260124186-cc/recipe-planner-service/cmd/server	[no test files]
?   	github.com/1260124186-cc/recipe-planner-service/internal/domain	[no test files]
--- FAIL: TestGeneratePlanOnlyUsesRecipesMatchingRequestedTag (0.00s)
    tag_filter_test.go:34: tagged plan selected "quick-pasta", want only warm-soup
FAIL
FAIL	github.com/1260124186-cc/recipe-planner-service/internal/service	1.540s
?   	github.com/1260124186-cc/recipe-planner-service/internal/store	[no test files]
ok  	github.com/1260124186-cc/recipe-planner-service/internal/transport	3.805s
FAIL
```

## 期望行为
使用“warm”等标签生成周计划时，计划中的每一餐都应来自该标签对应的菜谱。用户以不同大小写输入同一标签时，应得到一致的筛选结果。
