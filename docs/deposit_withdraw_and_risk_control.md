# Deposit, Withdraw and Risk Control

This document defines the relay account funding, charging, income distribution, and withdrawal flow across Relay and Relay Wallet.

## Two-Service Architecture

The system has two isolated services:

- Relay: public-facing coordination service.
- Relay Wallet: key-holding execution service with no inbound control plane from Relay.

Relay Wallet MUST pull work from Relay. Relay MUST NOT push withdrawal execution commands directly to Relay Wallet.

## Relay Account Ledger

Relay uses a single ledger model:

- `relay_accounts`: per-address current balance.
- `relay_account_events`: append-only balance change log.

Relay account events MUST be the only source used by Relay Wallet to reconstruct balances.

### Relay Account Event Types

| Type | Name | Balance Effect |
|------|------|----------------|
| 0 | `TaskIncome` | Increase |
| 1 | `DaoTaskShare` | Increase |
| 2 | `WithdrawFeeIncome` | Increase |
| 3 | `Deposit` | Increase |
| 4 | `TaskPayment` | Decrease |
| 5 | `TaskRefund` | Increase |
| 6 | `Withdraw` | Decrease |
| 7 | `WithdrawRefund` | Increase |

`TaskIncome` is for node-side task settlement income. `DaoTaskShare` is for DAO-side settlement share.

`relay_account_events` MUST contain the complete relay account balance history for all relay account event types defined in this document.

`task_quota_events` MUST NOT be treated as relay account history and MUST NOT affect relay account balances.

## Deposit

There is no API to initiate a deposit.

A deposit happens when a user sends a native token transfer to `relay_account.deposit_address` on a supported network.

Relay MUST:

1. Detect successful native transfers to the configured deposit address.
2. Enforce idempotency per transaction hash and network.
3. Create a `relay_account_events` record of type `Deposit`.
4. Create a `deposit_records` record for client query.
5. Increase the sender relay account balance.

## Task Charge, Refund and Income

### Task Charge

On task creation, Relay MUST:

1. Validate creator relay account balance is sufficient for `task_fee`.
2. Create a `TaskPayment` relay account event for the creator.
3. Decrease creator relay account balance by `task_fee`.

### Task Refund

If task reaches refund or abort terminal paths, Relay MUST:

1. Create a `TaskRefund` event for the creator.
2. Increase creator relay account balance by the refunded amount.

### Task Income Distribution

On successful settlement, Relay MUST:

1. Split task fee by DAO ratio.
2. Create one `TaskIncome` event for the node address.
3. Create one `DaoTaskShare` event for the DAO address.
4. Increase both recipient balances by their split amounts.

## Withdraw

A withdrawal converts relay account balance to on-chain native tokens.

### Relay-Side Flow

On `POST /v1/client/:address/withdraw`, Relay MUST:

1. Validate JWT, signature, amount, and benefit address.
2. Compute withdrawal fee based on configured policy. Set `withdrawal_fee` to zero when requester address equals `dao.task_fee_share_address` or `withdraw.withdrawal_fee_address`.
3. Create a `withdraw_records` row with `Status = Pending`.
4. Create a `Withdraw` relay account event for `amount + withdrawal_fee`.
5. Store the created relay account event ID into `withdraw_records.relay_account_event_id`.
6. Decrease requester relay account balance by `amount + withdrawal_fee`.

When wallet fulfills a withdrawal, Relay MUST:

1. Mark withdrawal `Status = Success`.
2. Create `WithdrawFeeIncome` event to fee collection address when fee is non-zero.
3. Increase fee collection address relay account balance by fee amount.

When wallet rejects a withdrawal, Relay MUST:

1. Mark withdrawal `Status = Failed`.
2. Create `WithdrawRefund` event for requester.
3. Increase requester relay account balance by `amount + withdrawal_fee`.

## Relay Wallet Synchronization Rules

