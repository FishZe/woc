#### 学校的任务罢了......

#### 因为我不会java，所以我只能写golang了

### 接口介绍

### 1. 登录

> POST /login

参数:

| 参数名      | 类型     | 说明        |
|----------|--------|-----------|
| name     | string | 用户名       |
| password | string | 密码        |

示例
```json
{
    "name": "24岁事学生",
    "password": "114514"
}
```

返回值:

| 参数名  | 类型  | 说明   |
|------|-----|------|
| code | int | 状态码  |
| data | obj | 相关信息 |

其中, `data`字段含义如下:

| 参数名  | 类型   | 说明     |
|------|------|--------|
| msg  | str  | 信息     |

当登陆成功时, `code`为`200`, `msg`为`login success`
登陆失败时, `msg`为登陆失败原因

### 2. 管理员相关接口

> POST /admin

**以下接口需要`role`为1时才能调用**

#### 否则会被中间件拦截哦

#### 2.1 新建用户

> POST /admin/new

参数:

| 参数名      | 类型     | 说明        |
|----------|--------|-----------|
| name     | string | 用户名       |
| password | string | 密码        |
| email    | string | 邮箱        |
| role     | string | 角色        |

示例
```json
{
    "name": "24岁不是学生",
    "password": "114514",
    "email": "114@514.com",
    "role": 0
}
```

返回值:

| 参数名  | 类型  | 说明   |
|------|-----|------|
| code | int | 状态码  |
| data | obj | 相关信息 |

其中, `data`字段含义如下:

| 参数名  | 类型   | 说明     |
|------|------|--------|
| msg  | str  | 信息     |

当创建成功时, `code`为`200`, `msg`为`insert user success`
登陆失败时, `msg`为登陆失败原因

#### 2.2 删除用户

> POST /admin/delete

参数:
| 参数名      | 类型     | 说明        |
|----------|--------|-----------|
| id      | int    | 用户id      |

示例
```json
{
    "id": 1
}
```
返回值

| 参数名  | 类型  | 说明   |
|------|-----|------|
| code | int | 状态码  |
| data | obj | 相关信息 |

其中, `data`字段含义如下:

| 参数名  | 类型   | 说明     |
|------|------|--------|
| msg  | str  | 信息     |

当删除成功时, `code`为`200`, `msg`为`delete user success`
登陆失败时, `msg`为登陆失败原因

#### 2.3 修改用户信息

> POST /admin/update

参数:
| 参数名      | 类型     | 说明        |
|----------|--------|-----------|
| id      | int    | 用户id      |
| name     | string | 用户名       |
| password | string | 密码        |
| email    | string | 邮箱        |
| role     | string | 角色        |

会修改`id`为传入`id`用户的信息

示例:
```json
{
    "id": 2,
    "name": "24岁是不是学生啊",
    "password": "1145144",
    "email": "514@144.com",
    "role": -1
}
```

返回值:

| 参数名  | 类型  | 说明   |
|------|-----|------|
| code | int | 状态码  |
| data | obj | 相关信息 |

其中, `data`字段含义如下:

| 参数名  | 类型   | 说明     |
|------|------|--------|
| msg  | str  | 信息     |

当创建成功时, `code`为`200`, `msg`为`update user success`
登陆失败时, `msg`为登陆失败原因

#### 2.4 搜索用户信息

> POST /admin/search

参数:
| 参数名      | 类型     | 说明        |
|----------|--------|-----------|
| id      | int    | 用户id      |
| name     | string | 用户名       |
| password | string | 密码        |
| email    | string | 邮箱        |
| role     | string | 角色        |

值不为空时, 会用此键进行搜索, 为空时不作为关键词

特殊的, `role`为`-2`时不进行搜索

示例:
```json
{
    "id": 0,
    "name": "",
    "password": "",
    "email": "514@144.com",
    "role": -2
}
```

返回值:

| 参数名  | 类型  | 说明   |
|------|-----|------|
| code | int | 状态码  |
| data | obj | 相关信息 |

当`code`为`200`时, `data`字段为数组, 其中每个元素含义如下:

| 参数名       | 类型      | 说明         |
|-----------|---------|------------|
| id        | int     | 用户id       |
| name      | string  | 用户名        |
| password  | string  | 密码         |
| email     | string  | 邮箱         |
| role      | string  | 角色         |

当`code`不为`200`时, `data`字段含义如下:

| 参数名  | 类型   | 说明     |
|------|------|--------|
| msg  | str  | 查询失败原因 |


#### 2.5 批量获取用户

> POST /admin/get

参数:

| 参数名      | 类型        | 说明     |
|----------|-----------|--------|
| from_id  | int       | 起始id   |
| sum      | int       | 数量     |

示例:
```json
{
    "from_id": 1,
    "sum": 10
}
```

返回值:

| 参数名  | 类型  | 说明   |
|------|-----|------|
| code | int | 状态码  |
| data | obj | 相关信息 |

当`code`为`200`时, `data`字段为数组, 其中每个元素含义如下:

| 参数名       | 类型      | 说明         |
|-----------|---------|------------|
| id        | int     | 用户id       |
| name      | string  | 用户名        |
| password  | string  | 密码         |
| email     | string  | 邮箱         |
| role      | string  | 角色         |

当`code`不为`200`时, `data`字段含义如下:

| 参数名  | 类型   | 说明     |
|------|------|--------|
| msg  | str  | 查询失败原因 |
