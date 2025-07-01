local opt = vim.opt

-- tabs & indentation
opt.tabstop = 4
opt.expandtab = false
opt.softtabstop = 4
opt.shiftwidth = 4
opt.autoindent = true

-- line numbers
opt.relativenumber = true
opt.number = true

-- line wrapping
opt.wrap = false

-- search casing
opt.ignorecase = true
opt.smartcase = true

-- colors
opt.termguicolors = true
opt.background = "dark"
opt.signcolumn = "yes"

-- backspace
opt.backspace = "indent,eol,start"

-- clipboard
opt.clipboard:append("unnamedplus")

-- window split
opt.splitright = true
opt.splitbelow = true


