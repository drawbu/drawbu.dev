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
          ./server
        '';
      in {
        formatter = pkgs.alejandra;

        devShell = pkgs.mkShell {
          inputsFrom = builtins.attrValues self.packages.${system};
          packages = [rundev];
        };

        packages = {
          default = self.packages.${system}.server;
          server = pkgs.buildGoModule {
            name = "server";
            src = ./.;
            vendorHash = null;
            ldflags = ["-X main.assetsDir=${placeholder "out"}/share/assets"];
            nativeBuildInputs = with pkgs; [templ tailwindcss];
            propagatedBuildInputs = with pkgs; [git];
            preBuild = ''
              templ generate
            '';
            postBuild = ''
              mkdir -p $out/share/assets
              tailwindcss -i ./assets/style.css -o $out/share/assets/style.css
            '';
          };
        };
      }
    );
}
