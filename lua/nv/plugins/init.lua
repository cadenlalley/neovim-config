return {
	{ "nvim-lua/plenary.nvim" },
	{
		"christoomey/vim-tmux-navigator",
		cond = function()
			return vim.env.HERDR_PANE_ID == nil
		end,
	},
}
