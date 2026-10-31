// Package fs provides filesystem-backed implementations of
// [tether.SessionStore] and [tether.DiffStore]. Each session is
// stored as a single file named by its ID - simple, dependency-free,
// and easy to inspect during development.
//
// Suitable for development, single-node deployments, and low-traffic
// applications. For production deployments with multiple nodes, use a
// shared store (Redis, SQLite, PostgreSQL) instead.
//
// Usage:
//
//	import "github.com/jpl-au/tether-adapter/fs"
//
//	tether.Stateful(app, tether.StatefulConfig[State]{
//	    SessionStore: fs.NewSessionStore("tmp/sessions"),
//	    DiffStore:    fs.NewDiffStore("tmp/diffs"),
//	})
package fs
