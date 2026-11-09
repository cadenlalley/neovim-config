return {
	"Cannon07/code-preview.nvim",
	lazy = false,
	opts = {
		diff = {
			layout = "tab",
			equalize = true,
			full_file = true,
			visible_only = false,
		},
		neo_tree = { enabled = false },
		keys = { close_all = "<leader>dq" },
	},
	config = function(_, opts)
		require("code-preview").setup(opts)

		local group = vim.api.nvim_create_augroup("CodePreviewNavigation", { clear = true })
		vim.api.nvim_create_autocmd({ "BufWinEnter", "WinEnter" }, {
			group = group,
			callback = function()
				-- The plugin enables diff mode after creating its scratch windows.
				vim.schedule(function()
					for _, win in ipairs(vim.api.nvim_list_wins()) do
						local buf = vim.api.nvim_win_get_buf(win)
						if vim.wo[win].diff and vim.bo[buf].buftype == "nofile" then
							vim.wo[win].wrap = false
							vim.wo[win].number = true
							vim.wo[win].relativenumber = false
							local function map(key, action, desc)
								vim.keymap.set("n", key, action, { buffer = buf, silent = true, desc = desc })
							end
							map("]c", function()
								vim.cmd("normal! " .. vim.v.count1 .. "]c")
								vim.cmd("normal! zz")
							end, "Next preview hunk")
							map("[c", function()
								vim.cmd("normal! " .. vim.v.count1 .. "[c")
								vim.cmd("normal! zz")
							end, "Previous preview hunk")
							map("<Tab>", "<C-w>w", "Switch preview side")
							map("<S-Tab>", "<C-w>W", "Switch preview side backward")
						end
					end
				end)
			end,
		})
	end,
}
