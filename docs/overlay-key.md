# Overlay

The intent of overlays is to add variation to a sequence at mathematically
determined points of the sequence. Overlays may be interesting and useful to
some people and over complex and not worth investing in to other people.
Either is fine, music is whatever we want it to be.

## Overview

Imagine two pieces of transparent paper one on top of the other. Each has a
tic-tac-toe board and each represents only one move. From above we can see
both the first move X and the next move O. By layering more transparencies we
can fill every section of the tic-tac-toe grid.

The principle is the same here. We want to create new sequences by layering
variations on top of the original sequence.

We determine when to apply these different layers with an overlay key. The
first overlay key that we see is `1/1`. This is the bottom overlay or root
overlay and it is applied on every iteration of the sequence.

An overlay key of `2/1` is applied every other cycle. An overlay key of `3/1`
is applied every third cycle.

Overlay keys can be complex, be most of the time they will be used simply.
Before describing the complexity, here are two tables that illustrate when each
overlay key is applied in a simple system.

### Every 1st, second, third and fourth iteration

| Key Cycle | `1/1` | `2/1` | `3/1` | `4/1` |
| --------- | ----- | ----- | ----- | ----- |
| 1         | X     |       |       |       |
| 2         | X     | X     |       |       |
| 3         | X     |       | X     |       |
| 4         | X     | X     |       | X     |
| 5         | X     |       |       |       |
| 6         | X     | X     | X     |       |
| 7         | X     |       |       |       |
| 8         | X     | X     |       | X     |

### The first of four, the second of four, the third of four, the fourth of four

| Key Cycle | `1/4` | `2/4` | `3/4` | `4/4` |
| --------- | ----- | ----- | ----- | ----- |
| 1         | X     |       |       |       |
| 2         |       | X     |       |       |
| 3         |       |       | X     |       |
| 4         |       |       |       | X     |
| 5         | X     |       |       |       |
| 6         |       | X     |       |       |
| 7         |       |       | X     |       |
| 8         |       |       |       | X     |

## Overlay Key Definition

The overlay key consists of four parts, two of which are initially visible.

The full key definition is:

```
1:1/1S1
```

The parts of which are:

```
1 - Shift - Segment of the interval to apply the overlay
: - Width Indicator
1 - Width - The duration of key cycles to apply the overlay once applied
/ - Interval Indicator
1 - Interval - The number of key cycles needed to reconsider application
S - Start Indicator
1 - Start - The number of key cycles to wait before allowing application of the overlay
```

Initially the overlay key is `1/1`. This is the root overlay and will be
applied to the first key cycle of every 1st key cycle. More simply, this
overlay will be applied every key cycle. If not using more specific overlays,
this will always be applied and no understanding of the overlay key is needed.

An overlay key of `1/4` describe an iteration length of 4 and a shift of 1, or
in other words, the first of every four. This will apply the overlay to the
first of every four key cycles.

If we have four overlay keys for each segment of the interval: `1/4`, `2/4`,
`3/4`, `4/4` then they would be applied according to the below table:

| Key Cycle | `1/4` | `2/4` | `3/4` | `4/4` |
| --------- | ----- | ----- | ----- | ----- |
| 1         | X     |       |       |       |
| 2         |       | X     |       |       |
| 3         |       |       | X     |       |
| 4         |       |       |       | X     |
| 5         | X     |       |       |       |
| 6         |       | X     |       |       |
| 7         |       |       | X     |       |
| 8         |       |       |       | X     |

If the interval is only 1, with a shift of 2 as in `2/1`, then it will be
applied every other cycle. A key of `4/2` is a fraction that gets reduced to
`2/1` and will also play every other cycle. `3/1` is applied every third cycle.
`4/1` is applied every fourth cycle and so on.

