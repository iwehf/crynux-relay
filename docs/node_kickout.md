# Node Kickout: Graduated Penalty System

## Problem

The current kickout mechanism removes a node permanently when 2 of its last 3 tasks are timeout failures. This is too aggressive in the presence of misbehaving applications (e.g., apps sending invalid tasks during testing). A burst of invalid tasks can cause honest, long-term stable nodes to be kicked from the network en masse, even though the timeouts were not their fault.

The core tension is:

- **Node protection**: An honest node should not be permanently removed due to a short burst of external failures.
- **Application quality**: Applications should not experience increased timeout rates because the network is too lenient toward genuinely bad nodes.

## Design Goals

1. A long-term stable node should tolerate occasional short-term errors without being removed.
2. A genuinely bad node should still be permanently removed, but based on sustained poor performance rather than a narrow 3-task window.
3. Application-side task success rates should not degrade — temporarily penalized nodes should be effectively excluded from task assignment.

## Overview

The new system replaces the binary kickout with a **temporary penalty and recovery mechanism**. Instead of permanently removing a node on short-term failures, the system temporarily reduces the node's chance of receiving tasks, then allows it to recover over time. Permanent kickout is reserved for nodes that demonstrate sustained poor performance over a much longer window.

### Temporary Penalty

Every time a task assigned to a node ends with a timeout, the node's task selection probability is **immediately and sharply reduced**. The reduction is multiplicative: each timeout multiplies the current probability by a penalty factor (0.3). So a single timeout reduces the probability by 70%. If the node times out again before recovering, the probability is reduced by another 70% of the already-reduced value. The effect is cumulative — consecutive timeouts compound rapidly.

When the probability drops low enough (below a threshold), the node is **completely excluded** from task selection — it receives zero tasks. From the application's perspective, this is identical to a kickout: the node is effectively invisible. This protects application quality during the penalty period.

### Recovery (Two Layers)

The penalty is temporary. The node's selection probability recovers through two complementary mechanisms:

**1. Passive time-based recovery.** Even if no tasks are assigned to the node, its selection probability slowly drifts back toward its normal value over time. This follows an exponential curve with a 30-minute time constant (~63% recovered after 30 min, ~86% after 60 min, ~95% after 90 min). This ensures that a penalized node is never stuck permanently — it always has a path back, even when it receives no tasks (which is exactly the situation when it's in the exclusion zone and can't earn tasks to prove itself).

**2. Active success-based recovery.** Every time the node completes a task successfully, its selection probability receives a discrete boost (+0.15). This is faster than passive recovery and serves as a proof-of-work mechanism — a node that actively demonstrates it can complete tasks recovers faster than one that simply waits. Once a penalized node starts receiving tasks again (after passive recovery lifts it out of the exclusion zone), a positive feedback loop kicks in: each success increases the probability of being selected, which leads to more tasks, which leads to more successes.

The two layers are complementary. Passive recovery handles the cold start problem (getting out of the exclusion zone where no tasks are available). Active recovery accelerates the rest. The time constant is deliberately set to be slow (30 min) so that the success-based recovery is a meaningful differentiator — a node that actively completes tasks recovers noticeably faster than one that simply waits.

### Permanent Kickout

Permanent removal is still possible, but based on the existing **QoS score** — the 50-task rolling average that already tracks each node's long-term performance. The QoS score is checked on every task completion. If it drops below a threshold (e.g., 1.0 out of 10), the node is permanently removed. This operates independently from the health multiplier — the two systems don't gate each other.

This reuses the existing QoS infrastructure rather than introducing a new tracking mechanism. It also has a useful built-in property: when an invalid app task causes ALL nodes in a group to timeout, the QoS scores for that group are set to NULL and excluded from the rolling average. This means app-caused failures do not drag down a node's QoS score at all — only genuine node-specific failures (where the node times out but other nodes in the same group succeed) affect the score.

### Summary

