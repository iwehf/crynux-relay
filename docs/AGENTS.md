## Documentation Index

| Document | Description |
|----------|-------------|
| [node_selection.md](./node_selection.md) | Hard filters, base weight, model locality boost, and weighted sampling for task-to-node assignment |
| [qos.md](./qos.md) | Long-term performance score (`Q_long`) and short-term reliability factor (`H`) that compose the runtime QoS |
| [task_version.md](./task_version.md) | Version matching rules between task requirements and node capabilities |
| [task_validation_and_slashing.md](./task_validation_and_slashing.md) | Validation task lifecycle, result comparison, and slashing conditions |
| [model_predownload.md](./model_predownload.md) | Pre-download scheduling, node notification, and model availability tracking |
| [deposit_withdraw_and_risk_control.md](./deposit_withdraw_and_risk_control.md) | Deposit and withdrawal lifecycle across Relay and Wallet, task fee balance management, and risk control checks |

## Doc Update Requirements

When updating documentation files:

1. Read the entire document first to understand its structure, sections, and flow
2. Find the most appropriate location to integrate new content based on:
   - Logical relationship with existing sections
   - Document flow and narrative
   - Where readers would naturally expect to find the information
3. Integrate new content naturally into existing sections when possible:
   - Add as a paragraph within a relevant section
   - Extend an existing list or table
   - Add as a subsection under an appropriate parent section
   - Distribute across multiple sections if a feature affects different parts of the document
4. Do NOT simply create a new top-level section and place all new content there
5. Only create a new section if the topic is truly distinct from all existing content

Write documentation as a specification.

Documentation MUST state clear, final decisions and requirements.

Documentation MUST NOT include:
- Recommendations or advice.
- Options or alternatives.
- Speculation or uncertainty.
- Future-facing placeholders.

Documentation MUST use definitive language that can be implemented and tested:
- Requirement keywords: MUST, MUST NOT, SHALL, SHOULD. Use SHOULD only when a requirement level is intended.
- Exact behavior, constraints, and interfaces.
