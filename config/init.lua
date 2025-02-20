local seq = require("seq")

local Channel = 10
local C1 = 36
local C3 = 60
local MessageType = "Note"

seq.addtemplate({
	name = "Drums",
	maxGateLength = 32,
	uistyle = "plain",
	lines = {
		{ Channel, MessageType, C1 },
		{ Channel, MessageType, C1 + 1 },
		{ Channel, MessageType, C1 + 2 },
		{ Channel, MessageType, C1 + 3 },
		{ Channel, MessageType, C1 + 4 },
		{ Channel, MessageType, C1 + 5 },
		{ Channel, MessageType, C1 + 6 },
		{ Channel, MessageType, C1 + 7 },
	},
})

local pianoLines = {}

local index = 1
for noteValue = C3, C1, -1 do
	pianoLines[index] = { 1, "Note", noteValue }
	index = index + 1
end

seq.addtemplate({
	name = "Piano2",
	maxGateLength = 32,
	uistyle = "blackwhite",
	lines = pianoLines,
})

seq.addinstrument({
	name = "Prophet 10",
	controlchanges = {
		{ 7, 120, "MASTER VOLUME" },
		{ 9, 120, "OSC B FREQUENCY" },
		{ 14, 127, "OSC B FINE TUNE" },
		{ 15, 1, "OSC A SAW ON/FF" },
		{ 20, 1, "OSC A SQUARE ON/OFF" },
		{ 21, 120, "OSC A PULSE WIDTH" },
		{ 22, 120, "OSC B PULSE WIDTH" },
		{ 23, 1, "OSC SYNC ON/OFF" },
		{ 24, 1, "OSC B LOW FREQ ON/OFF" },
		{ 25, 1, "OSC B KEYBOARD ON/OFF" },
		{ 26, 120, "GLIDE RATE" },
		{ 27, 120, "OSC A LEVEL" },
		{ 28, 120, "OSC B LEVEL" },
		{ 29, 120, "NOISE LEVEL" },
		{ 30, 1, "OSC B SAW ON/OFF" },
		{ 31, 120, "RESONANCE" },
		{ 35, 2, "FILTER KEYBOARD TRACK OFF/HALF/FULL" },
		{ 41, 1, "FILTER REV SELECT" },
		{ 46, 120, "LFO FREQUENCY" },
		{ 47, 120, "LFO INITIAL AMOUNT" },
		{ 52, 1, "OSC B TRI ON/OFF" },
		{ 53, 120, "LFO SOURCE MIX" },
		{ 54, 1, "LFO FREQ A ON/OFF" },
		{ 55, 1, "LFO FREQ B ON/OFF" },
		{ 56, 1, "LFO FREQ PW A ON/OFF" },
		{ 57, 1, "LFO FREQ PW B ON/OFF" },
		{ 58, 1, "LFO FILTER ON/OFF" },
		{ 59, 127, "POLY MOD FILT ENV AMOUNT" },
		{ 60, 120, "POLY MOD OSC B AMOUNT" },
		{ 61, 1, "POLY MOD FREQ A ON/OFF" },
		{ 62, 1, "POLY MOD PW ON/OFF" },
		{ 63, 1, "POLY MOD FILTER ON/OFF" },
		{ 70, 11, "PITCH WHEEL RANGE" },
		{ 71, 3, "RETRIGGER AND UNISON ASSIGN" },
		{ 73, 120, "CUTOFF" },
		{ 74, 127, "BRIGHTNESS" },
		{ 85, 127, "VINTAGE" },
		{ 86, 1, "PRESSURE FILTER" },
		{ 87, 1, "PRESSURE LFO" },
		{ 89, 120, "ENVELOPE FILTER AMOUNT" },
		{ 90, 1, "ENVELOPE FILTER VELOCITY ON/OFF" },
		{ 102, 1, "ENVELOPE VCA VELOCITY ON/OFF" },
		{ 103, 120, "ATTACK FILTER" },
		{ 104, 120, "ATTACK VCA" },
		{ 105, 120, "DECAY FILTER" },
		{ 106, 120, "DECAY VCA" },
		{ 107, 120, "SUSTAIN FILTER" },
		{ 108, 120, "SUSTAIN VCA" },
		{ 109, 120, "RELEASE FILTER" },
		{ 110, 120, "RELEASE VCA" },
		{ 111, 1, "RELEASE ON/OFF" },
		{ 112, 1, "UNISON ON/OFF" },
		{ 113, 10, "UNISON VOICE COUNT" },
		{ 114, 7, "UNISON DETUNE" },
		{ 116, 1, "OSC B SQUARE ON/OFF" },
		{ 117, 1, "LFO SAW ON/OFF" },
		{ 118, 1, "LFO TRI ON/OFF" },
		{ 119, 1, "LFO SQUARE ON/OFF" },
	},
})
