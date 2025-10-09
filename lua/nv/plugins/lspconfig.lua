return {
	'neovim/nvim-lspconfig',
	dependencies = { 'saghen/blink.cmp' },

	opts = {
		servers = {
			lua_ls = {
				on_attach = function(client, _)
					client.server_capabilities.documentHighlightProvider = false
				end
			},
			gopls = {
				on_attach = function(client, _)
					client.server_capabilities.documentHighlightProvider = false
				end,
				gofumpt = true,
				codelenses = {
					gc_details = false,
					generate = true,
					regenerate_cgo = true,
					run_govulncheck = true,
					test = true,
					tidy = true,
					upgrade_dependency = true,
					vendor = true,
				},
				hints = {
					assignVariableTypes = true,
					compositeLiteralFields = true,
					compositeLiteralTypes = true,
					constantValues = true,
					functionTypeParameters = true,
					parameterNames = true,
					rangeVariableTypes = true,
				},
				analyses = {
					nilness = true,
					unusedparams = true,
					unusedwrite = true,
					useany = true,
				},
				usePlaceholders = true,
				completeUnimported = true,
				staticcheck = true,
				directoryFilters = { "-.git", "-.vscode", "-.idea", "-.vscode-test", "-node_modules" },
				semanticTokens = true,
				hoverKind = "FullDocumentation",
			},
			terraformls = {
				cmd = { "terraform-ls", "serve" },
				filetypes = { "terraform", "terraform-vars", "tfvars" },
				flags = { debounce_text_changes = 100 },
				on_attach = function(client, _)
					client.server_capabilities.documentHighlightProvider = false
				end
			},
			ols = {

			},
			zls = {

			},
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
