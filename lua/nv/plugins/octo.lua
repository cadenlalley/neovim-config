return {
	'pwntester/octo.nvim',
	dependencies = {
		'nvim-lua/plenary.nvim',
		'folke/snacks.nvim',
		'nvim-tree/nvim-web-devicons',
	},
	opts = {
		picker = "snacks",
		enable_builtin = true,
		use_local_fs = true,
	}
}
