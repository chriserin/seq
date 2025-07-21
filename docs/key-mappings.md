# Key Mappings

This document lists all keyboard mappings for the sequencer application.

| Mapping | Key Binding | Description |
|---------|-------------|-------------|
| PlayStop | Space | Play the full arrangement once. If playing, stop |
| PlayOverlayLoop | ' + Space | Play the current overlay in a loop |
| PlayLoop | Alt + Space | Play the full arrangement in a loop |
| PlayPart | Ctrl + @ | Play current part in a loop |
| PlayRecord | : + Space | Play the full arrangement once and send a record message at the beginning |
| Increase | + / = | Increase value of current selection (tempo, beats, cycles, accents, etc.) or tempo by 5 if no specific selection |
| Decrease | - | Decrease value of current selection (tempo, beats, cycles, accents, etc.) or tempo by 5 if no specific selection |
| CursorLineStart | < | Move cursor to beginning of current line |
| CursorLineEnd | > | Move cursor to end of current line |
| CursorLastLine | b + l | Move cursor to last line |
| CursorFirstLine | b + f | Move cursor to first line |
| AccentIncrease | A | Increase accent value for current note. |
| AccentDecrease | a | Decrease accent value for current note. |
| ActionAddLineReset | s + s | Add line reset action to current line.  When the playback cursor reaches this action, the playback cursor will reset to the first beat. |
| ActionAddResetAll | s + S | Add reset action all to current line.  When the playback cursor reaches this action, all playback cursors will reset to the first beat.  |
| ActionAddLineBounce | s + b | Add line bounce action to current line.  When the playback cursor reaches this action it will reverse direction, and reverse again when reaching the line beginning creating a bouncing effect. |
| ActionAddLineBounceAll | s + B | Add line bounce all action to current line.  When the playback cursor reaches this action all playback cursors will reverse direction, and reverse again when reaching the line beginning creating a bouncing effect. |
| ActionAddSkipBeat | s + k | Add skip beat all action to current line.  When the playback cursor reaches this action, all the playback cursors will advance an additional beat.  |
| ActionAddSkipBeatAll | s + K | Add skip beat action to current line.  When the playback cursor reaches this action, it will advance an additional beat.  |
| ActionAddLineReverse | s + r | Add line reverse action to current line.  When the playback cursor reaches this action, the playback will reverse for this line. |
| ActionAddLineDelay | s + z | Add line delay action to current line. |
| ClearOverlay | C | Remove all notes and actions from the current overlay layer |
| GateIncrease | G | Increases the gate value for current note.  The gate corresponds to the length of the note. |
| GateDecrease | g | Decreases the gate value for current note.  The gate corresponds to the length of the note. |
| GateBigIncrease | E | Increase gate value for current note by 8, or 1 full beat. The gate corresponds to the length of the note. |
| GateBigDecrease | e | Decrease gate value for current note by 8, or 1 full beat. The gate corresponds to the length of the note. |
| RotateDown | J | Rotate pattern down. In the current column shift all notes down by one line, with a note in the bottom line moving to the top line |
| RotateUp | K | Rotate pattern up. In the current column shift all notes down by one line, with a note in the bottom line moving to the top line |
| RotateLeft | H | Rotate pattern left. On the current line shift all notes right of the cursor left by one beat.  A note at the cursor's beat will be moved to the last beat of the line. |
| RotateRight | L | Rotate pattern right. On the current line shift all notes right of the cursor right by one beat.  A note at the last beat will be moved to the cursor's beat.  |
| SelectKeyLine | Y | Selects the current line as the key line. The KeyCycle of the part is advanced when the cursor returns to the first beat. SEE KEYCYCLES |
| Mute | m | Mute the current line. Midi messages will not be sent from this line when the line is muted. |
| Solo | M | Solo the current line.  Only midi messages from this line or other soloed lines will be sent. |
| RatchetDecrease | r | Decrease the number of hits evenly divided within the span of 1 beat |
| RatchetIncrease | R | Increase the number of hits evenly divided within the span of 1 beat |
| WaitIncrease | W | Increase the wait value for current note. The wait value is the time between the playback of the note's beat and the sending of the midi message.  This is useful for creating a swing effect.  |
| WaitDecrease | w | Decrease the wait value for current note. The wait value is the time between the playback of the note's beat and the sending of the midi message.  This is useful for creating a swing effect. The initial value for a note will 0, in which case WaitDecrease will not have an effect.  |
| NextTheme | ] + c | Switch to next theme. A theme consists of the set of colors used to draw the seq application and the set of icons used to represent different accent levels. |
| PrevTheme | [ + c | Move to the previous theme. A theme consists of the set of colors used to draw the seq application and the set of icons used to represent different accent levels. |
| NextSection | ] + s | Move to the next section within the arrangement.  If the next section is a group, then this mapping will move to the first section within that group.  |
| PrevSection | [ + s | Move to the previous section. Move to the next section within the arrangement.  If the previous section is a group, then this mapping will move to the last section within that group. |
| NewSectionAfter | Ctrl + ] | Create new section after the current section |
| NewSectionBefore | Ctrl + p | Create new section before the current section |
| BeatInputSwitch | Ctrl + b | This selects the current part's beats which can be increased or decreased with +/-.  Using this key combination again will move through selections Beats, Start Beats, Cycles and Start Cycles.  | |
| AccentInputSwitch | Ctrl + e | This selects the controls that determine the accent values and target.  Use +/- to increase and decrease the selections. SEE ACCENT CONTROLS |
| OverlayInputSwitch | Ctrl + o | This selects the inputs that control the overlay period/key.  SEE OVERLAY KEY CONTROLS |
| SetupInputSwitch | Ctrl + s | Select the inputs that control the midi message for each line.  Pressing this key combo repeatedly will move through the channel, target and value inputs.  |
| TempoInputSwitch | Ctrl + t | Select the inputs that control the tempo and subdivision. Press once to select the tempo input, press again to select the subdivisions input.   |
| OverlayStackToggle | Ctrl + u | Toggle the behaviour of the current overlay layer between three options: No association, press up, press down.  SEE OVERLAYS |
| ChangePart | Ctrl + c | Change the part of the section to either an existing part or a new part |
| ToggleArrangementView | Ctrl + a | Open the arrangement view when closed.  Focus the arrangement view while unfocused and open.  Press enter to move focus back to the grid.  While open and focused, close the arrangement view.  |
| NewLine | Ctrl + l | Create a new line with a value 1 greater than the previous line |
| New | Ctrl + n | Create a new sequence using the same template as the current sequence |
| Save | Ctrl + v | Save the current sequence.  If not previously saved, you will be prompted to name the new file.  The file will be saved in the directory from which you opened seq |
| ToggleAccentMode | n + a | Start Pattern Mode - Accent.  Use the facilities of pattern mode to increase or decrease the accent values of the line. SEE PATTERN MODE |
| ToggleWaitMode | n + w | Start Pattern Mode - Wait.  Use the facilities of pattern mode to increase or decrease the wait values of the line. SEE PATTERN MODE |
| ToggleGateMode | n + g |  Start Pattern Mode - Gate.  Use the facilities of pattern mode to increase or decrease the gate values of the line. SEE PATTERN MODE |
| ToggleRatchetMode | n + r |  Start Pattern Mode - Ratchet.  Use the facilities of pattern mode to increase or decrease the ratchet values of the line. SEE PATTERN MODE |
|RatchetInputSwitch | Ctrl + y | Select the inputs that control the ratchets for the current note. Press again to select the Span input. |
| ClearLine | c | Remove all notes from the current line from the current cursor position to the end |
| NoteRemove | d | Remove note at current position, and remove it from any stacked overlays if the current overlay is higher than the overlay of the current note |
| OverlayNoteRemove | x | Remove note from overlay at current position, allowing notes in lower layers to show through |
| TogglePlayEdit | b + e | Toggle play edit mode.  Press while playing to ensure the current overlay/part does not change while editing.  Press again to allow changing. |
| NoteAdd | f | Add note at current position |
| ReloadFile | b + r | Reload current file, any changes since the last save will be lost |
| ActionAddSpecificValue | b + v | Add specific value note to the grid. When cursor is above this note, +/- will affect the specific value of the note |
| CursorLeft | h | Move cursor left |
| CursorDown | j | Move cursor down |
| CursorUp | k | Move cursor up |
| CursorRight | l | Move cursor right |
| ToggleChordMode | o | Toggle chord mode.  SEE CHORD MODE |
| Yank | y | Copy current selection to buffer.  Copies all values of a visual selection or the value under cursor if no visual selection. |
| Paste | p | Paste the buffer at the position of the cursor |
| Quit | q | Quit the application |
| Undo | u | Undo last action |
| Redo | U | Redo last undone action |
| ToggleVisualMode | v | Toggle visual selection.  SEE VISUAL MODE |
| NextOverlay | { | Move to next overlay SEE OVERLAYS|
| PrevOverlay | } | Move to previous overlay SEE OVERLAYS|
| Enter | Enter | Confirm current action, move focus to the grid when elsewhere, or escape from visual mode |
| Escape | Esc | Cancel current action or exit mode, move focus to the grid when elsewhere, or escape from visual mode  |


