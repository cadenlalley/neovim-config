return {
	{
		"mason-org/mason.nvim",
		lazy = false,
		opts = {},
		config = function()
			require("mason").setup()
		end
	},

	{ 
		"mason-org/mason-lspconfig.nvim", 
		lazy = false,
		config = function() 
		end 
	}
}
