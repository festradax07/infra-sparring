# ADR 0002 — spar llama al LLM con su propio loop

Estado: aceptado · 2026-10-12

## Contexto

Había dos formas de integrar un LLM:

- **A.** spar llama a la API del provider y el jugador solo usa la CLI.
- **B.** El jugador usa su propio agente (Claude Code, Codex, Cursor) conectado a spar por MCP.

B evita pagar una API key aparte, pero el agente con el que habla el jugador es el mismo que
generó la falla: conoce la solución y la puede revelar ante cualquier pregunta cercana.

## Decisión

A en el MVP. El provider se abstrae detrás de una interfaz chica, con un adaptador para
Anthropic primero y uno compatible con la API de OpenAI después (cubre OpenAI, Gemini y Ollama).

B queda para una fase posterior y solo para usos donde ocultar no importa: armar packs nuevos
o analizar informes.

## Consecuencias

- El jugador necesita una API key propia (ver ADR 0003).
- El costo por desafío es visible y se registra (RNF-09).
- La opacidad del desafío la garantiza el diseño, no la disciplina del jugador.
