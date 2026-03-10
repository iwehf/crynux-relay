# Deposit, Withdraw and Risk Control

This document describes the deposit and withdrawal system design, the trust model between Relay and Wallet, and the risk control strategy.

## Two-Service Architecture

The system is split across two independent services (and repositories): the Relay and the Relay Wallet.

The **Relay** is a public-facing service. It accepts deposit detection and withdrawal requests, tracks task fee balances, and serves as the coordination layer. Because it is publicly accessible, it is the more exposed attack surface.

The **Relay Wallet** holds the private keys to the system wallets that contain real funds. It runs in a restricted environment with no inbound network access. The Relay Wallet actively polls the Relay for pending work; the Relay never pushes to the Relay Wallet. This pull-only design means that even if the Relay is fully compromised, an attacker cannot directly instruct the Relay Wallet to send funds. The Relay Wallet independently validates every request before executing any transaction.

## Deposit

A deposit converts on-chain native tokens into off-chain task fee balance. There is no API to initiate a deposit. A user deposits by sending a native token transfer to a designated system address on any supported network. The Relay monitors the chain, detects the transfer, and credits the sender's task fee balance with the transferred amount. The deposit then propagates to the Relay Wallet as a task fee log entry like any other balance change.

## Withdraw

A withdrawal converts off-chain task fee balance back into on-chain native tokens. A withdrawal request charges a fee; the fee is waived for the fee collection address.

1. A node submits a signed withdrawal request to the Relay. The Relay validates the request and reserves the amount from the node's task fee balance.
2. The Relay Wallet periodically pulls pending requests from the Relay and independently validates each request against its own locally computed balance and on-chain state. If any request fails validation, the Relay Wallet halts the entire sync task and raises an alert for manual intervention -- it does not skip or reject individual requests.
3. For each accepted request, the Relay Wallet sends a blockchain transaction to transfer the amount to the node's benefit address.
4. The Relay Wallet reports the outcome back to the Relay. On success, the Relay finalizes the record. On failure, the Relay refunds the reserved amount back to the node's task fee balance.

## Task Fee Logs

Task fee logs are the mechanism that connects deposits, task rewards, and withdrawals. They are the sole input the Relay Wallet uses to build its local view of each node's balance. The Relay produces a task fee log entry for every balance-changing event. The Relay Wallet syncs these logs sequentially by ID from the Relay and applies them to its local balances. The Relay Wallet MUST NOT accept a withdrawal request until all task fee logs up to and including the request's `TaskFeeEventID` have been synced.

Each log entry has a type:

| Type | Name | Description |
|------|------|-------------|
| 0 | Task | Node reward for completing an inference task |
| 1 | Draw | DAO commission from the same task. One Task event and one Draw event are always created together per task completion |
| 2 | WithdrawalFee | Withdrawal fee credited to the fee collection address when a withdrawal is fulfilled |
| 3 | Bought | Balance deposited via on-chain token purchase |

## Trust Model

The Relay and the Relay Wallet do not trust each other's balances.

The Relay maintains its own task fee balance for each node address, updated as tasks complete and deposits arrive. When a withdrawal request arrives, the Relay checks its own balance and reserves the amount.

The Relay Wallet maintains a separate copy of each node's balance, built by independently syncing task fee logs from the Relay and applying them locally. Before executing a withdrawal, the Relay Wallet checks the request against its own locally computed balance. This means a corrupted or manipulated balance in the Relay database alone is not sufficient to cause an incorrect payout -- the Relay Wallet's independent balance must also agree.

### Record Integrity (Current: MAC, TODO: User Signature Forwarding)

The current implementation uses a MAC to protect against database tampering. The Relay computes an HMAC-SHA256 tag over the record's key fields (ID, address, benefit address, amount, network) using a secret key stored outside the database. The MAC is computed when the Relay's background job marks a record as processed, and verified when the Relay Wallet polls for records. If any covered field was modified after the MAC was computed, the verification fails and the record is excluded.

