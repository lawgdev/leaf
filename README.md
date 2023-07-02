# 🍃 Leaf

Leaf is a simple and lightweight CLI tool for generating, running, and managing your lawg projects and feeds.

## 🌿 Installation

### Homebrew

```bash
brew tap lawg/leaf
brew install leaf
```

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

