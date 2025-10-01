# PW Equipment Changer

A GUI application for automatically switching equipment sets in Perfect World using hotkeys. The application features a modern subscription-based system with email registration and hardware ID validation.

## 🎮 Features

- **Modern GUI Interface**: Built with Fyne framework for cross-platform compatibility
- **Subscription System**: Email-based registration with HWID (Hardware ID) validation
- **Equipment Set Switching**: Configure up to 11 items for automatic equipment changes
- **Hotkey Activation**: Press `Q` to switch between equipment sets
- **Configurable Timing**: Adjustable delay between item clicks
- **Persistent Configuration**: Saves user email and settings locally
- **Cross-Platform**: Supports Windows, macOS, and Linux

## 🏗️ Architecture

The application is structured as follows:

```
pw-equip-change/
├── main.go                 # Application entry point
├── equip/                  # Core package
│   ├── gui.go             # GUI interface and main application logic
│   ├── models.go          # Data structures (SetupEquip)
│   ├── subscription.go    # User registration and validation
│   ├── utils.go           # Utility functions and automation
│   └── config.go          # Configuration management
├── media/
│   └── icon.jpg           # Application icon
├── .github/workflows/
│   └── build_windows.yml  # GitHub Actions for Windows builds
└── build.sh               # Multi-platform build script
```

## 🚀 Installation & Usage

### Download Pre-built Binary

1. Download the latest release from the [Releases](https://github.com/your-repo/pw-equip-change/releases) page
2. Run the executable: `pw-equip-changer.exe` (Windows) or `pw-equip-changer` (macOS/Linux)

### Building from Source

#### Prerequisites
- Go 1.24.5 or later
- Fyne dependencies for your platform

#### Quick Build
```bash
# Clone the repository
git clone <repository-url>
cd pw-equip-change

# Build using the automated script
./build.sh
```

#### Manual Build
```bash
# Standard build
go build -o pw-equip-changer

# Windows (hide console window)
go build -ldflags="-H windowsgui" -o pw-equip-changer.exe

# Using Fyne packaging (recommended for distribution)
fyne package -os windows --name pw-equip-changer.exe -icon media/icon.jpg -release
```

## ⚙️ Configuration

### Initial Setup

1. **Email Registration**: Enter the email used for purchasing the program
2. **Equipment Configuration**:
   - **Number of Items**: How many equipment pieces to swap (1-11)
   - **Bar Change Key**: Key to switch skill bars (`v` or `` ` ``)
   - **Click Timing**: Delay between clicks in milliseconds (e.g., 200ms = 0.2 seconds)
   - **Item Keys**: Individual keys for each equipment piece

### Game Setup Instructions

1. **Prepare 3 Skill Bars**: Leave 3 bars available for rotation
2. **Main Bar**: Set up your skills/potions as desired
3. **Attack Set Bar**: Place attack equipment in the second bar
4. **Defense Set Bar**: Place defense equipment in the third bar
5. **Activation**: Press `Q` to switch between equipment sets

## 🔐 Subscription System

The application uses a robust subscription system:

### Registration Process
1. Enter your purchase email in the application
2. The app generates a unique HWID based on your hardware
3. Registration request is sent to `gamedevforge.ovh/register-user`
4. Subscription validation occurs via `gamedevforge.ovh/validate-user`

### HWID Management
- **Hardware ID**: Unique identifier based on system information
- **Reset Function**: Available if you need to change machines
- **Security**: Prevents unauthorized usage across multiple devices

### API Endpoints
- `POST /register-user?email={email}&hwid={hwid}` - Register user with HWID
- `GET /validate-user?email={email}&hwid={hwid}` - Validate subscription
- `PATCH /reset-hwid?email={email}&hwid={hwid}` - Reset HWID for new machine

## 🛠️ Development

### Dependencies

```go
// Core dependencies
fyne.io/fyne/v2 v2.6.3          // GUI framework
github.com/go-vgo/robotgo        // Keyboard automation
github.com/robotn/gohook         // Global hotkey capture
github.com/shirou/gopsutil/v4    // System information for HWID
```

### Development Tools (mise.toml)

```bash
# Run tests
mise run test

# Security checks
mise run sec

# Format code
mise run fmt

# Lint code
mise run lint

# Run all checks
mise run checks

# Build
mise run build
```

### Project Structure

- **main.go**: Simple entry point that initializes the GUI application
- **equip/gui.go**: Main application logic, UI components, and event handling
- **equip/models.go**: Data structures for equipment configuration
- **equip/subscription.go**: User authentication and subscription validation
- **equip/utils.go**: Automation functions and utility methods
- **equip/config.go**: Configuration persistence (saves to `~/.pw-equip-change/config.json`)

## 🔧 Technical Details

### Equipment Switching Logic

The application implements a two-set rotation system:

```go
// Set 1 → Set 2: Navigate to attack equipment bar
KeyChange → KeyChange → ClickItems → KeyChange

// Set 2 → Set 1: Navigate to defense equipment bar  
KeyChange → ClickItems → KeyChange → KeyChange
```

### HWID Generation

Hardware ID is generated using:
- Host ID
- Platform information
- Platform family
- System architecture
- MD5 hash for consistency

### Configuration Storage

User settings are stored in:
- **Windows**: `%USERPROFILE%\.pw-equip-change\config.json`
- **macOS/Linux**: `~/.pw-equip-change/config.json`

## 🚀 CI/CD

The project includes GitHub Actions for automated Windows builds:

- **Trigger**: Push to main branch
- **Environment**: Windows Latest with Go 1.24.5
- **Output**: Windows executable with Fyne packaging
- **Release**: Automatic GitHub release creation with version tagging

## 📝 Version History

Current version: **0.4** (see `version.txt`)

## 🤝 Support

For support and subscription activation:
1. Run the application to get your HWID
2. Contact support with your email and HWID
3. Visit [gamedevforge.ovh](https://gamedevforge.ovh) for purchases

## 📄 License

This is a commercial application with subscription-based licensing.

---

**Note**: This application is specifically designed for Perfect World gameplay automation and requires an active subscription to function.