However, the MAC is a workaround for a missing design: the user's ECDSA signature is not forwarded to the Relay Wallet. When a node submits a withdrawal request, it signs the full withdrawal details (address, amount, benefit address, network) with its private key. The Relay validates this signature at the API layer, but then discards it -- the signature is not stored in the database and not included in the data served to the Relay Wallet. As a result, the Relay Wallet has no way to verify what the user actually requested.

This matters primarily for the **amount** field. If an attacker tampers with the amount in the Relay's database (for example, changing a 1 ETH withdrawal to 50 ETH), the Relay Wallet would honor the tampered amount as long as it does not exceed the Relay Wallet's locally tracked balance. The benefit address is independently verified from chain by the Relay Wallet, so tampering that field alone would be caught regardless of the MAC. But the amount has no independent source of truth -- only the user's original signature can prove the intended value.

The correct approach is to store the user's ECDSA signature on the `WithdrawRecord` and forward it to the Relay Wallet. The Relay Wallet MUST then verify the signature against the record's fields before executing the transaction. This provides end-to-end integrity from the user to the Relay Wallet, covering all fields (address, amount, benefit address, network) without relying on a shared secret between Relay components. With signature forwarding in place, the MAC becomes redundant and can be removed.

## Risk Control Strategy

Risk controls are layered across both services. Each service enforces its own checks independently, so a single point of compromise does not lead to fund loss.

### Authentication and Authorization

Every withdrawal request from a node MUST carry a valid JWT token and an ECDSA signature over the withdrawal details. The signature proves the node's private key holder authorized the specific withdrawal (address, amount, benefit address, network). The signature timestamp MUST be within 60 seconds.

The Relay Wallet-facing APIs use a separate signature scheme. Every request from the Relay Wallet MUST be signed by the configured wallet address, preventing unauthorized callers from fulfilling or rejecting withdrawals.

### Balance Verification

Balances are checked at multiple points:

- **On request**: The Relay checks the in-memory task fee cache to confirm the node has enough balance for the amount plus the fee.
- **On reconciliation**: The Relay re-checks against the persisted database balance after all preceding events have been flushed.
- **On wallet sync**: The Relay Wallet checks its independently maintained balance for the node address.
- **On-chain check**: The Relay Wallet verifies the system wallet has sufficient on-chain funds for the aggregate pending withdrawals on each network.
- **Post-confirmation**: After the blockchain transaction is confirmed, the Relay Wallet re-checks the local balance before deducting, as a final safeguard.

### Benefit Address Verification

The benefit address in every withdrawal request is verified against the on-chain benefit address contract. This check happens both in the Relay (at request time) and in the Relay Wallet (before execution). An attacker who manipulates the Relay database to change a benefit address would be caught by the Relay Wallet's independent on-chain lookup.

### Task Fee Log Risk Controls

The Relay Wallet validates every batch of task fee logs before applying them:

- A single log entry MUST NOT exceed `MaxTaskFeeAmount`. Logs of type `Bought` are exempt from this limit.
- The number of log entries for a single address within one batch MUST NOT exceed `MaxAddressLogsCountInBatch`.
- The number of previously unseen addresses within one batch MUST NOT exceed `MaxNewAddressCountInBatch`.

If any check fails, the Relay Wallet MUST halt the sync task and raise an alert. Processing MUST NOT resume until manual intervention resolves the issue.

### Transaction Safety

The Relay Wallet limits the number of in-flight transactions per blockchain network. Failed transactions are retried up to a configurable maximum, with a delay between retries. Each withdrawal request has an overall processing deadline; if the deadline passes, the request is rejected and the funds are refunded on the Relay side.

## Implementation Details

### Deposit Flow

The Relay runs a per-network block listener that polls for new blocks. For each transaction in a block, it checks whether the transaction is a pure native token transfer (no calldata) to `BuyTaskFee.Address`. If so:

1. Fetches the transaction receipt and confirms `receipt.Status == ReceiptStatusSuccessful`.
2. Checks the transaction has not already been processed (idempotency via the `Reason` field: `"{type}-{txHash}-{network}"`).
3. Creates a `TaskFeeEvent` of type `Bought` with the sender's address and `tx.Value()` as the amount, with `Status = Pending`.
4. Creates a `DepositRecord` (address, amount, network, tx_hash) for user query.
5. Credits the amount to the sender's in-memory task fee balance.

