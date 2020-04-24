(window.webpackJsonp=window.webpackJsonp||[]).push([[20],{122:function(e,t,r){"use strict";r.r(t),r.d(t,"frontMatter",(function(){return i})),r.d(t,"metadata",(function(){return o})),r.d(t,"rightToc",(function(){return p})),r.d(t,"default",(function(){return s}));var n=r(1),a=r(6),c=(r(0),r(193)),i={title:"File Parameter [object]"},o={id:"opspec/reference/structure/op-directory/op/parameter/file",title:"File Parameter [object]",description:"An object defining a parameter which accepts a [file typed value](../../../../types/file.md).",source:"@site/docs/opspec/reference/structure/op-directory/op/parameter/file.md",permalink:"/docs/opspec/reference/structure/op-directory/op/parameter/file",editUrl:"https://github.com/opctl/opctl/edit/master/docs/docs/opspec/reference/structure/op-directory/op/parameter/file.md",lastUpdatedBy:"Chris Dostert",lastUpdatedAt:1585210706},p=[{value:"Properties:",id:"properties",children:[{value:"default",id:"default",children:[]},{value:"description",id:"description",children:[]},{value:"isSecret",id:"issecret",children:[]}]}],l={rightToc:p},u="wrapper";function s(e){var t=e.components,r=Object(a.a)(e,["components"]);return Object(c.b)(u,Object(n.a)({},l,r,{components:t,mdxType:"MDXLayout"}),Object(c.b)("p",null,"An object defining a parameter which accepts a ",Object(c.b)("a",Object(n.a)({parentName:"p"},{href:"/docs/opspec/reference/types/file"}),"file typed value"),"."),Object(c.b)("h2",{id:"properties"},"Properties:"),Object(c.b)("ul",null,Object(c.b)("li",{parentName:"ul"},"must have:",Object(c.b)("ul",{parentName:"li"},Object(c.b)("li",{parentName:"ul"},Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"#description"}),"description")))),Object(c.b)("li",{parentName:"ul"},"may have:",Object(c.b)("ul",{parentName:"li"},Object(c.b)("li",{parentName:"ul"},Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"#default"}),"default")),Object(c.b)("li",{parentName:"ul"},Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"#issecret"}),"isSecret"))))),Object(c.b)("h3",{id:"default"},"default"),Object(c.b)("p",null,"A relative or absolute path string to use as the default value of the parameter when no argument is provided."),Object(c.b)("p",null,"If the value is..."),Object(c.b)("ul",null,Object(c.b)("li",{parentName:"ul"},"an absolute path, the value is interpreted from the root of the op."),Object(c.b)("li",{parentName:"ul"},"a relative path, the value is interpreted from the current working directory at the time the op is called.",Object(c.b)("blockquote",{parentName:"li"},Object(c.b)("p",{parentName:"blockquote"},"relative path defaults are ignored when an op is called from an op as there is no current working directory.")))),Object(c.b)("h3",{id:"description"},"description"),Object(c.b)("p",null,"A ",Object(c.b)("a",Object(n.a)({parentName:"p"},{href:"/docs/opspec/reference/structure/op-directory/op/markdown"}),"markdown [string]")," defining a human friendly description of the parameter."),Object(c.b)("h3",{id:"issecret"},"isSecret"),Object(c.b)("p",null,"A boolean indicating if the value of the parameter is secret. This will cause it to be hidden in UI's for example. "))}s.isMDXComponent=!0},193:function(e,t,r){"use strict";r.d(t,"a",(function(){return s})),r.d(t,"b",(function(){return m}));var n=r(0),a=r.n(n);function c(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function i(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function o(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?i(Object(r),!0).forEach((function(t){c(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):i(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function p(e,t){if(null==e)return{};var r,n,a=function(e,t){if(null==e)return{};var r,n,a={},c=Object.keys(e);for(n=0;n<c.length;n++)r=c[n],t.indexOf(r)>=0||(a[r]=e[r]);return a}(e,t);if(Object.getOwnPropertySymbols){var c=Object.getOwnPropertySymbols(e);for(n=0;n<c.length;n++)r=c[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(a[r]=e[r])}return a}var l=a.a.createContext({}),u=function(e){var t=a.a.useContext(l),r=t;return e&&(r="function"==typeof e?e(t):o({},t,{},e)),r},s=function(e){var t=u(e.components);return(a.a.createElement(l.Provider,{value:t},e.children))},b="mdxType",d={inlineCode:"code",wrapper:function(e){var t=e.children;return a.a.createElement(a.a.Fragment,{},t)}},f=Object(n.forwardRef)((function(e,t){var r=e.components,n=e.mdxType,c=e.originalType,i=e.parentName,l=p(e,["components","mdxType","originalType","parentName"]),s=u(r),b=n,f=s["".concat(i,".").concat(b)]||s[b]||d[b]||c;return r?a.a.createElement(f,o({ref:t},l,{components:r})):a.a.createElement(f,o({ref:t},l))}));function m(e,t){var r=arguments,n=t&&t.mdxType;if("string"==typeof e||n){var c=r.length,i=new Array(c);i[0]=f;var o={};for(var p in t)hasOwnProperty.call(t,p)&&(o[p]=t[p]);o.originalType=e,o[b]="string"==typeof e?e:n,i[1]=o;for(var l=2;l<c;l++)i[l]=r[l];return a.a.createElement.apply(null,i)}return a.a.createElement.apply(null,r)}f.displayName="MDXCreateElement"}}]);