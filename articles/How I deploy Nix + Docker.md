---
title: How I deploy Nix + Docker
date: 2024-06-09
uri: deploy-nix-docker
author:
  name: drawbu
  email: contact@drawbu.dev
---

Hello everyone :)

I recently remade my 'portfolio' website for the fifth time and since I use and Nix as a package manager and to declare my developement environnements, it only made sense to me to avec a Nix flake at the root of my project. For the last year, it has become a habit that I have on every project I have so I never again have a problem with the dependencies I used in a old project.

I used Go + Templ to do a simple website and once I have been done I question how to deploy such a project. I'm a Nix advocate. I really like to have the possibility to define a declarative development environment with all the dependencies I need to build and dev a project. And of course curiosity hehe.

## Using nix build

So currently I have a Nix flake with a devShell that exports `go`, `templ` and `tailwind`.  Great.

```nix
devShell = pkgs.mkShell {
    packages = with pkgs; [go templ tailwindcss];
};
```

How do I build my project?
That's a tricky part. I several steps to achieve the build of the binary.

```sh
templ generate
mkdir -p /tmp/drawbu.dev
tailwindcss -i ./assets/style.css -o /tmp/drawbu.dev/style.css
go build -ldflags="-X 'main.assetsDir=/tmp/drawbu.dev'"
```

So...  I need to make a script to build process easier. Let's use our Nix flake!
We are gonna do a package named `app` so I can just `nix build` and all our build steps will be executed (with the advantages that has nix aka. declarative, reproducible & reliable builds).

Will just do a simple Go package and see were it goes.

```nix
app = pkgs.buildGoModule {
  name = "app";
  src = ./.;
  vendorHash = null;
};
```


If you are interested, here is the finished version:

```nix
app = pkgs.buildGoModule {
  name = "app";
  src = ./.;
  vendorHash = null;
  ldflags = ["-X main.assetsDir=${placeholder "out"}/share/assets"];
  nativeBuildInputs = with pkgs; [templ tailwindcss makeWrapper];
  preBuild = ''
    templ generate
  '';
  postBuild = ''
    mkdir -p $out/share/assets
    tailwindcss -i ./assets/style.css -o $out/share/assets/style.css
  '';
  postInstall = ''
    wrapProgram $out/bin/app \
      --prefix PATH : ${pkgs.lib.makeBinPath (with pkgs; [ git ])}
  '';
};
```

## Docker

So go is a compiled language and so I can make good use of that feature and just deploy the binary, right? Yeah... no. You see in modern day world, you almost never see just a random binary running your entire website on a disant VPS. Instead, most of the time we try to define a Docker image. It brings simplicity to the deployment, because what works on docker on your machine, should usually works the same on the production machine (we ain't putting your laptop in prod).

You need to isolate your server from the rest of the machine, because if something goes bad, it only takes place in the container, and does not propagate to the other services and the rest of the network.

I made some research and came around [this article](https://mitchellh.com/writing/nix-with-dockerfiles) that creates a Dockerfile where we use nix build to create binary a nix store, and create a development environment without the build time dependencies and only exports our app and it's run time dependencies.

That could work, I thought, but I still looked around and found a page on [nix.dev](https://nix.dev/tutorials/nixos/building-and-running-docker-images.html) that explain that we can declaratively declare Docker image with Nix! And ouptut them as a package!

You I made a second package exported by my flake called `docker` that uses `pkgs.dockerTools.buildImage`.

```nix
docker = pkgs.dockerTools.buildImage {
  name = "drawbu.dev";
  tag = "latest";
  copyToRoot = pkgs.buildEnv {
    name = "image-root";
    paths = [ self.packages.${system}.app ];
  };
  config.Cmd = ["app"];
};
```

Then I ran `nix build .#docker`, and... voilà after a few seconds nix just ouputed me a docker image with the completed build of my app with only its run time dependencies!

I then make a [GitHub workflow](https://github.com/drawbu/drawbu.dev/blob/main/.github/workflows/docker.yml) to build the image each time I push to the repo and publish it on ghcr.io.


## Finally deploying

So now I just have to make a ultra simple `docker-compose.yml` file that uses my pre-build image and I'll have everything I need!

```yml
services:
  drawbu.dev:
    image: ghcr.io/drawbu/drawbu.dev:latest
    ports:
        - "8080:8080"
```

That's it!

Thanks for reading me. This is one of the first time I do this, thanks you for your time, and I hope I see you around soon! :)