| Layer | Trigger | Effect | Recovery |
|-------|---------|--------|----------|
| Probability penalty | Each timeout failure | Selection probability *= 0.3 | Passive time decay (tau = 30 min) + active success boost (+0.15 per task) |
| Hard exclusion | Probability drops below threshold | Node receives zero tasks (same as kickout from app's view) | Automatic as probability recovers above threshold |
| Permanent kickout | QoS score below kickout threshold (checked on every task completion) | Node permanently removed | None (irreversible) |

## Health Multiplier

Each node carries a **health multiplier** H (range 0.0 to 1.0, default 1.0). This multiplier directly scales the node's task selection probability:

```
SelectionWeight = BaseProbability(staking, qos) * H
```

### On Timeout Failure

When a task assigned to the node ends with a timeout (same criteria as the current kickout check: `TaskEndAborted` + `TaskAbortTimeout` + node never submitted a result), the health multiplier is penalized:

```
H_new = H_effective * PenaltyFactor
```

Where `PenaltyFactor = 0.3`. This means:

- 1 timeout: H drops to 0.30 (70% reduction in selection probability)
- 2 consecutive timeouts: H drops to 0.09 (effectively excluded)
- 3 consecutive timeouts: H drops to 0.027

### On Successful Task Completion

When the node completes a task successfully (any non-timeout outcome), the health multiplier receives a boost:

```
H_new = min(1.0, H_effective + SuccessBoost)
```

Where `SuccessBoost = 0.15`. This provides active proof-of-work recovery — a node must demonstrate it can complete tasks to regain full health.

### Time-Based Recovery

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

Recovery works correctly across relay restarts because only the base value and timestamp are stored.

### How the Two Recovery Mechanisms Interact

Time-based recovery and success boost are both necessary. They serve complementary roles across different health ranges:

**H < 0.1 (exclusion zone):** The node is hard-excluded and receives zero tasks. Success boost is impossible — there are no tasks to succeed at. Time-based recovery is the **only** mechanism that works here. It slowly brings H back above 0.1, at which point the node re-enters the selection pool. This is intentional: the exclusion zone is a quarantine period designed to protect application quality, and the node must simply wait it out.

**H = 0.1 ~ 0.3 (low probability zone):** The node is back in the selection pool but receives very few tasks due to its low H. Time-based recovery is still the primary driver of recovery in this range. If a task does arrive and the node completes it, the success boost (+0.15) pushes H up significantly in relative terms (e.g., 0.2 to 0.35 is a 75% relative improvement). These boosts are rare but impactful when they happen.

**H = 0.3 ~ 0.7 (moderate zone):** This is where the success boost becomes the dominant force. The node receives tasks at a reasonable rate, and each success pushes H up meaningfully. There is a positive feedback loop: each success increases H, which increases the probability of being selected, which leads to more tasks, which leads to more successes. Recovery accelerates as the node climbs through this range.

**H > 0.7 (near-normal zone):** Both mechanisms bring H back to 1.0 quickly. The node is essentially recovered and operating normally.

The key takeaway: time-based recovery handles the cold start problem (getting out of the exclusion zone and the low-probability zone where tasks are rare). Success boost accelerates recovery once the node is receiving tasks again. Together they ensure that a penalized node recovers in a reasonable time regardless of task frequency, but a node that actively proves itself recovers faster.

Setting `tau` to a longer value (30 min instead of 15 min) makes the time-based recovery slower, which gives the success boost room to be a meaningful differentiator. A node that completes tasks successfully recovers noticeably faster than one that simply waits — the two mechanisms are not redundant.

## Hard Exclusion Threshold

When a node's effective health drops below the **exclusion threshold** (H < 0.1), it is completely excluded from task selection (selection weight = 0). This provides the same application-quality guarantee as the old kickout mechanism — no tasks are sent to the node during this period.

The difference from a permanent kickout is that the node automatically becomes eligible again as its health recovers above the threshold. After 2 consecutive timeouts (H ≈ 0.09), the node crosses back above 0.1 in approximately 1-2 minutes due to time recovery, and reaches 50% health in about 10 minutes.

## Permanent Kickout

Permanent removal is reserved for nodes that demonstrate sustained poor performance, using the existing **QoS score** (50-task rolling average, see [qos.md](qos.md)).

The QoS score is checked on every task completion. If it drops below `QoSKickoutThreshold` (e.g., 1.0 out of 10), the node is permanently kicked out (same execution as today: status set to quit, models deleted, unstake queued, `NodeKickedOutEvent` emitted). This check is independent of the health multiplier — the two systems operate in parallel, each handling its own concern.

This design has a useful built-in property: when an invalid app task causes ALL nodes in a validation group to timeout, the QoS scores for that group are set to NULL and **excluded from the rolling average**. App-caused failures do not drag down a node's QoS score. Only genuine node-specific failures — where the node times out but other nodes in the same group succeed (receiving a QoS score of 0) — affect the score. This means the QoS score is a reliable long-term signal for whether the node itself is the problem.

Note: QoS scores only cover grouped (validation) tasks, which are ~1% of all tasks. The temporary penalty mechanism (health multiplier) handles all tasks including single ones. The QoS-based kickout is the long-term backstop that catches nodes whose failures are consistently their own fault, not the app's.

## Worked Examples

### Honest Node Hit by 2 Invalid Tasks

```
State:          H = 1.0,  QoS = 7.5 (healthy)
Timeout #1:     H = 0.30                               selection prob reduced 70%
Timeout #2:     H = 0.09                               EXCLUDED (H < 0.1)
                QoS still ~7.5 (app-caused group                
                timeouts are NULL, not counted)         no kickout
  ... node is in exclusion zone, receives no tasks ...
  ... 20 minutes pass (time recovery) ...
Recovery:       H ≈ 0.37                               back in selection pool, low probability
  ... node receives a task and completes it ...
Success:        H ≈ 0.52                               success boost accelerates recovery
  ... 15 more minutes pass ...
Recovery:       H ≈ 0.71                               moderate probability
  ... another successful task ...
Success:        H ≈ 0.86                               nearly recovered
```

The node was excluded for ~20 minutes, then gradually returned to full capacity. Time recovery got it out of the exclusion zone, then the positive feedback loop (success -> higher H -> more tasks -> more successes) accelerated the rest. The QoS score was barely affected because the invalid tasks caused all-group timeouts (NULL scores, excluded from the rolling average). The node was never in danger of permanent removal.

### Genuinely Bad Node (Persistent Hardware Issue)

```
State:          H = 1.0, QoS = 5.0
Timeout #1:     H = 0.30                               reduced probability
                QoS drops (node got 0, others in
                group succeeded with 10 and 5)
  ... 30 min recovery, H ≈ 0.56 ...
Timeout #2:     H = 0.17, QoS drops further            further reduced
  ... 30 min recovery, H ≈ 0.48 ...
Timeout #3:     H = 0.14, QoS continues declining      barely above threshold
  ... pattern continues, each cycle: recover then fail ...
Eventually:     H drops below 0.1                      EXCLUDED
                QoS < 1.0                              PERMANENT KICKOUT
```

The genuinely bad node is permanently removed after its QoS score degrades below the threshold. Because QoS is a 50-task rolling average, this takes a sustained pattern of failures — each node-specific timeout pushes a 0 into the pool, gradually dragging the average down. The timeline is longer than the old system, but the tradeoff is that honest nodes are protected from false positives.

### Mixed Behavior (Occasional Failures)

A node that times out once every ~15 tasks sees its health dip briefly after each timeout (H = 0.3) but recovers within 30-60 minutes via time recovery plus success boosts from the tasks it completes in between. Its QoS score remains healthy because the occasional 0 is diluted by many good scores in the 50-task pool. This node is never excluded or kicked out — it simply experiences brief periods of reduced selection probability.

## Impact on Application Quality

| Scenario | Old system | New system |
|----------|-----------|------------|
| Node has 2 consecutive timeouts | Permanently kicked | Excluded ~10 min, then recovers |
| Node has 3+ consecutive timeouts | Permanently kicked | Excluded longer; if QoS also degraded, permanent kickout |
| During exclusion period | No tasks sent (kicked) | No tasks sent (H < 0.1, weight = 0) |
| After exclusion | Node gone forever | Node returns if healthy, or kicked if bad |

The key insight: **from the application's perspective, a hard-excluded node and a kicked-out node are identical** — neither receives tasks. The only difference is what happens after: a good node comes back, a bad node gets permanently removed. Application quality is preserved in both cases.

## Parameters

| Parameter | Value | Description |
|-----------|-------|-------------|
| `HealthPenaltyFactor` | 0.3 | Multiplier applied to health on each timeout |
| `HealthSuccessBoost` | 0.15 | Additive boost to health on each successful task |
| `HealthRecoveryTau` | 30 min | Time constant for exponential recovery toward 1.0 |
| `HealthExcludeThreshold` | 0.1 | Below this, node is hard-excluded from selection |
| `QoSKickoutThreshold` | 1.0 | QoS score (out of 10) below which a node is permanently kicked out |

## Relationship to Existing QoS Score

The health multiplier and the QoS score (50-task rolling pool) operate at different scopes and time scales:

| | QoS Score | Health Multiplier |
|---|---|---|
| **Scope** | Grouped (validation) tasks only (~1%) | All tasks (grouped and single) |
| **Time scale** | Long-term (50-task rolling average) | Short-term (immediate penalty + recovery over minutes) |
| **Purpose** | Reward fast execution among peers | Penalize timeout failures temporarily |
| **Effect on selection** | Part of base selection probability | Multiplier on top of base probability |
| **Effect on kickout** | QoS below threshold triggers permanent kickout | Health below threshold triggers exclusion + kickout check |
| **Recovery** | New task scores push out old ones | Time-based exponential decay (tau = 30 min) + success boost (+0.15 per task) |

Both feed into the final selection weight:

```
FinalWeight = BaseProbability(staking, qos) * HealthMultiplier
```

The QoS score serves two roles: it influences selection probability (higher QoS = more tasks) and it acts as the long-term signal for permanent kickout. The health multiplier is the short-term response mechanism that protects application quality immediately when a node starts failing, without waiting for the 50-task QoS average to catch up.
