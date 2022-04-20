#### 编译

```
git clone https://github.com/moiot/gravity.git

export GOOS=linux  
export GOARCH=amd64   

cd gravity && make

scp bin/gravity root@172.30.65.68:/data0/scripts/
```




#### mysql

MySQL 配置要求

```
[mysqld]
server_id=4
log_bin=mysql-bin
enforce-gtid-consistency=ON
gtid-mode=ON
binlog_format=ROW

```

gravity 账户权限如下所示

```
mysql -uroot -plhzq@123qwe -S /data0/mysql_3306//data/mysql.sock


CREATE USER _gravity IDENTIFIED BY '123456sl';
GRANT SELECT, RELOAD, LOCK TABLES, REPLICATION SLAVE, REPLICATION CLIENT, CREATE, INSERT, UPDATE, DELETE ON *.* TO '_gravity'@'%';
GRANT ALL PRIVILEGES ON _gravity.* TO '_gravity'@'%';


create database nihao;

use nihao;
create table a(id int auto_increment primary key, name varchar(20));

```



#### gravity 自动创建表，全量、增量数据。

配置ddl

[output.config]
enable-ddl = true



Gravity 在写入目标端 MySQL 的时会打上双向同步的内部标识（通过封装内部表事务的方式），
在源端配置好 ignore-bidirectional-data 就可以忽略 Gravity 内部的写流量。

./gravity -config 68_2_69_ddl_config.toml


```

# name 必填
name = "68_2_69_ddl_config"

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
host = "127.0.0.1"
username = "dba"
password = "a0d196ebf12968f065cdgzs103"
port = 3306

#
# Output 插件的定义，此处使用 mysql
#
[output]
type = "mysql"
[output.config.target]
host = "172.30.65.69"
username = "dba"
password = "a0d196ebf12968f065cdgzs103"
port = 3306

# 当前支持 create & alter table 语句。库表名会根据路由信息调整。
[output.config]
enable-ddl = true

[output.config.sql-engine-config]
type = "mysql-replace-engine"

[output.config.sql-engine-config.config]
tag-internal-txn = true


# 路由规则的定义
[[output.config.routes]]
match-schema = "nihao"
match-table = "*"
target-schema = "nihao2"
target-table = "*"


```


####  gravity 双向同步


配置2条链路，每条链路需要配置如下规则。

[input.config]
ignore-bidirectional-data = true
[output.config.sql-engine-config.config]
tag-internal-txn = true


./gravity -config 68_2_69_config.toml


```

# name 必填
name = "68_2_69_config"

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
host = "127.0.0.1"
username = "dba"
password = "a0d196ebf12968f065cdgzs103"
port = 3306

#
# Output 插件的定义，此处使用 mysql
#
[output]
type = "mysql"
[output.config.target]
host = "172.30.65.69"
username = "dba"
password = "a0d196ebf12968f065cdgzs103"
port = 3306

# 当前支持 create & alter table 语句。库表名会根据路由信息调整。
[output.config]
enable-ddl = true

[output.config.sql-engine-config]
type = "mysql-replace-engine"

[output.config.sql-engine-config.config]
tag-internal-txn = true


# 路由规则的定义
[[output.config.routes]]
match-schema = "nihao"
match-table = "*"
target-schema = "nihao2"
target-table = "*"

```

./gravity -config 69_2_68_config.toml

```

# name 必填
name = "69_2_68_config"

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
host = "172.30.65.69"
username = "dba"
password = "a0d196ebf12968f065cdgzs103"
port = 3306

#
# Output 插件的定义，此处使用 mysql
#
[output]
type = "mysql"

[output.config.target]
host = "172.30.65.68"
username = "dba"
password = "a0d196ebf12968f065cdgzs103"
port = 3306

# 当前支持 create & alter table 语句。库表名会根据路由信息调整。
[output.config]
enable-ddl = true

[output.config.sql-engine-config]
type = "mysql-replace-engine"

[output.config.sql-engine-config.config]
tag-internal-txn = true


# 路由规则的定义
[[output.config.routes]]
match-schema = "nihao2"
match-table = "*"
target-schema = "nihao"
target-table = "*"

```



####  gravity 双向同步test alter table

源：

```
{"level":"info","msg":"QueryEvent: database: nihao, sql: alter table b drop ce, position: {mysql-bin.000002 %!s(uint32=1805002) e1098580-b25f-11ec-bf63-002dac1e4144:1-3169}","time":"2022-04-06T16:06:52+08:00"}
{"level":"info","msg":"[output-mysql] executed ddl: alter table `nihao2`.`b` drop column `ce`","time":"2022-04-06T16:06:52+08:00"}
{"level":"info","msg":"[binlogTailer] ddl done with gtid: e1098580-b25f-11ec-bf63-002dac1e4144:1-3169, stmt: alter table b drop ce","time":"2022-04-06T16:06:52+08:00"}
```

目标：

```
{"level":"info","msg":"ignore internal ddl: /*gravityDDL*//*68_2_69_ddl_full_increment_config*/alter table `nihao2`.`b` drop column `ce`","time":"2022-04-06T16:06:52+08:00"}
```

#### 测试环境


源
172.30.65.68:3306 
nihao
目标
172.30.65.69:3306  
nihao2


GRANT ALL PRIVILEGES ON *.* TO 'test'@'%' IDENTIFIED BY '123456test001';
flush privileges;



#### 本地运行


go run cmd/gravity/main.go -config ./tiger_test/68_2_69_config.toml





