# wails

Runs [Tether](https://github.com/jpl-au/tether) pages inside a [Wails](https://wails.io) v3 window with `mode.IPC`. The page's connection goes to Go and back over the window's own Wails stream, as a WebSocket would over the network, so the app opens no port.

## Install

```bash
go get github.com/jpl-au/tether-adapter/wails
```

## Usage

Three things set it up: the handler's `Mode` is `mode.IPC`, its Layout carries `wails.Script()`, and the handler is both the app's asset handler and the one `wails.Serve` hands connections to.

```go
import (
    "github.com/jpl-au/tether"
    "github.com/jpl-au/tether-adapter/wails"
    "github.com/jpl-au/tether/mode"
    "github.com/wailsapp/wails/v3/pkg/application"
)

h := tether.Stateful(tether.App{}, tether.StatefulConfig[State]{
    Mode: mode.IPC,
    Layout: func(s State, content node.Node) node.Node {
        return html.New(head.New(wails.Script()), body.New(content))
    },
    // ...
})

app := application.New(application.Options{
    Assets: application.AssetOptions{Handler: h},
})
wails.Serve(app, h)
app.Window.NewWithOptions(application.WebviewWindowOptions{URL: "/"})
app.Run()
```

## How it works

- The page itself, its scripts and the connect ticket are ordinary requests, served by the handler through Wails' asset server.
- `wails.Script()` gives the page `window.Tether.ipc`: the Wails stream named `tether`, wrapped so it behaves as the WebSocket tether expects. Its first frame carries what a WebSocket's request would: the URL with the page's ticket, and the page's User-Agent.
- `wails.Serve` handles that stream. For each connection it reads the first frame and hands the connection to the handler's `ServeTransport`, where events and updates travel as they would over a WebSocket.
- Wails delivers every frame to the page as bytes. The page side turns them into text, so a binary wire format reaches tether as base64 text, as it does over SSE.

The connect request carries the page's URL and User-Agent, not its cookies, so an `Exclusive` key read from a cookie does not apply.

## Licence

MIT