The `syncTaskFeesToDB` background job later re-validates the on-chain receipt and promotes the event to `Processed`, at which point it is persisted to the `TaskFee` table and becomes visible to the Relay Wallet via the task fee log API.

The Relay Wallet has no special deposit handling. It processes `Bought` task fee logs the same as any other type, except that `Bought` logs are exempt from the per-log maximum amount check.

**TODO:** The Relay Wallet SHOULD independently verify deposit transactions on-chain. The current task fee log API does not include the transaction hash or network, so the Relay Wallet has no way to check on-chain state. This means a compromised Relay can fabricate `Bought` logs with arbitrary amounts and addresses, and the Relay Wallet will accept them as long as the batch risk controls pass. To close this gap:

1. The Relay MUST include the transaction hash and network in `Bought` task fee log entries.
2. The Relay Wallet MUST fetch the on-chain transaction receipt for each `Bought` log and verify:
   - The transaction exists and `receipt.Status == ReceiptStatusSuccessful`.
   - The transaction is a native token transfer to the configured `BuyTaskFee.Address`.
   - The `tx.Value()` matches the amount in the log entry.
   - The sender matches the address in the log entry.
3. The Relay Wallet MUST record every processed transaction hash and reject any `Bought` log whose transaction hash has already been processed, preventing replay of the same on-chain transaction across multiple log entries.

### Withdrawal Flow

### Phase 1: Client Request (Client to Relay)

1. The client sends `POST /client/:address/withdraw` with JWT authentication.
2. The Relay validates:
   - The JWT-authenticated user address matches the path parameter.
   - The ECDSA signature over the message `"Withdraw {amount} from {address} to {benefitAddress} on {network}"` is valid and the signer matches the address. The signature timestamp MUST be within 60 seconds of the current time.
   - The benefit address matches the on-chain benefit address for the node.
   - The amount is at least `min_withdrawal_amount`.
3. In a single database transaction, the Relay:
   - Computes the withdrawal fee (zero for the fee collection address and the DAO address).
   - Checks the in-memory task fee balance covers `amount + withdrawal_fee`.
   - Deducts `amount + withdrawal_fee` from the cache.
   - Records the latest `TaskFeeEventID` for the address.
   - Creates a `WithdrawRecord` with `Status = Pending` and `LocalStatus = Pending`.
4. Returns the request ID to the client.

### Phase 2: Relay-Side Reconciliation

The `syncTaskFeesToDB` background job runs every 10 seconds and processes pending withdrawal records in batches of 50:

1. Fetches `WithdrawRecord` rows where `LocalStatus = Pending`, ordered by ID.
2. For each record, validates:
   - The record's `TaskFeeEventID` MUST be less than or equal to the last processed task fee event ID.
   - The address MUST exist in the persisted `TaskFee` table.
   - The persisted task fee balance MUST be greater than or equal to the withdrawal amount (unless the record has `Status = Failed`, indicating a rejected record being re-processed).
3. Invalid records are marked `LocalStatus = Invalid`.
4. Valid records are processed in a database transaction: the task fee balance is deducted, a MAC is computed, and the record is promoted to `LocalStatus = Processed`.

### Phase 3: Relay Wallet Sync and Validation

The Relay Wallet runs two background sync tasks:

**Task Fee Log Sync** (every 5 seconds):
1. Fetches task fee logs from the Relay starting after the last synced log ID.
2. Runs task fee risk control checks on the batch.
3. Merges logs by address and updates local `RelayAccount` balances.

**Withdrawal Request Sync** (every 5 seconds):
1. Fetches processed withdrawal requests from the Relay using signature authentication.
2. The Relay only returns records where `LocalStatus = Processed` and MAC verification passes.
3. The Relay Wallet filters requests: only those whose `TaskFeeEventID` is already covered by the synced task fee logs are accepted.
4. The Relay Wallet runs withdrawal risk control checks on the batch.
5. Accepted requests are stored locally with `Status = Pending`.

