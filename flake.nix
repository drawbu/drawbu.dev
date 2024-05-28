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
          packages = with pkgs; [templ];
        };

        packages = {
          default = self.packages.${system}.server;
          server = pkgs.buildGoModule {
            name = "server";
            src = ./.;
            vendorHash = null;
            preBuild = ''
              ${pkgs.templ}/bin/templ generate
            '';
          };
        };
      }
    );
}
