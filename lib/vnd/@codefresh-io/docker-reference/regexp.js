'use strict';

/**
 * expression defines a full expression, where each regular expression must follow the previous.
 */

function expression() {
    for (var _len = arguments.length, regexps = Array(_len), _key = 0; _key < _len; _key++) {
        regexps[_key] = arguments[_key];
    }

    return new RegExp(regexps.map(function (re) {
        return re.source;
    }).join(''));
}

/**
 * group wraps the regexp in a non-capturing group.
 */
function group() {
    return new RegExp('(?:' + expression.apply(undefined, arguments).source + ')');
}

/**
 * repeated wraps the regexp in a non-capturing group to get one or more matches.
 */
function optional() {
    return new RegExp(group.apply(undefined, arguments).source + '?');
}

/**
 * repeated wraps the regexp in a non-capturing group to get one or more matches.
 */
function repeated() {
    return new RegExp(group.apply(undefined, arguments).source + '+');
}

/**
 * anchored anchors the regular expression by adding start and end delimiters.
 */
function anchored() {
    return new RegExp('^' + expression.apply(undefined, arguments).source + '$');
}

/**
 * capture wraps the expression in a capturing group.
 */
function capture() {
    return new RegExp('(' + expression.apply(undefined, arguments).source + ')');
}

var alphaNumericRegexp = /[a-z0-9]+/;
var separatorRegexp = /(?:[._]|__|[-]*)/;

var nameComponentRegexp = expression(alphaNumericRegexp, optional(repeated(separatorRegexp, alphaNumericRegexp)));

var domainComponentRegexp = /(?:[a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9])/;

var domainRegexp = expression(domainComponentRegexp, optional(repeated(/\./, domainComponentRegexp)), optional(/:/, /[0-9]+/));

var tagRegexp = /[\w][\w.-]{0,127}/;

var digestRegexp = /[a-zA-Z][a-zA-Z0-9]*(?:[-_+.][a-zA-Z][a-zA-Z0-9]*)*[:][a-fA-F0-9]{32,}/;

var anchoredDigestRegexp = anchored(digestRegexp);

var nameRegexp = expression(optional(domainRegexp, /\//), nameComponentRegexp, optional(repeated(/\//, nameComponentRegexp)));

var anchoredNameRegexp = anchored(optional(capture(domainRegexp), /\//), capture(nameComponentRegexp, optional(repeated(/\//, nameComponentRegexp))));

var referenceRegexp = anchored(capture(nameRegexp), optional(/:/, capture(tagRegexp)), optional(/@/, capture(digestRegexp)));

var identifierRegexp = /[a-f0-9]{64}/;

var anchoredIdentifierRegexp = anchored(identifierRegexp);

Object.assign(exports, {
    referenceRegexp: referenceRegexp,
    anchoredNameRegexp: anchoredNameRegexp,
    anchoredDigestRegexp: anchoredDigestRegexp,
    anchoredIdentifierRegexp: anchoredIdentifierRegexp
});