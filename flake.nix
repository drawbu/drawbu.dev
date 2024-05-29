{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = {self, ...} @ inputs:
    inputs.utils.lib.eachDefaultSystem (
      system: let
        pkgs = import inputs.nixpkgs {inherit system;};
      in {
        formatter = pkgs.alejandra;

        devShell = pkgs.mkShell {
          inputsFrom = builtins.attrValues self.packages.${system};
          # env.GOROOT = "${pkgs.go}/share/go";
        };

        packages = {
          default = self.packages.${system}.server;
          server = pkgs.buildGoModule {
            name = "server";
            src = ./.;
            vendorHash = null;
            ldflags = ["-X main.assetsDir=${placeholder "out"}/share/assets"];
            nativeBuildInputs = with pkgs; [templ];
            preBuild = ''
              templ generate
            '';
            postBuild = ''
              mkdir -p $out/share
              cp -r assets $out/share
            '';
          };
        };
      }
    );
}
