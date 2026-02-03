return {
	{
			"github/copilot.vim",
			event = "InsertEnter",
			cmd = "Copilot",
			init = function()
				if vim.g.copilot_enabled == nil then
					vim.g.copilot_enabled = 0
				end
			end,
			config = function()
				vim.g.copilot_no_tab_map = true
				vim.g.copilot_assume_mapped = true
			end,
		},
	}
