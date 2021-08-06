(window.webpackJsonp=window.webpackJsonp||[]).push([[32],{172:function(e,r,t){"use strict";t.r(r),t.d(r,"frontMatter",(function(){return c})),t.d(r,"metadata",(function(){return l})),t.d(r,"rightToc",(function(){return p})),t.d(r,"default",(function(){return u}));var n=t(1),a=t(9),o=(t(0),t(216)),c={title:"Parallel Loop Call [object]"},l={id:"reference/opspec/op-directory/op/call/parallel-loop",title:"Parallel Loop Call [object]",description:"An object defining a call loop in which all iterations happen in parallel (all at once without order).",source:"@site/docs/reference/opspec/op-directory/op/call/parallel-loop.md",permalink:"/docs/reference/opspec/op-directory/op/call/parallel-loop",editUrl:"https://github.com/opctl/opctl/edit/main/website/docs/reference/opspec/op-directory/op/call/parallel-loop.md",lastUpdatedBy:"Chris Dostert",lastUpdatedAt:1587672399,sidebar:"docs",previous:{title:"Op Call [object]",permalink:"/docs/reference/opspec/op-directory/op/call/op"},next:{title:"Predicate [object]",permalink:"/docs/reference/opspec/op-directory/op/call/predicate"}},p=[{value:"Properties",id:"properties",children:[{value:"range",id:"range",children:[]},{value:"run",id:"run",children:[]},{value:"vars",id:"vars",children:[]}]}],i={rightToc:p},b="wrapper";function u(e){var r=e.components,t=Object(a.a)(e,["components"]);return Object(o.b)(b,Object(n.a)({},i,t,{components:r,mdxType:"MDXLayout"}),Object(o.b)("p",null,"An object defining a call loop in which all iterations happen in parallel (all at once without order)."),Object(o.b)("h2",{id:"properties"},"Properties"),Object(o.b)("ul",null,Object(o.b)("li",{parentName:"ul"},"must have ",Object(o.b)("ul",{parentName:"li"},Object(o.b)("li",{parentName:"ul"},Object(o.b)("a",Object(n.a)({parentName:"li"},{href:"#range"}),"range")),Object(o.b)("li",{parentName:"ul"},Object(o.b)("a",Object(n.a)({parentName:"li"},{href:"#run"}),"run")))),Object(o.b)("li",{parentName:"ul"},"may have",Object(o.b)("ul",{parentName:"li"},Object(o.b)("li",{parentName:"ul"},Object(o.b)("a",Object(n.a)({parentName:"li"},{href:"#vars"}),"vars"))))),Object(o.b)("h3",{id:"range"},"range"),Object(o.b)("p",null,"A ",Object(o.b)("a",Object(n.a)({parentName:"p"},{href:"/docs/reference/opspec/op-directory/op/call/rangeable-value"}),"rangeable value")," to loop over."),Object(o.b)("h3",{id:"run"},"run"),Object(o.b)("p",null,"A ",Object(o.b)("a",Object(n.a)({parentName:"p"},{href:"/docs/reference/opspec/op-directory/op/call/index"}),"call [object]")," defining a call run on each iteration of the loop"),Object(o.b)("h3",{id:"vars"},"vars"),Object(o.b)("p",null,"A ",Object(o.b)("a",Object(n.a)({parentName:"p"},{href:"/docs/reference/opspec/op-directory/op/call/loop-vars"}),"loop-vars [object]")," binding iteration info to variables."))}u.isMDXComponent=!0},216:function(e,r,t){"use strict";t.d(r,"a",(function(){return u})),t.d(r,"b",(function(){return O}));var n=t(0),a=t.n(n);function o(e,r,t){return r in e?Object.defineProperty(e,r,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[r]=t,e}function c(e,r){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);r&&(n=n.filter((function(r){return Object.getOwnPropertyDescriptor(e,r).enumerable}))),t.push.apply(t,n)}return t}function l(e){for(var r=1;r<arguments.length;r++){var t=null!=arguments[r]?arguments[r]:{};r%2?c(Object(t),!0).forEach((function(r){o(e,r,t[r])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):c(Object(t)).forEach((function(r){Object.defineProperty(e,r,Object.getOwnPropertyDescriptor(t,r))}))}return e}function p(e,r){if(null==e)return{};var t,n,a=function(e,r){if(null==e)return{};var t,n,a={},o=Object.keys(e);for(n=0;n<o.length;n++)t=o[n],r.indexOf(t)>=0||(a[t]=e[t]);return a}(e,r);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(n=0;n<o.length;n++)t=o[n],r.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(a[t]=e[t])}return a}var i=a.a.createContext({}),b=function(e){var r=a.a.useContext(i),t=r;return e&&(t="function"==typeof e?e(r):l({},r,{},e)),t},u=function(e){var r=b(e.components);return(a.a.createElement(i.Provider,{value:r},e.children))},s="mdxType",d={inlineCode:"code",wrapper:function(e){var r=e.children;return a.a.createElement(a.a.Fragment,{},r)}},f=Object(n.forwardRef)((function(e,r){var t=e.components,n=e.mdxType,o=e.originalType,c=e.parentName,i=p(e,["components","mdxType","originalType","parentName"]),u=b(t),s=n,f=u["".concat(c,".").concat(s)]||u[s]||d[s]||o;return t?a.a.createElement(f,l({ref:r},i,{components:t})):a.a.createElement(f,l({ref:r},i))}));function O(e,r){var t=arguments,n=r&&r.mdxType;if("string"==typeof e||n){var o=t.length,c=new Array(o);c[0]=f;var l={};for(var p in r)hasOwnProperty.call(r,p)&&(l[p]=r[p]);l.originalType=e,l[s]="string"==typeof e?e:n,c[1]=l;for(var i=2;i<o;i++)c[i]=t[i];return a.a.createElement.apply(null,c)}return a.a.createElement.apply(null,t)}f.displayName="MDXCreateElement"}}]);