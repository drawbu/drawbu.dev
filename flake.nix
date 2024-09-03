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

        rundev = pkgs.writeShellApplication {
          name = "rundev";
          runtimeInputs = [
            pkgs.fd
            pkgs.entr
            (pkgs.writeShellApplication {
              name = "buildapp";
              text = ''
                ${self.defaultPackage.${system}.preBuild}
                go build
                ./app
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
                ${builtins.concatStringsSep "\n" (
                  genList (
                    i: "install -D ${font}/share/fonts/woff2/iosevka-comfy-fixed-${elemAt targets i}.woff2 static/"
                  ) (length targets)
                )}
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
