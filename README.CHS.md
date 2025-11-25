# 波场网络数据解析 TronETL

[![release](https://github.com/c0mm4nd/tronetl/actions/workflows/release.yml/badge.svg?branch=dev)](https://github.com/c0mm4nd/tronetl/actions/workflows/release.yml)

TRONETL 是一个用于将波场（TRON）网络中的区块链数据解析为 CSV 格式文件的命令行工具。

[English Version](./README.md)

## 准备工作

1.  运行一个波场全节点。您可以使用 [TronDeploy](https://github.com/c0mm4nd/trondeploy) 来搭建。
2.  确保全节点正在运行并且可以访问。
3.  安装 `tronetl`。

## 安装

### 使用 Docker (推荐)

```bash
git clone https://github.com/c0mm4nd/tronetl && cd tronetl
docker build -t tronetl .
docker run -it tronetl -h
```

### 使用 Go

```bash
git clone https://github.com/c0mm4nd/tronetl && cd tronetl
go install .
tronetl -h
```

## 用法

```
tronetl 是一个用于将波场（TRON）网络中的区块链数据解析为 CSV 格式文件的命令行工具。

用法:
  tronetl [command]

可用命令:
  completion                     为指定的 shell 生成自动补全脚本
  export_blocks_and_transactions 导出区块、交易和 TRC10 转账
  export_token_transfers         导出 TRC20 代币转账
  export_address_details         导出地址详情
  help                           显示任何命令的帮助信息
  server                         运行一个 Web 服务器以提供导出服务

标志:
  -h, --help   tronetl 的帮助信息

使用 "tronetl [command] --help" 获取更多关于命令的信息。
```

ETL 结果的模式（schema）请参考[此文档](./SCHEMA.CHS.md)。

### export_blocks_and_transactions

导出区块、交易和 TRC10 转账到 CSV 文件。

```
用法:
  tronetl export_blocks_and_transactions [flags]

标志:
      --blocks-output string         区块输出的 CSV 文件路径，使用 - 表示不输出 (默认 "blocks.csv")
      --end-block uint               结束区块高度
      --end-timestamp uint           结束区块的时间戳 (UTC)
  -h, --help                         export_blocks_and_transactions 的帮助信息
      --provider-uri string          TRON Full Node 的 URI (不带端口) (默认 "http://localhost")
      --start-block uint             起始区块高度
      --start-timestamp uint         起始区块的时间戳 (UTC)
      --transactions-output string   交易输出的 CSV 文件路径，使用 - 表示不输出 (默认 "transactions.csv")
      --trc10-output string          TRC10 转账输出的 CSV 文件路径，使用 - 表示不输出 (默认 "trc10.csv")
      --workers uint                 并行执行的 worker 数量
```

### export_token_transfers

导出 TRC20 代币转账、日志、内部交易和收据到 CSV 文件。

```
用法:
  tronetl export_token_transfers [flags]

标志:
      --contracts stringArray       仅输出指定合约的转账
      --end-block uint              结束区块高度
      --end-timestamp uint          结束区块的时间戳 (UTC)
  -h, --help                        export_token_transfers 的帮助信息
      --internal-tx-output string   内部交易输出的 CSV 文件路径，使用 - 表示不输出 (默认 "internal_transactions.csv")
      --logs-output string          日志输出的 CSV 文件路径，使用 - 表示不输出 (默认 "logs.csv")
      --provider-uri string         TRON Full Node 的 URI (不带端口) (默认 "http://localhost")
      --receipts-output string      收据输出的 CSV 文件路径，使用 - 表示不输出 (默认 "receipts.csv")
      --start-block uint            起始区块高度
      --start-timestamp uint        起始区块的时间戳 (UTC)
      --transfers-output string     代币转账输出的 CSV 文件路径，使用 - 表示不输出 (默认 "token_transfers.csv")
      --workers uint                并行执行的 worker 数量
```

### export_address_details

导出地址详情，包括账户信息、合约详情和代币信息。

```
用法:
  tronetl export_address_details [flags]

标志:
      --accounts-output string    账户信息输出的 CSV 文件路径，使用 - 表示不输出 (默认 "accounts.csv")
      --addrs stringArray         需要加载详情的地址列表
      --addrs-source string       包含地址列表的 CSV 文件路径 (默认 "-")
      --contracts-output string   合约账户详情输出的 CSV 文件路径，使用 - 表示不输出 (默认 "contract.csv")
  -h, --help                      export_address_details 的帮助信息
      --provider-uri string       TRON Full Node 的 URI (不带端口) (默认 "http://localhost")
      --tokens-output string      代币合约详情输出的 CSV 文件路径，使用 - 表示不输出 (默认 "tokens.csv")
```

### server

运行一个 Web 服务器，通过 REST API 提供导出功能。

服务器监听端口 `54173`。

**API 端点:**

*   `/export_blocks_and_transactions`:  导出区块和交易。
    *   查询参数: `start-block`, `end-block`, `start-timestamp`, `end-timestamp`
*   `/export_token_transfers`: 导出代币转账。
    *   查询参数: `start-block`, `end-block`, `start-timestamp`, `end-timestamp`, `contracts` (可重复)

输出为包含相应 CSV 文件的 zip 压缩包。