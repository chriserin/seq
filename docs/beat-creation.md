# Beat Pattern Creation

Creating beat patterns is simple and quick in sq with pattern mode. sq begins in pattern mode — fill, where number keys correspond to beat intervals.

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

To create some accents on the high hats, press `na` to enter Pattern Mode — Accent
and `shift+2` to increase the velocity of the high hat hit on every
other beat. Press `enter` to escape from this mode.

```
  BDK┤▧       ▧       ▧       ▧
  SN │    ▧       ▧       ▧       ▧
  H1 │▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧▤▧

```

That's 4 additional key strokes to create that accent pattern.

## Euclidean Rhythms

For more complex and mathematically-distributed patterns, you can use the Euclidean rhythm generator. This creates evenly-spaced note patterns that are useful for polyrhythms and interesting drum patterns.

### Creating a Euclidean Pattern

1. Position your cursor at the start of a line
2. Press `bu` to activate the Euclidean rhythm generator
3. Type the number of hits you want (e.g., "5" for 5 hits)
4. Press `Enter` to generate the pattern

For example, to create a 5-hit pattern over 16 beats on a percussion line:

```
  PC │
```

Press `bu`, type `5`, and press `Enter`:

```
  PC │▧  ▧  ▧  ▧  ▧
```

The hits are distributed as evenly as possible across the available beats.

### Using with Visual Selection

You can also apply Euclidean rhythms to a specific range using visual mode:

1. Press `v` to enter visual mode
2. Use `l` to expand your selection to the desired length
3. Press `bu` to activate the Euclidean generator
4. Type the number of hits
5. Press `Enter` to generate

This allows you to create Euclidean patterns of any length within your sequence.

### Common Euclidean Patterns

- **Euclidean(3, 8)**: Classic tresillo pattern used in Latin music
- **Euclidean(5, 8)**: Cuban cinquillo rhythm
- **Euclidean(5, 12)**: Common pattern for complex time signatures
- **Euclidean(7, 16)**: Asymmetric pattern for experimental rhythms

Euclidean rhythms work on the current overlay, so you can layer multiple Euclidean patterns with different hit counts to create complex polyrhythmic structures.
