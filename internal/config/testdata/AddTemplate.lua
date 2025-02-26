local seq = require("seq")

local Channel = 10
local C1 = 36
local MessageType = "Note"

seq.addtemplate({
	name = "Drums",
	maxgatelength = 32,
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
