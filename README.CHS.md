# 波场网络数据解析 TronETL

TRONETL 是 TRON Protocol 的 ETL 助手

## 编程语言

TronETL主程序使用golang语言，在对复杂嵌套的波场原始数据解析任务时，较EthereumETL的Python语言有极大的性能提升，且golang作为静态语言对数据质量有严格保障

此外有部分python脚本作为主程序功能上的补充。

## 准备工作

1. 使用 [TronDeploy](https://git.ngx.fi/c0mm4nd/trondeploy) 运行全节点！

2.确保下载了输出目录并启动了节点

3.安装tronetl
```bash
git clone https://git.ngx.fi/c0mm4nd/tronetl && cd tronetl

# 如果使用 docker（推荐）
docker build -t tronetl .
docker run -it tronetl -h

# 否则使用最新的 golang
go install .
tronetl -h
```

## 用法

```bash
tronetl 是一个 CLI 工具，用于将区块链数据从 tron 网络解析为 CSV 格式文件

用法：
   tronetl [命令]

可用命令：
   completion 为指定的 shell 生成自动完成脚本
   export_blocks_and_transactions 导出区块，以及区块的 trx 和 trc10 交易
   export_token_transfers 导出智能合约代币的转账
   help 有关任何命令的帮助
   服务器运行一个服务器来处理导出任务

标志：
   -h, --help 帮助 tronetl

使用“tronetl [command] --help”获取有关命令的更多信息。
```

ETL 结果的架构写在[此文档中](./SCHEMA.md)

### export_blocks_and_transactions

```bash
导出区块，包含区块的 trx 和 trc10 交易

用法：
   tronetl export_blocks_and_transactions [flags]

Flag：
       --blocks-output string 块输出的 CSV 文件，使用 - 省略（默认“blocks.csv”）
       --end-block uint 结束块号
       --end-timestamp uint 结束块的时间戳（UTC）
   -h, --help export_blocks_and_transactions 帮助
       --provider-uri 字符串 tron 全节点的基本 uri（无端口）（默认“http://localhost”）
       --start-block uint 起始块号
       --start-timestamp uint 起始块的时间戳（UTC）
       --transactions-output string 交易输出的CSV文件，使用-省略（默认“transactions.csv”）
       --trc10-output string trc10 输出的 CSV 文件，使用 - 省略（默认“trc10.csv”）
```

### export_token_transfers

```bash
导出智能合约代币的转账

用法：
   tronetl export_token_transfers [flags]

Flag：
       --contracts stringArray 只输出选定合约的转账
       --end-block uint 结束块号
       --end-timestamp uint 结束块的时间戳（UTC）
   -h, --help export_token_transfers 帮助
       --output string 用于令牌传输输出的 CSV 文件，使用 - 省略（默认“token_transfer.csv”）
       --provider-uri 字符串 tron 全节点的基本 uri（无端口）（默认“http://localhost”）
       --start-block uint 起始块号
       --start-timestamp uint 起始块的时间戳（UTC）
```

### server

启动一个REST服务器，从http请求得到解析结果



