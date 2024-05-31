# Flicker

Flicker is a live reloading tool for Flutter projects, designed to automate the hot reload process whenever changes are detected in your project files. This tool is inspired by Air for Go and aims to streamline your Flutter development workflow.

## Features

- **Automatic Hot Reload**: Automatically triggers Flutter's hot reload whenever a file change is detected.
- **Configurable Watch Directories**: Specify which directories to watch for changes.
- **Device Selection**: Choose which device to run your Flutter application on.
- **Verbose Logging**: Optional verbose logging for detailed output.
- **Clean Exit**: Automatically exits when the Flutter application is closed.

## Installation

### Prerequisites

- Go 1.16 or later
- Flutter installed and configured

### Steps

1. **Clone the repository**:
   ```sh
   git clone https://github.com/ZLUN73L_thdgit/flicker.git
   cd flicker
   ```

2. **Install dependencies**:
   ```sh
   go get -u github.com/fsnotify/fsnotify
   go get -u github.com/pelletier/go-toml
   ```

3. **Build the binary**:
   ```sh
   go build -o flicker main.go
   ```

4. **Move the binary to a directory in your PATH** (optional but recommended):
   ```sh
   mv flicker /usr/local/bin/
   ```

## Configuration

Create a `flicker.toml` file in the root of your Flutter project with the following structure:

```toml
# flicker.toml

[watch]
# Directories to watch;
directories = ["lib"]

# Flutter target device (e.g., chrome, ios, android)
device = "chrome"
```

### Configuration Options

- `directories`: List of directories to watch for changes. If empty, Flicker will watch all directories in the Flutter project.
- `device`: The target device to run the Flutter application on. Examples include `chrome`, `ios`, `android`.

## Usage

Generating the Configuration File

To generate a standard flicker.toml configuration file with the default values of lib and chrome, run the following command:

```sh
flicker -generate-config
```

Navigate to your Flutter project directory and run Flicker:

```sh
flicker
```

Flicker will start the Flutter application on the specified device and automatically perform hot reloads whenever a change is detected in the watched directories. It will also exit cleanly when the Flutter application is closed.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request with your changes.

## Acknowledgements

- Inspired by [Air for Go](https://github.com/cosmtrek/air)
- Uses [fsnotify](https://github.com/fsnotify/fsnotify) for file system notifications
- Uses [go-toml](https://github.com/pelletier/go-toml) for TOML configuration

---

Happy coding with Flicker! 🚀
