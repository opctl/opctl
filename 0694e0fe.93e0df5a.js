(window.webpackJsonp=window.webpackJsonp||[]).push([[3],{193:function(e,t,r){"use strict";r.d(t,"a",(function(){return s})),r.d(t,"b",(function(){return O}));var n=r(0),a=r.n(n);function i(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function c(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function o(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?c(Object(r),!0).forEach((function(t){i(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):c(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function l(e,t){if(null==e)return{};var r,n,a=function(e,t){if(null==e)return{};var r,n,a={},i=Object.keys(e);for(n=0;n<i.length;n++)r=i[n],t.indexOf(r)>=0||(a[r]=e[r]);return a}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(n=0;n<i.length;n++)r=i[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(a[r]=e[r])}return a}var p=a.a.createContext({}),b=function(e){var t=a.a.useContext(p),r=t;return e&&(r="function"==typeof e?e(t):o({},t,{},e)),r},s=function(e){var t=b(e.components);return(a.a.createElement(p.Provider,{value:t},e.children))},f="mdxType",u={inlineCode:"code",wrapper:function(e){var t=e.children;return a.a.createElement(a.a.Fragment,{},t)}},m=Object(n.forwardRef)((function(e,t){var r=e.components,n=e.mdxType,i=e.originalType,c=e.parentName,p=l(e,["components","mdxType","originalType","parentName"]),s=b(r),f=n,m=s["".concat(c,".").concat(f)]||s[f]||u[f]||i;return r?a.a.createElement(m,o({ref:t},p,{components:r})):a.a.createElement(m,o({ref:t},p))}));function O(e,t){var r=arguments,n=t&&t.mdxType;if("string"==typeof e||n){var i=r.length,c=new Array(i);c[0]=m;var o={};for(var l in t)hasOwnProperty.call(t,l)&&(o[l]=t[l]);o.originalType=e,o[f]="string"==typeof e?e:n,c[1]=o;for(var p=2;p<i;p++)c[p]=r[p];return a.a.createElement.apply(null,c)}return a.a.createElement.apply(null,r)}m.displayName="MDXCreateElement"},99:function(e,t,r){"use strict";r.r(t),r.d(t,"frontMatter",(function(){return c})),r.d(t,"metadata",(function(){return o})),r.d(t,"rightToc",(function(){return l})),r.d(t,"default",(function(){return s}));var n=r(1),a=r(6),i=(r(0),r(193)),c={title:"File"},o={id:"opspec/reference/types/file",title:"File",description:"File typed values are a filesystem file entry.",source:"@site/docs/opspec/reference/types/file.md",permalink:"/docs/opspec/reference/types/file",editUrl:"https://github.com/opctl/opctl/edit/master/docs/docs/opspec/reference/types/file.md",lastUpdatedBy:"Chris Dostert",lastUpdatedAt:1583255634},l=[{value:"Initialization",id:"initialization",children:[]},{value:"Coercion",id:"coercion",children:[]}],p={rightToc:l},b="wrapper";function s(e){var t=e.components,r=Object(a.a)(e,["components"]);return Object(i.b)(b,Object(n.a)({},p,r,{components:t,mdxType:"MDXLayout"}),Object(i.b)("p",null,"File typed values are a filesystem file entry."),Object(i.b)("p",null,"Files..."),Object(i.b)("ul",null,Object(i.b)("li",{parentName:"ul"},"are mutable, i.e. making changes to a file results in the file being changed everywhere it's referenced"),Object(i.b)("li",{parentName:"ul"},"can be passed in/out of ops via ",Object(i.b)("a",Object(n.a)({parentName:"li"},{href:"/docs/opspec/reference/structure/op-directory/op/parameter/file"}),"file parameters")),Object(i.b)("li",{parentName:"ul"},"can be initialized via ",Object(i.b)("a",Object(n.a)({parentName:"li"},{href:"#initialization"}),"file initialization")),Object(i.b)("li",{parentName:"ul"},"are coerced according to ",Object(i.b)("a",Object(n.a)({parentName:"li"},{href:"#coercion"}),"file coercion"))),Object(i.b)("h3",{id:"initialization"},"Initialization"),Object(i.b)("p",null,"File typed values can be constructed from a literal string or templated string; see ",Object(i.b)("a",Object(n.a)({parentName:"p"},{href:"/docs/opspec/reference/types/string#initialization"}),"string initialization"),"."),Object(i.b)("h3",{id:"coercion"},"Coercion"),Object(i.b)("p",null,"File typed values are coercible to:"),Object(i.b)("ul",null,Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",Object(n.a)({parentName:"li"},{href:"/docs/opspec/reference/types/boolean"}),"boolean")," (files which are empty, all ",Object(i.b)("inlineCode",{parentName:"li"},'"0"'),", or (case insensitive) ",Object(i.b)("inlineCode",{parentName:"li"},'"f"')," or ",Object(i.b)("inlineCode",{parentName:"li"},'"false"')," coerce to false; all else coerce to true)"),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",Object(n.a)({parentName:"li"},{href:"/docs/opspec/reference/types/array"}),"array")," (if value of file is an array in JSON format)"),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",Object(n.a)({parentName:"li"},{href:"/docs/opspec/reference/types/number"}),"number")," (if value of file is numeric)"),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",Object(n.a)({parentName:"li"},{href:"/docs/opspec/reference/types/object"}),"object")," (if value of file is an object in JSON format)"),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",Object(n.a)({parentName:"li"},{href:"/docs/opspec/reference/types/string"}),"string"))))}s.isMDXComponent=!0}}]);