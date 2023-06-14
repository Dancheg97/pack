<h2 align="center">Pack registry</h2>

This library provides interfaces required to embed pack registry to foreign projects written in go. You have to implement 2 interfaces, to get `http.Handler` which will be able to function with your system:

- Source from emails to GPG keys
- Database former, that will be able to add apckages to resulting pacman database

By default, pack uses a directory with GPG keys, also with nested child directories, to validate pushed packages, before adding them to database, but it default behaviour can be adjusted for any database/web paradigm.

Example of usage:

```go

import "fmnx.su/core/pack/registry"

func main() {
	d := local.DirStorage{
		Dir: "/dir/with/packages",
	}

	k := local.LocalKeyDir{
		Dir: "/dir/with/gpg/keys",
	}

	r := registry.Registry{
		TmpDir:      "/tmp",
		Dbname:      "domain.com",
		FileStorage: &d,
		KeyReader:   &k,
	}

	router := mux.NewRouter()

	router.HandleFunc(p.Endpoint+"/push", r.Push)
	router.HandleFunc(p.Endpoint+"/{owner}/{file}", r.Get)
	router.HandleFunc(p.Endpoint+"/{file}", r.Get)

	msg := fmt.Sprintf("Starting registry %s on port %s", p.Name, p.Port)

	tmpl.Amsg(p.Stdout, msg)

	return http.ListenAndServe(":"+p.Port, router)
}

```
