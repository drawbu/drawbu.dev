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
          tailwindcss -i static/style.css -o static/generated.css && \
          templ generate                                          && \
          go build                                                && \
          ./app
        '';
        rundev = pkgs.writeShellScriptBin "rundev" ''
          if [ -z "$1" ]; then
            ${devbuild}/bin/devbuild
          elif [ "$1" = "--watch" ]; then
            ${pkgs.fd}/bin/fd | ${pkgs.entr}/bin/entr -c -r ${devbuild}/bin/devbuild
          else
            echo "Usage: $0 [--watch]"
          fi
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
            tags = ["production"];
            nativeBuildInputs = with pkgs; [templ tailwindcss makeWrapper];
            preBuild = ''
              templ generate
              tailwindcss -i static/style.css -o static/generated.css
            '';
          };

          docker = pkgs.dockerTools.buildImage {
            name = "drawbu.dev";
            tag = "latest";
            created = "now";
            copyToRoot = pkgs.buildEnv {
              name = "image-root";
              paths = [packages.app];
            };
            config.Cmd = ["app"];
          };
        };

        defaultPackage = packages.app;
      }
    );
}
