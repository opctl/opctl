/**
 * Copyright (c) 2017-present, Facebook, Inc.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

// See https://docusaurus.io/docs/site-config for all the possible
// site configuration options.

// List of projects/orgs using your project for the users page.
const users = [
  {
    caption: 'Expedia',
    // You will need to prepend the image path with your baseUrl
    // if it is not '/', like: '/test-site/img/image.jpg'.
    image: '/img/expedia-logo.svg',
    infoLink: 'https://expedia.com',
    pinned: true,
  },
  {
    caption: 'helloera',
    // You will need to prepend the image path with your baseUrl
    // if it is not '/', like: '/test-site/img/image.jpg'.
    image: '/img/helloera-logo.svg',
    infoLink: 'https://helloera.co',
    pinned: true,
  },
  {
    caption: 'Nintex',
    // You will need to prepend the image path with your baseUrl
    // if it is not '/', like: '/test-site/img/image.jpg'.
    image: '/img/nintex-logo.svg',
    infoLink: 'https://nintex.com',
    pinned: true,
  },
  {
    caption: 'ProKarma',
    // You will need to prepend the image path with your baseUrl
    // if it is not '/', like: '/test-site/img/image.jpg'.
    image: '/img/prokarma-logo.svg',
    infoLink: 'https://prokarma.com',
    pinned: true,
  },
  {
    caption: 'Remitly',
    // You will need to prepend the image path with your baseUrl
    // if it is not '/', like: '/test-site/img/image.jpg'.
    image: '/img/remitly-logo.png',
    infoLink: 'https://remitly.com',
    pinned: true,
  },
  {
    caption: 'Samsung SDS',
    // You will need to prepend the image path with your baseUrl
    // if it is not '/', like: '/test-site/img/image.jpg'.
    image: '/img/samsung-logo.png',
    infoLink: 'https://www.samsungsds.com',
    pinned: true,
  },
];

const slackUrl = 'https://opctl-slackin.herokuapp.com'
const githubUrl = 'https://github.com/opctl/opctl'

const siteConfig = {
  title: 'opctl', // Title for your website.
  tagline: 'Distributed operation control system',
  url: 'https://opctl.io', // Your website URL
  baseUrl: '/', // Base URL for your project */
  // For github.io type URLs, you would set the url and baseUrl like:
  url: 'https://opctl.io',
  //   baseUrl: '/test-site/',

  // Used for publishing and more
  projectName: 'opctl',
  organizationName: 'opctl',
  cname: "opctl.io" /* the CNAME for your website */,
  // For top-level user or org sites, the organization is still the same.
  // e.g., for the https://JoelMarcey.github.io site, it would be set like...
  //   organizationName: 'JoelMarcey'

  // For no header links in the top nav bar -> headerLinks: [],
  headerLinks: [
    { doc: 'setup', label: 'Docs' },
    { doc: "run-a-react-app", label: "Examples" },
    { href: githubUrl, label: "GitHub" },
    { href: slackUrl, label: 'Slack' }
  ],

  // If you have users set above, you add it here:
  users,

  /* path to images for header/footer */
  headerIcon: 'img/logo.svg',
  footerIcon: 'img/logo.svg',
  favicon: 'img/logo.svg',

  editUrl: "https://github.com/opctl/opctl/edit/master/docs/",
  /* Colors for website */
  colors: {
    primaryColor: "#222",
    secondaryColor: "#333",
  },

  /* Custom fonts for website */
  /*
  fonts: {
    myFont: [
      "Times New Roman",
      "Serif"
    ],
    myOtherFont: [
      "-apple-system",
      "system-ui"
    ]
  },
  */

  // This copyright info is used in /core/Footer.js and blog RSS/Atom feeds.
  copyright: `Copyright © ${new Date().getFullYear()} opctl`,

  highlight: {
    // Highlight.js theme to use for syntax highlighting in code blocks.
    theme: 'default',
  },

  // Add custom scripts here that would be placed in <script> tags.
  scripts: ['https://buttons.github.io/buttons.js'],

  scrollToTop: true,

  // On page navigation for the current documentation page.
  onPageNav: 'separate',
  // No .html extensions for paths.
  cleanUrl: true,

  // Open Graph and Twitter card images.
  ogImage: 'img/undraw_online.svg',
  twitterImage: 'img/undraw_tweetstorm.svg',

  // Show documentation's last contributor's name.
  enableUpdateBy: true,

  // Show documentation's last update time.
  enableUpdateTime: true,

  // You may provide arbitrary config keys to be used as needed by your
  // template. For example, if you need your repo's URL...
  githubUrl,
  slackUrl
};

module.exports = siteConfig;
