'use strict';

var _slicedToArray = function () { function sliceIterator(arr, i) { var _arr = []; var _n = true; var _d = false; var _e = undefined; try { for (var _i = arr[Symbol.iterator](), _s; !(_n = (_s = _i.next()).done); _n = true) { _arr.push(_s.value); if (i && _arr.length === i) break; } } catch (err) { _d = true; _e = err; } finally { try { if (!_n && _i["return"]) _i["return"](); } finally { if (_d) throw _e; } } return _arr; } return function (arr, i) { if (Array.isArray(arr)) { return arr; } else if (Symbol.iterator in Object(arr)) { return sliceIterator(arr, i); } else { throw new TypeError("Invalid attempt to destructure non-iterable instance"); } }; }();

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function _possibleConstructorReturn(self, call) { if (!self) { throw new ReferenceError("this hasn't been initialised - super() hasn't been called"); } return call && (typeof call === "object" || typeof call === "function") ? call : self; }

function _inherits(subClass, superClass) { if (typeof superClass !== "function" && superClass !== null) { throw new TypeError("Super expression must either be null or a function, not " + typeof superClass); } subClass.prototype = Object.create(superClass && superClass.prototype, { constructor: { value: subClass, enumerable: false, writable: true, configurable: true } }); if (superClass) Object.setPrototypeOf ? Object.setPrototypeOf(subClass, superClass) : subClass.__proto__ = superClass; }

var _require = require('./regexp'),
    referenceRegexp = _require.referenceRegexp,
    anchoredNameRegexp = _require.anchoredNameRegexp,
    anchoredIdentifierRegexp = _require.anchoredIdentifierRegexp;

var _require2 = require('./digest'),
    validateDigest = _require2.validateDigest,
    isDigest = _require2.isDigest;

var _require3 = require('./reference'),
    Reference = _require3.Reference;

var NAME_MAX_LENGTH = 255;

var InvalidReferenceFormatError = function (_Error) {
    _inherits(InvalidReferenceFormatError, _Error);

    function InvalidReferenceFormatError() {
        _classCallCheck(this, InvalidReferenceFormatError);

        var _this = _possibleConstructorReturn(this, (InvalidReferenceFormatError.__proto__ || Object.getPrototypeOf(InvalidReferenceFormatError)).call(this, 'invalid reference format'));

        _this.name = 'InvalidReferenceFormatError';
        return _this;
    }

    return InvalidReferenceFormatError;
}(Error);

var NameContainsUppercaseError = function (_Error2) {
    _inherits(NameContainsUppercaseError, _Error2);

    function NameContainsUppercaseError() {
        _classCallCheck(this, NameContainsUppercaseError);

        var _this2 = _possibleConstructorReturn(this, (NameContainsUppercaseError.__proto__ || Object.getPrototypeOf(NameContainsUppercaseError)).call(this, 'repository name must be lowercase'));

        _this2.name = 'NameContainsUppercaseError';
        return _this2;
    }

    return NameContainsUppercaseError;
}(Error);

var EmptyNameError = function (_Error3) {
    _inherits(EmptyNameError, _Error3);

    function EmptyNameError() {
        _classCallCheck(this, EmptyNameError);

        var _this3 = _possibleConstructorReturn(this, (EmptyNameError.__proto__ || Object.getPrototypeOf(EmptyNameError)).call(this, 'repository name must have at least one component'));

        _this3.name = 'EmptyNameError';
        return _this3;
    }

    return EmptyNameError;
}(Error);

var NameTooLongError = function (_Error4) {
    _inherits(NameTooLongError, _Error4);

    function NameTooLongError() {
        _classCallCheck(this, NameTooLongError);

        var _this4 = _possibleConstructorReturn(this, (NameTooLongError.__proto__ || Object.getPrototypeOf(NameTooLongError)).call(this, 'repository name must not be more than ' + NAME_MAX_LENGTH + ' characters'));

        _this4.name = 'NameTooLongError';
        return _this4;
    }

    return NameTooLongError;
}(Error);

var DEFAULT_DOMAIN = 'docker.io';
var LEGACY_DEFAULT_DOMAIN = 'index.docker.io';
var OFFICIAL_REPOSITORY_NAME = 'library';

exports.parseQualifiedName = function (name) {
    var matches = referenceRegexp.exec(name);

    if (!matches) {
        if (name === '') {
            throw new EmptyNameError();
        }

        if (referenceRegexp.test(name.toLowerCase())) {
            throw new NameContainsUppercaseError();
        }

        throw new InvalidReferenceFormatError();
    }

    if (matches[1].length > NAME_MAX_LENGTH) {
        throw new NameTooLongError();
    }

    var reference = void 0;

    var nameMatch = anchoredNameRegexp.exec(matches[1]);
    if (nameMatch && nameMatch.length === 3) {
        reference = {
            domain: nameMatch[1],
            repository: nameMatch[2]
        };
    } else {
        reference = {
            domain: '',
            repository: matches[1]
        };
    }

    reference.tag = matches[2];

    if (matches[3]) {
        validateDigest(matches[3]);
        reference.digest = matches[3];
    }

    return new Reference(reference);
};

function splitDockerDomain(name) {
    var domain = void 0;
    var reminder = void 0;

    var indexOfSlash = name.indexOf('/');
    if (indexOfSlash === -1 || !(name.lastIndexOf('.', indexOfSlash) !== -1 || name.lastIndexOf(':', indexOfSlash) !== -1 || name.startsWith('localhost/'))) {

        domain = DEFAULT_DOMAIN;
        reminder = name;
    } else {
        domain = name.substring(0, indexOfSlash);
        reminder = name.substring(indexOfSlash + 1);
    }

    if (domain === LEGACY_DEFAULT_DOMAIN) {
        domain = DEFAULT_DOMAIN;
    }

    if (domain === DEFAULT_DOMAIN && !reminder.includes('/')) {
        reminder = OFFICIAL_REPOSITORY_NAME + '/' + reminder;
    }

    return [domain, reminder];
}

exports.parseFamiliarName = function (name) {
    if (anchoredIdentifierRegexp.test(name)) {
        throw new TypeError('invalid repository name (' + name + '),' + 'cannot specify 64-byte hexadecimal strings');
    }

    var _splitDockerDomain = splitDockerDomain(name),
        _splitDockerDomain2 = _slicedToArray(_splitDockerDomain, 2),
        domain = _splitDockerDomain2[0],
        remainder = _splitDockerDomain2[1];

    var remoteName = void 0;
    var tagSeparatorIndex = remainder.indexOf(':');
    if (tagSeparatorIndex > -1) {
        remoteName = remainder.substring(0, tagSeparatorIndex);
    } else {
        remoteName = remainder;
    }

    if (remoteName.toLowerCase() !== remoteName) {
        throw new TypeError('invalid reference format: repository name must be lowercase');
    }

    return exports.parseQualifiedName(domain + '/' + remainder);
};

exports.parseAll = function (name) {
    if (anchoredIdentifierRegexp.test(name)) {
        return new Reference({ digest: 'sha256:' + name });
    }

    if (isDigest(name)) {
        return new Reference({ digest: name });
    }

    return exports.parseFamiliarName(name);
};