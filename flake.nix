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

        mkWOFF2From =
          name: pkg: ext:
          pkgs.stdenvNoCC.mkDerivation {
            name = "${name}-woff2";
            nativeBuildInputs = [
              pkgs.fontforge
              pkg
            ];
            dontInstall = true;
            unpackPhase = ''
              WOFF2_DIR="$out/share/fonts/woff2/"
              mkdir -p "$WOFF2_DIR"
              for file in ${pkg}/share/fonts/truetype/*.${ext}; do
              	NAME="$(basename $file .${ext})"
              	fontforge --lang=ff \
              		-c 'Open($1); Generate($2);' \
              		"$file" \
              		"$WOFF2_DIR/$NAME.woff2" &
              done
              wait
            '';
          };

        iosevka-comfy-woff2 = mkWOFF2From "iosevka-comfy-fixed" pkgs.iosevka-comfy.comfy-fixed "ttf";

        rundev = pkgs.writeShellApplication {
          name = "rundev";
          runtimeInputs = [
            pkgs.fd
            pkgs.entr
            (pkgs.writeShellApplication {
              name = "buildapp";
              text = ''
                function buildapp() {
                  ${self.defaultPackage.${system}.preBuild}
                  go build
                }
                buildapp && ./app
              '';
            })
          ];
          text = ''
            if [ $# -eq 0 ]; then
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
              install -D ${iosevka-comfy-woff2}/share/fonts/woff2/iosevka-comfy-fixed-regular.woff2 static/
              install -D ${iosevka-comfy-woff2}/share/fonts/woff2/iosevka-comfy-fixed-bold.woff2    static/
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
