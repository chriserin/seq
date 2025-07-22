# Arrangement

Arrangements consist of parts, sections and groups. The Arrangement View is
shown when pressing `ctrl+a`.

## Sections and Parts

An arrangement consists of sections.  These sections are sequenced one after
the other as listed in the Arrangement View.  Each section contains a part that
could be the same or different than other sections.

Initially, there is one section and one part.  In the arrangement view that
looks like this:

```
Section ────────────────────⟳ Amount─────⟳ Start──Start Beat──────⟳ Keep
   Part 1  ● ───────────────────────1───────────1───────────0───────────-

```

To create another section press `ctrl+]` and accept the "Choose Part" prompt
with `enter`.  Now there are two sections in the arrangement view:

```
 Section ────────────────────⟳ Amount─────⟳ Start──Start Beat──────⟳ Keep
   Part 1                           1           1           0           -
   Part 2  ● ───────────────────────1───────────1───────────0───────────-

```

"Part 1" will play for 1 key cycle and "Part 2" will also play for one key
cycle.  You can add another section but choose an already existing part by
pressing `ctlr+]` and selecting an existing part with `+` before accepting
with `enter`.  This will create an arrangement that looks like this:

```
 Section ────────────────────⟳ Amount─────⟳ Start──Start Beat──────⟳ Keep
   Part 1                           1           1           0           -
   Part 2                           1           1           0           -
   Part 1  ● ───────────────────────1───────────1───────────0───────────-

```

Sections are just a thin wrapper around parts that determines how many key
cycles a part will play along with a set of attributes described below.

Move the section selection down and up with `j` and `k`.

## Groups

A group can be used to iterate over a sequence of parts multiple times.  When
the Arrangement View is focused, press `g` to create a group that consists of
the currently focused section and the section after it, if a following section
exists.

When "Part 2" of our working example is currently selected, pressing `g` will
alter an arrangement to look like this:

```
 Section ────────────────────⟳ Amount─────⟳ Start──Start Beat──────⟳ Keep
   Part 1                           1           1           0           -
   Group ────────────────────────── 1────────────────────────────────────
   ├─ Part 2  ●                     1           1           0           -
   └─ Part 1                        1           1           0           -
```

The only attribute of the group is the "Cycles Amount", change this value with `+/-`.

Creating a group has the effect of creating an arrangement tree.

## Section Attributes

Each section has four attributes, Cycle Amount, Cycle Start, Start Beat and
Cycle Keep.  Move left and right between these attributes with `h` and `l`.


### Cycles Amount  /  ⟳ Amount

The number of Key Cycles a section will play is determined by the Cycles
Amount.  This defaults to 1 as a way to force a more considered choice for the
sequence.  The Cycles Amount is the left most attribute and is selected by
default when the arrangement view is focused initially, allowing quick
selections of higher Cycle Amount values with the `+` key.

### Cycles Start  /  ⟳ Start

The current number of Key Cycles determines which Overlay is active.  A section
can differ from another section playing the same part by using a different Key
Cycle at the start.  By starting with a key cycle of 2 rather than 1, a part
will begin with overlay `2/1`.

### Start Beat

A part does not need to start at the first beat.  Using the start beat
attribute, a section can start at any beat of the part.  This can be used to
inject oddly timed breaks between sections or to create fills at the starts of
sections.

### Keep Cycles

If a section is within a group and the group iterates multiple times, it is
possible to maintain the Key Cycles between plays within a group.  A "✔" as a
Keep Cycles value indicates that Key Cycles for that section will persist
between group plays.

For the following arrangement: 

```
 Section ────────────────────⟳ Amount─────⟳ Start──Start Beat──────⟳ Keep
   Part 1                           1           1           0           -
   Group                            2
   ├─ Part 2  ● ─────────────────── 1───────────1───────────0───────────✔
   └─ Part 1                        1           1           0           -

```

The group will play twice for a sequence of "Part 2", "Part 1", "Part 2", "Part
1".  The first time "Part 2" plays the Key Cycles start with a value of 1.  The
second time "Part 2" plays the Key Cycles start with a value of 2.


### Moving Sections 

When there are multiple sections it is possible to move the currently selected
section up or down within the arrangement tree with `J` or `K`.  When a group
follows a section that you want to move down then pressing `J` will move the
section _into_ the group.  To arrange the section so that the section follows
the group keep press `J` until the section is below and exists at the same tree
level.

For instance, given the following arrangement:

```
Section ────────────────────⟳ Amount─────⟳ Start──Start Beat──────⟳ Keep
   Part 1  ● ────────────────────── 1───────────1───────────0───────────-
   Group                            2
   ├─ Part 2                        1           1           0           ✔
   └─ Part 1                        1           1           0           -

```

Pressing `J` will alter the arrangement to be:

```
  Group                            2
   ├─ Part 1  ● ─────────────────── 1───────────1───────────0───────────-
   ├─ Part 2                        1           1           0           ✔
   └─ Part 1                        1           1           0           -

```

Pressing `J` again will alter the arrangement to be:

```
 Group                            2
   ├─ Part 2                        1           1           0           ✔
   ├─ Part 1  ● ─────────────────── 1───────────1───────────0───────────-
   └─ Part 1                        1           1           0           -

```

Pressing `J` again will alter the arrangement to be:

```
Section ────────────────────⟳ Amount─────⟳ Start──Start Beat──────⟳ Keep
   Group                            2
   ├─ Part 2                        1           1           0           ✔
   ├─ Part 1                        1           1           0           -
   └─ Part 1  ● ─────────────────── 1───────────1───────────0───────────-

```

Pressing `J` a final time will alter the arrangement to be:

```
 Section ────────────────────⟳ Amount─────⟳ Start──Start Beat──────⟳ Keep
   Group                            2
   ├─ Part 2                        1           1           0           ✔
   └─ Part 1                        1           1           0           -
   Part 1  ● ────────────────────── 1───────────1───────────0───────────-

```

Our first section has now been moved to the bottom of the arrangement and is
now the final section.

Pressing `K` 4 times will move the selected section up four times so that it
once again is the beginning section of the sequence.
