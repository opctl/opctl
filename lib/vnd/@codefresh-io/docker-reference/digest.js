'use strict';

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function _possibleConstructorReturn(self, call) { if (!self) { throw new ReferenceError("this hasn't been initialised - super() hasn't been called"); } return call && (typeof call === "object" || typeof call === "function") ? call : self; }

function _inherits(subClass, superClass) { if (typeof superClass !== "function" && superClass !== null) { throw new TypeError("Super expression must either be null or a function, not " + typeof superClass); } subClass.prototype = Object.create(superClass && superClass.prototype, { constructor: { value: subClass, enumerable: false, writable: true, configurable: true } }); if (superClass) Object.setPrototypeOf ? Object.setPrototypeOf(subClass, superClass) : subClass.__proto__ = superClass; }

var _require = require('./regexp'),
    anchoredDigestRegexp = _require.anchoredDigestRegexp;

var InvalidDigestFormatError = function (_Error) {
    _inherits(InvalidDigestFormatError, _Error);

    function InvalidDigestFormatError() {
        _classCallCheck(this, InvalidDigestFormatError);

        var _this = _possibleConstructorReturn(this, (InvalidDigestFormatError.__proto__ || Object.getPrototypeOf(InvalidDigestFormatError)).call(this, 'invalid digest format'));

        _this.name = 'InvalidDigestFormatError';
        return _this;
    }

    return InvalidDigestFormatError;
}(Error);

var UnsupportedAlgorithmError = function (_Error2) {
    _inherits(UnsupportedAlgorithmError, _Error2);

    function UnsupportedAlgorithmError() {
        _classCallCheck(this, UnsupportedAlgorithmError);

        var _this2 = _possibleConstructorReturn(this, (UnsupportedAlgorithmError.__proto__ || Object.getPrototypeOf(UnsupportedAlgorithmError)).call(this, 'unsupported digest algorithm'));

        _this2.name = 'UnsupportedAlgorithmError';
        return _this2;
    }

    return UnsupportedAlgorithmError;
}(Error);

var InvalidDigestLengthError = function (_Error3) {
    _inherits(InvalidDigestLengthError, _Error3);

    function InvalidDigestLengthError() {
        _classCallCheck(this, InvalidDigestLengthError);

        var _this3 = _possibleConstructorReturn(this, (InvalidDigestLengthError.__proto__ || Object.getPrototypeOf(InvalidDigestLengthError)).call(this, 'invalid checksum digest length'));

        _this3.name = 'InvalidDigestLengthError';
        return _this3;
    }

    return InvalidDigestLengthError;
}(Error);

var algorithmsSizes = {
    sha256: 32,
    sha384: 48,
    sha512: 64
};

function checkDigest(digest, handleError) {
    var indexOfColon = digest.indexOf(':');
    if (indexOfColon < 0 || indexOfColon + 1 === digest.length || !anchoredDigestRegexp.test(digest)) {
        return handleError(InvalidDigestFormatError);
    }

    var algorithm = digest.substring(0, indexOfColon);
    if (!algorithmsSizes[algorithm]) {
        return handleError(UnsupportedAlgorithmError);
    }

    if (algorithmsSizes[algorithm] * 2 !== digest.length - indexOfColon - 1) {
        return handleError(InvalidDigestLengthError);
    }

    return true;
}

exports.validateDigest = function (digest) {
    checkDigest(digest, function (ErrorType) {
        throw new ErrorType();
    });
};

exports.isDigest = function (digest) {
    return checkDigest(digest, function () {
        return false;
    });
};