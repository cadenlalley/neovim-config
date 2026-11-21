return {
	"sindrets/diffview.nvim",
	cmd = { "DiffviewOpen", "DiffviewClose", "DiffviewFileHistory" },
	keys = {
		{ "<leader>rd", "<cmd>DiffviewOpen origin/HEAD...HEAD --imply-local<cr>", desc = "review PR diff" },
		{ "<leader>rq", "<cmd>DiffviewClose<cr>", desc = "review close" },
		{ "<leader>rh", "<cmd>DiffviewFileHistory %<cr>", desc = "review file history" },
	},
	opts = {},
}
