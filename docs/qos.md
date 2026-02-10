# Quality of Service (QoS) Implementation

This document describes the QoS implementation in the Crynux Relay codebase.

## Overview

The QoS system incentivizes nodes to perform well by:

1. Assigning **Task Scores** to individual tasks based on execution speed.
2. Maintaining a **Node QoS Score** as a rolling average of recent task scores.
3. Using the Node QoS Score to influence **node selection probability** for future tasks.
4. **Permanently kicking out** nodes whose QoS score drops below a threshold.
5. Distributing **rewards proportionally** to task scores within grouped tasks.

For the short-term **health multiplier** system that temporarily penalizes nodes on timeout failures, see [node_health.md](node_health.md).

## Task Grouping (Validation Tasks)

Not all tasks go through QoS scoring. Only **grouped tasks** (validation tasks) receive QoS scores.

A VRF (Verifiable Random Function) determines whether a task requires validation. The task creator generates a VRF proof, and the relay verifies it. The selection logic is:

```go
// utils/vrf.go
func VrfNeedValidation(vrfNumber []byte) bool {
    number := big.NewInt(0).SetBytes(vrfNumber)
    r := big.NewInt(0).Mod(number, big.NewInt(100)).Uint64()
    return r == 0
}
```

This means approximately **1% of tasks** are selected for validation grouping. A grouped task is executed by **3 different nodes** simultaneously. Single (non-grouped) tasks are not scored for QoS.

## Task Score

When a grouped task is validated, each of the 3 tasks in the group receives a **Task Score** based on execution speed.

### Execution Time Measurement

Execution time is measured as:

```
ExecutionTime = ScoreReadyTime - StartTime
```

- `StartTime`: when the task was assigned to the node and started.
- `ScoreReadyTime`: when the node submitted its result score (hash).

If either timestamp is missing, the execution time is set to the maximum possible duration (effectively infinite).

### Score Assignment

Tasks within a group are sorted by execution time (fastest first). If two tasks have the same execution time, they are sorted by database ID (lower ID first). The fixed score values are:

| Completion Order | Task Score |
|-----------------|------------|
| 1st (fastest)   | 10         |
| 2nd             | 5          |
| 3rd (slowest)   | 2          |

These values are defined as:

```go
// service/qos.go
TASK_SCORE_REWARDS [3]uint64 = [3]uint64{10, 5, 2}
```

Special cases:
- Tasks that were **aborted** before the group validation receive a score of **0**.
- If **all 3 tasks** in a group are aborted, QoS scores are set to NULL (not valid) and are **not included** in the node's rolling average.

## Node QoS Score

Each node maintains a QoS score that represents its recent performance. This is a **rolling average** of the task scores from its most recent tasks.

### Rolling Pool Mechanism

The node QoS score is calculated using an in-memory rolling pool:

- **Pool size**: 50 tasks (`NODE_QOS_SCORE_POOL_SIZE = 50`)
- The pool is stored per node address in a concurrent-safe map (`NodeQosScorePool`).
- When a new task score arrives, it is appended to the pool. If the pool exceeds 50 entries, the oldest entry is removed.
- The node's QoS score is the **arithmetic mean** of all scores in the pool.

### Pool Initialization

When a node's pool does not yet exist in memory (e.g., after a relay restart), the pool is initialized as follows:

- If the node already has a non-zero `QOSScore` in the database, the pool is pre-filled with 49 copies of that existing score, then the new score is appended (total = 50).
- If the node has no existing score, the pool starts empty and the new score is the first entry.

This ensures that the rolling average transitions smoothly from the persisted score rather than jumping abruptly.

### When the Node QoS Score is Updated

The node QoS score is updated in the following task status transitions:

1. **TaskGroupValidated** (`SetTaskStatusGroupValidated`): The "winning" task in the group. Node QoS is updated with the task's QoS score.
2. **TaskEndGroupRefund** (`SetTaskStatusEndGroupRefund`): The other valid tasks in the group that are refunded. Node QoS is updated with the task's QoS score.
3. **TaskEndAborted** (`SetTaskStatusEndAborted`): If the task has a valid QoS score (i.e., it was part of a group that was validated), the node QoS is updated.

The `updateNodeQosScore` function calls `getNodeTaskQosScore` to compute the new rolling average and persists it to the database:

```go
// service/node.go
func updateNodeQosScore(ctx context.Context, db *gorm.DB, node *models.Node, qos uint64) error {
    qosScore, err := getNodeTaskQosScore(node, qos)
    if err != nil {
        return err
    }
    return node.Update(ctx, db, map[string]interface{}{
        "qos_score": qosScore,
    })
}
```

