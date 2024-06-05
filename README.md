# drawbu.dev

Hello there! :)

This is the repository for my personal website, [drawbu.dev](https://drawbu.dev).
This is a simple project that I made to have some fun and showcase my projects
and "articles".


## Build

So this is a Go project, using tailwind and templ. I highly recommend you to
check out the nix flake that I have for this project, it makes everything much
easy to build and run.

```bash
# Build the project (using nix)
nix build

# Build the project (using go)
templ generate
mkdir -p /tmp/drawbu.dev
tailwindcss -i ./assets/style.css -o /tmp/drawbu.dev/style.css
go build -ldflags="-X 'main.assetsDir=/tmp/drawbu.dev'"
```


## Development

You have all the dependencies in the nix flake, so you can just run the
`nix develop` command and you will have everything you need. You will have
also available the `rundev` command for fast iteration on the project.

The project depends on Go (with a few mods), Tailwind and Templ.


## Docker

Yeah you have no way to avoid Nix hehe

```bash
nix build .#docker
docker load < result
```
```bash
docker run --run -p 8080:8080 drawbu.dev
```
