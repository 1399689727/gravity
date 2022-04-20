## 架构

master:

172.30.75.250:3307
172.30.75.21:3307

auto_increment_offset=1
auto_increment_increment=7

master:

172.30.75.22:3307

auto_increment_offset=2
auto_increment_increment=7


./gravity -config 21_2_22_config.toml
./gravity -config 22_2_21_config.toml


## MySQL 配置

```
[mysqld]
server_id=4
log_bin=mysql-bin
enforce-gtid-consistency=ON
gtid-mode=ON
binlog_format=ROW

```



#### 开启gtid

https://cloud.tencent.com/developer/article/1418501


在实践online升级之前，我们需要了解MySQL 5.7版本的GTID_MODE 的含义:

OFF            :不产生GTID,Slave只接受不带GTID的事务
OFF_PERMISSIVE :不产生GTID,Slave即接受不带GTID的事务,也接受带GTID的事务
ON_PERMISSIVE  :产生GTID,Slave即接受不带GTID的事务,也接受带GTID的事务
ON             :产生GTID,Slave只能接受带GTID的事务。


2.1 在主从复制结构中所有的实例中执行

set global ENFORCE_GTID_CONSISTENCY = WARN;

在正常运行的业务系统数据库中，设置ENFORCEGTIDCONSISTENCY为WARN，目的是观察err log是否有不满足要求的sql出现。如果有发现任何warning，需要通知应用进行调整相关sql，直到不出现warning为止。GTID 使用限制如下:

1.不支持非事务引擎。
2.不支持create table ... select 语句(在主库执行时直接报错)。
3.不允许一个SQL同时更新一个事务引擎和非事务引擎的表。
4.不支持create temporary table和drop temporary语句。
如果没有任何warning 出现，则在所有实例上执行:

set global ENFORCE_GTID_CONSISTENCY = ON;


2.2 在主从复制结构中所有实例中执行:

set global GTID_MODE = OFF_PERMISSIVE;

让主库不产生GTID,Slave实例即接受不带GTID的事务,也接受带GTID的事务。确保一定要在所有实例中执行完该命令之后再执行接下来的步骤。


2.3 在主从复制结构中所有实例中执行:

set global GTID_MODE = ON_PERMISSIVE;

主库开始产生GTID,Slave即接受不带GTID的事务,也接受带GTID的事务。


2.4 在主从复制结构中所有的实例中执行:

在各个实例节点上执行如下命令检查匿名事务是否消耗完毕，最好多检查几次，以便确认该参数的值是0.

[RW][TEST:3316]>SHOW STATUS LIKE 'ONGOING_ANONYMOUS_TRANSACTION_COUNT';
+-------------------------------------+-------+
| Variable_name                       | Value |
+-------------------------------------+-------+
| Ongoing_anonymous_transaction_count | 0     |
+-------------------------------------+-------+
1 row in set (0.00 sec)
如果在从库上检查只需要一次满足为0 即可。


2.5 确保第四步之前的binlog全部为应用。

确保操作之前的所有binlog都已经被其他服务器应用了，因为匿名的GTID必须确保已经复制应用成功，才可以进行下一步操作。如何检查呢？ 其实最简单的方式是在从库库执行show slave status检查应用位点的情况。
如果追上了，则可以继续。否则需要等待从库应用完binlog之后在进行下一步。

2.5 在主从复制结构中所有的实例中执行:
set global GTID_MODE = ON;

该参数的功能是让系统产生GTID ,Slave只能接受带GTID的事务。

2.6 在从库上执行:
设置slave 复制中MASTER_AUTO_POSITION=1。

stop slave;
CHANGE MASTER TO MASTER_AUTO_POSITION = 1;
start slave;

至此，将基于位点的复制关系升级为GTID模式。结束了吗？还没呢，记得修改my.cnf 添加

enforce-gtid-consistency=ON
gtid-mode=ON


#### 创建用户

CREATE USER _gravity IDENTIFIED BY '123456sl';
GRANT SELECT, RELOAD, LOCK TABLES, REPLICATION SLAVE, REPLICATION CLIENT, CREATE, INSERT, UPDATE, DELETE ON *.* TO '_gravity'@'%';
GRANT ALL PRIVILEGES ON _gravity.* TO '_gravity'@'%';



## 配置双向链路


./gravity -config 21_2_22_config.toml


```

# name 必填
name = "21_2_22_config"

# 内部用于保存位点、心跳等事项的库名，默认为 _gravity
internal-db-name = "_gravity"

#
# Input 插件的定义，此处定义使用 mysql
#
[input]
type = "mysql"
mode = "replication"

[input.config]
ignore-bidirectional-data = true

[input.config.source]
host = "172.30.75.21"
port = 3307
username = "_gravity"
password = "123456sl"

#
# Output 插件的定义，此处使用 mysql
#
[output]
type = "mysql"
[output.config.target]

host = "172.30.75.22"
port = 3307
username = "_gravity"
password = "123456sl"

# 当前支持 create & alter table 语句。库表名会根据路由信息调整。
[output.config]
enable-ddl = true

[output.config.sql-engine-config]
type = "mysql-replace-engine"

[output.config.sql-engine-config.config]
tag-internal-txn = true


# 路由规则的定义
[[output.config.routes]]
match-schema = "stock_user"
match-table = "*"
target-schema = "stock_user"
target-table = "*"

```

./gravity -config 22_2_21_config.toml

```


# name 必填
name = "22_2_21_config"

# 内部用于保存位点、心跳等事项的库名，默认为 _gravity
internal-db-name = "_gravity"

#
# Input 插件的定义，此处定义使用 mysql
#
[input]
type = "mysql"
mode = "replication"

[input.config]
ignore-bidirectional-data = true

[input.config.source]
host = "172.30.75.22"
port = 3307
username = "_gravity"
password = "123456sl"

#
# Output 插件的定义，此处使用 mysql
#
[output]
type = "mysql"
[output.config.target]

host = "172.30.75.21"
port = 3307
username = "_gravity"
password = "123456sl"

# 当前支持 create & alter table 语句。库表名会根据路由信息调整。
[output.config]
enable-ddl = true

[output.config.sql-engine-config]
type = "mysql-replace-engine"

[output.config.sql-engine-config.config]
tag-internal-txn = true


# 路由规则的定义
[[output.config.routes]]
match-schema = "stock_user"
match-table = "*"
target-schema = "stock_user"
target-table = "*"


```

