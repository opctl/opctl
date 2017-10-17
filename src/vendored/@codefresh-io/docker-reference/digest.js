'use strict';

const { anchoredDigestRegexp } = require('./regexp');

class InvalidDigestFormatError extends Error {
    constructor() {
        super('invalid digest format');
        this.name = 'InvalidDigestFormatError';
    }
}

class UnsupportedAlgorithmError extends Error {
    constructor() {
        super('unsupported digest algorithm');
        this.name = 'UnsupportedAlgorithmError';
    }
}

class InvalidDigestLengthError extends Error {
    constructor() {
        super('invalid checksum digest length');
        this.name = 'InvalidDigestLengthError';
    }
}

const algorithmsSizes = {
    sha256: 32,
    sha384: 48,
    sha512: 64,
};

function checkDigest(digest, handleError) {
    const indexOfColon = digest.indexOf(':');
    if (indexOfColon < 0 ||
        indexOfColon + 1 === digest.length ||
        !anchoredDigestRegexp.test(digest)) {
        return handleError(InvalidDigestFormatError);
    }

    const algorithm = digest.substring(0, indexOfColon);
    if (!algorithmsSizes[algorithm]) {
        return handleError(UnsupportedAlgorithmError);
    }

    if (algorithmsSizes[algorithm] * 2 !== (digest.length - indexOfColon - 1)) {
        return handleError(InvalidDigestLengthError);
    }

    return true;
}

exports.validateDigest = (digest) => {
    checkDigest(digest, (ErrorType) => {
        throw new ErrorType();
    });
};

exports.isDigest = (digest) => {
    return checkDigest(digest, () => false);
};
