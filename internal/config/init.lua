local seq = require("seq")
local instruments_path = CONFIG_DIR .. "/instruments.lua"
local chunk = loadfile(instruments_path)
if chunk then
	chunk()
end

local Channel = 10
local C1 = 36
local C3 = 60
local MessageType = "Note"

seq.addtemplate({
	name = "Drums",
	maxgatelength = 32,
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

seq.addtemplate({
	name = "VMNA",
	maxgatelength = 1,
	uistyle = "plain",
	lines = {
		{ Channel, MessageType, C1, "BD" },
		{ Channel, MessageType, 48, "D1" },
		{ Channel, MessageType, 41, "D2" },
		{ Channel, MessageType, 58, "MU" },
		{ Channel, MessageType, 40, "SN" },
		{ Channel, MessageType, 49, "H1" },
		{ Channel, MessageType, 51, "O2" },
		{ Channel, MessageType, 42, "H2" },
		{ Channel, MessageType, 44, "O2" },
		{ Channel, MessageType, 39, "CL" },
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
	seqtype = "polyphony",
	maxgatelength = 32,
	uistyle = "blackwhite",
	lines = pianoLines,
})
