<h2 align="center">Pack server</h2>

This library provides interfaces required to embed pack registry to foreign projects written in go. You have to implement 2 interfaces, to get `http.Handler` which will be able to function with your system:

- Source from emails to GPG keys
- Database former, that will be able to add apckages to resulting pacman database

By default, pack uses a directory with GPG keys, also with nested child directories, to validate pushed packages, before adding them to database, but it default behaviour can be adjusted for any database/web paradigm.

Example of usage:

```go

import "fmnx.su/core/pack/server"

func main() {
	d := server.LocalDirDb{
		Dir:    "/var/cache/pacman/pkg",
		DbName: "localhost",
	}

	k := server.LocalGpgDir{
		GpgDir: "/home/user/gpg",
	}

	s := server.Pusher{
		Stdout:          os.Stdout,
		Stderr:          os.Stderr,
		GPGVireivicator: &k,
		DbFormer:        &d,
	}

	fs := http.FileServer(http.Dir(p.Dir))
	http.Handle("/api/pack", http.StripPrefix("/api/pack", fs))
	http.HandleFunc("/api/pack/push", s.Push)

	return http.ListenAndServe(":"+p.Port, http.DefaultServeMux)
}

```