Relay Wallet MUST synchronize in event-order:

1. Sync `relay_account_events` in ascending ID order.
2. Verify signatures and integrity constraints before apply.
3. Apply each event to local account table using event type balance effect.
4. Sync withdrawal records and only execute records whose `relay_account_event_id` is not greater than the last applied relay account event ID.

Relay Wallet MUST reject or halt on ordering or integrity violations.

Relay account events MUST be retained as a complete ledger for audit, including `Withdraw` and `WithdrawRefund` events.

Relay Wallet should skip applying `Withdraw` and `WithdrawRefund` to its local balance if withdrawal deduction and rollback are handled by withdrawal processing flow.

When Relay Wallet skips applying an event type, it MUST still:

1. Verify event integrity.
2. Keep event-order continuity.
3. Advance sync checkpoint by event ID.

`withdraw_records.relay_account_event_id` MUST be treated as a synchronization watermark.

## Integrity and Authentication

- Client withdrawal requests MUST carry valid JWT and user signature.
- Wallet-facing APIs MUST use wallet signature authentication.
- Relay MUST generate MAC for processed ledger and withdrawal records.
- Relay Wallet MUST verify MAC before acceptance.

## Risk Control

### Relay Controls

- Relay MUST validate sufficient relay account balance before task charge and withdrawal charge.
- Relay MUST validate benefit address from on-chain source before withdrawal acceptance.
- Relay MUST enforce idempotency on deposit events.

### Relay Wallet Controls

- Relay Wallet MUST enforce batch-level risk limits on synced relay account events.
- Relay Wallet MUST verify local account sufficiency before sending withdrawal transactions.
- Relay Wallet MUST re-check local balance before marking a fulfilled withdrawal as deducted.

## Data Models

### Relay `withdraw_records`

| Field | Type | Description |
|-------|------|-------------|
| `address` | string | User address |
| `benefit_address` | string | Destination on-chain address |
| `amount` | BigInt | Requested withdraw amount |
| `network` | string | Target blockchain network |
| `status` | enum | `Pending`, `Success`, `Failed` |
| `local_status` | enum | `Pending`, `Processed`, `Invalid` |
| `relay_account_event_id` | uint | Ledger ordering anchor |
| `tx_hash` | string nullable | Fulfillment transaction hash |
| `withdrawal_fee` | BigInt | Fee amount |
| `mac` | string | Integrity tag |

### Relay `deposit_records`

| Field | Type | Description |
|-------|------|-------------|
| `address` | string | Depositor address |
| `amount` | BigInt | Deposit amount |
| `network` | string | Source network |
| `tx_hash` | string | Deposit transaction hash |

## Configuration

### Relay `relay_account`

| Key | Type | Description |
|-----|------|-------------|
| `deposit_address` | string | Native token deposit target |

### Relay `withdraw`

| Key | Type | Description |
|-----|------|-------------|
| `relay_wallet_address` | string | Wallet service signer address |
| `min_withdrawal_amount` | uint64 | Minimum withdraw amount |
| `withdrawal_fee` | uint64 | Withdraw fee. This fee is waived for `dao.task_fee_share_address` and `withdrawal_fee_address` |
| `withdrawal_fee_address` | string | Relay operator fee income address |

## API Endpoints

### Client APIs

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/v1/client/:address/withdraw` | Create withdrawal request |
| `GET` | `/v1/client/:address/withdraw/list` | Query withdrawals |
| `GET` | `/v1/client/:address/deposit/list` | Query deposits |
| `GET` | `/v1/client/:address/task_fee` | Query task fee records |

### Wallet APIs

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/v1/relay_account/logs` | Query relay account events |
| `GET` | `/v1/withdraw/list` | Query pending withdrawals |
| `POST` | `/v1/withdraw/:id/fulfill` | Mark fulfilled |
| `POST` | `/v1/withdraw/:id/reject` | Mark rejected |