### Phase 4: Transaction Execution

For each pending local withdrawal record:

1. A processing goroutine is spawned with a deadline (configurable timeout from record creation time).
2. A blockchain transaction is queued to send ETH to the benefit address on the target network.
3. The transaction sender manages nonce allocation and submits the transaction, respecting the per-network in-flight transaction limit.
4. The transaction confirmer waits for a receipt. On failure, it retries up to the configured maximum.
5. On confirmation, the Relay Wallet re-checks the local balance, deducts the amount, and marks the record as `Success`.
6. If the transaction fails permanently or the deadline passes, the record is marked as `Failed`.

### Phase 5: Reporting

- **On success**: The Relay Wallet reports the transaction hash to the Relay. The Relay records the hash and credits the withdrawal fee to the fee collection address.
- **On failure or timeout**: The Relay Wallet reports rejection to the Relay. The Relay refunds `amount + withdrawal_fee` back to the node's task fee balance.

### Risk Control Checks Reference

#### Relay-Side Controls

| Control | Rule |
|---------|------|
| JWT authentication | The requesting user MUST own the address |
| ECDSA signature | The withdrawal message MUST be signed by the address owner; timestamp within 60 seconds |
| Benefit address verification | The benefit address MUST match the on-chain benefit address |
| Minimum withdrawal amount | `amount >= min_withdrawal_amount` |
| Task fee balance check (cache) | The in-memory task fee balance MUST cover `amount + withdrawal_fee` |
| Task fee balance check (DB) | The persisted task fee balance MUST cover the withdrawal amount |
| Task fee event ordering | The record's `TaskFeeEventID` MUST NOT exceed the last processed event ID |
| MAC integrity | Each processed record is MAC-tagged; MAC is verified before serving to the Relay Wallet |
| Relay Wallet signature authentication | All Relay Wallet API requests MUST be signed by the configured wallet address |

#### Relay Wallet-Side Withdrawal Controls

| Check | Rule |
|-------|------|
| Status validation | The request status MUST be `Pending` |
| Amount validation | The amount MUST be a valid decimal string and at least `min_withdrawal_amount` |
| Address existence | The address MUST exist in the Relay Wallet's local `RelayAccount` table |
| Balance per address | The sum of pending withdrawal amounts per address MUST NOT exceed the local balance |
| System wallet balance | The sum of pending withdrawal amounts per network MUST NOT exceed the system wallet's on-chain balance |
| Benefit address | The benefit address MUST match the on-chain benefit address |
| Post-confirmation re-check | Before deducting after blockchain confirmation, the local balance MUST cover the amount |

#### Relay Wallet-Side Task Fee Log Controls

| Check | Rule |
|-------|------|
| Per-log maximum | A single log amount MUST NOT exceed the configured maximum (token purchases are exempt) |
| Amount validation | The amount MUST be a valid decimal string |
| Per-address batch limit | The number of logs per address in a single batch MUST NOT exceed the configured limit |
| New address batch limit | The number of previously unseen addresses in a single batch MUST NOT exceed the configured limit |

#### Error Handling

- A risk control check failure causes the respective background task to stop and triggers an alert, requiring manual intervention.
- Non-risk-control errors (network failures, temporary issues) are retried after a delay.
- Each withdrawal record's processing goroutine is independent; a failure in one does not affect others.

## Data Models

### WithdrawRecord (Relay)

| Field | Type | Description |
|-------|------|-------------|
| `Address` | string | Node address |
| `BenefitAddress` | string | Destination address for the withdrawal |
| `Amount` | BigInt | Withdrawal amount in wei |
| `Network` | string | Target blockchain network |
| `Status` | WithdrawStatus | `Pending` (0), `Success` (1), `Failed` (2) |
| `LocalStatus` | WithdrawLocalStatus | `Pending` (0), `Processed` (1), `Invalid` (2) |
| `TaskFeeEventID` | uint | ID of the latest task fee event at creation time |
| `TxHash` | NullString | Blockchain transaction hash (set on fulfillment) |
| `WithdrawalFee` | BigInt | Fee amount in wei |
| `MAC` | string | HMAC-SHA256 integrity tag |

