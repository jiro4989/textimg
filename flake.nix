{
  description = "textimg is command to convert from color text (ANSI or 256) to image.";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    # 複数のシステム(Linux, Macなど)に簡単に対応するためのライブラリ
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "textimg";
          version = "3.2.0";
          src = ./.;
          vendorHash = "sha256-esTGTsar8qahGw625WjTjlFoVmeEL/72yLiHYvhOQi8=";

          # 開発用のツールの scripts/width を除外するため
          subPackages = ["."];
        };

        devShells.default = pkgs.mkShell {
          packages = [
            pkgs.go_1_26
            pkgs.gopls
            pkgs.markdown-toc
          ];

          shellHook = ''
            echo "go development environment was activated"
          '';
        };
      }
    );
}
