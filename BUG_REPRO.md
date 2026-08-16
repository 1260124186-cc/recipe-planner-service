# 修复前故障复现（Docker）

## 项目与标准命令
该项目提供菜谱创建、库存管理、用餐计划和采购清单服务。在仓库根目录可使用以下标准命令：

```bash
go build ./...
go run ./cmd/server
go test ./...
```

本基线状态下，`go build ./...` 可以成功；`go test ./...` 会触发下述故障。

## 环境构建与编译
已分别在 linux/arm64 和 linux/amd64 完成镜像构建、容器内编译和健康检查：

```bash
docker build --platform linux/arm64 -f benzhi.Dockerfile -t recipe-planner-service:benzhi .
docker run --rm --platform linux/arm64 --entrypoint go recipe-planner-service:benzhi build ./...
docker build --platform linux/amd64 -f benzhi.Dockerfile -t recipe-planner-service:benzhi .
docker run --rm --platform linux/amd64 --entrypoint go recipe-planner-service:benzhi build ./...
```

两个平台的镜像构建和容器内 `go build ./...` 均成功；目标故障在下节命令中触发。

## 故障触发步骤
在仓库根目录执行：

```bash
docker build --platform linux/arm64 -f benzhi.Dockerfile -t recipe-planner-service:benzhi .
docker run --rm --platform linux/arm64 --entrypoint go recipe-planner-service:benzhi test ./...
```

## 实际错误输出

```text
?   	github.com/1260124186-cc/recipe-planner-service/cmd/server	[no test files]
?   	github.com/1260124186-cc/recipe-planner-service/internal/domain	[no test files]
--- FAIL: TestStoredRecipeDoesNotChangeWhenCallerReusesSlices (0.00s)
    workflow_test.go:119: stored recipe changed with caller slices: domain.Recipe{ID:"shared-slices", Name:"Shared Slices", Tags:[]string{"slow"}, Steps:[]string{"Discard ingredients"}, Ingredients:[]domain.IngredientNeed{domain.IngredientNeed{Name:"tomato", Portions:99}}}
FAIL
FAIL	github.com/1260124186-cc/recipe-planner-service/internal/service	0.003s
?   	github.com/1260124186-cc/recipe-planner-service/internal/store	[no test files]
ok  	github.com/1260124186-cc/recipe-planner-service/internal/transport	0.003s
FAIL

命令退出状态：1
```

## 期望行为
使用相同输入创建菜谱后，后续对最初提交数据的本地处理不应改变已保存菜谱的标签、制作步骤或食材份数。查询、排期和采购清单应继续使用创建成功时的菜谱内容。
