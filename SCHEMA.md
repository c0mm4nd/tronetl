
# Data Structure

In the design of the data structure, the compatibility with the CSV format output by the tronetl project is first guaranteed.
And based on this, special parameters in the Tron network are appendec to the tail.

## blocks.csv

The block structure is as follows:

| Column            | Type                                         |
| ----------------- | -------------------------------------------- |
| number            | uint64                                       |
| hash              | hex_string                                   |
| parent_hash       | hex_string                                   |
| nonce             | always `""`                                  |
| sha3_uncles       | always `""`                                  |
| logs_bloom        | always `"0"*512`                             |
| transactions_root | hex_string                                   |
| state_root        | always `""`                                  |
| receipts_root     | always `""`                                  |
| miner             | address, actually stands for witness_address |
| difficulty        | always `""`                                  |
| total_difficulty  | always `""`                                  |
| size              | uint64                                       |
| extra_data        | always `""`                                  |
| gas_limit         | bigint                                       |
| gas_used          | bigint                                       |
| timestamp         | uint64, the unit is second                   |
| transaction_count | int                                          |
| base_fee_per_gas  | always `""`                                  |

Increase:

| Column            | Type                      |
| ----------------- | ------------------------- |
| witness_signature | hex_string, PoS signature |

---

## transactions.csv

The TRX transaction data structure is as follows, here it mainly follows the transaction data that comes with the rpc interface block in java-tron

| Column                   | Type                             |
| ------------------------ | -------------------------------- |
| hash                     | hex_string                       |
| nonce                    | always empty                     |
| block_hash               | hex_string                       |
| block_number             | uint64                           |
| transaction_index        | uint                             |
| from_address             | address                          |
| to_address               | address                          |
| value                    | bigint                           |
| gas                      | bigint, = Energy Consumed        |
| gas_price                | bigint, meaningless              |
| input                    | hex_string                       |
| block_timestamp          | timestamp in second              |
| max_fee_per_gas          | always empty                     |
| max_priority_fee_per_gas | always empty                     |
| transaction_type         | string, refer to the event types |

Increase:

| Column                 | Type                       |
| ---------------------- | -------------------------- |
| transaction_timestamp  | uint64, unit is second     |
| transaction_expiration | uint64, the unit is second |
| fee_limit              | bigint                     |


For details about event types, please refer to the [system-contracts chapter in the official document](https://tronprotocol.github.io/documentation-en/mechanism-algorithm/system-contracts/), and functions can be added later as needed.

---

## trc10.csv

The TRC10 transaction data structure (from `TransferAssetContract` and `TransferContract` events) is as follows:

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


---


## token_transfers.csv

The token transaction data structure is as follows:

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

## logs.csv

The log data structure is as follows:

| Column           | Type                |
| ---------------- | ------------------- |
| block_number     | uint64              |
| transaction_hash | hex_string          |
| log_index        | int                 |
| address          | address             |
| topics           | topics joint by `;` |
| data             | hex_string          |

---

## internal_transactions.csv

The txinfo's internal transaction data structure is as follows:

| Column             | Type                                      |
| ------------------ | ----------------------------------------- |
| transaction_hash   | hex_string                                |
| internal_index     | uint                                      |
| caller_address     | address                                   |
| transferTo_address | address                                   |
| call_info_index    | uint, index of the call info              |
| call_token_id      | uint, token id (empty means TRX)          |
| call_value         | int64, the amount of the transfered token |
| note               | hex_string                                |
| rejected           | bool                                      |

---

## receipts.csv

The tx receipt data structure is as follows:

| Column              | Type                                                   |
| ------------------- | ------------------------------------------------------ |
| transaction_hash    | hex_string                                             |
| transaction_index   | uint                                                   |
| block_number        | hex_string                                             |
| contract_address    | address (the called one, not newly created one as eth) |
| energy_fee          | int64                                                  |
| origin_energy_usage | int64                                                  |
| energy_usage_total  | int64                                                  |
| net_usage           | int64                                                  |
| net_fee             | int64                                                  |
| result              | string                                                 |

---

## addresses.csv

The tx receipt data structure is as follows:

| Column  | Type    |
| ------- | ------- |
| address | address |

---


## Note

For all `address` types, it is parsed into T-addr format (that is, a base58 string starting with T).
The `boolean` type indicates a value: `true` or `false` (all lowercase).