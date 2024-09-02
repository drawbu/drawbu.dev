{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { self, ... }@inputs:
    inputs.utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import inputs.nixpkgs { inherit system; };

        rundev = pkgs.writeShellApplication {
          name = "rundev";
          runtimeInputs = [
            pkgs.fd
            pkgs.entr
            (pkgs.writeShellApplication {
              name = "buildapp";
              text = ''
                tailwindcss -i static/style.css -o static/generated.css && \
                templ generate                                          && \
                go build                                                && \
                ./app
              '';
            })
          ];
          text = ''
            if [ -z "$1" ]; then
              buildapp
            elif [ "$1" = "--watch" ]; then
              fd | entr -c -r buildapp
            else
              echo "Usage: $0 [--watch]"
            fi
          '';
        };

      in
      rec {
        formatter = pkgs.alejandra;

        devShell = pkgs.mkShell {
          inputsFrom = builtins.attrValues self.packages.${system};
          packages = [ rundev ];
        };

        packages = {
          app = pkgs.buildGoModule {
            name = "app";
            src = ./.;
            vendorHash = null;
            tags = [ "production" ];
            nativeBuildInputs = with pkgs; [
              templ
              tailwindcss
              makeWrapper
            ];
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
              paths = [ packages.app ];
            };
            config =
              let
                healthcheck = pkgs.writeShellApplication {
                  name = "healthcheck";
                  runtimeInputs = [ pkgs.curl ];
                  text = ''
                    test "$(curl --fail localhost:8080/health)" = "OK"
                  '';
                };
              in
              {
                Healthcheck.Test = [
                  "CMD"
                  "${pkgs.lib.getExe healthcheck}"
                ];
                Entrypoint = [ "app" ];
              };
          };
        };

        defaultPackage = packages.app;
      }
    );
}
