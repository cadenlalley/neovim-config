local keymap = vim.keymap

----------------------
-- plugin imports
----------------------

local explorer = require("snacks.explorer")
local picker = require("snacks.picker")
local git = require("snacks.git")
local gitbrowse = require("snacks.gitbrowse")
local lazygit = require("snacks.lazygit")
local flash = require("flash")
local harpoon = require("harpoon")
local treesittercontext = require("treesitter-context")

----------------------
-- misc
----------------------

keymap.set("n", "<Esc>", "<cmd>noh<cr>", { desc = "clear search" })

----------------------
-- buffer management
----------------------

keymap.set("n", "<Tab>", "<cmd>bnext<cr>", { desc = "buffer next" })
keymap.set("n", "<S-Tab>", "<cmd>bprev<cr>", { desc = "buffer prev" })
keymap.set("n", "<leader>x", "<cmd>bdelete<cr>", { desc = "buffer delete" })
keymap.set("n", "<leader>n", "<cmd>enew<cr>", { desc = "new buffer" })

----------------------
-- window management
----------------------

keymap.set("n", "<leader>wsv", "<C-w>v", { desc = "window split vertically" })
keymap.set("n", "<leader>wsh", "<C-w>h", { desc = "window split horizontally" })
keymap.set("n", "<leader>wse", "<C-w>=", { desc = "window size equal" })
keymap.set("n", "<leader>wsm", "<cmd>MaximizerToggle<cr>", { desc = "window toggle maximized" })

----------------------
-- explorer
----------------------

keymap.set("n", "<leader>ee", explorer.open, { desc = "explorer toggle" })
keymap.set("n", "<leader>er", explorer.reveal, { desc = "explorer reveal" })

----------------------
-- picker
----------------------

keymap.set("n", "<leader>ff", picker.files, { desc = "find files" })
keymap.set("n", "<leader>fW", picker.grep, { desc = "find word globally" })
keymap.set("n", "<leader>fw", picker.lines, { desc = "find word" })
keymap.set("n", "<leader>fb", picker.buffers, { desc = "find buffer" })
keymap.set("n", "<leader>fs", picker.lsp_symbols, { desc = "find symbols" })
keymap.set("n", "<leader>fS", function() picker.lsp_workspace_symbols({ live = true }) end, { desc = "find symbols" })
keymap.set("n", "<leader>fgd", picker.git_diff, { desc = "find git diff" })
keymap.set("n", "<leader>fgs", picker.git_status, { desc = "find git status" })
keymap.set("n", "<leader>fgt", picker.git_stash, { desc = "find git stash" })
keymap.set("n", "<leader>f'", picker.registers, { desc = "find registers" })
keymap.set("n", "<leader>f/", picker.search_history, { desc = "find search history" })
keymap.set("n", "<leader>fa", picker.autocmds, { desc = "find autocmds" })
keymap.set("n", "<leader>fc", picker.command_history, { desc = "find command history" })
keymap.set("n", "<leader>fC", picker.commands, { desc = "find commands" })
keymap.set("n", "<leader>fd", picker.diagnostics, { desc = "find diagnostics" })
keymap.set("n", "<leader>fD", picker.diagnostics_buffer, { desc = "find buffer diagnostics" })
keymap.set("n", "<leader>fh", picker.help, { desc = "find help pages" })
keymap.set("n", "<leader>fj", picker.jumps, { desc = "find jumps" })
keymap.set("n", "<leader>fk", picker.keymaps, { desc = "find keymaps" })
keymap.set("n", "<leader>fl", picker.loclist, { desc = "find location list" })
keymap.set("n", "<leader>fM", picker.man, { desc = "find man pages" })
keymap.set("n", "<leader>fm", picker.marks, { desc = "find marks" })
keymap.set("n", "<leader>fR", picker.resume, { desc = "find resume" })
keymap.set("n", "<leader>fq", picker.qflist, { desc = "find quickfix list" })
keymap.set("n", "<leader>fu", picker.undo, { desc = "find undotree" })

----------------------
-- source control
----------------------

keymap.set("n", "<leader>sb", git.blame_line, { desc = "git toggle blame" })
keymap.set("n", "<leader>sob", gitbrowse.open, { desc = "git open in browser" })
keymap.set("n", "<leader>soc", gitbrowse.get_url, { desc = "git get url" })
keymap.set("n", "<leader>slo", lazygit.open, { desc = "git open lazygit" })
keymap.set("n", "<leader>sll", lazygit.log, { desc = "git open lazygit log" })

----------------------
-- lsp
----------------------

keymap.set("n", "<leader>gd", picker.lsp_definitions, { desc = "lsp go to definitions" })
keymap.set("n", "<leader>gD", picker.lsp_declarations, { desc = "lsp to go declarations" })
keymap.set("n", "<leader>gI", picker.lsp_implementations, { desc = "lsp go to implementations" })
keymap.set("n", "<leader>gr", picker.lsp_references, { desc = "lsp go to references" })
keymap.set("n", "<leader>lf", vim.diagnostic.open_float, { desc = "lsp open diagnostic" })
keymap.set("n", "<leader>ca", vim.lsp.buf.code_action, { desc = "lsp code actions" })
keymap.set("n", "<leader>rn", vim.lsp.buf.rename, { desc = "lsp rename" })

----------------------
-- flash
----------------------

keymap.set({ "n", "x", "o" }, "s", flash.jump, { desc = "flash jump" })

----------------------
-- harpoon
----------------------

keymap.set("n", "<leader>ha", function() harpoon:list():add() end, { desc = "harpoon add mark" })
keymap.set("n", "<leader>hl", function() harpoon.ui:toggle_quick_menu(harpoon:list()) end, { desc = "harpoon show list" })
keymap.set("n", "<leader>h[", function() harpoon:list():prev() end, { desc = "harpoon previous mark" })
keymap.set("n", "<leader>h]", function() harpoon:list():next() end, { desc = "harpoon next mark" })

----------------------
-- treesitter context
----------------------

keymap.set("n", "[c", treesittercontext.go_to_context, { desc = "go to context" })
