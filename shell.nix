{
  pkgs ? import sources.nixpkgs { },
  sources ? import nix/sources.nix,
}:
pkgs.mkShell {
  nativeBuildInputs = with pkgs; [
    bazelisk
    nodePackages.aws-cdk
    nodejs-18_x
  ];
}
