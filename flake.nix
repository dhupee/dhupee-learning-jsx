{
  description = "dhupee testing uv2nix";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-25.05";
    nixpkgs-unstable.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    pyproject-nix = {
      url = "github:pyproject-nix/pyproject.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    uv2nix = {
      url = "github:pyproject-nix/uv2nix";
      inputs.pyproject-nix.follows = "pyproject-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    pyproject-build-systems = {
      url = "github:pyproject-nix/build-system-pkgs";
      inputs.pyproject-nix.follows = "pyproject-nix";
      inputs.uv2nix.follows = "uv2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = {
    nixpkgs,
    nixpkgs-unstable,
    flake-utils,
    pyproject-nix,
    uv2nix,
    pyproject-build-systems,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
        };

        pkgs-unstable = import nixpkgs-unstable {
          inherit system;
          config.allowUnfree = true;
        };

        # Load workspace and create overlays
        workspace = uv2nix.lib.workspace.loadWorkspace {
          # always point the root towards your UV project root
          workspaceRoot = ./backend;
        };
        overlay = workspace.mkPyprojectOverlay {sourcePreference = "wheel";};
        editableOverlay = workspace.mkEditablePyprojectOverlay {root = "$REPO_ROOT";};

        # Create Python package set
        pythonSet =
          (pkgs.callPackage pyproject-nix.build.packages {
            python = pkgs.python312;
          }).overrideScope (
            pkgs.lib.composeManyExtensions [
              pyproject-build-systems.overlays.wheel
              overlay
            ]
          );

        # Create editable Python set for dev shell
        editablePythonSet = pythonSet.overrideScope editableOverlay;
      in {
        devShells = {
          default = pkgs.mkShell {
            packages =
              [
                # For UV
                (editablePythonSet.mkVirtualEnv "dh-learning-jsx-dev-env" workspace.deps.all)
                pkgs.uv
              ]
              # Other packages with stable packages
              ++ (with pkgs; [
                # Javascripts
                nodejs_22

                # Dev Utils
                go-task
              ])
              # Other packages with unstable packages
              ++ (
                with pkgs-unstable; [
                  # Dependencies
                  yt-dlp
                  ffmpeg
                ]
              );
            env = {
              UV_NO_SYNC = "1";
              UV_PYTHON = editablePythonSet.python.interpreter;
              UV_PYTHON_DOWNLOADS = "never";
            };
            shellHook = ''
              unset PYTHONPATH
              export REPO_ROOT=$(git rev-parse --show-toplevel)
              export SHELL=$(which ${pkgs.zsh})
              echo "Development Shell Initialized"
              exec zsh
            '';
          };

          bruno = pkgs.mkShell {
            packages = with pkgs; [
              bruno
            ];
            shellHook = ''
              echo "Bruno Initialized"
            '';
          };
        };

        packages.default = pythonSet.mkVirtualEnv "dh-learning-jsx-env" workspace.deps.default;
      }
    );
}
