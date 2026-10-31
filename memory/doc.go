// Package memory provides in-process implementations of
// [tether.SessionStore] and [tether.DiffStore] backed by a plain
// map. Data lives in memory only and does not survive process
// restarts.
//
// Use this for testing, development, and short-lived deployments
// where persistence beyond the process lifetime is not needed.
//
// Usage:
//
//	import "github.com/jpl-au/tether-adapter/memory"
//
//	tether.Stateful(app, tether.StatefulConfig[State]{
//	    SessionStore: memory.NewSessionStore(),
//	    DiffStore:    memory.NewDiffStore(),
//	})
package memory
