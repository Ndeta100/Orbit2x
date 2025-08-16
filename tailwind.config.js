/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./views/**/*.html","./**/*.templ", "./**/go", "./templates/**/*.html",
        "./static/**/*.js",],
    theme: {
        extend: {},
    },
    plugins: [ require('@tailwindcss/forms'),
        require('@tailwindcss/typography')],
}
