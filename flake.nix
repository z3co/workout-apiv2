{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = inputs @ { flake-parts, ... }: 
  flake-parts.lib.mkFlake {inherit inputs; }
  {
    systems = [
      "x86_64-linux" "x86_64-darwin"
    ];
    perSystem = {pkgs, ... }: {
      devShells.default = pkgs.mkShell {
        name = "go";
        packages = with pkgs; [
          zsh go gopls postgresql docker-compose sqlc
        ];
        shellHook = ''
          export SHELL=${pkgs.zsh}/bin/zsh
          exec ${pkgs.zsh}/bin/zsh
        '';
        PGUSER = "postgres";
        PGPASSWORD = "postgres";
        PGDATABASE = "main";
        PGHOST = "localhost";
        PGPORT = 5432;
      };
    };
  };
}
