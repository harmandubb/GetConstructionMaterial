/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      backgroundColor: {
        'brand': '#FFEA9C',
        'logo-blue': '#069AF0',
        'logo-red-dark': '#E5191C',
        'logo-red-lights': '#FF3129',
      },
      spacing: {
        '120': '700px',
      }, 
    },
  },
  plugins: [],
}
