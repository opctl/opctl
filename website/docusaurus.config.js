module.exports = {
  title: 'Opctl',
  tagline: 'Automate operating your project; use containers as building blocks.',
  url: 'https://opctl.io',
  baseUrl: '/',
  favicon: 'img/favicon.ico',
  organizationName: 'opctl', // Usually your GitHub org/user name.
  projectName: 'opctl', // Usually your repo name.
  themeConfig: {
    disableDarkMode: true,    
    googleAnalytics: {
      trackingID: 'UA-94109316-1',
    },
    navbar: {
      title: 'Opctl',
      logo: {
        alt: 'opctl Logo',
        src: 'img/logo.svg',
      },
      links: [
        { to: 'docs/zero-to-hero/hello-world', label: 'Docs', position: 'left' },
        {
          href: 'https://opctl-slackin.herokuapp.com/',
          label: 'Slack',
          position: 'right',
        },
        {
          href: 'https://github.com/opctl/opctl',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {
              label: 'Docs',
              to: 'docs/zero-to-hero/hello-world',
            },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'Slack',
              href: 'https://opctl-slackin.herokuapp.com/',
            },
          ],
        },
        {
          title: 'Social',
          items: [
            {
              label: 'Github',
              href: 'https://github.com/opctl/opctl'
            }
          ]
        }
      ],
      copyright: `Copyright © ${new Date().getFullYear()} opctl.io`,
    },
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          editUrl: "https://github.com/opctl/opctl/edit/master/docs/",
          sidebarPath: require.resolve('./sidebars.js'),
          // Equivalent to `enableUpdateBy`.
          showLastUpdateAuthor: true,
          // Equivalent to `enableUpdateTime`.
          showLastUpdateTime: true,
        }
      },
    ],
  ],
};
