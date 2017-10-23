'use strict';

const typesTemplates = {
    'digest': ref => `${ref.digest}`,
    'canonical': ref => `${ref.repositoryUrl}@${ref.digest}`,
    'repository': ref => `${ref.repositoryUrl}`,
    'tagged': ref => `${ref.repositoryUrl}:${ref.tag}`,
    'dual': ref => `${ref.repositoryUrl}:${ref.tag}@${ref.digest}`
};

exports.Reference = class {
    constructor(options) {
        if (!options.repository && !options.domain) {
            if (options.digest) {
                this.digest = options.digest;
                this.type = 'digest';
            } else {
                throw new TypeError('Empty Reference');
            }
        } else if (!options.tag) {
            this.domain     = options.domain;
            this.repository = options.repository;
            if (options.digest) {
                this.digest = options.digest;
                this.type = 'canonical';
            } else {
                this.type = 'repository';
            }
        } else if (!options.digest) {
            this.domain     = options.domain;
            this.repository = options.repository;
            this.tag        = options.tag;
            this.type       = 'tagged';
        } else {
            this.domain     = options.domain;
            this.repository = options.repository;
            this.tag        = options.tag;
            this.digest     = options.digest;
            this.type       = 'dual';
        }
    }

    toString() {
        return typesTemplates[this.type](this);
    }

    get repositoryUrl() {
        return `${this.domain}/${this.repository}`;
    }
};
