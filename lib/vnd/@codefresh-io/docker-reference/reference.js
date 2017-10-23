'use strict';

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

var typesTemplates = {
    'digest': function digest(ref) {
        return '' + ref.digest;
    },
    'canonical': function canonical(ref) {
        return ref.repositoryUrl + '@' + ref.digest;
    },
    'repository': function repository(ref) {
        return '' + ref.repositoryUrl;
    },
    'tagged': function tagged(ref) {
        return ref.repositoryUrl + ':' + ref.tag;
    },
    'dual': function dual(ref) {
        return ref.repositoryUrl + ':' + ref.tag + '@' + ref.digest;
    }
};

exports.Reference = function () {
    function _class(options) {
        _classCallCheck(this, _class);

        if (!options.repository && !options.domain) {
            if (options.digest) {
                this.digest = options.digest;
                this.type = 'digest';
            } else {
                throw new TypeError('Empty Reference');
            }
        } else if (!options.tag) {
            this.domain = options.domain;
            this.repository = options.repository;
            if (options.digest) {
                this.digest = options.digest;
                this.type = 'canonical';
            } else {
                this.type = 'repository';
            }
        } else if (!options.digest) {
            this.domain = options.domain;
            this.repository = options.repository;
            this.tag = options.tag;
            this.type = 'tagged';
        } else {
            this.domain = options.domain;
            this.repository = options.repository;
            this.tag = options.tag;
            this.digest = options.digest;
            this.type = 'dual';
        }
    }

    _createClass(_class, [{
        key: 'toString',
        value: function toString() {
            return typesTemplates[this.type](this);
        }
    }, {
        key: 'repositoryUrl',
        get: function get() {
            return this.domain + '/' + this.repository;
        }
    }]);

    return _class;
}();