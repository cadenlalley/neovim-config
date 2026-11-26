local M = {}
local cache = {}

local function run(cmd)
	local res = vim.system(cmd, { text = true }):wait()
	if res.code ~= 0 then
		return nil
	end
	return vim.trim(res.stdout)
end

local function git(args)
	return run(vim.list_extend({ "git" }, args))
end

function M.base()
	local cwd = vim.fn.getcwd()
	if cache[cwd] then
		return cache[cwd]
	end
	local name = run({ "gh", "pr", "view", "--json", "baseRefName", "-q", ".baseRefName" })
	if not name or name == "" then
		name = "HEAD"
	end
	cache[cwd] = "origin/" .. name
	return cache[cwd]
end

function M.start()
	cache[vim.fn.getcwd()] = nil
	git({ "fetch", "--quiet", "origin" })
	local base = M.base()
	local sha = git({ "merge-base", base, "HEAD" })
	if not sha then
		vim.notify("No merge base with " .. base, vim.log.levels.ERROR)
		return
	end
	require("gitsigns").change_base(sha, true)
	vim.notify("Reviewing against " .. base)
end

function M.stop()
	require("gitsigns").reset_base(true)
end

function M.files()
	Snacks.picker.git_diff({ base = M.base() })
end

function M.diffview()
	vim.cmd("DiffviewOpen " .. M.base() .. "...HEAD --imply-local")
end

return M
