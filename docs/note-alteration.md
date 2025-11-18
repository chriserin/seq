# Note Alteration

When `f` (fill) is pressed a note will be created at the cursors current
position.  This note will have default properties for accent and gate that are
ideal for triggering a drum on either a hardware or software drum machine.

A default note will have an accent value of 60 and a gate value of 20ms.  In
other words, a midi note on message will be sent with a velocity of 60 and 20ms
later, a note off message will be sent.

There are 4 different values that can be changed for each note:  Accent, Gate, Ratchet and Wait.

## Accent

Accent by default corresponds to the velocity of a midi note on message.  There
are by default 8 different accent levels. Pressing `A` will increase the accent
level for the note at the cursor position.  Pressing `a` will decrease the
accent level for this note.

There are different unicode characters for each accent level to help visually
recognize patterns being created and to give visual feedback when changing the
accent level of a note.  These unicode characters will be different for every
theme. See [Themes](key-mappings.md) (NextTheme/PrevTheme mappings).

## Gate

The gate attribute corresponds to the amount of time between the midi note on
message and the midi note off message.  By default there are 8 different gate
levels.  Pressing `G` will increase the gate level for the note at the cursor
position.  Pressing `g` will decrease the gate level for this note.

The initial gate level will be the shortest possible gate, which is an absolute
time of 20ms.  Each gate level above this level will create a gate value that
is a percentage of the interval created by the tempo and subdivision.

For instance, a gate level of 5 corresponds to a percentage value of 0.5.  When
the tempo is 120 and the subdivision value is 2, then there are 240
subdivisions per minute and 250ms between each beat.  The gate length given
these factors is 50% of 250ms or 125ms.

When a gate is increased or decreased a small symbol below the note will
change.

### Long Gates

Pressing `E` will increase the gate level above the length of the current Beat.
Pressing `e` will similarly decrease the gate level.  Gates up to 16 beats are
currently possible.

When the gate length is longer than the beat interval a bar extending from the
note will be drawn to indicate it's length.

## Ratchet

Ratchet corresponds to the number of midi note on messages that will be sent
for the note.  There are 8 different Ratchet levels and the initial ratchet
level for a note will be 1 corresponding to 1 midi note on message (a hit).  Pressing
`R` will increase the ratchet level for the note at the cursor position.
Pressing `r` will decrease the ratchet level for this note.

For each ratchet above level 2 the length of the gate will be 20ms.  The
ratchet hits will be evenly distributed through the beat interval.

When the ratchet level is increased or decreased a small symbol corresponding
to the ratchet level will change.

## Wait

The wait attribute corresponds to the amount of time between the beat and the
midi note on message.  By default there are 8 different wait levels.  Each
level corresponds to a percentage of the beat interval, from 0 to 54.  Pressing
`W` will increase the wait level for the note at the cursor position.  Pressing
`w` will decrease the wait level for this note.

Combined with pattern mode this can be useful for creating swing effects.
