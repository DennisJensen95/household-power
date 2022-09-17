{ pkgs ? import (fetchTarball "http://nixos.org/channels/nixos-21.05/nixexprs.tar.xz") {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.which
    pkgs.htop
    pkgs.zlib
  ];

  shellHook = ''echo Entering household-power development environment'';
}