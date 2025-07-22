# Beat Pattern Creation

Creating beat patterns is very simple and quick in seq with pattern mode. 

## Example

Here I have a blank pattern with some named lines that correspond to a bass drum
(BD), a snare (SN) and a high hat (H1).

```
 BDK┤                                  
 SN │                                  
 H1 │                                  
```

With the cursor at the first beat of the bass drum I can press `8` to create a note every 8 beats.

```
  BDK┤▧       ▧       ▧       ▧         
  SN │
  H1 │

```

I can move the cursor down to the snare line with `j` and press `4` and `8` to create notes on every 8th beat starting at the 5 beat.

```
  BDK│▧       ▧       ▧       ▧         
  SN ┤    ▧       ▧       ▧       ▧
  H1 │
```

For the high hat I can move down one line with `j` and press `1` to create a high hat hit on every beat.

```
  BDK│▧       ▧       ▧       ▧         
  SN │    ▧       ▧       ▧       ▧
  H1 ┤▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧

```

6 keystrokes to create a very simple beat.

To create some accents on the high hats I can enter Pattern Mode - Accent with
`na` and `shift+2` twice to increase the velocity of the high hat hit on every
other beat.  Press `enter` to escape from this mode.

```
  BDK┤▧       ▧       ▧       ▧         
  SN │    ▧       ▧       ▧       ▧
  H1 │▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧

```

It's subtle but you can see a different unicode character in the high hat row to represent the accent.

That's 5 additional key strokes to create that accent pattern.  And now you
have a beat you can send to the drum machine of your choice.
