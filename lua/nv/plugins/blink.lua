return {
  'saghen/blink.cmp',
  dependencies = { 
    'rafamadriz/friendly-snippets'
  },
  version = '1.*',
  opts = {
	keymap = {
		preset = "enter",
		["<Enter>"] = { "select_and_accept", "fallback" },
		["<Tab>"] = false,
	},

    appearance = {
      nerd_font_variant = 'mono'
    },

    completion = { documentation = { auto_show = false } },

    snippets = { preset = 'luasnip' },

    sources = {
      default = { 'lsp', 'path', 'snippets', 'buffer' },
    },

    fuzzy = { implementation = "prefer_rust_with_warning" }
  },
  opts_extend = { "sources.default" }
}
