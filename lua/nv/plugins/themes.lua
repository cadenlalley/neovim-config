return {
	"navarasu/onedark.nvim",
	priority = 1000,
	opts = {},
	config = function()
		require('onedark').setup {
			style = 'dark',
			transparency = true,
			term_colors = true,
			ending_tildes = true
		}
		require('onedark').load()
	end
}
