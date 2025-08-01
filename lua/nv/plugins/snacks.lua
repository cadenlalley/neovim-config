return {
	"folke/snacks.nvim",
	priority = 1000,
	lazy = false,

	opts = {
		bigfile = { enabled = true },
		dashboard = {
			enabled = true,
			preset = {
				header = [[
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣈⣻⣿⣿⣿⣶⣶⣦⣄⣀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⢀⣤⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣦⡀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠛⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠿⠿⣿⣿⣿⣦⡀⠀⠀⠀
⠀⠀⠀⠀⠀⣀⣀⣰⣿⣿⣿⣿⣿⣿⣿⣿⣿⡟⣡⣶⣿⣷⣶⣍⠻⣿⣿⣄⠀⠀
⠀⠀⠀⠀⠀⠀⠤⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⢰⣿⣿⣿⣿⣿⣿⣷⡙⣿⣿⡆⠀
⠀⠀⠀⠤⠤⠤⣤⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⢸⣿⣿⡟⢻⣿⣿⣿⣇⢹⣿⣿⡄
⠀⢰⡄⠀⢀⣤⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣆⢻⣿⣿⣿⣿⣿⣿⡟⣸⣿⣿⣇
⣷⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣦⣙⠿⣿⣿⣿⠟⣱⣿⣿⣿⣿
⣭⣻⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣶⣶⣶⣿⣿⡿⢟⣫⣷
⠈⢻⡿⣶⣾⣭⣽⣛⣛⣿⣿⣿⠿⠿⠿⠿⣿⣿⣿⣛⣛⣻⣭⣽⣷⣶⣿⣿⣿⡏
⠀⠈⠁⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠁
⠀⠀⠈⠉⠉⠉⠉⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠏⠀
⠀⠀⠀⠀⠈⢉⣉⣩⣭⣭⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠋⠀⠀
⠀⠀⠀⠀⠀⠈⠉⠛⠛⣛⣛⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠟⠁⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⣭⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠟⠁⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⣿⣿⣿⠿⠿⠛⠋⠉⠀⠀⠀⠀⠀⠀⠀⠀]],
			},
			sections = {
				{ section = "header", header = [[Test]] },
				{ section = "keys", gap = 1, padding = 1 },
				{ pane = 2, icon = " ", title = "Recent Files", section = "recent_files", indent = 2, padding = 1 },
				{ pane = 2, icon = " ", title = "Projects", section = "projects", indent = 2, padding = 1 },
				{
					pane = 2,
					icon = " ",
					title = "Git Status",
					section = "terminal",
					enabled = function()
						return Snacks.git.get_root() ~= nil
					end,
					cmd = "git status --short --branch --renames",
					height = 5,
					padding = 1,
					ttl = 5 * 60,
					indent = 3,
				},
				{ section = "startup" },
			},
		},
		explorer = { enabled = true, sources = { files = { hidden = true, } } },
		indent = { enabled = true },
		input = { enabled = true },
		picker = { enabled = true, sources = { files = { hidden = true, ignored = true } }, hidden = true, ignored = true },
		notifier = { enabled = true },
		quickfile = { enabled = true },
		scope = { enabled = true },
		scroll = { enabled = true },
		statuscolumn = { enabled = true },
		words = { enabled = true },
		terminal = { enabled = true },
		mage = { enabled = true },
		gitbrowse = { enabled = true },
		git = { enabled = true },
		lazygit = { enabled = true },
	},
}
