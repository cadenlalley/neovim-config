local M = {}

local COPILOT_PLUGIN = "copilot.vim"

local function notify(state)
	local msg = state and "Copilot enabled" or "Copilot disabled"
	vim.notify(msg, vim.log.levels.INFO, { title = "Copilot" })
end

local function ensure_loaded()
	if vim.g.loaded_copilot == 1 then
		return true
	end

	local ok, lazy = pcall(require, "lazy")
	if ok then
		lazy.load({ plugins = { COPILOT_PLUGIN } })
	end

	return vim.g.loaded_copilot == 1
end

local function current_state()
	if ensure_loaded() and vim.fn.exists("*copilot#Enabled") == 1 then
		local ok, enabled = pcall(function()
			return vim.fn["copilot#Enabled"]()
		end)
		if ok then
			return enabled == 1
		end
	end
	return vim.g.copilot_enabled == 1
end

local function start_client()
	if vim.fn.exists("*copilot#Client") == 1 then
		pcall(function()
			vim.fn["copilot#Client"]()
		end)
	end
end

local function clear_suggestion()
	if vim.fn.exists("*copilot#Dismiss") == 1 then
		pcall(function()
			vim.fn["copilot#Dismiss"]()
		end)
	end
end

local function run_command(cmd, state)
	if not ensure_loaded() then
		vim.notify("Copilot plugin is not loaded", vim.log.levels.ERROR, { title = "Copilot" })
		return false
	end
	local ok, err = pcall(vim.cmd, cmd)
	if not ok then
		vim.notify(("Copilot command failed: %s"):format(err), vim.log.levels.ERROR, { title = "Copilot" })
	end
	if ok then
		vim.g.copilot_enabled = state and 1 or 0
		notify(state)
	end
	return ok
end

function M.enable()
	if not current_state() then
		if run_command("Copilot enable", true) then
			start_client()
		end
	end
end

function M.disable()
	if current_state() then
		if run_command("Copilot disable", false) then
			clear_suggestion()
		end
	end
end

function M.toggle()
	if current_state() then
		M.disable()
	else
		M.enable()
	end
end

return M
