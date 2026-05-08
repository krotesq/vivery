# Rules

## Error Handling

- Repository-Methoden geben bei Fehler `nil, err` zurück, nicht ein leeres Struct.
- `cmdTag.RowsAffected()` darf erst nach einem `err != nil`-Check aufgerufen werden.
- Handler geben keine internen Fehlermeldungen an den Client weiter. Jeder Handler definiert eigene, generische Error-Messages.
- Fehler aus DB-Operationen dürfen nicht mit `_` ignoriert werden. Mindestens loggen.
- Errors werden mit `fmt.Errorf("context: %w", err)` gewrappt, um Kontext zu erhalten und `errors.Is`/`errors.As` zu ermöglichen.
- Sentinel Errors (z.B. `ErrNotFound`, `ErrConflict`) verwenden, um im Handler zwischen Client-Fehler (4xx) und Server-Fehler (5xx) zu unterscheiden.
- Panic nur in unrecoverable Situationen. Für erwartbare Fehler immer `error` zurückgeben.

## Security

- Cookies setzen immer `HttpOnly: true` und `Secure: true`. `Secure` wird per Environment konfiguriert.
- Cookies setzen `MaxAge` passend zur Token-Lifetime.
- Authentifizierungs-Endpoints haben Rate Limiting.
- Passwörter werden vor dem Hashen validiert: min. 8 Zeichen, min. 1 Zahl, min. 1 Sonderzeichen, darf den Username nicht enthalten.
- Alle Login-Fehler geben dieselbe generische Fehlermeldung zurück. Kein Unterschied zwischen "User nicht gefunden", "deaktiviert" oder "falsches Passwort".
- JWTs enthalten `iat` (issued at) und `jti` (unique token ID) Claims.
- Request Bodies haben ein Size Limit (z.B. 1 MB via `io.LimitReader`).
- HTTP Server definiert `ReadTimeout`, `WriteTimeout` und `IdleTimeout`.
- `X-Forwarded-For` wird nicht blind vertraut. Nur bei bekanntem Reverse Proxy auswerten, sonst `RemoteAddr` verwenden.
- URLs aus User-Input werden vor dem Speichern auf interne Adressen geprüft (SSRF-Schutz).
- Account-Lockout ist zeitbasiert. Gesperrte Accounts werden nach einer definierten Dauer automatisch entsperrt.
- `sslmode` für die DB-Verbindung ist konfigurierbar. In Produktion `require` oder `verify-full`.
- DSN-Komponenten (User, Passwort, DB-Name) werden URL-encoded.
- Secrets (DB-Passwörter, JWT-Secrets, API-Keys) werden nie hardcoded. Immer über Environment-Variablen oder Secret-Manager laden.
- CORS-Origins werden per Environment konfiguriert, nicht hardcoded.
- State-ändernde Endpoints (POST, PUT, PATCH, DELETE) sind bei Cookie-basierter Auth gegen CSRF geschützt (z.B. `SameSite=Strict` oder CSRF-Token).
- SQL-Queries verwenden immer Parameterized Queries (`$1`, `$2`). Kein String-Concatenation für Queries.
- User-Input wird validiert bevor er verarbeitet wird: Typ, Länge, Format, erlaubte Zeichen.
- Sensible Daten (Credentials, Keys, Tokens) werden in API-Responses nicht exponiert oder nur maskiert zurückgegeben.
- Security-relevante Aktionen (Login, fehlgeschlagene Logins, Passwortänderungen, Account-Sperrungen) werden geloggt.

## HTTP / API

- Korrekte HTTP-Statuscodes verwenden: `201` für Erstellung, `204` für Löschung ohne Body, `400` für Client-Fehler, `500` für Server-Fehler, `409` für Konflikte.
- Responses setzen immer `Content-Type: application/json`.
- Pointer in Request-DTOs vor Zugriff auf `nil` prüfen, um Panics zu vermeiden.
- Path-Parameter (z.B. IDs) werden vor der Weitergabe an den Service validiert (Format, Länge).
- Listen-Endpoints unterstützen Pagination. Nie unbegrenzt viele Datensätze zurückgeben.
- Der Server fährt graceful herunter: laufende Requests werden abgeschlossen, neue abgelehnt.

## Datenbank

- `SELECT *` vermeiden. Explizite Spaltenlisten verwenden, um bei Schema-Änderungen nicht zu brechen.
- Transaktionen verwenden wenn mehrere zusammenhängende Writes stattfinden.
- Connection Pools konfigurieren: `MaxConns`, `MinConns`, `MaxConnLifetime`.
- `context.Context` an alle DB-Aufrufe durchreichen, damit Timeouts und Cancellation greifen.

## Go Conventions

- Receiver-Namen sind kurz: ein oder zwei Buchstaben (z.B. `h`, `s`, `r`), nicht der volle Typname.
- Variablen-Namen dürfen nicht den Typnamen shadwen (z.B. nicht `var dto dto`).
- Parameter verwenden camelCase, kein snake_case.
- Struct-Namen sind nicht identisch mit dem Package-Namen (nicht `source.source`).
- Keine Leading-Underscores für Typen (nicht `_type`). Stattdessen beschreibende Namen verwenden.
- Funktionen mit mehr als 4 Parametern verwenden ein Config- oder Options-Struct.
- `time.Now()` nicht direkt in Repositories aufrufen. Für Testbarkeit als Parameter übergeben oder eine Clock-Abstraktion verwenden.
- Unused Parameters und Imports entfernen.
- Structured Logging verwenden (`slog` oder vergleichbar), nicht `log.Println` mit freien Strings.

## Serialisierung

- Leere Slices werden als `[]` serialisiert, nicht als `null`. Slices mit `make([]T, 0)` initialisieren.
- JSON-Tags verwenden einheitlich camelCase.
