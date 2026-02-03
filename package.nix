{
  perSystem = { pkgs, ... }: {
    packages.default = pkgs.buildGoModule {
      pname = "workout-apiv2";
      version = "0.1.0";
      src = ./.;
      vendorHash = "sha256-QO/MohQ+hqhSqbnGjCLqSMovNHNH6zEhKGznHV96sBc=";
    };
  };
}
