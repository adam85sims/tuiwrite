# Testing Spell Check Highlighting

## How to Test

1. Run tuiwrite with the test file:
   ```bash
   ./tuiwrite spell-test.md
   ```

2. In Read Mode, press `:` to enter command mode

3. Type `spell` or `spellcheck` and press Enter to enable spell-checking
   - The status bar will show "Spell checking enabled (uk)"
   - This will download the UK English dictionary if not already present

4. You should now see misspelled words highlighted with a RED background

5. Misspelled words in the test file:
   - recieve (should be receive)
   - occured (should be occurred)  
   - seperate (should be separate)
   - definately (should be definitely)
   - thier (should be their)
   - mispelled (should be misspelled)

## Commands

- `:spell` or `:spellcheck` - Toggle spell-checking on/off
- `:spell -uk` - Switch to UK English (default)
- `:spell -us` - Switch to US English
- `:spell -ca` - Switch to Canadian English
- `:spell -au` - Switch to Australian English

## Visual Appearance

- **Misspelled words**: Red background (#f38ba8) with dark text
- **Cursor**: Maroon background (#eba0ac) 
- The highlighting uses the Catppuccin Mocha theme colors

Press Ctrl+C to quit.
