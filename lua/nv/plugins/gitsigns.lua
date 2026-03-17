return {
	"lewis6991/gitsigns.nvim",
	event = "BufEnter",
	config = function()
		require("gitsigns").setup({
			current_line_blame = true,
			current_line_blame_opts = {
				virt_text = true,
				virt_text_pos = 'eol',
				delay = 0,
				ignore_whitespace = true
			}
		})
	end
}
