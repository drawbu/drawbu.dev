{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = {self, ...} @ inputs:
    inputs.utils.lib.eachDefaultSystem (
      system: let
        pkgs = import inputs.nixpkgs {inherit system;};

        rundev = pkgs.writeShellScriptBin "rundev" ''
          templ generate                                                 && \
          mkdir -p /tmp/drawbu.dev                                       && \
          tailwindcss -i ./assets/style.css -o /tmp/drawbu.dev/style.css && \
          go build -ldflags="-X 'main.assetsDir=/tmp/drawbu.dev'"        && \
          ./app
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

          docker = pkgs.dockerTools.buildImage {
            name = "drawbu.dev";
            tag = "latest";
            created = "now";
            copyToRoot = pkgs.buildEnv {
              name = "image-root";
              paths = [ packages.app ];
            };
            config = {
              Cmd = ["app"];
              Env = ["SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"]; # wut
            };
          };
        };

        defaultPackage = packages.app;
      }
    );
}
