# infra-sparring

Your sparring partner for Linux certs: an AI agent that breaks your lab on command so you learn to fix it.

`spar` is a trainer for Linux and Kubernetes certifications. An LLM generates fresh faults in a
VM lab, you diagnose them from the terminal, and at the end you get a report on how you solved them.

> Status: in design. First pack: LFCS.

## How to play

```bash
spar auth login                 # once: stores the API key in the system keyring
spar start                      # you get a symptom, not the cause
vagrant ssh ubuntu              # diagnose it like a real incident
spar hint                       # up to two hints if you get stuck
spar ready                      # verifies the fix and generates the report
```

## Documentation

- [Design decisions (ADRs, in Spanish)](docs/adr/)
