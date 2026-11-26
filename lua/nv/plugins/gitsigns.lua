return {
	"lewis6991/gitsigns.nvim",
	event = "BufEnter",
	keys = {
		{ "<leader>rs", function() require("nv.review").start() end, desc = "review start (diff vs PR base)" },
		{ "<leader>rx", function() require("nv.review").stop() end, desc = "review stop" },
		{ "<leader>rf", function() require("nv.review").files() end, desc = "review changed hunks" },
		{ "<leader>rp", function() require("gitsigns").preview_hunk_inline() end, desc = "review preview hunk" },
		{ "]h", function() require("gitsigns").nav_hunk("next") end, desc = "next hunk" },
		{ "[h", function() require("gitsigns").nav_hunk("prev") end, desc = "prev hunk" },
	},
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
