/* to disable code block CSS from tailwind/typography, we use another code highlighter */
const disabledCss = {
	'code::before': false,
	'code::after': false,
	'blockquote p:first-of-type::before': false,
	'blockquote p:last-of-type::after': false,
	pre: false,
	code: false,
	'pre code': false,
	'code::before': false,
	'code::after': false,
}

module.exports = {
	content: ["internal/views/**/*.templ"],
	theme: {
		extend: {
			/* disable code block CSS */
			typography: {
				DEFAULT: { css: disabledCss },
				sm: { css: disabledCss },
				lg: { css: disabledCss },
				xl: { css: disabledCss },
				'2xl': { css: disabledCss },
			},
		},
	},
	/* @tailwind/typography plugin should be befire daisyui */
	plugins: [require("@tailwindcss/typography"), require("daisyui")],
	daisyui: {
		themes: [
			{
				light: {
					"primary": "#06b6d4",
					"secondary": "#9FB798",
					"accent": "#d946ef",
					"neutral": "#4B4744",
					"base-100": "#f1f5f9",
					"info": "#4ECDC4",
					"warning": "#FFBA08",
					"error": "#E84855",
					"success": "35FF69",
					"--rounded-box": "0.15rem",
					"--rounded-btn": "0.15rem"
				},
				dark: {
					"primary": "#22d3ee",
					"secondary": "#9FB798",
					"accent": "#d946ef",
					"neutral": "#7B9EA8",
					"base-100": "#1B4965",
					"info": "#4ECDC4",
					"warning": "#FFBA08",
					"error": "#E84855",
					"success": "35FF69",
					"--rounded-box": "0.15rem",
					"--rounded-btn": "0.15rem"
				}
			},
		],
	},
};
