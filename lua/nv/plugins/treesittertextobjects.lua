return {
	"nvim-treesitter/nvim-treesitter-textobjects",
	dependencies = { "nvim-treesitter/nvim-treesitter" },
	config = function()
		require('nvim-treesitter').setup {
			textobjects = {
				move = {
					enable = true,
					set_jumps = true,
					goto_next_start = {
						["]f"] = "@function.outer",
					},
					goto_previous_start = {
						["[f"] = "@function.outer",
					},
				},
			},
		}
	end
}
