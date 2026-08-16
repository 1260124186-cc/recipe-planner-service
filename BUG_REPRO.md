# 修复前故障复现（Docker）

## 项目与标准命令
Recipe Planner Service 用于管理菜谱、食材库存和每日用餐计划，并根据未完成的计划计算采购缺口。服务使用内存保存运行期间的数据，默认监听 `:8080`。

在仓库根目录执行：

```bash
go build ./...
go run ./cmd/server
go test ./...
```

## 环境构建与编译
在仓库根目录执行以下命令：

```bash
docker build --platform linux/amd64 -f benzhi.Dockerfile -t recipe-planner-service:benzhi .
docker run --rm --platform linux/amd64 --entrypoint go recipe-planner-service:benzhi build ./...
docker build --platform linux/arm64 -f benzhi.Dockerfile -t recipe-planner-service:benzhi .
docker run --rm --platform linux/arm64 --entrypoint go recipe-planner-service:benzhi build ./...
```

以上两个平台的镜像构建和容器内编译均成功。两个平台的标准健康检查命令也均成功：

```bash
./build_benzhi_docker.sh linux/amd64
./build_benzhi_docker.sh linux/arm64
```

## 故障触发步骤
在修复前源码目录根目录执行：

```bash
docker run --rm --platform linux/arm64 --entrypoint go recipe-planner-service:benzhi test ./...
```

## 实际错误输出
```text
?   	github.com/1260124186-cc/recipe-planner-service/cmd/server	[no test files]
?   	github.com/1260124186-cc/recipe-planner-service/internal/domain	[no test files]
--- FAIL: TestRestockRejectsInvalidBatchWithoutChangingPantry (0.00s)
    workflow_test.go:84: invalid restock changed pantry: map[string]int{"tomato":2}
FAIL
FAIL	github.com/1260124186-cc/recipe-planner-service/internal/service	0.006s
--- FAIL: TestRestockDoesNotApplyValidPrefixWhenBatchContainsInvalidItem (0.00s)
    memory_test.go:28: failed batch partially changed pantry: map[string]int{"tomato":2}
FAIL
FAIL	github.com/1260124186-cc/recipe-planner-service/internal/store	0.004s
ok  	github.com/1260124186-cc/recipe-planner-service/internal/transport	0.002s
FAIL
exit code: 1
```

## 期望行为
当一次补货请求包含无效条目时，服务应返回失败，并且本次请求不应改变库存。正常的有效补货请求仍应成功合并同名食材；后续采购清单应根据未完成计划和实际库存正确计算缺口。
