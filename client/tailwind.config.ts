import type { Config } from 'tailwindcss'

const config: Config = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          primary: '#3B82F6',  // blue-500
          secondary: '#1D4ED8', // blue-700
          accent: '#EF4444',   // red-500
        },
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}

export default config