### WithdrawRecord (Relay Wallet)

| Field | Type | Description |
|-------|------|-------------|
| `RemoteID` | uint | ID of the record on the Relay |
| `Address` | string | Node address |
| `BenefitAddress` | string | Destination address |
| `Amount` | BigInt | Withdrawal amount in wei |
| `Network` | string | Target blockchain network |
| `Status` | WithdrawStatus | `Pending`, `Success`, `Failed`, `Finished` |
| `BlockchainTransactionID` | NullInt64 | Associated blockchain transaction |

### DepositRecord (Relay)

| Field | Type | Description |
|-------|------|-------------|
| `Address` | string | Depositor address |
| `Amount` | BigInt | Deposit amount in wei |
| `Network` | string | Source blockchain network |
| `TxHash` | string | On-chain transaction hash |

### RelayAccount (Relay Wallet)

| Field | Type | Description |
|-------|------|-------------|
| `Address` | string | Node address (unique) |
| `Balance` | BigInt | Locally tracked balance from task fee logs |

## Configuration

### Relay (`buy_task_fee` section)

| Key | Type | Description |
|-----|------|-------------|
| `address` | string | Target address for native token deposits |

### Relay (`withdraw` section)

| Key | Type | Description |
|-----|------|-------------|
| `address` | string | Wallet service address for signature authentication |
| `min_withdrawal_amount` | uint64 | Minimum withdrawal amount (in ether, converted to wei) |
| `withdrawal_fee` | uint64 | Per-withdrawal fee (in ether, converted to wei) |
| `withdrawal_fee_address` | string | Address that receives withdrawal fees |

### Relay (`mac` section)

| Key | Type | Description |
|-----|------|-------------|
| `secret_key` | string | HMAC-SHA256 secret key for MAC computation |

### Relay Wallet (`tasks.sync_withdrawal_requests`)

| Key | Type | Description |
|-----|------|-------------|
| `interval_seconds` | uint | Polling interval |
| `batch_size` | uint | Maximum requests per fetch |
| `min_withdrawal_amount` | uint64 | Minimum withdrawal amount (in ether) |

### Relay Wallet (`tasks.process_withdrawal_requests`)

| Key | Type | Description |
|-----|------|-------------|
| `interval_seconds` | uint | Processing interval |
| `batch_size` | uint | Maximum records per processing cycle |
| `timeout` | uint | Seconds before a withdrawal request is rejected due to timeout |

### Relay Wallet (`tasks.sync_task_fee_logs`)

| Key | Type | Description |
|-----|------|-------------|
| `interval_seconds` | uint | Polling interval |
| `batch_size` | uint | Maximum logs per fetch |
| `max_task_fee_amount` | uint | Maximum single task fee log amount (in ether) |
| `max_address_logs_count_in_batch` | uint | Maximum logs per address in a single batch |
| `max_new_address_count_in_batch` | uint | Maximum new addresses in a single batch |

### Relay Wallet (`blockchains.<network>`)

| Key | Type | Description |
|-----|------|-------------|
| `max_retries` | uint8 | Maximum retry attempts for a failed transaction |
| `retry_interval` | uint64 | Seconds between retries |
| `receipt_wait_time` | uint64 | Seconds to wait for a receipt before marking as failed |
| `sent_transaction_count_limit` | uint64 | Maximum in-flight transactions per network |

## API Endpoints

### Client APIs (JWT-protected)

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/client/:address/withdraw` | Create a withdrawal request |
| `GET` | `/client/:address/withdraw/list` | List withdrawal records for a user |
| `GET` | `/client/:address/deposit/list` | List deposit records for a user |

### Relay Wallet APIs (Signature-protected)

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/v1/withdraw/list` | List pending withdrawal requests |
| `POST` | `/v1/withdraw/:id/fulfill` | Mark a withdrawal as fulfilled with a transaction hash |
| `POST` | `/v1/withdraw/:id/reject` | Reject a withdrawal request |
