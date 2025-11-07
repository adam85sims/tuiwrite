# Terminal Setup for Catppuccin Colors

TUIWrite uses the beautiful Catppuccin Mocha color palette with 24-bit true color support. For the best experience, your terminal needs to be configured to support true color (16 million colors).

## Testing Your Terminal

Run this command to test your terminal's color support:

```bash
./tuiwrite --test-colors
```

This will display all Catppuccin colors. If they appear as intended (vibrant pastels), your terminal is properly configured! If they look incorrect, monochrome, or washed out, follow the setup instructions below.

## Quick Test

Run this in your terminal:
```bash
echo $COLORTERM
```

If it outputs `truecolor` or `24bit`, you're good to go! If not, follow the setup for your terminal below.

---

## Terminal Configuration Guides

### GNOME Terminal (Linux)

GNOME Terminal supports true color but may override it with theme settings.

1. Open GNOME Terminal
2. Go to **Preferences** → **Profiles** → Select your profile
3. Click the **Colors** tab
4. **Uncheck** "Use colors from system theme"
5. Set **Built-in schemes** to "Custom" or "Tango Dark"
6. Restart your terminal

**Set COLORTERM environment variable** (add to `~/.bashrc` or `~/.zshrc`):
```bash
export COLORTERM=truecolor
```

### Konsole (KDE)

Konsole has excellent true color support by default.

1. Just ensure your profile isn't forcing a specific color scheme
2. Edit → Current Profile → Appearance → Get New...
3. Search for "Catppuccin" if you want a matching terminal theme

**Already set by default**, but you can verify:
```bash
echo $COLORTERM  # Should output 'truecolor'
```

### Alacritty (Cross-platform)

Alacritty has perfect true color support out of the box. No configuration needed!

**Optional**: Add Catppuccin theme to your `~/.config/alacritty/alacritty.toml`:
```toml
[colors.primary]
background = '#1e1e2e'
foreground = '#cdd6f4'
```

### Kitty (Cross-platform)

Kitty has excellent true color support by default.

**Optional**: Add to `~/.config/kitty/kitty.conf`:
```conf
# Ensure true color support
term xterm-kitty
```

### iTerm2 (macOS)

iTerm2 fully supports true color.

1. Open iTerm2 → **Preferences** → **Profiles** → **Terminal**
2. Under "Terminal Emulation", ensure **Report Terminal Type** is `xterm-256color` or higher
3. iTerm2 will automatically use true color

**Optional**: Install Catppuccin theme for iTerm2:
```bash
curl -O https://raw.githubusercontent.com/catppuccin/iterm/main/colors/catppuccin-mocha.itermcolors
```
Then import via Preferences → Profiles → Colors → Color Presets → Import

### Windows Terminal

Windows Terminal has true color support since version 1.0.

1. Open Windows Terminal
2. Settings (Ctrl+,) → Your Profile → Appearance
3. Color scheme: Any modern scheme works (Catppuccin is available!)
4. Make sure "Use acrylic" is OFF for best color accuracy

**Install Catppuccin theme**:
1. Download theme from: https://github.com/catppuccin/windows-terminal
2. Add to `settings.json` under `"schemes"`

### Warp (macOS)

Warp has excellent true color support by default. No configuration needed!

### Wezterm (Cross-platform)

Wezterm has perfect true color support.

Add to `~/.wezterm.lua`:
```lua
return {
  term = "wezterm",
  -- Optional: Use Catppuccin theme
  color_scheme = "Catppuccin Mocha",
}
```

### Standard Terminal.app (macOS)

macOS Terminal has limited color support and **does not support true color**. For the best TUIWrite experience on macOS, we recommend:
- **iTerm2** (recommended)
- **Wezterm**
- **Alacritty**

### tmux Users

If you use tmux, add this to your `~/.tmux.conf`:

```conf
# Enable true color support
set -g default-terminal "tmux-256color"
set -ga terminal-overrides ",*256col*:Tc"

# Or for more compatibility:
set -g default-terminal "screen-256color"
set -ga terminal-overrides ",xterm-256color:Tc"
```

Then restart tmux:
```bash
tmux kill-server
tmux
```

### SSH Sessions

Colors should work over SSH if:
1. Your **local** terminal supports true color
2. Your `COLORTERM` environment variable is forwarded

Add to `~/.ssh/config`:
```
Host *
    SendEnv COLORTERM
```

And on the server, ensure `/etc/ssh/sshd_config` has:
```
AcceptEnv COLORTERM
```

---

## Troubleshooting

### Colors look washed out or wrong
- Check if your terminal theme is overriding colors
- Try the `--test-colors` flag to verify
- Ensure COLORTERM=truecolor is set

### Colors are completely monochrome
- Your terminal may not support true color
- Try a modern terminal emulator (Alacritty, Kitty, iTerm2, Windows Terminal)

### Colors work in vim/neovim but not TUIWrite
- Check if vim is using a different `$TERM` value
- Run `:echo $TERM` in vim and compare to `echo $TERM` in your shell

### Still having issues?
Run these diagnostics and check the debug log:

```bash
echo "TERM: $TERM"
echo "COLORTERM: $COLORTERM"
./tuiwrite --test-colors
cat ~/.config/tuiwrite/debug.log
```

---

## Recommended Terminals

For the absolute best TUIWrite experience, we recommend:

**Linux**: Alacritty, Kitty, or Konsole  
**macOS**: iTerm2, Wezterm, or Alacritty  
**Windows**: Windows Terminal or Wezterm  
**Cross-platform**: Alacritty or Wezterm

All of these support true color out of the box with minimal configuration.
