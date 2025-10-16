return {
	'L3MON4D3/LuaSnip',
	version = "v2.*",
	build = "make install_jsregexp",
	dependencies = {
		'rafamadriz/friendly-snippets',
	},
	config = function()
		local ls = require('luasnip')
		local s = ls.snippet
		local t = ls.text_node
		local i = ls.insert_node

		-- Load snippets from friendly-snippets
		require('luasnip.loaders.from_vscode').lazy_load()
		-- Custom Go error handling snippet

		ls.add_snippets("go", {
			s("iferr", {
				t("if err != nil {"),
				t({ "", "\treturn " }),
				i(1),
				t({ "", "}" }),
				i(0)
			})
		})

		ls.add_snippets("go", {
			s("fmterr", {
				t("if err != nil {"),
				t({ "", "\treturn fmt.Errorf(\"" }),
				i(1, ""), -- cursor will jump here
				t({ " :%w\", err)", "", "}" }),
			}),
		})
	end,
}

