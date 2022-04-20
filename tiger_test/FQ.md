

#### 汇总问题

1, 不支持创建索引，影响不大。

{"level":"info","msg":"[output-mysql] ignore unsupported ddl: create index ind_a on a(nihao)","time":"2022-04-13T16:08:18+08:00"}


"level":"info","msg":"QueryEvent: database: nihao, sql: create index ind_a on a(nihao), position: {mysql-bin.000002 116731300 e1098580-b25f-11ec-bf63-002dac1e4144:1-205493}","time":"2022-04-13T16:08:18+08:00"}
{"level":"info","msg":"[output-mysql] ignore unsupported ddl: create index ind_a on a(nihao)","time":"2022-04-13T16:08:18+08:00"}
{"level":"info","msg":"[binlogTailer] ddl done with gtid: e1098580-b25f-11ec-bf63-002dac1e4144:1-205493, stmt: create index ind_a on a(nihao)","time":"2022-04-13T16:08:18+08:00"}
{"level":"info","msg":"ignore internal ddl: /*gravityDDL*//*69_2_68_config*/alter table `nihao`.`a` add column `nihao3` varchar(20)","time":"2022-04-13T16:14:14+08:00"}


2，进程自动退出


{"level":"info","msg":"[binlogTailer] quit for context canceled, currentBinlogGTIDSet: cb588668-58e7-11ec-9913-002dac1e4b15:1-71,66ed7226-58e8-11ec-b194-002dac1e4b16:1-298965, currentBinlogFileName: mysql-bin.000006, currentBinlogFilePosition: 185188198","time":"2022-04-13T19:54:51+08:00"}
[2022/04/13 19:54:51] [info] binlogsyncer.go:172 syncer is closing...
[2022/04/13 19:54:51] [error] binlogsyncer.go:639 connection was bad
[2022/04/13 19:54:51] [error] binlogstreamer.go:77 close sync with err: Sync was closed
