{
  description = "A Go-based CLI to generate .dockerignore files, inspired by gibo.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "dibo";
          version = "0.1.0";
          src = ./.;
          
          vendorHash = "sha256-AhGRCTU795Rwv2K435aSl8GCNwZX3M5mohYW5NipPkY="; 

          proxyVendor = true;
          
          preBuild = ''
            export GOTOOLCHAIN=auto
          '';

          preConfigure = ''
            export GOTOOLCHAIN=auto
          '';

          ldflags = [ "-s" "-w" ];

          meta = with pkgs.lib; {
            description = "A Go-based CLI to generate .dockerignore files, inspired by gibo.";
            homepage = "https://github.com/anton-fuji/dibo";
            license = licenses.mit;
            maintainers = [ "anton-fuji" ];
            platforms = platforms.all;
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
          ];
          shellHook = ''
            export GOTOOLCHAIN=auto
          '';
        };
      }
     );
}
