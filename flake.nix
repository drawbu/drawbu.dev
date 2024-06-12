{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = {self, ...} @ inputs:
    inputs.utils.lib.eachDefaultSystem (
      system: let
        pkgs = import inputs.nixpkgs {inherit system;};

        devbuild = pkgs.writeShellScriptBin "devbuild" ''
          rm -rf /tmp/drawbu.dev                                                && \
          cp -r static /tmp/drawbu.dev                                          && \
          tailwindcss -i /tmp/drawbu.dev/style.css -o /tmp/drawbu.dev/style.css && \
          templ generate                                                        && \
          go build -ldflags="-X 'main.staticDir=/tmp/drawbu.dev'"               && \
          ./app
        '';
        rundev = pkgs.writeShellScriptBin "rundev" ''
          if [ -z "$1" ]; then
            exec ${devbuild}/bin/devbuild
          fi
          if [ "$1" = "--watch" ]; then
            exec ${pkgs.fd}/bin/fd | ${pkgs.entr}/bin/entr -c -r ${devbuild}/bin/devbuild
          fi
          echo "Usage: $0 [--watch]"
        '';
      in rec {
        formatter = pkgs.alejandra;

        devShell = pkgs.mkShell {
          inputsFrom = builtins.attrValues self.packages.${system};
          packages = [rundev];
        };

        packages = {
          app = pkgs.buildGoModule {
            name = "app";
            src = ./.;
            vendorHash = null;
            ldflags = [
              "-X main.staticDir=${placeholder "out"}/share/static"
              "-X main.articlesDir=${placeholder "out"}/share/articles"
              "-X main.prod=true"
            ];
            nativeBuildInputs = with pkgs; [templ tailwindcss makeWrapper];
            preBuild = ''
              templ generate
            '';
            postBuild = ''
              mkdir -p $out/share
              cp -r articles $out/share/articles
              cp -r static $out/share/static
              tailwindcss -i $out/share/static/style.css -o $out/share/static/style.css
            '';
          };

          docker = pkgs.dockerTools.buildImage {
            name = "drawbu.dev";
            tag = "latest";
            created = "now";
            copyToRoot = pkgs.buildEnv {
              name = "image-root";
              paths = [packages.app] ++ (with pkgs.dockerTools; [caCertificates]);
            };
            config.Cmd = ["app"];
          };
        };

        defaultPackage = packages.app;
      }
    );
}
