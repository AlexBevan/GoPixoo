# GoPixoo

A Go CLI tool for controlling Divoom Pixoo64 LED displays over your local network.

## Features

- 📸 **Send images & animated GIFs** — Push static images (PNG, JPEG, BMP) or multi-frame GIFs to the display
- 📝 **Scrolling text** — Display customizable text with fonts, colors, speed, and alignment
- ⏱️ **Built-in tools** — Timer, stopwatch, scoreboard, noise meter
- 🎨 **Drawing primitives** — Pixel-level control, fill, clear
- 🔧 **Device management** — Get device info, set time, reboot
- 🌈 **Channel switching** — Clock, cloud, visualizer, custom modes
- 🔌 **Brightness & display control** — Adjust brightness, turn display on/off
- 🎮 **Raw JSON commands** — Direct API access for unsupported commands
- 📦 **Reusable library** — `internal/pixoo` package for Go projects

## Installation

### From GitHub Releases

Download the latest binary from [releases](https://github.com/alexbevan/gopixoo/releases) and add to your `$PATH`.

### Build from Source

```bash
git clone https://github.com/alexbevan/gopixoo.git
cd gopixoo
go install ./cmd/gopixoo
```

Or install directly:

```bash
go install github.com/alexbevan/gopixoo@latest
```

## Quick Start

### 1. Set Your Device IP

Find your Pixoo64's IP address on your network. Then choose one of these methods:

**Option A: CLI flag**
```bash
gopixoo --ip 192.168.1.100 send animation.gif
```

**Option B: Environment variable**
```bash
export GOPIXOO_IP=192.168.1.100
gopixoo send animation.gif
```

**Option C: Config file** (recommended)
```bash
mkdir -p ~/.config/gopixoo
cat > ~/.config/gopixoo/config.yaml << EOF
device:
  ip: "192.168.1.100"
EOF

# Now you can omit --ip on all commands
gopixoo send animation.gif
```

### 2. Send Your First Image

```bash
# Send an animated GIF
gopixoo send ~/Pictures/animation.gif

# Send a static PNG
gopixoo send ~/Pictures/logo.png --resize fit

# Send with custom frame delay (milliseconds)
gopixoo send animation.gif --speed 150
```

### 3. Display Text

```bash
# Simple scrolling text
gopixoo text send "Hello World!"

# Text with color and font
gopixoo text send "Alert!" --color "#FF0000" --font 3

# Clear text
gopixoo text clear
```

### 4. Control the Device

```bash
# Turn display on/off
gopixoo display on
gopixoo display off

# Get/set brightness (0-100)
gopixoo brightness get
gopixoo brightness set 75

# Switch channels
gopixoo channel clock
gopixoo channel visualizer

# View device info
gopixoo device info
```

## Configuration

### Config File

Create `~/.config/gopixoo/config.yaml`:

```yaml
device:
  ip: "192.168.1.100"
  # port: 80  (default, rarely changed)

defaults:
  brightness: 80
  text_color: "#FFFFFF"
  text_font: 2
```

### Environment Variables

Set `GOPIXOO_IP` to override the config file:

```bash
export GOPIXOO_IP=192.168.1.100
gopixoo device info
```

### Configuration Priority

1. **CLI flags** (`--ip 192.168.1.100`) — highest priority
2. **Environment variables** (`GOPIXOO_IP=...`)
3. **Config file** (`~/.config/gopixoo/config.yaml`)
4. **Defaults** — lowest priority

## Commands Reference

### `send` — Push Image or GIF

Send static images (PNG, JPEG, BMP) or animated GIFs to the display. For animated GIFs,
all frames are extracted and sent sequentially. For static images, a single frame is sent.
The `--speed` flag only applies to animated GIFs.

```bash
gopixoo send <file> [flags]
```

**Flags:**
- `--ip <IP>` — Device IP address
- `-s, --speed <ms>` — Frame delay in milliseconds (default: 100)
- `-r, --resize <mode>` — Resize mode: `fit` (default), `fill`, `stretch`, `none`
  - `fit` — Scale with aspect ratio, pad with black
  - `fill` — Scale with aspect ratio, crop to anchor point
  - `stretch` — Ignore aspect ratio, fill 64x64
  - `none` — No resize, crop to 64x64
- `--anchor <position>` — Crop anchor for `fill` mode (default: center)
  - `center` — Crop from center
  - `top` — Crop from top
  - `bottom` — Crop from bottom
  - `left` — Crop from left
  - `right` — Crop from right
- `--size <px>` — Target size (default: 64 for Pixoo64)

**Examples:**
```bash
# Basic send
gopixoo send birthday.gif

# High-speed animation (50ms per frame)
gopixoo send animation.gif --speed 50

# Large image, scale to fit
gopixoo send poster.png --resize fit

# Stretch to fill (no aspect ratio)
gopixoo send pattern.jpg --resize stretch

# Fill mode, crop from top
gopixoo send portrait.jpg --resize fill --anchor top

# Fill mode, crop from left
gopixoo send landscape.jpg --resize fill --anchor left

# Send to custom size (e.g., Pixoo16)
gopixoo send small.gif --size 16
```

### `text` — Send or Clear Text

Display scrolling or static text on the display.

```bash
gopixoo text send <message> [flags]
gopixoo text clear
```

**Send Flags:**
- `--x <px>` — X position (default: 0)
- `--y <px>` — Y position (default: 0)
- `--font <id>` — Font index 0-7 (default: 0)
- `--color <hex>` — Text color as hex `#RRGGBB` (default: #FFFFFF)
- `--speed <ms>` — Scroll speed in milliseconds (default: 100)
- `--dir <dir>` — Direction: 0=left (default), 1=right
- `--align <align>` — Alignment: 1=left, 2=center (default), 3=right
- `--id <id>` — Text ID for managing multiple texts (default: 1)
- `--width <px>` — Text area width (default: 64)

**Examples:**
```bash
# Simple text
gopixoo text send "Hi there!"

# Red text, centered, fast scroll
gopixoo text send "Alarm!" --color "#FF0000" --align 2 --speed 50

# Static text at position (10, 20)
gopixoo text send "Status: OK" --x 10 --y 20 --speed 5000

# Clear all text
gopixoo text clear
```

### `brightness` — Get or Set Display Brightness

Control display brightness (0-100).

```bash
gopixoo brightness get          # Show current brightness
gopixoo brightness set 75       # Set to 75%
```

### `display` — Screen Power Control

Turn the display on or off.

```bash
gopixoo display on      # Turn on
gopixoo display off     # Turn off
```

### `channel` — Switch Display Channel

Switch between built-in channels.

```bash
gopixoo channel clock           # Clock display
gopixoo channel cloud           # Cloud/online content
gopixoo channel visualizer      # Music/sound visualizer
gopixoo channel custom          # Custom/user content
```

### `clock` — Get or Set Clock Face

Configure the clock display.

```bash
gopixoo clock get           # Show current clock face ID
gopixoo clock set <id>      # Set clock face by ID
```

**Examples:**
```bash
gopixoo clock set 0
gopixoo clock set 5
gopixoo clock get
```

### `tool` — Built-in Device Tools

Access Pixoo's built-in tools: timer, stopwatch, scoreboard, noise meter.

#### Timer

```bash
gopixoo tool timer <minutes> <seconds>     # Start timer
gopixoo tool timer <m> <s> --stop          # Stop timer
```

**Examples:**
```bash
gopixoo tool timer 5 0          # Start 5-minute timer
gopixoo tool timer 1 30         # Start 1m 30s timer
gopixoo tool timer 0 30 --stop  # Stop timer
```

#### Stopwatch

```bash
gopixoo tool stopwatch start    # Start stopwatch
gopixoo tool stopwatch stop     # Pause stopwatch
gopixoo tool stopwatch reset    # Reset stopwatch
```

#### Scoreboard

```bash
gopixoo tool scoreboard <blue> <red>
```

**Examples:**
```bash
gopixoo tool scoreboard 10 8
gopixoo tool scoreboard 0 0
```

#### Noise Meter

```bash
gopixoo tool noise start        # Start noise meter
gopixoo tool noise stop         # Stop noise meter
```

### `draw` — Low-Level Drawing Primitives

Draw directly to the display at the pixel level.

#### Pixel

```bash
gopixoo draw pixel <x> <y> [--color <hex>]
```

**Examples:**
```bash
gopixoo draw pixel 32 32                    # White pixel at center
gopixoo draw pixel 10 10 --color "#FF0000" # Red pixel
```

#### Fill

```bash
gopixoo draw fill [--color <hex>]
```

**Examples:**
```bash
gopixoo draw fill                           # Red fill (default)
gopixoo draw fill --color "#00FF00"        # Green fill
```

#### Clear

```bash
gopixoo draw clear
```

Clear the entire display (fill with black).

### `device` — Device Information and Management

Get device status and control.

#### Info

```bash
gopixoo device info
```

Shows device settings: brightness, clock face, temperature mode, rotation, mirror mode, etc.

#### Time

```bash
gopixoo device time
```

Show device's current time.

#### Reboot

```bash
gopixoo device reboot
```

Reboot the Pixoo64 device.

### `raw` — Send Raw JSON

For unsupported commands, send raw JSON directly to the device.

```bash
gopixoo raw '<JSON>'
```

**Examples:**
```bash
# Get all channel configuration
gopixoo raw '{"Command": "Channel/GetAllConf"}'

# Custom payload
gopixoo raw '{"Command": "Device/GetDeviceSettings"}'
```

The response is printed as pretty-printed JSON.

## Global Flags

Available on all commands:

- `--ip <IP>` — Pixoo64 device IP address
- `--config <path>` — Path to config file (default: `~/.config/gopixoo/config.yaml`)
- `--verbose` — Enable verbose output (shows HTTP requests, frame counts, etc.)
- `-h, --help` — Show command help

**Examples:**
```bash
gopixoo --verbose send animation.gif
gopixoo --config /etc/gopixoo/custom.yaml brightness get
```

## Examples

### Send an Animated Birthday GIF

```bash
gopixoo send ~/birthdays/confetti.gif --speed 100 --resize fit
```

### Display a Live Scoreboard

```bash
gopixoo channel custom
gopixoo tool scoreboard 42 38
gopixoo text send "Game Score" --color "#00FFFF"
```

### Set Up a Clock with Adjusted Brightness

```bash
gopixoo channel clock
gopixoo clock set 3
gopixoo brightness set 60
```

### Show System Status Text

```bash
gopixoo text send "CPU: 45% RAM: 62%" --font 1 --color "#00FF00" --speed 200
```

### Cycle Through Channels

```bash
gopixoo channel clock
sleep 10
gopixoo channel visualizer
sleep 10
gopixoo channel cloud
```

### Create a Simple Animation Loop

```bash
#!/bin/bash
for gif in ~/animations/*.gif; do
  echo "Sending $(basename "$gif")..."
  gopixoo send "$gif" --resize fit --speed 80
  sleep 5
done
```

## Troubleshooting

### Device Not Responding

1. **Verify device IP**
   ```bash
   ping 192.168.1.100
   ```

2. **Check device is on network**
   - Ensure Pixoo64 is powered on and connected to WiFi
   - Look for the device's IP in your router's admin panel

3. **Try specifying IP directly**
   ```bash
   gopixoo --ip 192.168.1.100 device info
   ```

### Image Doesn't Display Correctly

- Try `--resize fill` for better center crop
- Use `--speed 200` for slower, clearer animation
- Ensure image is in supported format (PNG, JPEG, GIF, BMP)

### Text Appears Garbled

- Try different `--font` values (0-7)
- Reduce `--speed` for slower scrolling
- Use `--align 2` for centered text

### Verbose Output for Debugging

```bash
gopixoo --verbose send animation.gif
```

Shows detailed HTTP requests and frame processing.

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature`)
3. Commit your changes (`git commit -am 'Add your feature'`)
4. Push to the branch (`git push origin feature/your-feature`)
5. Open a Pull Request

### Development

Build and test locally:

```bash
go build -o gopixoo ./cmd/gopixoo
go test ./...
```

## Architecture

GoPixoo follows a clean separation of concerns:

- **`cmd/`** — Cobra command definitions (CLI layer)
- **`internal/pixoo/`** — Pixoo64 HTTP client library
  - `client.go` — HTTP transport
  - `commands.go` — Command builders
  - `draw.go`, `animation.go`, `channel.go`, `device.go`, `tool.go` — Command-specific builders
- **`internal/imaging/`** — Image processing utilities
  - `convert.go` — Pixel encoding, resizing
- **`internal/config/`** — Configuration management (future)

The `internal/pixoo` package is a standalone library and can be imported by other Go projects.

## License

MIT License — see [LICENSE](LICENSE) for details.

## Support

For issues, feature requests, or questions:

- 📝 [GitHub Issues](https://github.com/alexbevan/gopixoo/issues)
- 💬 [GitHub Discussions](https://github.com/alexbevan/gopixoo/discussions)

## See Also

- [Divoom Pixoo64 Official](https://www.divoom.com/products/pixoo64)
- [Pixoo API Reference](https://github.com/divoom/pixoo-api) (third-party documentation)
