/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./static/*.css", "./pkg/**/*.templ"],
  theme: {
    extend: {},
    fontFamily: {
      body: ['Iosevka Comfy', 'sans-serif'],
      display: ['Iosevka Comfy', 'sans-serif'],
      mono: ['Iosevka Comfy', 'sans-serif'],
      sans: ['Iosevka Comfy', 'sans-serif'],
      serif: ['Iosevka Comfy', 'sans-serif'],
    },

  },
  plugins: [],
}
