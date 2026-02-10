# Node Health Multiplier

## Purpose

The health multiplier protects application quality when nodes start failing, without permanently removing nodes that may be experiencing temporary issues. It addresses two competing goals:

- **Application quality**: When a node times out, it should be immediately excluded from receiving further tasks so that applications are not affected by unreliable nodes.
- **Node protection**: An otherwise healthy node should not be permanently removed from the network due to a short burst of failures (e.g., caused by invalid application tasks or transient network issues). It should be given a chance to recover and prove itself.

The health multiplier achieves both by sharply reducing a failing node's selection probability (protecting applications), while allowing automatic recovery over time and through successful task completions (protecting nodes).

## Overview

### Penalty

Every time a task assigned to a node ends with a timeout, the node's task selection probability is **immediately and sharply reduced**. The reduction is multiplicative: each timeout multiplies the current probability by a penalty factor (0.3). A single timeout reduces the probability by 70%. If the node times out again before recovering, the probability is reduced by another 70% of the already-reduced value. Consecutive timeouts compound rapidly.

When the probability drops low enough (below a threshold), the node is **completely excluded** from task selection — it receives zero tasks. From the application's perspective, this is identical to a kickout: the node is effectively invisible.

### Recovery (Two Layers)

The penalty is temporary. The node's selection probability recovers through two complementary mechanisms:

