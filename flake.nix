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
        inherit (pkgs) lib;

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
              	${lib.getExe reduceFont} "$file" "$WOFF2_DIR/$NAME.woff2" &
              done
              wait
            '';
          };

        reduceFont =
          pkgs.writers.writePython3Bin "reduce-font" { libraries = with pkgs.python3.pkgs; [ fontforge ]; }
            ''
              import fontforge
              import sys

              if len(sys.argv) != 3:
                  print("Usage: python input.ttf output.woff2")
                  sys.exit(1)

              input_font = sys.argv[1]
              output_font = sys.argv[2]

              font = fontforge.open(input_font)

              ascii_glyphs = range(0, 256)
              for glyph in font.glyphs():
                  if glyph.unicode not in ascii_glyphs:
                      font.removeGlyph(glyph)
                  else:
                      glyph.simplify()

              font.autoHint()
              font.removeOverlap()

              font.generate(output_font, flags=("tfm"))

              font.close()
            '';
      in
      rec {
        formatter = pkgs.nixfmt-rfc-style;

        devShells.default = pkgs.mkShell {
          inputsFrom = builtins.attrValues packages;
          packages = [ ];
        };

        packages = {
          default = packages.app;
          app = pkgs.buildGoModule {
            name = "app";
            src = ./.;
            vendorHash = null;
            tags = [ "production" ];
            nativeBuildInputs = with pkgs; [
              templ
              tailwindcss_4
              makeWrapper
            ];
            preBuild =
              let
                targets = [
                  "regular"
                  "bold"
                  "extrabold"
                ];

                font = mkWOFF2From "iosevka-comfy-fixed" pkgs.iosevka-comfy.comfy-fixed "ttf";
                inherit (builtins) length genList elemAt;
              in
              ''
                ldflags+=" -X main.rev=${self.shortRev or "dev"}"

                templ generate

                mkdir -p static/dist
                cp static/robots.txt static/dist
                ${builtins.concatStringsSep "\n" (
                  genList (
                    i: "install -D ${font}/share/fonts/woff2/iosevka-comfy-fixed-normal${elemAt targets i}upright.woff2 static/dist"
                  ) (length targets)
                )}
                tailwindcss -i static/style.css -o static/dist/style.css -m
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
                  "${lib.getExe healthcheck}"
                ];
                Entrypoint = [ "app" ];
              };
          };

          dev = pkgs.writeShellApplication {
            name = "rundev";
            runtimeInputs = [
              pkgs.fd
              pkgs.entr
              (pkgs.writeShellApplication {
                name = "buildapp";
                text = ''
                  ${packages.default.preBuild}
                  go build -ldflags "$ldflags" -o app
                  ./app
                '';
              })
            ];
            text = ''
              if [ $# -eq 0 ]; then
                buildapp
              elif [ "$1" = "--watch" ]; then
                fd | entr -c -r -s buildapp
              else
                echo "Usage: $0 [--watch]"
              fi
            '';
          };

        };
      }
    );
}
