# Key Cycles

Every time a playback cursor returns to the start of the Key line, the key cycles
for this section are increased.  Because seq has playback cursors for each line, one line must
 be the Key line. This is signified by the `K` next to the line number/name.

Key Cycles are central to both the Arrangement View and the Overlay Key.  For
Arrangement, we use the count of Key Cycles to determine when to move to the
next part.  For Overlays, we use the Key Cycles to determine which Overlay to
use and when to use it.

Without any actions, a playback cursor will reach the end of the pattern and
return to the start of the line, increasing the Key Cycles.  If the Arrangement
has a setting of just 1 Cycle, which is the default, then this section is done
and the arrangement moves to the next section, resetting the Key Cycles.  Or if
there are no more sections, then the sequence is over and playback stops.

The amount of key cycles for a section can be increased with `ctrl+k` which
selects the key cycle inputs.

There are various actions (SEE ACTIONS) that you can place in any line to manipulate the
playback cursor for that line.  For instance, if you add the line reset action
(`s`) to the line with the Key indicator at beat 5, then when playback occurs
only the first 4 beats of that line will play before the playback cursor for
that line resets to zero which will increase the amount of Key Cycles.
