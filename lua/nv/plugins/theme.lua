return {
	"olimorris/onedarkpro.nvim",
	priority = 1000,
	opts = {},
	config = function()
		require('onedarkpro').setup({})
		vim.cmd("colorscheme onedark")
	end
}
