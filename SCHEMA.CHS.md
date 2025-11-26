# 数据结构

在数据结构的设计上，首先保证了与 tronetl 项目输出的 CSV 格式的兼容性。
并在此基础上，将 Tron 网络中的特殊参数附加到尾部。

## blocks.csv

区块结构如下：

| 列名 | 类型 |
| --- | --- |
| number | uint64 |
| hash | hex_string |
| parent_hash | hex_string |
| nonce | 始终为 `""` |
| sha3_uncles | 始终为 `""` |
| logs_bloom | 始终为 `"0"*512` |
| transaction_root | hex_string |
| state_root | 始终为 `""` |
| receipts_root | 始终为 `""` |
| miner | address，实际上是 witness_address |
| difficulty | 始终为 `""` |
| total_difficulty | 始终为 `""` |
| size | uint64 |
| extra_data | 始终为 `""` |
| gas_limit | bigint |
| gas_used | bigint |
| timestamp | uint64，单位为秒 |
| transaction_count | int |
| base_fee_per_gas | 始终为 `""` |

增加：

| 列名 | 类型 |
| --- | --- |
| witness_signature | hex_string, PoS 签名 |

---

## transactions.csv

TRX 交易数据结构如下，这里主要遵循 java-tron 中 rpc 接口块附带的交易数据

| 列名 | 类型 |
| --- | --- |
| hash | hex_string |
| nonce | 始终为空 |
| block_hash | hex_string |
| block_number | uint64 |
| transaction_index | uint |
| from_address | address |
| to_address | address |
| value | bigint |
| gas | bigint, = 消耗的能量 |
| gas_price | bigint, 无意义 |
| input | hex_string |
| data | string，`raw_data.data`（交易备注/附加数据，原样输出） |
| block_timestamp | 时间戳（秒） |
| max_fee_per_gas | 始终为空 |
| max_priority_fee_per_gas | 始终为空 |
| transaction_type | 字符串，请参阅事件类型 |
| status | 字符串，可以是 SUCCESS 或 REVERT |

增加：

| 列名 | 类型 |
| --- | --- |
| transaction_timestamp | int64, 单位为秒 |
| transaction_expiration | int64, 单位为秒 |
| fee_limit | bigint |


有关事件类型的详细信息，请参阅[官方文档中的系统合约章节](https://tronprotocol.github.io/documentation-en/mechanism-algorithm/system-contracts/)，后续可根据需要添加功能。

---

## trc10.csv

TRC10 交易数据结构（来自 `TransferAssetContract` 和 `TransferContract` 事件）如下：

| 列名 | 类型 |
| --- | --- |
| block_number | uint64 |
| block_hash | hex_string |
| transaction_hash | hex_string |
| transaction_index | int |
| contract_call_index | int |
| asset_name | string |
| from_address | address |
| to_address | address |
| value | bigint |


---


## token_transfers.csv

代币交易数据结构如下：

| 列名 | 类型 |
| --- | --- |
| block_number | uint64 |
| transaction_hash | hex_string |
| log_index | int |
| token_address | address |
| from_address | address |
| to_address | address |
| value | bigint |

---

## logs.csv

日志数据结构如下：

| 列名 | 类型 |
| --- | --- |
| block_number | uint64 |
| transaction_hash | hex_string |
| log_index | int |
| address | address |
| topics | topics 以分号连接 |
| data | hex_string |

---

## internal_transactions.csv

txinfo 的内部交易数据结构如下：

| 列名 | 类型 |
| --- | --- |
| block_number | uint64 |
| transaction_hash | hex_string |
| internal_index | uint |
| internal_hash | hex_string |
| caller_address | address |
| transferTo_address | address |
| call_info_index | uint, 调用信息的索引 |
| call_token_id | uint, 代币 id (空表示 TRX) |
| call_value | int64, 转移的代币数量 |
| note | hex_string |
| rejected | bool |

---

## receipts.csv

交易收据数据结构如下：

| 列名 | 类型 |
| --- | --- |
| transaction_hash | hex_string |
| transaction_index | uint |
| block_number | uint64 |
| contract_address | address (被调用的地址，而不是像 eth 那样新创建的地址) |
| energy_usage | int64 |
| energy_fee | int64 |
| origin_energy_usage | int64 |
| energy_usage_total | int64 |
| net_usage | int64 |
| net_fee | int64 |
| result | string |

---

## accounts.csv

账户数据结构如下：

| 列名 | 类型 |
| --- | --- |
| account_name | string |
| address | address |
| type | string |
| create_time | int64 |

---

## contracts.csv

合约数据结构如下：

| 列名 | 类型 |
| --- | --- |
| address | address |
| bytecode | string |
| function_sighashes | string |
| is_erc20 | bool |
| is_erc721 | bool |
| block_number | uint64 |
| contract_name | string |
| consume_user_resource_percent| int |
| origin_address | address |
| origin_energy_limit | int64 |

---

## tokens.csv

代币数据结构如下：

| 列名 | 类型 |
| --- | --- |
| address | address |
| symbol | string |
| name | string |
| decimals | uint64 |
| total_supply | uint64 |
| block_number | uint64 |

---

## 注意

对于所有 `address` 类型，它被解析为 T-addr 格式（即以 T 开头的 base58 字符串）。
`boolean` 类型表示一个值：`true` 或 `false`（全小写）。
