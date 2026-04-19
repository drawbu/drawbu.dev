# drawbu.dev

Hello there! :)

This is the repository for my personal website, [drawbu.dev](https://drawbu.dev).
This is a simple project that I made to have some fun and showcase my projects
and "articles".

This is a standard Astro blog, powered by htmx+preload to mimic SPA behavior.
Every pages is prerendered, so this is just incremental improvement! The website
works the same without JS enabled.


## Build

This is a `pnpm` project, so I highly recommend you to use `corepack`. 
I encourage you to use the Nix flake if you have the Nix package manager
installed, as it should provide everything you need. Anyway, the environment
is very simple.

```bash
pnpm app:build
# or
cd app
pnpm build
```
