import '@mdi/font/css/materialdesignicons.css'
import "vuetify/styles"
import { createVuetify, ThemeDefinition } from 'vuetify'
import { aliases, mdi } from 'vuetify/iconsets/mdi'

const imeslight: ThemeDefinition = {
  dark: false,
  colors: {
    primary: '#1976D2',
    secondary: '#424242',
    accent: '#82B1FF',
    error: '#FF5252',
    info: '#2196F3',
    success: '#4CAF50',
    warning: '#FB8C00',
  }
}

const imesdark: ThemeDefinition = {
  dark: true,
  colors: {
    primary: '#2196F3',
    secondary: '#424242',
    accent: '#FF4081',
    error: '#FF5252',
    info: '#9769F3',
    success: '#4CAF50',
    warning: '#FB8C00',
  }
}

export default createVuetify({
  defaults: {
    global: {
      ripple: false
    },
    icons: {
      defaultSet: 'mdi',
      aliases,
      sets: {
        mdi,
      }
    },
    theme: {
      defaultTheme: 'dark',
      themes: {
        imeslight,
        imesdark,
      },
    },
    VSheet: {
      elevation: 4,
    },
  }
})