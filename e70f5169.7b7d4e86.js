(window.webpackJsonp=window.webpackJsonp||[]).push([[65],{205:function(e,t,r){"use strict";r.r(t),r.d(t,"frontMatter",(function(){return i})),r.d(t,"metadata",(function(){return c})),r.d(t,"rightToc",(function(){return l})),r.d(t,"default",(function(){return s}));var n=r(1),o=r(9),a=(r(0),r(216)),i={title:"Boolean"},c={id:"reference/opspec/types/boolean",title:"Boolean",description:"Boolean typed values are a boolean i.e. `true` or `false`.",source:"@site/docs/reference/opspec/types/boolean.md",permalink:"/docs/reference/opspec/types/boolean",editUrl:"https://github.com/opctl/opctl/edit/main/website/docs/reference/opspec/types/boolean.md",lastUpdatedBy:"Chris Dostert",lastUpdatedAt:1587672399,sidebar:"docs",previous:{title:"Array",permalink:"/docs/reference/opspec/types/array"},next:{title:"Dir",permalink:"/docs/reference/opspec/types/dir"}},l=[{value:"Initialization",id:"initialization",children:[]},{value:"Coercion",id:"coercion",children:[]}],p={rightToc:l},b="wrapper";function s(e){var t=e.components,r=Object(o.a)(e,["components"]);return Object(a.b)(b,Object(n.a)({},p,r,{components:t,mdxType:"MDXLayout"}),Object(a.b)("p",null,"Boolean typed values are a boolean i.e. ",Object(a.b)("inlineCode",{parentName:"p"},"true")," or ",Object(a.b)("inlineCode",{parentName:"p"},"false"),"."),Object(a.b)("p",null,"Booleans..."),Object(a.b)("ul",null,Object(a.b)("li",{parentName:"ul"},"are immutable, i.e. assigning to a boolean results in a copy of the original boolean"),Object(a.b)("li",{parentName:"ul"},"can be passed in/out of ops via ",Object(a.b)("a",Object(n.a)({parentName:"li"},{href:"/docs/reference/opspec/op-directory/op/parameter/boolean"}),"boolean parameters")),Object(a.b)("li",{parentName:"ul"},"can be initialized via ",Object(a.b)("a",Object(n.a)({parentName:"li"},{href:"#initialization"}),"boolean initialization")),Object(a.b)("li",{parentName:"ul"},"are coerced according to ",Object(a.b)("a",Object(n.a)({parentName:"li"},{href:"#coercion"}),"boolean coercion"))),Object(a.b)("h3",{id:"initialization"},"Initialization"),Object(a.b)("p",null,"Boolean typed values can be constructed from a literal boolean."),Object(a.b)("h3",{id:"coercion"},"Coercion"),Object(a.b)("p",null,"Boolean typed values are coercible to:"),Object(a.b)("ul",null,Object(a.b)("li",{parentName:"ul"},Object(a.b)("a",Object(n.a)({parentName:"li"},{href:"/docs/reference/opspec/types/file"}),"file")," (will be serialized to JSON)"),Object(a.b)("li",{parentName:"ul"},Object(a.b)("a",Object(n.a)({parentName:"li"},{href:"/docs/reference/opspec/types/string"}),"string")," (will be serialized to JSON)")))}s.isMDXComponent=!0},216:function(e,t,r){"use strict";r.d(t,"a",(function(){return s})),r.d(t,"b",(function(){return O}));var n=r(0),o=r.n(n);function a(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function i(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function c(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?i(Object(r),!0).forEach((function(t){a(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):i(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function l(e,t){if(null==e)return{};var r,n,o=function(e,t){if(null==e)return{};var r,n,o={},a=Object.keys(e);for(n=0;n<a.length;n++)r=a[n],t.indexOf(r)>=0||(o[r]=e[r]);return o}(e,t);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(n=0;n<a.length;n++)r=a[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(o[r]=e[r])}return o}var p=o.a.createContext({}),b=function(e){var t=o.a.useContext(p),r=t;return e&&(r="function"==typeof e?e(t):c({},t,{},e)),r},s=function(e){var t=b(e.components);return(o.a.createElement(p.Provider,{value:t},e.children))},u="mdxType",f={inlineCode:"code",wrapper:function(e){var t=e.children;return o.a.createElement(o.a.Fragment,{},t)}},d=Object(n.forwardRef)((function(e,t){var r=e.components,n=e.mdxType,a=e.originalType,i=e.parentName,p=l(e,["components","mdxType","originalType","parentName"]),s=b(r),u=n,d=s["".concat(i,".").concat(u)]||s[u]||f[u]||a;return r?o.a.createElement(d,c({ref:t},p,{components:r})):o.a.createElement(d,c({ref:t},p))}));function O(e,t){var r=arguments,n=t&&t.mdxType;if("string"==typeof e||n){var a=r.length,i=new Array(a);i[0]=d;var c={};for(var l in t)hasOwnProperty.call(t,l)&&(c[l]=t[l]);c.originalType=e,c[u]="string"==typeof e?e:n,i[1]=c;for(var p=2;p<a;p++)i[p]=r[p];return o.a.createElement.apply(null,i)}return o.a.createElement.apply(null,r)}d.displayName="MDXCreateElement"}}]);