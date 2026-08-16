# Recipe Planner Service

Recipe Planner Service 是一个本地运行的菜谱与周计划服务。用户可以录入菜谱和食材库存，按标签生成每日用餐计划，标记已完成的餐食，并查看尚缺食材的采购清单。

服务提供以下相互关联的能力：

- 创建和浏览菜谱；每份菜谱包含标签、制作步骤及每餐所需食材份数。
- 批量补充本地食材库存，并在烹饪时原子地扣减对应份数。
- 根据起始日期、天数和可选标签生成计划，随后生成采购清单或将某日餐食标记为已烹饪。

项目不依赖数据库或外部服务，运行期间的数据保存在内存中。HTTP 服务默认监听 `:8080`，可通过 `RECIPE_PLANNER_ADDR` 环境变量修改监听地址。

## API

- `GET /health`：返回服务健康状态。
- `POST /recipes`、`GET /recipes`：创建和查询菜谱。
- `POST /pantry/restock`：补充食材库存。
- `POST /plans/generate`、`GET /plans/{id}`：生成和查看用餐计划。
- `POST /plans/{id}/meals/{date}/cook`：完成某天的一餐并扣减库存。
- `GET /shopping?plan_id={id}`：计算计划中未完成餐食缺少的食材。

例如，服务启动后可检查健康状态：

```bash
curl http://127.0.0.1:8080/health
```
