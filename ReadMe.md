## Ceph Quota Setter API 说明
Ceph Quota Setter 一共提供两个基于HTTP的API接口:

### 【/set 接口】

/set接口使用POST请求，用来设置目标CephFS目录的配额值。

#### 请求Bady样例：
``` json
{
    "target": "/mnt/cephfs-ht-test2/temp/",
    "max_bytes": 20000000,
    "max_files": 10000
}
```
请求的Bady为Json格式，其中target为字符串类型的目标目录；max_bytes和max_files都为int类型的配额数值。
如上请求表示将target所指目录的属性ceph.quota.max_bytes和ceph.quota.max_files设置成max_bytes和max_files中所设置的值。
其中max_bytes和max_files至少设置一个，若都不设置则返回错误。

#### 结果解析：
```json
{
    "code":20,
    "success":true,
    "description":"Set quotas success."
}
```
结果也以Json格式返回，其中code为int类型返回码；success为布尔变量，成功为true、失败为false；description为字符串类型的结果描述。下表是返回码及描述意义：

|code|success| description|可能原因|
|:--:| :--: | :--: |:--: |
| 10 |false |Can not decode data.   |Json字段类型不匹配|
| 11 |false |Target path is not exist. |目标目录挂载参数不一致|
| 12 |false |Missing parameters or Invalid parameters. |没有配额参数或配置值都为负|
| 13 |false |(返回系统标准错误) |设置max_files出错|
| 14 |false |(返回系统标准错误) |设置max_bytes出错|
| 20 |true  |Set quotas success. | 配额设置成功|


### 【/get 接口】
/get接口使用POST请求，用来查询目标CephFS目录的配额值。

#### 请求Bady样例：
``` json
{
    "target": "/mnt/cephfs-ht-test2/temp/"
}
```
请求的Bady为Json格式，其中target为字符串类型的目标目录。

#### 结果解析：
```json
{
    "code":21,
    "success":true,
    "description":"Get quotas success.",
    "max_bytes":20000000,
    "max_files":10000
}
```
结果也以Json格式返回，其中code为int类型返回码；success为布尔变量，成功为true、失败为false；description为字符串类型的结果描述；max_bytes和max_files都为int类型的配额数值。下表是返回码及描述意义：

|code|success| description|可能原因|
|:--:| :--: | :--: |:--: |
| 10 |false |Can not decode data.   |Json字段类型不匹配|
| 11 |false |Target path is not exist. |目标目录挂载参数不一致|
| 15 |false |(返回系统标准错误) |获取max_files出错|
| 16 |false |(返回系统标准错误) |获取max_bytes出错|
| 21 |true  |Get quotas success. | 获取配额成功|
