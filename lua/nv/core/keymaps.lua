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
keymap.set("n", "<leader>fw", picker.grep, { desc = "find word" })
keymap.set("n", "<leader>fb", picker.buffers, { desc = "find buffer" })

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
keymap.set('n', '<leader>ca', vim.lsp.buf.code_action, { desc = "lsp code actions" })
keymap.set('n', '<leader>rn', vim.lsp.buf.rename, { desc = "lsp rename" })

----------------------
-- flash
----------------------

keymap.set({ "n", "x", "o" }, "<leader>b", flash.jump, { desc = "flash Jump" })
keymap.set({ "n", "x", "o" }, "<leader>S", flash.treesitter, { desc = "flash treesitter" })
keymap.set({ "n", "x", "o" }, "<leader>r", flash.remote, { desc = "flash remote" })
keymap.set({ "n", "x", "o" }, "<leader>R", flash.treesitter_search, { desc = "flash treesitter search" })
