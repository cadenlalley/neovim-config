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
			ts_ls = {
				on_attach = function(client, _)
					client.server_capabilities.documentHighlightProvider = false
				end,
				filetypes = {
					"javascript",
					"javascriptreact",
					"javascript.jsx",
					"typescript",
					"typescriptreact",
					"typescript.tsx",
				},
				root_markers = { "package.json", "tsconfig.json", "jsconfig.json", ".git" },
			},
			eslint = {
				on_attach = function(client, _)
					client.server_capabilities.documentHighlightProvider = false
				end
			},
			rust_analyzer = {
				on_attach = function(client, _)
					client.server_capabilities.documentHighlightProvider = false
				end
			},
			dockerls = {
				on_attach = function(client, _)
					client.server_capabilities.documentHighlightProvider = false
				end,
				filetypes = { 'dockerfile' },
				root_markers = { 'Dockerfile', '.git', 'docker-compose.yml' },
			}
		}
	},

	config = function(_, opts)
		local cmp = require('blink.cmp')
		for server, server_opts in pairs(opts.servers) do
			local config = vim.deepcopy(server_opts)
			config.capabilities = cmp.get_lsp_capabilities(config.capabilities)
			vim.lsp.config(server, config)
			vim.lsp.enable(server)
		end
	end
}
