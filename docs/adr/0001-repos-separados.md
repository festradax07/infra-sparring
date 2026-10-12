# ADR 0001 — spar es un repositorio aparte de lfcs.d

Estado: aceptado · 2026-10-12

## Contexto

spar nació dentro de `lfcs.d`, el laboratorio personal para preparar la LFCS. Pero
spar apunta a varias certificaciones y a que otras personas puedan usarlo, mientras que
`lfcs.d` es el recorrido de una persona: sus soluciones, su configuración y sus informes.

## Decisión

Dos repositorios. `infra-sparring` (binario `spar`) contiene la herramienta y los packs, y publica releases con el
binario. `lfcs.d` consume ese binario como cualquier usuario y guarda `spar.yaml`,
`reports/` y las soluciones escritas a mano.

## Consecuencias

- spar no puede depender de rutas ni archivos de `lfcs.d`. Lo que necesite del pack
  (temario, ejemplos, allowlist) vive en `packs/` dentro de este repo.
- Los escenarios `ec-*` de `lfcs.d` se copian a `packs/lfcs/examples/` como ejemplos para el
  generador; no se referencian entre repos.
- Hace falta versionar el contrato del desafío, porque `lfcs.d` puede quedar con una versión
  vieja del binario.
