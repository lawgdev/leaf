# 🍃 Leaf

Leaf is a simple and lightweight CLI tool for generating, running, and managing your lawg projects and feeds.

## 🌿 Installation

### Homebrew

```bash
brew install lawgdev/tap/leaf
```

### Windows

1. Download the latest release from the [releases page](https://github.com/lawgdev/leaf/releases)
2. Extract the zip file
3. Move the `leaf.exe` file to a directory in your `PATH` environment variable or add the directory to your `PATH` environment variable. (We recommend putting it in your User directory)

### Manual

```bash
git clone
cd leaf
make install
```

## 🌿 Usage

### Authenticate

```bash
leaf login
```

### Link an application to a feed.

```bash
leaf connect
```

### Listening for logs with that application.

```bash
leaf listen
```

### Automatically start on system boot (launchd/systemctl)

### Supported platforms

| Platform             | Supported |
| -------------------- | --------- |
| macOS                | ✅        |
| Linux (systemctl)    | ✅        |
| Windows              | ❌        |



```bash
leaf upstart
```

