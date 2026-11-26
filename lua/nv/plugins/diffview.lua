return {
	"sindrets/diffview.nvim",
	cmd = { "DiffviewOpen", "DiffviewClose", "DiffviewFileHistory" },
	keys = {
		{ "<leader>rd", function() require("nv.review").diffview() end, desc = "review PR diff" },
		{ "<leader>rq", "<cmd>DiffviewClose<cr>", desc = "review close" },
		{ "<leader>rh", "<cmd>DiffviewFileHistory %<cr>", desc = "review file history" },
	},
	opts = {},
}
