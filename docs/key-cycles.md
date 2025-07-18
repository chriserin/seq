# Key Cycles

Key Cycles are an integral concept that holds together both Arrangements and
Overlays.  When the sequence plays, a playback cursor scrolls across the screen
as the beat advances.  When the playback cursor gets to the end it returns to
the start of the line.  Whenever the playback cursor returns to the
beginning of the line the number of key cycles increases.

However, there is a playback cursor for each line, which necessitates choosing
a line as the line that will increase the key cycles when the playback cursor
returns to the start.
