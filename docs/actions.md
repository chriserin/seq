# Actions

Each line has a playback cursor. There are various actions that can be added
to a line that will manipulate the playback cursor in different ways. Each of
these actions can be added to a line with a keycombo begining with `s`.

## Line Reset

`ss` — Add a Line Reset Action

When the playback cursor reaches this action, it will reset to the first beat.

## Line Reset All

`sS` — Add a Line Reset All Action

When the playback cursor reaches this action, it will reset all playback cursors
to the first beat.

## Line Bounce

`sb` — Add a Line Bounce Action

When the playback cursor reaches this action, it will go backwards until
reaching the first beat. When reaching the first beat it will change directions
and go forwards until again reaching the bounce action. Unless some other
action intervenes it will bounce back and forth between the first beat and the
bounce action.

## Line Bounce All

`sB` — Add a Line Bounce All Action

When the playback cursor reaches this action, all playback cursors will go
backwards until reaching the first beat. When each playback cursor reaches the
first beat they will change directions and go forwards until the line with the
playback cursor again reaches the Line Bounce All action. Unless some other
action intervenes each playback cursor will bounce back and forth between the
first beat and this action.

## Line Skip Beat

`sk` — Add a Line Skip Action

When the playback cursor reaches this action, it will skip this beat. This will
place this line's playback cursor ahead of other playback cursors.

## Line Delay

`sz` — Add a Line Delay Action

When the playback cursor reaches this action it will pause on the beat before
this action and play the note at that location repeatedly until either
interrupted by another action or that part or overlay changes.

## Line Reverse

`sr` — Add a Line Reverse

When the playback cursor reaches this action the playback cursor moves in the
opposite direction. When the playback cursor reaches the start of the line it
will reset to the location of the action and will continue to move backwards.