## Pattern Mode Mappings

The keys of pattern mode have two different behaviours, value and fill.  The
default mode is Pattern Mode - Fill.

### PATTERN MODE - Fill

This is the default mode.  When seq starts, it starts in Pattern Mode - Fill.

Numbers will a note every X beats from the cursor to the end of the line.  If a
note already exists in that location the note will be removed.

EXAMPLE: With the cursor at the start of the line `1` will add a note at every
beat.  `1` will then remove a note at every beat.


| Mapping | Key Binding | Description |
|---------|-------------|-------------|
| NumberPattern | 1 | Add/remove a note every beat  |
| NumberPattern | 2 | Add/remove a note every 2nd beat |
| NumberPattern | 3 | Add/remove a note every 3rd beat |
| NumberPattern | 4 | Add/remove a note every 4th beat |
| NumberPattern | 5 | Add/remove a note every 5th beat |
| NumberPattern | 6 | Add/remove a note every 6th beat |
| NumberPattern | 7 | Add/remove a note every 7th beat |
| NumberPattern | 8 | Add/remove a note every 8th beat |
| NumberPattern | 9 |  Add/remove a note every 9th beat |


### PATTERN MODE - Value (Accent, Gate, Ratchet, Wait)

To enter pattern mode for a value, type `na` for accent, `nw` for wait, `nr`
for ratchet or `ng` for gate.


