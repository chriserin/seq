# Beat Pattern Creation

Creating beat patterns is simple and quick in seq with pattern mode.  seq begins in pattern mode - fill, where number keys correspond to beat intervals.  

## Example

Here is a blank pattern with some named lines that correspond to a bass drum
(BD), a snare (SN) and a high hat (H1).

```
 BDK┤                                  
 SN │                                  
 H1 │                                  
```

With the cursor at the first beat of the bass drum, press `8` to create a note every 8 beats.

```
  BDK┤▧       ▧       ▧       ▧         
  SN │
  H1 │

```

Move the cursor down to the snare line with `j` and press `4` and `8` to create notes on every 8th beat starting at the 5 beat.

```
  BDK│▧       ▧       ▧       ▧         
  SN ┤    ▧       ▧       ▧       ▧
  H1 │
```

For the high hat, move down one line with `j` and press `1` to create a high hat hit on every beat.

```
  BDK│▧       ▧       ▧       ▧         
  SN │    ▧       ▧       ▧       ▧
  H1 ┤▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧

```

6 keystrokes to create a very simple beat.

To create some accents on the high hats, press `na` to enter Pattern Mode -
Accent and `shift+2` to increase the velocity of the high hat hit on every
other beat.  Press `enter` to escape from this mode.

```
  BDK┤▧       ▧       ▧       ▧         
  SN │    ▧       ▧       ▧       ▧
  H1 │▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧

```

That's t
additional key strokes to create that accent pattern.