## Node Selection Probability

The QoS score directly influences a node's probability of being selected for new tasks. The selection probability combines two factors: **Staking Score** and **QoS Score**.

### Staking Score

```
StakingScore = sqrt(staking / maxStaking)
```

- `staking`: the node's staked amount.
- `maxStaking`: the maximum staked amount among all nodes in the network (tracked globally and refreshed on node join/quit).

### QoS Score Normalization

```
QoSProb = nodeQoSScore / maxQoSScore
```

- `nodeQoSScore`: the node's current rolling average QoS score.
- `maxQoSScore`: a **fixed constant** equal to `TASK_SCORE_REWARDS[0]` = **10** (the maximum possible task score).

**Important**: If the calculated `QoSProb` is 0 (e.g., for a new node with no score yet), it defaults to **0.5** as a baseline.

### Combined Probability (Harmonic Mean)

The final selection weight combines staking and QoS using the **harmonic mean formula**:

```
ProbWeight = StakingScore * QoSProb / (StakingScore + QoSProb)
```

If either component is 0, the combined probability is 0.

### Model Locality Boost

After computing base probabilities, nodes are further boosted based on whether they already have the required models locally. A task may require multiple models (stored as the `ModelIDs` list) -- for example, an SD task might need a base model plus one or more LoRA models. For LLM tasks, this is typically just a single model.

The boost logic works as follows:

- If the node's currently **in-use models match exactly** with the task's required models: **2x boost**.
- If the node has **some (but not all)** of the required models locally: boost by `1 + matchCount / totalRequired` (between 1x and 2x).
- If **at least one node** has matching local models, then **only** those nodes with local models are considered (nodes without local models are excluded from the selection pool).
- If **no nodes** have any matching local models, the boost step is skipped entirely, and selection falls back to the **full candidate list** using only the base staking + QoS probabilities.

### Model Pre-download (Download Task)

To ensure that in-demand models are available on enough nodes, the relay triggers a **model pre-download mechanism** every time a task starts. This proactively spreads models to additional nodes when fewer than 3 available nodes have a required model, so that future tasks are more likely to benefit from the model locality boost described above.

For a detailed description of this mechanism, including the spreading logic, node selection process, and relevant code locations, see [Model Pre-download Mechanism](model_predownload.md).

### Weighted Random Selection

The actual node selection uses weighted random sampling (`gonum/stat/sampleuv.NewWeighted`), where each node's weight is its computed probability. This is a probabilistic selection, not deterministic - higher-weighted nodes are more likely to be selected, but any eligible node can be chosen.

## Permanent Kickout

The permanent kickout mechanism removes nodes whose QoS score demonstrates sustained poor performance. It uses the same 50-task rolling average QoS score described above.

### Kickout Criteria

The kickout check (`ShouldPermanentKickout` in `service/qos.go`) evaluates two conditions, both of which must be true:

1. The node's QoS score is below the configured `qos.kickout_threshold` (default: 2.0).
2. The QoS score pool has accumulated enough samples (equal to `qos.score_pool_size`, default: 50).

The second condition prevents premature kickout of nodes that have only completed a few grouped tasks — the system waits until there is a statistically meaningful sample before making a permanent removal decision.

### When the Kickout Check Runs

The check happens inside `nodeFinishTask`, which is called every time a node completes processing a task — whether the outcome is success, group refund, or abort.

### Kickout Execution

If the kickout condition is met:

1. The node's status is set to quit.
2. All of the node's local model records are deleted.
3. An unstake transaction is queued on the blockchain (the node is **not slashed** — its stake is returned).
4. A `NodeKickedOutEvent` is emitted.

### Interaction with Application-Caused Failures

When an invalid application task causes ALL nodes in a validation group to timeout, the QoS scores for that group are set to NULL and **excluded from the rolling average**. This means application-caused failures do not drag down a node's QoS score. Only genuine node-specific failures — where the node times out but other nodes in the same group succeed (receiving a QoS score of 0) — affect the score.

Since QoS scores only cover grouped (validation) tasks (~1% of all tasks), the permanent kickout is a long-term backstop that catches nodes whose failures are consistently their own fault. For short-term timeout penalty and temporary exclusion, see the [health multiplier](node_health.md).

## Reward Distribution for Grouped Tasks

When a grouped task succeeds (the result is uploaded), rewards are distributed proportionally to the task QoS scores rather than equally.

### Calculation

For each valid task in the group (status `TaskGroupValidated` or `TaskEndGroupRefund`):

