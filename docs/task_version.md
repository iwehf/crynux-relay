# Task Version Matching

This document explains how the relay uses task version to select compatible nodes and how worker version is handled.

For the full node selection pipeline that applies these version filters, see [node_selection.md](./node_selection.md).

At a high level, task version defines the minimum compatible runtime for execution. The relay validates version formats at task creation and node join or update, then selects nodes using a strict compatibility gate where `major` must match exactly and node `minor.patch` must be greater than or equal to task `minor.patch`. The same compatibility check is performed again right before moving a task to `Started` to prevent race-condition mismatches. Worker version is tracked separately for worker count APIs and does not affect task dispatch or node selection.

## Scope

- Task version matching is implemented for node selection and task start validation.
- Worker version is tracked as runtime count data and is not used in task dispatch.

## Version Format

Both task version and node version use semantic version style with three numeric parts.

- Task version format: `major.minor.patch`
- Node version format: `major.minor.patch`

Validation entry points:

- Task creation validates `task_version` in `api/v1/inference_tasks/create_task.go`
- Node join validates `version` in `api/v1/nodes/join.go` and `api/v2/nodes/join.go`
- Node version update validates `version` in `api/v1/nodes/version.go`

## Data Model

- Task version is stored as string in `models.InferenceTask.TaskVersion`
- Node version is stored as numeric fields in `models.Node`
  - `MajorVersion`
  - `MinorVersion`
  - `PatchVersion`

The task version parser is `InferenceTask.VersionNumbers` in `models/inference_task.go`.

## Node Compatibility Rule

The compatibility rule used by relay is:

- Node major version must equal task major version
- Node minor and patch must be greater than or equal to task minor and patch

Equivalent comparison:

- `node.major == task.major`
- `node.minor > task.minor`
- or `node.minor == task.minor` and `node.patch >= task.patch`

This means matching is strict on major version and forward compatible on minor and patch.

## Where Matching Happens

### Step 1: Candidate Filtering During Node Selection

When dispatcher selects a node for a queued task, the filtering happens in:

- `service/select_nodes.go`
  - `filterNodesByGPU`
  - `filterNodesByVram`

Database filter conditions include:

- `status = available`
- hardware constraints
- `major_version = task.major`
- `minor_version > task.minor OR minor_version = task.minor AND patch_version >= task.patch`

Selection call chain:

- `service/start_task.go` -> `TaskDispatcher.Dispatch`
- `service/select_nodes.go` -> `selectNodeForInferenceTask`

### Step 2: Safety Recheck Before Task Starts

Before final state change to `TaskStarted`, relay checks version again in:

- `service/task_status.go`
  - `isNodeVersionValidForTask`
  - used by `SetTaskStatusStarted`

If version is not compatible, task start fails with:

- `node version is not compatible with task`

This prevents race conditions where node version may have changed after initial selection.

## Worker Version and Runner Matching

In this codebase, runner version is represented by worker version endpoints.

- `api/v1/worker/worker.go` stores counts by exact `WorkerVersion` string
- `models/worker.go` persists `WorkerVersion` and `Count`

Current behavior:

- Worker version data is only used for join, quit, and count query APIs
- No code path uses worker version to filter nodes
- No code path uses task version to pick worker version
- No task dispatch logic depends on worker version

## Practical Outcome

Given a task version `A.B.C`:

- Relay will only consider nodes with major version `A`
- Among those nodes, relay accepts versions `A.B.C` and newer within major `A`
- Relay rejects any node with major version not equal to `A`
- Worker version has no effect on this decision

## Related Source Files

- `models/inference_task.go`
- `models/node.go`
- `service/select_nodes.go`
- `service/start_task.go`
- `service/task_status.go`
- `api/v1/inference_tasks/create_task.go`
- `api/v1/nodes/join.go`
- `api/v2/nodes/join.go`
- `api/v1/nodes/version.go`
- `api/v1/worker/worker.go`
- `models/worker.go`
