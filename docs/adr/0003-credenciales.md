# ADR 0003 — Credenciales del provider fuera de los archivos de configuración

Estado: aceptado · 2026-10-12

## Contexto

`spar.yaml` vive dentro de repos del usuario (por ejemplo, `lfcs.d`) y termina en GitHub.
Una API key escrita ahí se filtra.

## Decisión

La configuración solo nombra la variable de entorno (`llm.api_key_env`). La key se busca en
este orden:

1. Esa variable de entorno.
2. El llavero del sistema (Secret Service en Linux, Keychain en macOS), donde la guarda
   `spar auth login` después de validarla con una llamada mínima.
3. `~/.config/spar/credentials`, solo con permisos `600`. Con otros permisos spar
   se niega a leerlo, como `ssh` con las claves privadas.

`spar auth status` muestra de qué fuente sale la key; `spar auth logout` la borra del llavero.

## Consecuencias

- Dependencia de `zalando/go-keyring` y de que haya un llavero disponible; en sesiones sin
  escritorio se usa el archivo.
- Los logs y errores redactan cualquier valor que parezca una key.
- `spar auth login` lee la key sin eco en la terminal.