```
payment = taskFee * taskQoSScore / totalScore
```

Where:
- `taskFee`: the fee associated with each individual task in the group.
- `taskQoSScore`: the task's QoS score (10, 5, or 2).
- `totalScore`: the sum of all valid tasks' QoS scores in the group.

Any remainder from integer division is accumulated and added to the **last valid task's** payment to ensure no tokens are lost.

### Example

For a group of 3 tasks with fee = 100 each, all completed successfully:

| Node | Score | Share | Payment |
|------|-------|-------|---------|
| 1st  | 10    | 10/17 | 58      |
| 2nd  | 5     | 5/17  | 29      |
| 3rd  | 2     | 2/17  | 13 (includes remainder) |

## Task Validation Logic

The validation process in `ValidateTaskGroup` determines the outcome for each task in a group:

### Score Comparison

- For **SD/SD-FT-LoRA tasks**: Scores (perceptual hashes) are compared using Hamming distance. Each 8-byte segment must have a Hamming distance below a configurable threshold.
- For **LLM tasks**: Scores must be exactly equal.

### Validation Outcomes

Given 3 tasks in a group, the system checks pairwise similarity:

| Scenario | Task 1 | Task 2 | Task 3 |
|----------|--------|--------|--------|
| All 3 match | GroupValidated | GroupRefund | GroupRefund |
| 1 & 2 match, 3 differs | GroupValidated | GroupRefund | Invalidated (slashed) |
| 1 & 3 match, 2 differs | GroupValidated | Invalidated (slashed) | GroupRefund |
| 2 & 3 match, 1 differs | Invalidated (slashed) | GroupValidated | GroupRefund |
| None match | All Aborted | All Aborted | All Aborted |

The **first matching task** (by execution time order) becomes `GroupValidated` (the "winner" whose result is used). The other matching task(s) become `GroupRefund` (they get refunded but still receive QoS scores). Invalidated tasks result in the node being **slashed** (staked tokens confiscated).

Note: The `GroupValidated` task is the one that eventually has its result uploaded by the task creator (triggering `SetTaskStatusEndSuccess` -> `SetTaskStatusEndGroupSuccess`). At that point, rewards are distributed to all valid tasks in the group proportionally to their QoS scores.

## Data Flow Summary

```
Task Created
    |
    v
Task Started (assigned to node)
    |
    v
Node submits result -> TaskScoreReady (records ScoreReadyTime)
    |
    v
Task Creator validates group (VRF proof + score comparison)
    |
    v
+-- GroupValidated (winner) -----> updateNodeQosScore(node, taskScore)
|
+-- GroupRefund (matching) -------> updateNodeQosScore(node, taskScore)
|
+-- Invalidated (cheater) -------> nodeSlash (no QoS update, kicked & slashed)
|
+-- Aborted (all failed) --------> QoS score set to NULL (ignored)
    |
    v
Task Creator uploads result -> EndGroupSuccess
    |
    v
Rewards distributed proportionally to task QoS scores
```

## Key Constants and Config

| Constant / Config | Value | Description |
|-------------------|-------|-------------|
| `TASK_SCORE_REWARDS` | [10, 5, 2] | Task scores for 1st, 2nd, 3rd place |
| `qos.score_pool_size` | 50 (default) | Rolling pool size for node QoS calculation |
| `qos.kickout_threshold` | 2.0 (default) | QoS score below which a node is permanently kicked out |
| `maxQoSScore` | 10.0 | Fixed normalization denominator for QoS score |
| Default QoS Prob | 0.5 | Fallback when QoS probability is 0 |

## Relevant Source Files

| File | Description |
|------|-------------|
| `service/qos.go` | Core QoS logic: task score assignment, rolling pool, permanent kickout check |
| `service/health.go` | Health multiplier: penalty, boost, effective health calculation |
| `service/selecting_prob.go` | Selection probability calculation (staking + QoS) |
| `service/validate_task.go` | Task group validation and QoS score assignment |
| `service/task_status.go` | Task state transitions, QoS updates, health penalty/boost |
| `service/select_nodes.go` | Node selection using QoS-weighted probability * health multiplier |
| `service/node.go` | Node management, QoS persistence, kickout execution |
| `models/inference_task.go` | Task model with QOSScore field and ExecutionTime method |
| `models/node.go` | Node model with QOSScore, HealthBase, HealthUpdatedAt fields |
| `models/network.go` | NetworkNodeData with QoS field (synced periodically) |
| `tasks/sync_network.go` | Background task syncing node QoS to network statistics |
| `utils/vrf.go` | VRF validation check for task grouping |
