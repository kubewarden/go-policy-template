package main

// This file contains all the data structures that need
// to serialized and deserialized to JSON.
//
// For each type we need to generate the code that implements
// the easyjson interfaces.
//
// Every time this file changes, the Makefile will regenerate the `types_easyjson.go`
// file.
//
// Note: the `easyjson` cli cannot process files that are inside of the `main`
// package. This is a known limitation. The Makefile target from above has a
// workaround for that.
//
// Important: limit the number of imports inside of this file. Also, don't
// try to use interface types (or anything making use of them). This isn't
// going to play out well with TinyGo at **runtime** due to its limited
// support of Go reflection.

type Settings struct {
	DeniedNames []string `json:"denied_names"`
}
