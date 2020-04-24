(window.webpackJsonp=window.webpackJsonp||[]).push([[82],{184:function(e,t,r){"use strict";r.r(t),r.d(t,"frontMatter",(function(){return o})),r.d(t,"metadata",(function(){return c})),r.d(t,"rightToc",(function(){return p})),r.d(t,"default",(function(){return b}));var n=r(1),a=r(6),i=(r(0),r(193)),o={title:"File Parameter [object]"},c={id:"reference/opspec/op-directory/op/parameter/file",title:"File Parameter [object]",description:"An object defining a parameter which accepts a [file typed value](../../../../types/file.md).",source:"@site/docs/reference/opspec/op-directory/op/parameter/file.md",permalink:"/docs/reference/opspec/op-directory/op/parameter/file",editUrl:"https://github.com/opctl/opctl/edit/master/docs/docs/reference/opspec/op-directory/op/parameter/file.md",lastUpdatedBy:"Chris Dostert",lastUpdatedAt:1587672399,sidebar:"docs",previous:{title:"Dir Parameter [object]",permalink:"/docs/reference/opspec/op-directory/op/parameter/dir"},next:{title:"Number Parameter [object]",permalink:"/docs/reference/opspec/op-directory/op/parameter/number"}},p=[{value:"Properties:",id:"properties",children:[{value:"default",id:"default",children:[]},{value:"description",id:"description",children:[]},{value:"isSecret",id:"issecret",children:[]}]}],l={rightToc:p},u="wrapper";function b(e){var t=e.components,r=Object(a.a)(e,["components"]);return Object(i.b)(u,Object(n.a)({},l,r,{components:t,mdxType:"MDXLayout"}),Object(i.b)("p",null,"An object defining a parameter which accepts a ",Object(i.b)("a",Object(n.a)({parentName:"p"},{href:"../../../../types/file.md"}),"file typed value"),"."),Object(i.b)("h2",{id:"properties"},"Properties:"),Object(i.b)("ul",null,Object(i.b)("li",{parentName:"ul"},"must have:",Object(i.b)("ul",{parentName:"li"},Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",Object(n.a)({parentName:"li"},{href:"#description"}),"description")))),Object(i.b)("li",{parentName:"ul"},"may have:",Object(i.b)("ul",{parentName:"li"},Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",Object(n.a)({parentName:"li"},{href:"#default"}),"default")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",Object(n.a)({parentName:"li"},{href:"#issecret"}),"isSecret"))))),Object(i.b)("h3",{id:"default"},"default"),Object(i.b)("p",null,"A relative or absolute path string to use as the default value of the parameter when no argument is provided."),Object(i.b)("p",null,"If the value is..."),Object(i.b)("ul",null,Object(i.b)("li",{parentName:"ul"},"an absolute path, the value is interpreted from the root of the op."),Object(i.b)("li",{parentName:"ul"},"a relative path, the value is interpreted from the current working directory at the time the op is called.",Object(i.b)("blockquote",{parentName:"li"},Object(i.b)("p",{parentName:"blockquote"},"relative path defaults are ignored when an op is called from an op as there is no current working directory.")))),Object(i.b)("h3",{id:"description"},"description"),Object(i.b)("p",null,"A ",Object(i.b)("a",Object(n.a)({parentName:"p"},{href:"/docs/reference/opspec/op-directory/op/markdown"}),"markdown [string]")," defining a human friendly description of the parameter."),Object(i.b)("h3",{id:"issecret"},"isSecret"),Object(i.b)("p",null,"A boolean indicating if the value of the parameter is secret. This will cause it to be hidden in UI's for example. "))}b.isMDXComponent=!0},193:function(e,t,r){"use strict";r.d(t,"a",(function(){return b})),r.d(t,"b",(function(){return m}));var n=r(0),a=r.n(n);function i(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function o(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function c(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?o(Object(r),!0).forEach((function(t){i(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):o(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function p(e,t){if(null==e)return{};var r,n,a=function(e,t){if(null==e)return{};var r,n,a={},i=Object.keys(e);for(n=0;n<i.length;n++)r=i[n],t.indexOf(r)>=0||(a[r]=e[r]);return a}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(n=0;n<i.length;n++)r=i[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(a[r]=e[r])}return a}var l=a.a.createContext({}),u=function(e){var t=a.a.useContext(l),r=t;return e&&(r="function"==typeof e?e(t):c({},t,{},e)),r},b=function(e){var t=u(e.components);return(a.a.createElement(l.Provider,{value:t},e.children))},s="mdxType",d={inlineCode:"code",wrapper:function(e){var t=e.children;return a.a.createElement(a.a.Fragment,{},t)}},f=Object(n.forwardRef)((function(e,t){var r=e.components,n=e.mdxType,i=e.originalType,o=e.parentName,l=p(e,["components","mdxType","originalType","parentName"]),b=u(r),s=n,f=b["".concat(o,".").concat(s)]||b[s]||d[s]||i;return r?a.a.createElement(f,c({ref:t},l,{components:r})):a.a.createElement(f,c({ref:t},l))}));function m(e,t){var r=arguments,n=t&&t.mdxType;if("string"==typeof e||n){var i=r.length,o=new Array(i);o[0]=f;var c={};for(var p in t)hasOwnProperty.call(t,p)&&(c[p]=t[p]);c.originalType=e,c[s]="string"==typeof e?e:n,o[1]=c;for(var l=2;l<i;l++)o[l]=r[l];return a.a.createElement.apply(null,o)}return a.a.createElement.apply(null,r)}f.displayName="MDXCreateElement"}}]);