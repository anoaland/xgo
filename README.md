# XGO

Package xgo is for internal use only and is currently hosted in a private repository.

To use this private Go library in your project, you need to configure your environment to allow access to the private repository.

Follow these steps:

1. Set the `GOPRIVATE` environment variable to include the private repository. Add the following environment variable to your shell configuration file (e.g., `~/.bashrc`, `~/.zshrc`, or `~/.profile`):

   ```sh
   export GOPRIVATE=github.com/anoaland/*
   ```

2. Configure Git to use SSH instead of HTTPS for GitHub repositories:

   ```sh
   git config --global url."ssh://git@github.com/".insteadOf "https://github.com/"
   ```

   Alternatively, modify the `~/.gitconfig` file to include:

   ```sh
   [url "ssh://git@github.com/"]
   	 insteadOf = https://github.com/
   ```

   This configuration ensures that Git uses SSH authentication when accessing GitHub repositories.

3. Use `go get` to fetch the package in your project:

   ```sh
   go get github.com/anoaland/xgo
   ```

Note: Ensure that your Git configuration is set up to use the correct credentials for accessing the private repository.
