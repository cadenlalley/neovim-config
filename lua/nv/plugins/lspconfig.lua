return {
	'neovim/nvim-lspconfig',
	dependencies = { 'saghen/blink.cmp' },

	opts = {
		servers = {
			lua_ls = {
				on_attach = function(client, bufnr)
					client.server_capabilities.documentHighlightProvider = false
				end
			},
			gopls = {
				on_attach = function(client, bufnr)
					client.server_capabilities.documentHighlightProvider = false
				end,
				usePlaceholders = true,
				completeUnimported = true,
				staticcheck = true,
				semanticTokens = true,
				analyses = {
					unusedparams = true,
					nilness = true,
					shadow = true,
				},
				codelenses = {
					test = true,
					tidy = true,
					regenerate_cgo = true,
				},
				hints = {
					parameterNames = true,
				},
				buildFlags = { "-tags=integration" },
				experimentalPostfixCompletions = true,
				hoverKind = "FullDocumentation",
			}
		}
	},

	config = function(_, opts)
		local lspconfig = require('lspconfig')
		for server, config in pairs(opts.servers) do
			config.capabilities = require('blink.cmp').get_lsp_capabilities(config.capabilities)
			lspconfig[server].setup(config)
		end
	end
}
