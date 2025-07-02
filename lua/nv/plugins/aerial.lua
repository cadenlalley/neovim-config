return {
	"stevearc/aerial.nvim",
	event = "BufEnter",
	config = function()
		require("aerial").setup()
	end
}
