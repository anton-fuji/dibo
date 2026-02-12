{
  description = "A Go-based CLI to generate .dockerignore files, inspired by gibo.";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        packages.dibo = pkgs.buildGoModule {
          pname = "dibo";
          version = "0.1.0";

          src = ./.;

          vendorHash = null;

          go = pkgs.go_1_26;

          ldflags = [ "-s" "-w" ];

          meta = with pkgs.lib; {
            description = "A Go-based CLI to generate .dockerignore files, inspired by gibo.";
            homepage = "https://github.com/anton-fuji/dibo";
            license = licenses.mit;
            platforms = platforms.linux ++ platforms.darwin;
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go_1_25
            gopls
            golangci-lint
          ];
        };
      });
}
