{
  description = "GOST flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {inherit system;};
      in {
        packages.default = pkgs.buildGoPackage {
          name = "GOST";
          src = ./.;
          goPackagePath = "github.com/grig-iv/gost";
          meta.mainProgram = "gost";
        };
      }
    );
}
