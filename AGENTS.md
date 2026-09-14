# AGENTS.md — Cavaridge Go SDK local instructions

Fleet-wide engineering-agent governance, approval boundaries, KISS/quiet mode, and the canonical ten-role topology are defined by `Cavaridge-LLC/cavaridge-engineering-control`. This SDK repository may implement an approved public API contract but may not define that contract, product authority, or security policy by itself.

**Shipping** owns SDK engineering delivery. **Product Partner** owns approved API/product requirements. **Security & Compliance** owns authentication/authorization and security assurance. Model/vendor tools are subordinate implementation workers.

## SDK authority

- The canonical API contract must come from the current accepted Cavaridge platform/API specification and verified implementation, not from this README or a legacy monorepo path.
- `Cavaridge-LLC/cavaridge-platform` is the canonical shared-platform implementation reference.
- If this SDK, documentation site, sandbox endpoint, or legacy API spec conflicts with current accepted platform contracts, stop the affected change and reconcile the contract rather than guessing.
- Do not publish or version SDK methods for planned/unimplemented endpoints as if they are production-supported.

## Security

- Never commit or print real API keys, bearer tokens, customer identifiers, or credentials.
- Examples use environment-variable placeholders or clearly fictional values only.
- Authentication examples do not authorize creation of live keys/tokens.
- Mutating SDK methods must preserve the platform's authorization, idempotency, tenant context, audit, approval, and error semantics.
- Agent/model identity never grants API, customer-data, deployment, or publishing authority.

## Delivery

- Preserve Go API compatibility according to the approved SDK/versioning contract.
- Keep examples, README, and generated/client code synchronized with the actual API surface.
- KISS applies to SDK docs and examples.
- Public release/tagging, production API changes, credential operations, publishing, or other consequential actions require applicable approval.
- Security-critical changes require independent review before `accepted` unless an authorized-human exception is explicitly recorded.

Do not treat links into the legacy `cavaridge` monorepo as current authority without reconciling them against the canonical platform/API contract.
