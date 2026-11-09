return {
	"nvim-treesitter/nvim-treesitter",
	lazy = false,
	build = ":TSUpdate",
	config = function()
		require("nvim-treesitter").setup()

		vim.api.nvim_create_autocmd("FileType", {
			callback = function(args)
				pcall(vim.treesitter.start, args.buf)
			end,
		})

		vim.api.nvim_create_autocmd("User", {
			pattern = "TSUpdate",
			once = true,
			callback = function()
				local parsers = {
					"lua",
					"go",
					"rust",
					"dart",
					"yaml",
					"dockerfile",
					"gitignore",
					"json",
					"python",
				}
				local installed = require("nvim-treesitter.config").get_installed()
				local to_install = vim.tbl_filter(function(p)
					return not vim.list_contains(installed, p)
				end, parsers)
				if #to_install > 0 then
					vim.cmd("TSInstall " .. table.concat(to_install, " "))
				end
			end,
		})
	end,
}
