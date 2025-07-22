# Actions

Each line has a playback cursor.  There are various actions that can be added
to a line that will manipulate the playback cursor in different ways. Each of
these actions can be added to a line with a keycombo begining with `s`.

## Line Reset

Add a Line Reset action at the cursor location with `ss`.  When the playback
cursor reaches this action, it will reset to the first beat.

## Line Reset All

Add a Line Reset All action at the cursor location with `sS`.  When the playback
cursor reaches this action, it will reset all playback cursors to the first beat.

## Line Bounce

Add a Line Bounce action at the cursor location with `sb`.  When the playback
cursor reaches this action, it will go backwards until reaching the first beat.
When reaching the first beat it will change directions and go forwards until
again reaching the bounce action.  Unless some other action intervenes it will
bounce back and forth between the first beat and the action.

## Line Bounce All

Add a Line Bounce All action at the cursor location with `sB`.  When the
playback cursor reaches this action, all playback cursors will go backwards
until reaching the first beat. When each playback cursor reaches the first beat
they will change directions and go forwards until the line with the playback
cursor again reaches the Line Bounce All action.  Unless some other action
intervenes each playback cursor will bounce back and forth between the first
beat and this action.

## Line Skip Beat 

Add a Line Skip Beat action at the cursor location with `sk`.  When the
playback cursor reaches this action, it will skip this beat.  This will place
this line's playback cursor ahead of other playback cursors.

## Line Skip Beat All

Add a Line Skip Beat All action at the cursor location with `sK`.  When the
playback cursor reaches this action, each playback cursor will move forward one
beat.

## Line Delay

Add a Line Delay action at the cursor location with `sz`.  When the playback
cursor reaches this action it will pause on the beat before this action and
play the note at that location repeatedly until either interrupted by another
action or that part or overlay changes.

## Line Reverse

Add a Line Reverse action at the cursor location with `sr`.  When the playback
cursor reaches this action the playback cursor moves in this opposite
direction.  When the playback cursor reaches the start of the line it will
reset to the location of the action and will continue to move backwards.
