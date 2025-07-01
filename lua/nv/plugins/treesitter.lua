return {
	"nvim-treesitter/nvim-treesitter",
	event = "BufEnter",
	branch = 'master',
	lazy = false,
	build = ":TSUpdate",
	config = function()
		require("nvim-treesitter.configs").setup{
			ensure_installed = {
				"lua",
				"go",
				"rust",
				"dart",
				"yaml",
				"dockerfile",
				"gitignore",
				"json",
			},
			sync_install = true,
			auto_install = true,
		}
	end
}
