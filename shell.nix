{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    python311Packages.mkdocs-material
    python311Packages.mkdocs-material-extensions
    python311Packages.pymdown-extensions
  ];
  shellHook = ''
  '';
}