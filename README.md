# yze-go-errlast

A [`yze`](https://github.com/gomatic/yze) analyzer (category `errors`) enforcing the gomatic Go idiom that `error` is always the **last** return value. The convention is a contract on any signature returning an error, so it is checked on function and method declarations, interface method signatures, function literals, and function-typed definitions alike.

- **Rule:** `yze/errlast`
- **Library:** exports `Analyzer` and `Registration` for the [`yze`](https://github.com/gomatic/yze) aggregator and [`stickler`](https://github.com/gomatic/stickler) runner.
- **Binary:** `cmd/yze-go-errlast` runs it standalone (`text`/`-json`, and as a `go vet -vettool`).

Built on the [`go-yze`](https://github.com/gomatic/go-yze) framework.