**1. Passive time-based recovery.** Even if no tasks are assigned to the node, its selection probability slowly drifts back toward its normal value over time. This follows an exponential curve with a 30-minute time constant (~63% recovered after 30 min, ~86% after 60 min, ~95% after 90 min). This ensures that a penalized node is never stuck permanently — it always has a path back, even when it receives no tasks (which is exactly the situation when it's in the exclusion zone).

**2. Active success-based recovery.** Every time the node completes a task successfully, its selection probability receives a discrete boost (+0.15). This is faster than passive recovery and serves as a proof-of-work mechanism — a node that actively demonstrates it can complete tasks recovers faster than one that simply waits.

The two layers are complementary. Passive recovery handles the cold start problem (getting out of the exclusion zone where no tasks are available). Active recovery accelerates the rest. The time constant is deliberately slow (30 min) so that success-based recovery is a meaningful differentiator — a node that completes tasks successfully recovers noticeably faster than one that simply waits.

### Permanent Kickout

Permanent removal is handled separately by the [QoS score system](qos.md#permanent-kickout). The health multiplier only applies temporary penalties; permanent kickout is based on sustained poor performance over the 50-task QoS rolling average.

| Layer | Trigger | Effect | Recovery |
|-------|---------|--------|----------|
| Probability penalty | Each timeout failure | Selection probability *= 0.3 | Passive time decay (tau = 30 min) + active success boost (+0.15 per task) |
| Hard exclusion | Probability drops below 0.1 | Node receives zero tasks | Automatic as probability recovers above threshold |
| Permanent kickout | QoS score below threshold | Node permanently removed ([see qos.md](qos.md#permanent-kickout)) | None (irreversible) |

## Health Multiplier

Each node carries a **health multiplier** H (range 0.0 to 1.0, default 1.0). This multiplier directly scales the node's task selection probability:

```
SelectionWeight = BaseProbability(staking, qos) * H
```

## On Timeout Failure

When a task assigned to the node ends with a timeout (`TaskEndAborted` + `TaskAbortTimeout` + node never submitted a result), the health multiplier is penalized:

```
H_new = H_effective * PenaltyFactor
```

Where `PenaltyFactor = 0.3`. This means:

- 1 timeout: H drops to 0.30 (70% reduction in selection probability)
- 2 consecutive timeouts: H drops to 0.09 (effectively excluded)
- 3 consecutive timeouts: H drops to 0.027

## On Successful Task Completion

When the node completes a task successfully, the health multiplier receives a boost:

```
H_new = min(1.0, H_effective + SuccessBoost)
```

Where `SuccessBoost = 0.15`. This provides active proof-of-work recovery — a node must demonstrate it can complete tasks to regain full health.

Success boost is applied in:

- `SetTaskStatusEndSuccess` — single and group task success
- `SetTaskStatusGroupValidated` — winning task in a group
- `SetTaskStatusEndGroupRefund` — other matching tasks in a group

## Time-Based Recovery

Health recovers automatically over time via exponential decay toward 1.0:

```
H_effective(t) = H_base + (1 - H_base) * (1 - exp(-(t - t_base) / tau))
```

Where:

- `H_base` is the stored health value at the time of the last update.
- `t_base` is the timestamp of the last update.
- `tau = 30 minutes` is the recovery time constant.

This formula is computed **lazily** — there is no background timer. The effective health is calculated on-the-fly whenever the value is needed (during node selection or penalty application). This means:

- After 1 tau (30 min): ~63% of the gap to 1.0 is recovered.
- After 2 tau (60 min): ~86% recovered.
- After 3 tau (90 min): ~95% recovered.

Recovery works correctly across relay restarts because only the base value and timestamp are stored in the database.

## How the Two Recovery Mechanisms Interact

Time-based recovery and success boost serve complementary roles across different health ranges:

**H < 0.1 (exclusion zone):** The node is hard-excluded and receives zero tasks. Success boost is impossible — there are no tasks to succeed at. Time-based recovery is the **only** mechanism that works here. It slowly brings H back above 0.1, at which point the node re-enters the selection pool.

**H = 0.1 ~ 0.3 (low probability zone):** The node is back in the selection pool but receives very few tasks due to its low H. Time-based recovery is still the primary driver. If a task does arrive and the node completes it, the success boost (+0.15) pushes H up significantly in relative terms (e.g., 0.2 to 0.35 is a 75% relative improvement).

**H = 0.3 ~ 0.7 (moderate zone):** This is where the success boost becomes the dominant force. The node receives tasks at a reasonable rate, and each success pushes H up meaningfully. There is a positive feedback loop: each success increases H, which increases the probability of being selected, which leads to more tasks, which leads to more successes.

**H > 0.7 (near-normal zone):** Both mechanisms bring H back to 1.0 quickly. The node is essentially recovered and operating normally.

The key takeaway: time-based recovery handles the cold start problem (getting out of the exclusion zone where tasks are unavailable). Success boost accelerates recovery once the node is receiving tasks again.

## Hard Exclusion Threshold

When a node's effective health drops below the **exclusion threshold** (H < 0.1), it is completely excluded from task selection (selection weight = 0). No tasks are sent to the node during this period.

The node automatically becomes eligible again as its health recovers above the threshold. After 2 consecutive timeouts (H ~ 0.09), the node crosses back above 0.1 in approximately 1-2 minutes due to time recovery, and reaches 50% health in about 10 minutes.

## Health Reset on Join

When a node joins (or re-joins) the network, its health is reset to 1.0. Returning nodes start fresh.

## Worked Examples

### Node Hit by Invalid Application Tasks

```
State:          H = 1.0,  QoS = 7.5 (healthy)
Timeout #1:     H = 0.30                               selection prob reduced 70%
Timeout #2:     H = 0.09                               EXCLUDED (H < 0.1)
                QoS still ~7.5 (app-caused group
                timeouts are NULL, not counted)         no kickout
  ... node is in exclusion zone, receives no tasks ...
  ... 20 minutes pass (time recovery) ...
Recovery:       H ~ 0.37                               back in selection pool, low probability
  ... node receives a task and completes it ...
Success:        H ~ 0.52                               success boost accelerates recovery
  ... 15 more minutes pass ...
Recovery:       H ~ 0.71                               moderate probability
  ... another successful task ...
Success:        H ~ 0.86                               nearly recovered
```

The node is excluded for ~20 minutes, then gradually returns to full capacity. The QoS score is barely affected because the invalid tasks caused all-group timeouts (NULL scores, excluded from the rolling average). The node is never in danger of permanent removal.

### Genuinely Bad Node (Persistent Hardware Issue)

```
State:          H = 1.0, QoS = 5.0
Timeout #1:     H = 0.30                               reduced probability
                QoS drops (node got 0, others in
                group succeeded with 10 and 5)
  ... 30 min recovery, H ~ 0.56 ...
Timeout #2:     H = 0.17, QoS drops further            further reduced
  ... 30 min recovery, H ~ 0.48 ...
Timeout #3:     H = 0.14, QoS continues declining      barely above threshold
  ... pattern continues, each cycle: recover then fail ...
Eventually:     H drops below 0.1                      EXCLUDED
                QoS < 2.0                              PERMANENT KICKOUT (see qos.md)
```

The genuinely bad node is permanently removed after its QoS score degrades below the threshold. The health multiplier keeps it mostly out of the task pool in the meantime, protecting application quality while the 50-task QoS average catches up.

### Mixed Behavior (Occasional Failures)

A node that times out once every ~15 tasks sees its health dip briefly after each timeout (H = 0.3) but recovers within 30-60 minutes via time recovery plus success boosts. Its QoS score remains healthy because the occasional 0 is diluted by many good scores in the 50-task pool. This node is never excluded or kicked out — it simply experiences brief periods of reduced selection probability.

## Parameters

All parameters are configurable via the `node_health` section in the config file.

| Parameter | Config Key | Default | Description |
|-----------|-----------|---------|-------------|
| Penalty Factor | `node_health.penalty_factor` | 0.3 | Multiplier applied to health on each timeout |
| Success Boost | `node_health.success_boost` | 0.15 | Additive boost to health on each successful task |
| Recovery Tau | `node_health.recovery_tau_minutes` | 30 | Time constant (minutes) for exponential recovery toward 1.0 |
| Exclude Threshold | `node_health.exclude_threshold` | 0.1 | Below this health value, node is hard-excluded from selection |

## Relationship to QoS Score

The health multiplier and the QoS score operate at different scopes and time scales:

| | QoS Score | Health Multiplier |
|---|---|---|
| **Scope** | Grouped (validation) tasks only (~1%) | All tasks (grouped and single) |
| **Time scale** | Long-term (50-task rolling average) | Short-term (immediate penalty + recovery over minutes) |
| **Purpose** | Reward fast execution among peers | Penalize timeout failures temporarily |
| **Effect on selection** | Part of base selection probability | Multiplier on top of base probability |
| **Effect on removal** | QoS below threshold triggers permanent kickout | Health below threshold triggers temporary exclusion |
| **Recovery** | New task scores push out old ones | Time-based exponential decay + success boost |

Both feed into the final selection weight:

```
FinalWeight = BaseProbability(staking, qos) * HealthMultiplier
```

## Relevant Source Files

| File | Description |
|------|-------------|
| `service/health.go` | Core health logic: effective health calculation, penalty, boost |
| `service/select_nodes.go` | Health multiplier applied during node selection |
| `service/task_status.go` | Health penalty on timeout, health boost on success |
| `service/node.go` | Health reset on node join |
| `models/node.go` | Node model with `HealthBase` and `HealthUpdatedAt` fields |
| `config/app_config.go` | `NodeHealth` config section |
