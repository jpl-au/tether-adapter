// Package wails runs tether pages inside a Wails v3 window over the
// window's own stream, with no network port. The page's connection
// goes to Go and back through Wails, as a WebSocket would over the
// network.
//
// Three things set it up. Serve the tether handler as the app's asset
// handler with its Mode set to mode.IPC, put Script in the handler's
// Layout, and hand the handler to Serve:
//
//	h := tether.Stateful(tether.App{}, tether.StatefulConfig[State]{
//	    Mode:   mode.IPC,
//	    Layout: func(s State, content node.Node) node.Node {
//	        return html.New(head.New(wails.Script()), body.New(content))
//	    },
//	    // ...
//	})
//	app := application.New(application.Options{
//	    Assets: application.AssetOptions{Handler: h},
//	})
//	wails.Serve(app, h)
package wails

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/jpl-au/fluent/html5/script"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// stream is the name of the Wails stream every page connects over. The
// page side in ipc.js opens it by the same name.
const stream = "tether"

// ipc is the page side of Serve. See Script.
//
//go:embed ipc.js
var ipc string

// Script gives the page its socket: the Wails stream, wrapped so it
// behaves as the WebSocket tether expects. Put it in the head of the
// Layout of every page served with mode.IPC. It is a module, so it has
// run by the time the page connects.
func Script() node.Node {
	return script.InlineModule(ipc)
}

// TransportServer is a tether handler, which runs a connection the host
// hands it. Every [tether.Handler] is one.
type TransportServer interface {
	ServeTransport(r *http.Request, t tether.Transport)
}

// hello is the first frame a page sends: what a WebSocket would carry
// in its request, the URL with the page's one-time ticket and the
// page's User-Agent, which the ticket is bound to.
type hello struct {
	URL       string `json:"url"`
	UserAgent string `json:"userAgent"`
}

// Serve hands every connection a page opens over the app's tether
// stream to h, until the page goes away or the app stops.
func Serve(app *application.App, h TransportServer) {
	app.HandleStream(stream, func(c *application.StreamConn) {
		defer c.Close()
		frame, err := c.Receive()
		if err != nil {
			return
		}
		var first hello
		if err := json.Unmarshal(frame, &first); err != nil {
			slog.Warn("tether-adapter/wails: the page's first frame did not name its connection", "error", err)
			return
		}
		r, err := http.NewRequestWithContext(c.Context(), http.MethodGet, first.URL, nil)
		if err != nil {
			slog.Warn("tether-adapter/wails: the page asked to connect to an unusable URL", "url", first.URL, "error", err)
			return
		}
		r.Header.Set("User-Agent", first.UserAgent)
		h.ServeTransport(r, transport{conn: c})
	})
}

// transport is a [tether.Transport] over one Wails stream connection.
type transport struct {
	conn *application.StreamConn
}

// Send queues data for the page. Wails keeps the slice it is given
// until the frame is delivered, so it is handed a copy tether can never
// reuse.
func (t transport) Send(data []byte) error {
	return t.conn.Send(bytes.Clone(data))
}

// ReceiveEvent reads the next event the page sent. A closed stream is
// the page going away or the app stopping, which ends the connection
// cleanly.
func (t transport) ReceiveEvent() (tether.Event, error) {
	frame, err := t.conn.Receive()
	if errors.Is(err, application.ErrStreamClosed) {
		return tether.Event{}, io.EOF
	}
	if err != nil {
		return tether.Event{}, err
	}
	var ev tether.Event
	if err := json.Unmarshal(frame, &ev); err != nil {
		return tether.Event{}, err
	}
	return ev, nil
}

// Close ends the connection. Wails makes it safe to call more than once
// and from any goroutine.
func (t transport) Close() error {
	return t.conn.Close()
}
