return {
	"nvim-neotest/neotest",
	dependencies = {
		"nvim-lua/plenary.nvim",
		"antoinemadec/FixCursorHold.nvim",
		"nvim-neotest/neotest-go",
	},
	config = function()
		local neotest = require("neotest")
		neotest.setup({
			adapters = {
				require("neotest-go")({
					experimental = { test_table = true },
				}),
			},
		})
	end,
	keys = {
		{ "<leader>tn", function() require("neotest").run.run() end,                        desc = "run nearest test" },
		{ "<leader>tf", function() require("neotest").run.run(vim.fn.expand("%")) end,      desc = "run file tests" },
		{ "<leader>tw", function() require("neotest").watch.toggle(vim.fn.expand("%")) end, desc = "test watch" },
		{ "<leader>to", function() require("neotest").output.open() end,                    desc = "test output" },
		{ "<leader>tp", function() require("neotest").output_panel.toggle() end,            desc = "test output panel" },
		{ "<leader>ts", function() require("neotest").summary.toggle() end,                 desc = "toggle summary" },
	},
}
