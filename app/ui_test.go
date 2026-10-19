package main

import "github.com/pocketbase/pocketbase/ui"

// apis.NewRouter binds a superuser UI extensions OnServe hook on every call.
// The request helpers in this package build a router per request against a
// shared test app, so the second request would register GET /_/extensions.js
// twice and panic. The tests never use the superuser UI, so leave it
// unbundled, the same as building with the no_ui tag.
func init() {
	ui.DistDirFS = nil
}
