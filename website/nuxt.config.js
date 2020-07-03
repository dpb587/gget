export default {
  mode: 'universal',
  target: 'server',
  head: {
    title: 'gget',
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
    ]
  },
  loading: {
    color: '#48BB78',
    height: '3px'
  },
  buildModules: [
    '@nuxtjs/tailwindcss',
    '@nuxtjs/fontawesome',
    [
      '@nuxtjs/google-analytics',
      {
        id: process.env.GA_TRACKING_ID
      }
    ]
  ],
  env: {
    baseUrl: process.env.BASE_URL || 'http://127.0.0.1:3000/'
  },
  fontawesome: {
    component: "fa",
    icons: {
      solid: ['faArrowAltCircleDown', 'faBeer', 'faCheckDouble', 'faChevronDown', 'faCode', 'faCodeBranch', 'faCopy', 'faExchangeAlt', 'faFileDownload', 'faGraduationCap', 'faHeart', 'faHome', 'faLifeRing', 'faLock', 'faPlay', 'faPlayCircle', 'faRandom', 'faServer', 'faSignature', 'faTerminal', 'faTools', 'faUserAlt'],
      brands: ['faApple', 'faDocker', 'faGitAlt', 'faGithub', 'faLinux', 'faTwitter', 'faWindows'],
    }
  }
}
