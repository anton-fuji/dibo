{
  description = "A Go-based CLI to generate .dockerignore files, inspired by gibo.";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-23.11";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
      goVersion = "go1_25";
    in {
      packages.dibo = pkgs.buildGoModule {
        pname = "dibo";
        version = "0.1.0"; 

        src = ./.;

        vendorHash = null;

        proxyVendor = true;

        buildFlags = [ "-mod=vendor" ];

        ldflags = [ "-s" "-w" ];

        go = pkgs.${goVersion};

        meta = with pkgs.lib;
          {
            description = "A Go-based CLI to generate .dockerignore files, inspired by gibo.";
            homepage = "https://github.com/anton-fuji/dibo";
            license = licenses.mit;
            maintainers = [ "anton-fuji" ]; 
            platforms = platforms.linux ++ platforms.darwin;
          };
      };

      # nix-shell
      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          gopls
          go-tools
          golangci-lint
        ];
      };
    });
}
