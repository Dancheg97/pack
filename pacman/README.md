<h2 align="center">Go wrapper for arch package manager</h2>

This library aims to provide concurrent, stable and extensible interface to interact with arch package manager - pacman. Library has 0 dependencies and written in pure go with only a few packages from stdlib.

Some of the default options for functions can contain sudo, if you don't need it you can manually disable it.

Functions:

- `Sync` - syncronize packages

```go
import "fmnx.su/dancheg97/pacman"

func main() {
	err := pacman.Sync("nano")
	fmt.Println(err)
}
```

- `Search` - search for packages in pacman databases

```go
import "fmnx.su/dancheg97/pacman"

func main() {
	r, err := pacman.Search("vim")
	fmt.Println(r)
	fmt.Println(err)
}
```

- `Upgrade` - install packages from local files

```go
import "fmnx.su/dancheg97/pacman"

func main() {
	err := pacman.Upgrade("nvim-1-1-any.pkg.tar.zst")
	fmt.Println(err)
}
```

- `Query` - list installed packages

```go
import "fmnx.su/dancheg97/pacman"

func main() {
	r, err := pacman.Query()
	fmt.Println(r)
	fmt.Println(err)
}
```

- `Makepkg` - build package

```go
import "fmnx.su/dancheg97/pacman"

func main() {
	err := pacman.Makepkg()
	fmt.Println(err)
}
```

- `Remove` - remove installed packages

```go
import "fmnx.su/dancheg97/pacman"

func main() {
	err := pacman.Remove("emacs")
	fmt.Println(err)
}
```

- `RepoAdd` - add package to local repository

```go
import "fmnx.su/dancheg97/pacman"

func main() {
	err := pacman.RepoAdd("localhost.db.tar.gz", "vscodium-1-1-x86_64.pkg.tar.zst")
	fmt.Println(err)
}
```
