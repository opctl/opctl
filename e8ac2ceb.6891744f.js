(window.webpackJsonp=window.webpackJsonp||[]).push([[67],{207:function(e,t,r){"use strict";r.r(t),r.d(t,"frontMatter",(function(){return c})),r.d(t,"metadata",(function(){return o})),r.d(t,"rightToc",(function(){return p})),r.d(t,"default",(function(){return s}));var n=r(1),a=r(9),i=(r(0),r(216)),c={title:"File Parameter [object]"},o={id:"reference/opspec/op-directory/op/parameter/file",title:"File Parameter [object]",description:"An object defining a parameter which accepts a [file typed value](../../../types/file.md).",source:"@site/docs/reference/opspec/op-directory/op/parameter/file.md",permalink:"/docs/reference/opspec/op-directory/op/parameter/file",editUrl:"https://github.com/opctl/opctl/edit/main/website/docs/reference/opspec/op-directory/op/parameter/file.md",lastUpdatedBy:"=",lastUpdatedAt:1628636355,sidebar:"docs",previous:{title:"Dir Parameter [object]",permalink:"/docs/reference/opspec/op-directory/op/parameter/dir"},next:{title:"Number Parameter [object]",permalink:"/docs/reference/opspec/op-directory/op/parameter/number"}},p=[{value:"Properties:",id:"properties",children:[{value:"default",id:"default",children:[]},{value:"isSecret",id:"issecret",children:[]}]}],l={rightToc:p},u="wrapper";function s(e){var t=e.components,r=Object(a.a)(e,["components"]);return Object(i.b)(u,Object(n.a)({},l,r,{components:t,mdxType:"MDXLayout"}),Object(i.b)("p",null,"An object defining a parameter which accepts a ",Object(i.b)("a",Object(n.a)({parentName:"p"},{href:"/docs/reference/opspec/types/file"}),"file typed value"),"."),Object(i.b)("h2",{id:"properties"},"Properties:"),Object(i.b)("ul",null,Object(i.b)("li",{parentName:"ul"},"may have:",Object(i.b)("ul",{parentName:"li"},Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",Object(n.a)({parentName:"li"},{href:"#default"}),"default")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",Object(n.a)({parentName:"li"},{href:"#issecret"}),"isSecret"))))),Object(i.b)("h3",{id:"default"},"default"),Object(i.b)("p",null,"A ",Object(i.b)("a",Object(n.a)({parentName:"p"},{href:"/docs/reference/opspec/types/file#initialization"}),"file initializer")," to use as the value of the parameter when no argument is provided."),Object(i.b)("p",null,"If the value is a relative path it will be resolved from the current working directory of the caller. If no current working directory exists, such as when the caller is an op or web UI, the default will be ignored."),Object(i.b)("h3",{id:"issecret"},"isSecret"),Object(i.b)("p",null,"A boolean indicating if the value of the parameter is secret. This will cause it to be hidden in UI's for example. "))}s.isMDXComponent=!0},216:function(e,t,r){"use strict";r.d(t,"a",(function(){return s})),r.d(t,"b",(function(){return m}));var n=r(0),a=r.n(n);function i(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function c(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function o(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?c(Object(r),!0).forEach((function(t){i(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):c(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function p(e,t){if(null==e)return{};var r,n,a=function(e,t){if(null==e)return{};var r,n,a={},i=Object.keys(e);for(n=0;n<i.length;n++)r=i[n],t.indexOf(r)>=0||(a[r]=e[r]);return a}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(n=0;n<i.length;n++)r=i[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(a[r]=e[r])}return a}var l=a.a.createContext({}),u=function(e){var t=a.a.useContext(l),r=t;return e&&(r="function"==typeof e?e(t):o({},t,{},e)),r},s=function(e){var t=u(e.components);return(a.a.createElement(l.Provider,{value:t},e.children))},f="mdxType",b={inlineCode:"code",wrapper:function(e){var t=e.children;return a.a.createElement(a.a.Fragment,{},t)}},d=Object(n.forwardRef)((function(e,t){var r=e.components,n=e.mdxType,i=e.originalType,c=e.parentName,l=p(e,["components","mdxType","originalType","parentName"]),s=u(r),f=n,d=s["".concat(c,".").concat(f)]||s[f]||b[f]||i;return r?a.a.createElement(d,o({ref:t},l,{components:r})):a.a.createElement(d,o({ref:t},l))}));function m(e,t){var r=arguments,n=t&&t.mdxType;if("string"==typeof e||n){var i=r.length,c=new Array(i);c[0]=d;var o={};for(var p in t)hasOwnProperty.call(t,p)&&(o[p]=t[p]);o.originalType=e,o[f]="string"==typeof e?e:n,c[1]=o;for(var l=2;l<i;l++)c[l]=r[l];return a.a.createElement.apply(null,c)}return a.a.createElement.apply(null,r)}d.displayName="MDXCreateElement"}}]);