// +build !prod

package main

import "net/http"

var assets http.FileSystem = http.Dir("assets")
