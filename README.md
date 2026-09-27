# WebCacher 2

Proxy HTTP/HTTPS con caché en disco, cola de peticiones offline y estadísticas.
Reescritura en Go de WebCacher 1, manteniendo **compatibilidad con su carpeta `.cache`**:
los archivos se nombran con el mismo algoritmo, así que una caché creada por la
versión 1 se puede reutilizar tal cual.


## Cómo funciona

```
                                            no
Cliente ──▶ goproxy (:8092) ──> En cache? ────▶ OnRequest ───────▶ ¿hay Internet?
                                   │                                       │    no
                                   |                                       │────────▶ queue.GQueue ──▶ workers ──▶ se reintentan al volver la red
                                   │ Si                                    │
                                   |                                    si |
                                   ▼                                       ▼
                            cache.Global.Pop                            red real ──▶ OnResponse ──▶ cache.Global.Push
                                   |                                                    |
                                   |                                                    |
                                   |                                                    ▼
                                   ────────────────────────────────────────────────> browser

```

1. **MITM**: `ConnectHandler` fuerza `MitmConnect` para poder interceptar HTTPS
   (necesario para cachear por URL, no por host).
2. **OnRequest**: si la petición es `GET` y es cacheable, se busca en la caché. Si
   hay acierto se devuelve la respuesta con la cabecera `webcacher: true`. Si no hay
   acierto y no hay internet, la petición se encola.
3. **OnResponse**: se guarda en disco todo `GET` con status 200. Las respuestas que
   vienen de caché (marcadas con `webcacher: true`) solo cuentan estadísticas.
4. **Workers**: en bucle constante, sacan elementos de la cola y los reintentan
   apuntando de vuelta al propio proxy (cabecera `webcacher-queue: true` para no
   reencolar en bucle). Cada 10 s se guardan cola, caché y stats en disco.

## Estructura

| Ruta | Contenido |
| --- | --- |
| `proxy.go` | `main`: carga config, stats, cola y caché, y arranca el proxy |
| `proxy/responses.go` | Handlers goproxy (`OnRequest`, `OnResponse`, MITM) |
| `proxy/works.go` | Workers de la cola: reintento de peticiones offline |
| `proxy/stats.go` | Contadores y de historial, persistidos en `stats.json` |
| `proxy/utils.go` | Chequeo de internet, tamaño de caché, helper de URL |
| `cache/cache.go` | Caché en disco: `Push`, `Pop`, carga/guardado del índice |
| `queue/queue.go` | Cola doblemente enlazada + persistencia en `queue.json` |
| `config/` | Lectura de `webcacher.conf` y flags de CLI |
| `urlutils/` | Clave de caché, extensión y candidatos de búsqueda |
| `debug/` | Log a `/tmp/webacher.log` |
| `public/proxy-ca.crt` | CA de goproxy, necesaria para MITM |

## Formato de la caché

Cada respuesta se guarda con `httputil.DumpResponse` en:

```
.cache/<sha256(METHOD + URL)>[.<ext>].phttp
```

- `ext` sale de `URL.Path` (no de la URL completa) y se descarta si mide más de 10
  caracteres, replicando la regla de WebCacher 1.
- `Pop` prueba varios candidatos en orden de prioridad (clave exacta, sin extensión,
  clave completa con args, y las variantes "legacy" que usaba WebCacher 2 antes de la
  compatibilidad), para que la caché de la v1 se resuelva siempre.
- El índice (`cache.json`) guarda `Counts` (hits por clave), `Hashes` y `MemKeys`
  (ruta de archivo → clave).

## Configuración

`webcacher.conf` es un archivo de listas, una por sección:

```
CacheArgs:          # args que sí se incluyen en la clave (ej: t=...)
NoCacheExt:         # extensiones nunca cacheadas
Syncs:              # endpoints de sincronización
Pproxy:             # proxies padre
IgnoreQueue:        # hosts que no se encolan offline
NoCacheSite:        # hosts nunca cacheados
NoArgs:             # hosts cuya query string se descarta de la clave
```

Flags de CLI (vía `config.ParseArgs`):

| Flag | Efecto |
| --- | --- |
| `--no-args` | Activa el modo de descarte de argumentos (igual que la v1) |
| `--no-all` | Solo cachea archivos con extensión reconocida |
| `--no-cache` | Ignora la caché en todas las peticiones |
| `--no-queue` | Desactiva la cola offline (y sus workers) |
| `--no-mem-cache` | Solo usa caché en disco |
| `--help` | Muestra la ayuda |

## Uso

```bash
go build -o webcacher .
./webcacher [--no-args] [--no-queue] ...

# configurar el proxy del sistema hacia 127.0.0.1:8092
# y confiar el certificado de public/proxy-ca.crt en el almacén del sistema
```

En `main` el orden de arranque es: internet checker → config → stats → tamaño de
caché → cola → caché → proxy.

## Estado / pendientes

Ver `TODO.md`: sustitución de argumentos en la URL (ej. `t=12` → timestamp) y una
GUI web con gráficas para las estadísticas.

Otros puntos abiertos en el código:

- `stats.json` se reescribe completo cada 10 s sin escritura atómica.
- `Work` reintenta la petición pero no propaga ni registra errores de `client.Do`.
- `proxy/utils.go` inicializa `Internet` con un `DialTimeout` síncrono en el
  arranque, antes de que arranque el checker en background.
- El índice `MemKeys` se escribe en disco en cada `Save` completo de `cache.json`.
