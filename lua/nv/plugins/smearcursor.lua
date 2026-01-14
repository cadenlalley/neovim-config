return {
	"sphamba/smear-cursor.nvim",
	event = "VeryLazy",
	config = function()
		local color = "#ffffff"
		require("smear_cursor").setup({
			cursor_color = color, -- white for snow, use a red/orange hex for fire
			cursor_color_insert_mode = color,
			smear_insert_mode = true,
			stiffness_insert_mode = 0.65,
			max_length_insert_mode = 30,
			distance_stop_animating_vertical_bar = 0.3,
			gradient_exponent = 0,
			particles_enabled = true,
			particle_spread = 1,
			particles_per_second = 100,
			particles_per_length = 50,
			particle_max_lifetime = 1000,
		})
	end,
}
