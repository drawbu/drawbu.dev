/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./assets/*.css", "./pkg/**/*.templ"],
  theme: {
    extend: {},
    fontFamily: {
      body: ['Iosevka Aile Iaso', 'sans-serif'],
      display: ['Iosevka Aile Iaso', 'sans-serif'],
      mono: ['Iosevka Aile Iaso', 'sans-serif'],
      sans: ['Iosevka Aile Iaso', 'sans-serif'],
      serif: ['Iosevka Aile Iaso', 'sans-serif'],
    },

  },
  plugins: [],
}
