# Deposit and Withdraw E2E Flow

## Scope

This document defines the execution flow for validating deposit and withdraw across Relay and Relay Wallet in e2e containers.

## Authoritative References

- Relay e2e guide: `tests/e2e/AGENTS.md` in the Relay repository.
- Relay Wallet e2e guide: `tests/e2e/AGENTS.md` in the Relay Wallet repository.
- Ledger and withdrawal behavior specification: `docs/deposit_withdraw_and_risk_control.md`.

## Execution Flow

### 1. Build e2e container images

- Build Relay e2e image from the Relay repository.
- Build Relay Wallet e2e image from the Relay Wallet repository.
- Confirm both images are available for compose startup.

### 2. Prepare mount folders and base config files

- Follow mount-folder and base-config preparation from:
  - Relay: `tests/e2e/AGENTS.md`.
  - Relay Wallet: `tests/e2e/AGENTS.md` in the Relay Wallet repository.

### 3. Configure runtime integration in mounted config files

- Update mounted Relay and Relay Wallet config files to align:
  - Database connectivity.
  - Relay API connectivity from Relay Wallet.
  - Shared blockchain network selection.
  - Required account and address fields used by deposit and withdraw flows.
- Set the following config values:
  - Relay `withdraw.min_withdrawal_amount` MUST be set to `5`.
  - Relay `withdraw.withdrawal_fee` MUST be set to `1`.
  - Relay Wallet minimum accepted deposit amount MUST be set to `5`.
  - Relay Wallet minimum accepted withdrawal amount MUST be set to `5`.

### 4. Configure wallets and keys with Crynux MCP

- Generate all private keys with Crynux MCP and use the fixed business names below.
- Prepare accounts with the following table:
| Name | Type | Purpose | Required mapping | Minimum test-token funding for 10 runs |
|------|------|---------|------------------|----------------------------------------|
| `relay_chain_system` | private key | Relay blockchain signer key | Relay `blockchains.dymension.account.address` and `blockchains.near.account.address` MUST be derived from this key. Relay `config/secrets/blockchain_system_private_key.txt` MUST store this key. | `0` |
| `relay_wallet_chain_system` | private key | Relay Wallet system payout key | Relay Wallet `blockchains.dymension.account.address` and `blockchains.near.account.address` MUST be derived from this key. Relay Wallet `/app/config/blockchain_privkey.txt` MUST store this key. | `50` (payout `4` per run x `10` + gas reserve `10`) |
| `relay_wallet_relay_api` | private key | Relay Wallet authentication key for Relay withdraw APIs | Relay Wallet `relay.api.private_key_file` MUST point to the file storing this key. Relay `withdraw.relay_wallet_address` MUST equal the address derived from this key. | `0` |
| `relay_credits_api_auth` | private key | Relay credits API authorization signer key | Relay `credits.api_auth_address` MUST equal the address derived from this key. | `0` |
| `client_e2e_user` | private key | Client account key for connect-wallet authentication, withdraw request signing, deposit transfer, and deposit record ownership checks | This key MUST be used by the e2e client actor for authenticated API calls and MUST be used as the sender in the deposit transaction of this scenario. | `60` (deposit `5` per run x `10` + gas reserve `10`) |
| `relay_account.deposit_address` | address | Deposit target address | Relay `relay_account.deposit_address` MUST be set to this address. | `0` |
| `withdraw.withdrawal_fee_address` | address | Withdrawal fee receiver address | Relay `withdraw.withdrawal_fee_address` MUST be set to this address. | `0` |
| `dao.task_fee_share_address` | address | DAO task fee share receiver address | Relay `dao.task_fee_share_address` MUST be set to this address. | `0` |
- Top up balances before container startup:
  - `relay_wallet_chain_system` account MUST hold enough native tokens on each configured network for gas and withdrawal payouts.
  - `client_e2e_user` account MUST hold enough native tokens for deposit transfer amount, gas, and any signed on-chain actions used by the test workflow.

### 5. Configure Relay secret keys

- Generate Relay secret keys locally as random strings.
- Prepare secrets with the following table:

| Target file | Type | Purpose | Required mapping |
|------|------|---------|------------------|
| `config/secrets/jwt_secret_key.txt` | random string | JWT signing secret for Relay client authentication tokens | This file MUST contain a single-line random string with no trailing whitespace. Relay `http.jwt.secret_key_file` MUST point to this file. |
| `config/secrets/mac_secret_key.txt` | random string | MAC signing secret for Relay record integrity protection | This file MUST contain a single-line random string with no trailing whitespace. Relay `mac.secret_key_file` MUST point to this file. |

### 6. Create e2e compose definition

- Create one compose file dedicated to this scenario.
- Include required services:
  - MySQL for Relay and Relay Wallet databases.
  - Relay service container.
  - Relay Wallet service container.
- Ensure service dependencies and volume mounts use the prepared host mount folders for Relay and Relay Wallet config and data.
- Host port mappings in the compose file MUST be unique across all services. Relay HTTP `8080` MAY be exposed once, and Relay Wallet MUST NOT expose an additional HTTP service port in this scenario.

### 7. Start Relay and Relay Wallet containers

- Start all compose services.
- Verify Relay and Relay Wallet are healthy and task loops are running.
- Confirm no startup errors that block deposit or withdrawal processing.

### 8. Execute deposit scenario

- Client deposit transaction amount MUST be `5` test token units.
- Send an on-chain transfer to Relay deposit address using Crynux MCP.
- After sending the transfer transaction, wait about `10` seconds before checking confirmation and relay ingestion.
- Use Crynux MCP to check relay account balance for the client account before and after deposit ingestion.
- Use Crynux MCP to check relay deposit records for the client account.
- Verify deposit record fields, deposited amount, and expected balance increase.

### 9. Execute withdraw scenario

- Client withdrawal request amount MUST be `5` test token units.
- With `withdrawal_fee = 1`, Relay Wallet MUST pay out `4` test token units to the user.
- Submit a withdrawal request for the same client account using Crynux MCP.
- After submitting the withdrawal request, wait about `10` seconds before checking withdraw status.
- Use Crynux MCP to check relay account balance for the client account before and after withdrawal processing.
- Use Crynux MCP to check withdrawal records for the client account.
- Verify withdrawal record fields, success status, payout transaction hash, payout amount `4`, and expected balance decrease.
- Verify final relay account balance with the fixed formula for this scenario: `balance_after_deposit - 5`.

### 10. Teardown and cleanup

- Stop and remove e2e containers and temporary runtime resources.
- Keep artifacts needed for troubleshooting when verification fails.
