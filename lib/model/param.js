/**
* Typed param of an op
* @typedef {Object} Param
* @property {ArrayParam|null} [array] array typed param
* @property {BooleanParam|null} [boolean] boolean typed param
* @property {DirParam|null} [dir] dir typed param
* @property {FileParam|null} [file] string typed param
* @property {NumberParam|null} [number] number typed param
* @property {ObjectParam|null} [object] object typed param
* @property {SocketParam|null} [socket] socket typed param
* @property {StringParam|null} [string] string typed param
*/

/**
* Param of type object
* @typedef {Object} ArrayParam
* @property {Object} [constraints]
* @property {Array<*>|null} [default] default value for param; if any
* @property {string} [description]
* @property {boolean} [isSecret]
*/

/**
* Param of type boolean
* @typedef {Object} BooleanParam
* @property {boolean|null} [default] default value for param; if any
* @property {string} [description]
*/

/**
* Param of type dir
* @typedef {Object} DirParam
* @property {string|null} [default] default value for param; if any
* @property {string} [description]
* @property {boolean} [isSecret]
*/

/**
* Param of type file
* @typedef {Object} FileParam
* @property {string|null} [default] default value for param; if any
* @property {string} [description]
* @property {boolean} [isSecret]
*/

/**
* Param of type number
* @typedef {Object} NumberParam
* @property {Object} [constraints]
* @property {number|null} [default] default value for param; if any
* @property {string} [description]
* @property {boolean} [isSecret]
*/

/**
* Param of type object
* @typedef {Object} ObjectParam
* @property {Object} [constraints]
* @property {Object|null} [default] default value for param; if any
* @property {string} [description]
* @property {boolean} [isSecret]
*/

/**
* Param of type socket
* @typedef {Object} SocketParam
* @property {string|null} [default] default value for param; if any
* @property {boolean} [isSecret]
*/

/**
* Param of type string
* @typedef {Object} StringParam
* @property {Object} [constraints]
* @property {string|null} [default] default value for param; if any
* @property {string} [description]
* @property {boolean} [isSecret]
*/
