```
 .ooooo.         .o   
d88'   `8.     .d88   
Y88..  .8'   .d'888   
 `88888b.  .d'  888   
.8'  ``88b 88ooo888oo 
`8.   .88P      888   
 `boood8'      o888o  
                      
```

# 84 - A Retro Markdown Editor

84 is a lightweight, terminal-based Markdown editor with a 1980s-inspired interface, designed for simplicity and integration with the nbgo note-taking system. It provides a distraction-free editing experience with function key navigation and Markdown-specific shortcuts, perfect for quick note and bookmark editing.

## Features

Retro UI: Function key menu (F1–F5, F10) and Ctrl-based shortcuts inspired by nano and emacs.

Markdown Editing: Supports headers, lists, bold, italic, links, and inline code via shortcuts.

Mouse Support: Click to set the cursor position.

Search: Find text with case-insensitive matching (F3).

Preview: Render Markdown with a dark theme (F5).

Help Modal: View keybindings (F1).

Command-Line Usage: Edit files with 84 <filename> (creates .md files if needed).

Integration: Works with nbgo for note and bookmark editing.

## Usage

With nbgo

Switch to a notebook:nbgo use work

Run the TUI:nbgo

Navigate with up/down or mouse.
a: Add a note.
b: Add a bookmark.
v: View in Glow.
e: Edit in 84.
q or Ctrl+C: Quit.

## Standalone

Edit a Markdown file directly:

```bash
84 newnote
```

Creates/edits newnote.md in the current directory.

Keybindings in 84

F1: Show help modal (lists all keybindings).
F2: Save and exit.
F3: Search text (Enter to find, Esc to cancel).
F4: Toggle function key menu.
F5: Toggle Markdown preview.
F10/Esc: Quit without saving.
Ctrl+H: Insert header (# ) or increment level.
Ctrl+L: Insert list item (- ).
Ctrl+B: Insert bold markers (****, cursor between).
Ctrl+I: Insert italic markers (**, cursor between).
Ctrl+K: Insert link template ([text](url), cursor on text).
Ctrl+M: Insert inline code (`code`, cursor inside).
Mouse Click: Set cursor position.

## Notes

Terminal Support: Requires a modern terminal (e.g., iTerm2, Kitty, Alacritty) for mouse and function keys. Test on compact keyboards (e.g., Fn+5 for F5).
Limitations: No mouse-based text selection or syntax highlighting (by design). Search moves to the first match; cycling requires future enhancements.
Integration: Designed for nbgo, storing notes/bookmarks in

```bash
~/.nbgo/<notebook> as .md or .bookmark.md files.
```
