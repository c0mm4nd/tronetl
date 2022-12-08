
# 数据结构

在数据结构的设计上首先保证与tronetl项目输出CSV格式的兼容性，在此基础上结尾增加Tron网络中特殊参数。

## blocks.csv

区块结构如下：

| Column            | Type                             |
| ----------------- | -------------------------------- |
| number            | uint64                           |
| hash              | hex_string                       |
| parent_hash       | hex_string                       |
| nonce             | 始终为`""`                       |
| sha3_uncles       | 始终为`""`                       |
| logs_bloom        | 始终为`"0"*512`                  |
| transactions_root | hex_string                       |
| state_root        | 始终为`""`                       |
| receipts_root     | 始终为`""`                       |
| miner             | address，实际代表witness_address |
| difficulty        | 始终为`""`                       |
| total_difficulty  | 始终为`""`                       |
| size              | uint64                           |
| extra_data        | 始终为`""`                       |
| gas_limit         | bigint                           |
| gas_used          | bigint                           |
| timestamp         | uint64，单位为ms                 |
| transaction_count | int                              |
| base_fee_per_gas  | 始终为`""`                       |

增加：

| Column            | Type                |
| ----------------- | ------------------- |
| witness_signature | hex_string，PoS签名 |

---

## transactions.csv

常规交易数据结构如下，此处主要遵循java-tron中rpc接口区块自带交易数据

| Column                   | Type       |
| ------------------------ | ---------- |
| hash                     | hex_string |
| nonce                    | bigint     |
| block_hash               | hex_string |
| block_number             | bigint     |
| transaction_index        | bigint     |
| from_address             | address    |
| to_address               | address    |
| value                    | numeric    |
| gas                      | bigint     |
| gas_price                | bigint     |
| input                    | hex_string |
| block_timestamp          | bigint     |
| max_fee_per_gas          | bigint     |
| max_priority_fee_per_gas | bigint     |
| transaction_type         | bigint     |

增加：

| Column                 | Type                |
| ---------------------- | ------------------- |
| transaction_timestamp  | uint64，单位为ms    |
| transaction_expiration | uint64，单位为ms    |
| fee_limit              | bigint              |
| contract_calls         | int, 合约调用事件数 |


---

## trc10.csv

TRC10交易数据结构如下，也可以视为以太坊中的traces：

| Column              | Type       |
| ------------------- | ---------- |
| block_number        | uint64     |
| block_hash          | hex_string |
| transaction_hash    | hex_string |
| transaction_index   | int        |
| contract_call_index | int        |
| asset_name          | string     |
| from_address        | address    |
| to_address          | address    |
| value               | bigint     |


trc10交易以外事件类型详见[官方文档中system-contracts章节](https://tronprotocol.github.io/documentation-en/mechanism-algorithm/system-contracts/)，后续根据需要可增加功能。

---


## token_transfers.csv

代币交易数据结构如下：

| Column           | Type       |
| ---------------- | ---------- |
| block_number     | uint64     |
| transaction_hash | hex_string |
| log_index        | int        |
| token_address    | address    |
| from_address     | address    |
| to_address       | address    |
| value            | bigint     |

---

## 备注

对所有 `address` 类型，都解析为T-addr格式（即T开头的base58字符串）。
`boolean` 类型表示值为: `true` 或 `false` （全小写）。
