{
  perSystem = { pkgs, ... }: {
    packages.default = pkgs.buildGoModule {
      pname = "workout-apiv2";
      version = "0.1.0";
      src = ./.;
      vendorHash = "sha256-zXugligcuW4EutdXx8dSpjeIS15yuxP2af+Ykv7I3kU=";
      doCheck = false;
    };
  };
}
