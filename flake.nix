{
  description = "Go development shell";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      perSystem =
        { pkgs, ... }:
        let
          go-migrate-mysql = pkgs.go-migrate.overrideAttrs (_oldAttrs: {
            tags = [ "mysql" ];
          });
        in
        {
          devShells.default = pkgs.mkShell {
            packages = with pkgs; [
              go
              gopls
              gotools
              golangci-lint
              delve
              gofumpt
              revive
              goctl
              go-migrate-mysql

              # gRPC / Protobuf tooling
              protobuf
              buf
              protoc-gen-go
              protoc-gen-go-grpc

              grpcurl
              grpcui
              httpyac

              # Frontend tooling
              nodejs_22
              pnpm
            ];
          };
        };
    };
}
