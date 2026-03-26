## How to run [deposit_withdraw_test.md](./deposit_withdraw_test.md)

Run the E2E test in [tests/e2e/deposit_withdraw_test.md](./deposit_withdraw_test.md):

- Use dymension blockchain.
- The Relay source code is in this repository.
- The Relay Wallet source code is at <_>.
- Use <_> as the root directory for all test files.
- Use <_> in Crynux MCP to transfer tokens to test accounts.

If an error occurs during testing, automatically fix and continue only when the error is in a config file created during the test run (for example, URL mismatch or account mismatch); if the error is in the code itself, stop immediately and explain the error without attempting a fix.

When diagnosing DB migration errors:
During test execution, migrations always start from a fresh database. After the first failure, the container may restart and trigger additional migration errors. Focus only on the first error, because it cannot be caused by an inconsistent DB state. Identify the migration error from the initial run and determine why it was triggered on a fresh database.
