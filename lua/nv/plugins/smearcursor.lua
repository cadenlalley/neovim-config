return {
	"sphamba/smear-cursor.nvim",
	event = "VeryLazy",
	config = function()
		local color = "#ffffff"
		require("smear_cursor").setup({
			cursor_color = color,
			cursor_color_insert_mode = color,
			stiffness = 0.75,
			trailing_stiffness = 0.55,
			damping = 0.92,
			distance_stop_animating = 0.4,
			time_interval = 12,
			max_length = 40,
			never_draw_over_target = true,
			particles_enabled = true,
			particles_per_second = 280,
			particles_per_length = 6,
			particle_max_lifetime = 550,
			particle_max_initial_velocity = 14,
			particle_velocity_from_cursor = 0.35,
			particle_damping = 0.17,
			particle_gravity = -15,
			min_distance_emit_particles = 0.5,
		})
	end,
}
