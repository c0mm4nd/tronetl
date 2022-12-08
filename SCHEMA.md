
# Data Structure

In the design of the data structure, the compatibility with the CSV format output by the tronetl project is first guaranteed.
And based on this, special parameters in the Tron network are added at the end.

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
| timestamp         | uint64, the unit is ms                       |
| transaction_count | int                                          |
| base_fee_per_gas  | always `""`                                  |

Increase:

| Column            | Type                      |
| ----------------- | ------------------------- |
| witness_signature | hex_string, PoS signature |

---

## transactions.csv

The TRX transaction data structure is as follows, here it mainly follows the transaction data that comes with the rpc interface block in java-tron

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

Increase:

| Column                 | Type                                    |
| ---------------------- | --------------------------------------- |
| transaction_timestamp  | uint64, unit is ms                      |
| transaction_expiration | uint64, the unit is ms                  |
| fee_limit              | bigint                                  |
| contract_calls         | int, the number of contract call events |


---

## trc10.csv

The TRC10 transaction data structure is as follows, which can also be regarded as traces in Ethereum:

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


For details about event types other than trc10 transactions, please refer to the [system-contracts chapter in the official document](https://tronprotocol.github.io/documentation-en/mechanism-algorithm/system-contracts/), and functions can be added later as needed.

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

## Note

For all `address` types, it is parsed into T-addr format (that is, a base58 string starting with T).
The `boolean` type indicates a value: `true` or `false` (all lowercase).