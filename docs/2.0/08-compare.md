目前有许多的数据复制组件，比如 

- [canal/otter](https://github.com/alibaba/canal)
- [dtle](https://github.com/actiontech/dtle)
- [syncer](https://pingcap.com/docs-cn/tools/syncer/)

摩拜 DRC 与它们最大的不同就在于对多种数据源、多种输出的支持，以及对 Kubernetes 的原生支持。


|   | 多种数据源  | 多种数据输出  | MySQL 双向复制  | 数据清洗 | 集群版本 |多种语言二次开发|
|---|---|---|---|---|---|---|
| canal/otter  |  😩 | 🙂 | 😩 | 😩 | 😩 | 😩 |
|  dtle | 😩  | 😩  |  ✅ | 😩 | ✅  | 🙂 |
|  datapipeline | ✅ ✅  | ✅ ✅   |  😩 | ✅  | ✅  | ✅ |
|  DRC | ✅ |  ✅  | ✅  | ✅ |✅  |✅|
