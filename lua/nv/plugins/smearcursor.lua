return {
	"sphamba/smear-cursor.nvim",
	event = "VeryLazy",
	config = function()
		local color = "#ffffff"
		require("smear_cursor").setup({
			cursor_color = "#ffffff", -- white for snow, use a red/orange hex for fire
			gradient_exponent = 0,
			particles_enabled = true,
			particle_spread = 1,
			particles_per_second = 100,
			particles_per_length = 50,
			particle_max_lifetime = 1500,
		})
	end,
}
