module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'twitter': {
          '50': '#f4fafe', 
          '100': '#e8f6fe', 
          '200': '#c7e8fc', 
          '300': '#a5d9fa', 
          '400': '#61bdf6', 
          '500': '#1da1f2', 
          '600': '#1a91da', 
          '700': '#1679b6', 
          '800': '#116191', 
          '900': '#0e4f77'
        },
        'facebook': {
          '50': '#f3f8fe', 
          '100': '#e8f1fe', 
          '200': '#c5ddfc', 
          '300': '#a3c9fa', 
          '400': '#5da0f6', 
          '500': '#1877f2', 
          '600': '#166bda', 
          '700': '#1259b6', 
          '800': '#0e4791', 
          '900': '#0c3a77'
        }
      }
    },
    fontFamily: {
      'sans': ['"Outfit"', 'system-ui', 'sans-serif'],
    }
  },
  daisyui: {
    themes: false,
    logs: false
  },
  plugins: [
    require('daisyui')
  ],
}