| Key Cycle | `1/1` | `2/1` | `3/1` | `4/1` |
| --------- | ----- | ----- | ----- | ----- |
| 1         | X     |       |       |       |
| 2         | X     | X     |       |       |
| 3         | X     |       | X     |       |
| 4         | X     | X     |       | X     |
| 5         | X     |       |       |       |
| 6         | X     | X     | X     |       |
| 7         | X     |       |       |       |
| 8         | X     | X     |       | X     |

For keys that are fractions not evenly divisible like `3/2` then the interval
is doubled until it is greater than the shift. `3/2` is the same as `3/4`.

## WIDTH

By default, an overlay key will only be applied for one cycle before we
determine if it should be applied again. To apply the overlay for two cycles
we can increase the width.

`1:1/1` is the same as `1/1`, the number after the `:` indicates the width. For
simple keys this might be redundant. `2:2/1` will be applied every cycle the
same as `1/1`. For other keys it has more meaning. `3:2/4` will be applied
for the third and fourth out of every four cycles.

| Key Cycle | `1:2/4` | `2:2/4` | `3:2/4` | `4:2/4` |
| --------- | ------- | ------- | ------- | ------- |
| 1         | X       |         |         |         |
| 2         | X       | X       |         |         |
| 3         |         | X       | X       |         |
| 4         |         |         | X       | X       |
| 5         | X       |         |         | X       |
| 6         | X       | X       |         |         |
| 7         |         | X       | X       |         |
| 8         |         |         | X       | X       |

## START

By default, an overlay key will be applied at the first cycle that matches.
`3/4` will be applied at the third cycle. This is the same as `3/4S1`. The `S`
indicates start. We can delay the application of the overlay by providing a
start value that is higher than the first matching cycle. `3/4S4` will be
applied at the 7th cycle because the program will not check if the key matches
until the 4th cycle.

| Key Cycle | `1:2/4S4` | `2:2/4S4` | `3:2/4S4` | `4:2/4S4` |
| --------- | --------- | --------- | --------- | --------- |
| 1         |           |           |           |           |
| 2         |           |           |           |           |
| 3         |           |           |           |           |
| 4         |           |           |           |           |
| 5         | X         |           |           |           |
| 6         | X         | X         |           |           |
| 7         |           | X         | X         |           |
| 8         |           |           | X         | X         |

Another useful example is `2:1S8` which applies every other cycle but only after
the 8th cycle.

## Order And Stacking

Without any stacking options, only one overlay will be applied at any given
time. If two overlays both match, like `1/1` and `2/1` on key cycle 2, then
the overlay that matches less often will be on top and the overlay that matches
the most will be on the bottom. In this case only `2/1` will be applied for
key cycle 2.

The **stack** options of the overlay determine how that overlay relates to
other overlays. Overlays can either **press up**, **press down** or **stand
alone**. Cycle through stack options with `ctrl+u`

- **press up** - ↑̅ - Other overlays will be stacked on top of this overlay if they both match.
- **press down** - ↓̲ - This matching overlay will be stacked on top of the overlay underneath it, even if it doesn't match.
- **stand alone** - This overlay does not affect whether other overlays are applied.

By default the root overlay, `1/1`, always starts with the option **press up**.
All other overlays will applied on top of the root overlay.

The **press down** option is useful when we want to aggregate the variations as
we move through the cycles.

| Key Cycle | `1/4` ↓̲ | `2/4` ↓̲ | `3/4` ↓̲ | `4/4` ↓̲ |
| --------- | ------- | ------- | ------- | ------- |
| 1         | X       |         |         |         |
| 2         | X       | X       |         |         |
| 3         | X       | X       | X       |         |
| 4         | X       | X       | X       | X       |
| 5         | X       |         |         |         |
| 6         | X       | X       |         |         |
| 7         | X       | X       | X       |         |
| 8         | X       | X       | X       | X       |

As you can see, as we move through the cycles the **press down** option has an
accumulative effect. With the **press down** option an overlay can be applied
by an overlay directly above it if the above overlay has the **press down**
option even if the overlay does not match.
