vim.g.mapleader = " "

-- automatically reload files changed outside Neovim
vim.o.autoread = true

-- trigger checktime when gaining focus, switching buffers, or idling
vim.api.nvim_create_autocmd({ "FocusGained", "BufEnter", "CursorHold", "CursorHoldI" }, {
	callback = function()
		vim.cmd("checktime")
	end,
})

require("nv.lazy")
require("nv.core")