| Mapping | Key Binding | Description |
|---------|-------------|-------------|
| NumberPattern | shift+1 / ! | Increase value every beat |
| NumberPattern | shift+2 / @ | Increase value every 2nd beat |
| NumberPattern | shift+3 / # | Increase value every 3rd beat |
| NumberPattern | shift+4 / $ | Increase value every 4th beat |
| NumberPattern | shift+5 / % | Increase value every 5th beat |
| NumberPattern | shift+6 / ^ | Increase value every 6th beat |
| NumberPattern | shift+7 / & | Increase value every 7th beat |
| NumberPattern | shift+8 / * | Increase value every 8th beat |
| NumberPattern | shift+9 / ( | Increase value every 9th beat |
| NumberPattern | 1 | Decrease value every beat  |
| NumberPattern | 2 | Decrease value every 2nd beat  |
| NumberPattern | 3 | Decrease value every 3rd beat  |
| NumberPattern | 4 | Decrease value every 4th beat  |
| NumberPattern | 5 | Decrease value every 5th beat  |
| NumberPattern | 6 | Decrease value every 6th beat  |
| NumberPattern | 7 | Decrease value every 7th beat  |
| NumberPattern | 8 | Decrease value every 8th beat  |
| NumberPattern | 9 | Decrease value every 9th beat  |


## Chord Mode Mappings

Chord mode allows users to create and manipulate chords with a set of key mappings.

Example:  `tM` create a Major Triad Chord.  `7M` will add a major seventh to that
chord.  `]i` will invert the chord once, placing the root note 12 steps higher
at the top of the chord.

Some mappings exist in Pattern Fill mode as well.  `L` will move the entire
chord to the right one beat.  `A` will increase the accent value of every note
in the chord.

Pattern mode (value) also works differently in chord mode. Enter pattern
mode for accents with `na` and then use `shift+2` to increase the accent
value for every 2nd note in the chord.

Chord notes can be doubled.  Pressing `]d` once will double the first note of
the second.  Pressing `]d` again will double the second note of the while the
first note remains doubled.  If there are no more notes to double, `]d` will
remove all doubled notes.  `[d` behaves the same way but reversed.

Arpeggiate the notes of the chord with `]p` or `[p`.  These mappings will cycle
through the available arpeggiated patterns, at the moment there are only two
patterns: up and down.


| Mapping | Key Binding | Description |
|---------|-------------|-------------|
| MajorTriad | t + M | Add major triad chord |
| MinorTriad | t + m | Add minor triad chord |
| DiminishedTriad | t + d | Add diminished triad chord |
| AugmentedTriad | t + a | Add augmented triad chord |
| MinorSeventh | 7 + m | Add minor seventh |
| MajorSeventh | 7 + M | Add major seventh |
| AugFifth | 5 + a | Add augmented fifth |
| DimFifth | 5 + d | Add diminished fifth |
| PerfectFifth | 5 + p | Add perfect fifth |
| MinorSecond | 2 + m | Add minor second |
| MajorSecond | 2 + M | Add major second |
| MinorThird | 3 + m | Add minor third |
| MajorThird | 3 + M | Add major third |
| PerfectFourth | 4 + p | Add perfect fourth |
| MajorSixth | 6 + M | Add major sixth |
| Octave | 8 + p | Add octave |
| MinorNinth | 9 + m | Add minor ninth |
| MajorNinth | 9 + M | Add major ninth |
| DecreaseInversions | [ + i | Decrease chord |
| IncreaseInversions | ] + i | Increase chord |
| OmitRoot | 1 + o | Omit root note from chord |
| OmitSecond | 2 + o | Omit second note from chord |
| OmitThird | 3 + o | Omit third note from chord |
| OmitFourth | 4 + o | Omit fourth note from chord |
| OmitFifth | 5 + o | Omit fifth note from chord |
| OmitSixth | 6 + o | Omit sixth note from chord |
| OmitSeventh | 7 + o | Omit seventh note from chord |
| OmitOctave | 8 + o | Omit eighth note from chord |
| OmitNinth | 9 + o | Omit ninth note from chord |
| RemoveChord | D | Remove chord at current position |
| NextArpeggio | ] + p | Next arpeggio pattern |
| PrevArpeggio | [ + p | Previous arpeggio pattern |
| NextDouble | ] + d | Next double pattern |
| PrevDouble | [ + d | Previous double pattern |
| ConvertToNotes | n + n | Convert chord to individual notes |
| RotateRight | L | Move chord to the right |
| RotateLeft | H | Move chord to the left |
| RotateUp | K | Move chord up |
| RotateUp | K | Move chord down |


## Arrangement Mappings

Once having opened and focused the arrangement view with `ctrl+a` a new set of
mappings are available.  Some mappings are duplicated between the grid and the
arrangement view.


| Mapping | Key Binding | Description |
|---------|-------------|-------------|
| CursorUp |  k| Move the cursor to the previous arrangement section   |
| CursorDown | j | Move the cursor to the next arrangement section   |
| CursorLeft | h | Move the cursor to the section attribute to the left  |
| CursorRight | l | Move the cursor to the section attribute to the right   |
| Increase | + | Increase the value of the currently selected section attribute   |
| Decrease | - | Decrease the value of the currently selected section attribute   |
| GroupNodes | g | Group one or two parts together   |
| DeleteNode | d | Remove the current section attribute   |
| MovePartDown | J | Move the current section below the next section   |
| MovePartUp | K | Move the current section above the next section   |
| RenamePart | R | Rename the current part   |
| Escape | esc / enter  | Move focus back to the grid   |


## Overlay Key Mappings

Once having focused the overlay key inputs with `ctrl+o` a new set of mappings are available.


| Mapping | Key Binding | Description |
|---------|-------------|-------------|
| FocusWidth | :  | Select the width attribute   |
| FocusInterval | /  |  Select the interval attribute  |
| FocusShift | ^  |  Select the shift attribute |
| FocusStart | S |  Select the start attribute  |
| RemoveStart | s |  Remove the start attribute (set to zero)  |
| Increase | + |  Increase the selected value  |
| Decrease | - |  Decrease the selected value  |
| Escape | esc / enter  |  Return focus to the grid  |

