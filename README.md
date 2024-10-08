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
nix build # Or use `nix run` directly
```


## Development

You have all the dependencies in the nix flake, so you can just run the
`nix develop` command and you will have everything you need. You will have
also available the `nix run .#dev` command for fast iteration on the project.

The project depends on Go (with a few deps), Tailwind and Templ.


## Docker

You can use `docker-compose` to run the project using pre-built images:
```docker-compose
services:
  drawbu.dev:
    image: ghcr.io/drawbu/drawbu.dev:latest
    ports:
        - "8080:8080"
```

I don't have a real Dockerfile, so if you want to build the image, you'll need
to use Nix, or just use the one built from GitHub actions.
```bash
nix build .#docker
docker load < result
```
```bash
docker run --run -p 8080:8080 drawbu.dev
```
