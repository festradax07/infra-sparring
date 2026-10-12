# ADR 0004 — El LLM no ejecuta: devuelve datos y el pipeline ejecuta

Estado: aceptado · 2026-10-12

## Contexto

El diseño inicial le daba al agente herramientas para actuar sobre el laboratorio (`exec`,
`snapshot`, `restore`, `check`) expuestas como server MCP. Con el loop propio (ADR 0002) eso
no hace falta: generar un desafío no requiere explorar la VM, sino producir artefactos.

## Decisión

El LLM devuelve el desafío como datos estructurados: síntoma, `setup.sh`, `check.sh`,
`solution.sh`, evidencia y pistas. El pipeline de validación, que es código determinista,
analiza esos artefactos y los ejecuta en las VMs. Si algo falla, el error vuelve al LLM como
contexto para regenerar.

El LLM nunca recibe herramientas que ejecuten comandos.

## Consecuencias

- La superficie de riesgo baja: nada de lo que el LLM decida se ejecuta sin pasar antes por la
  allowlist y la verificación de recursos protegidos.
- El MVP no necesita un server MCP. Si se implementa la integración B del ADR 0002, el server
  MCP será una capa delgada sobre `internal/lab` y `internal/validate`.
- La revisión del informe también es datos: el LLM recibe la traza como texto y devuelve el análisis.
- El documento de diseño original muestra un "server MCP del lab"; este ADR lo reemplaza.
