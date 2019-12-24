/**
 * Copyright (c) 2017-present, Facebook, Inc.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

import React from 'react';
import classnames from 'classnames';
import Layout from '@theme/Layout';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import useBaseUrl from '@docusaurus/useBaseUrl';
import styles from './styles.module.css';

const features = [
  {
    title: <>Composable</>,
    imageUrl: 'img/fa_puzzle-piece.svg',
    description: (
      <>
        Build ops out of other ops.
      </>
    ),
  },
  {
    title: <>Containerized</>,
    imageUrl: 'img/fa_square.svg',
    description: (
      <>
        Say goodbye to operations-induced <a href="https://whatis.techtarget.com/definition/yak-shaving">yak shaving</a>.
      </>
    ),
  },
  {
    title: <>Distributable</>,
    imageUrl: 'img/fa_exchange.svg',
    description: (
      <>
        Define once, use everywhere.
      </>
    ),
  },
  {
    title: <>Versionable</>,
    imageUrl: 'img/fa_history.svg',
    description: (
      <>
        Version in standard source control.
      </>
    ),
  },
];

function Feature({ imageUrl, title, description }) {
  const imgUrl = useBaseUrl(imageUrl);
  return (
    <div className={classnames('col col--3', styles.feature)}>
      {imgUrl && (
        <div className="text--center">
          <img
            style={{
              maxWidth: '80px'
            }}
            src={imgUrl}
            alt={title}
          />
        </div>
      )}
      <h3>{title}</h3>
      <p>{description}</p>
    </div>
  );
}

function Home() {
  const context = useDocusaurusContext();
  const { siteConfig = {} } = context;
  return ( 
    <Layout
      title={siteConfig.title}
      description={siteConfig.tagline}>
      <header className={classnames('hero hero--primary', styles.heroBanner)}>
        <div className="container">
          <div className="project-logo">
            <img style={{ height: '20rem' }} src={'img/op-backdrop.svg'} alt="opctl" />
          </div>
          <h1 className="hero__title">{siteConfig.title}</h1>
          <p className="hero__subtitle">{siteConfig.tagline}</p>
          <div className={styles.buttons}>
            <Link
              className={classnames(
                'button button--outline button--secondary button--lg',
                styles.getStarted,
              )}
              to={useBaseUrl('docs/zero-to-hero/hello-world')}>
              Get Started
            </Link>
          </div>
        </div>
      </header>
      <main>
        {features && features.length && (
          <section className={styles.features}>
            <div className="container">
              <div className="row">
                {features.map((props, idx) => (
                  <Feature key={idx} {...props} />
                ))}
              </div>
            </div>
          </section>
        )}
      </main>
    </Layout>
  );
}

export default Home;
