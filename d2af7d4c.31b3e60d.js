(window.webpackJsonp=window.webpackJsonp||[]).push([[58],{198:function(e,t,r){"use strict";r.r(t),r.d(t,"frontMatter",(function(){return a})),r.d(t,"metadata",(function(){return i})),r.d(t,"rightToc",(function(){return p})),r.d(t,"default",(function(){return s}));var n=r(1),o=r(9),c=(r(0),r(216)),a={sidebar_label:"create",title:"opctl op create"},i={id:"reference/cli/op/create",title:"opctl op create",description:"```sh",source:"@site/docs/reference/cli/op/create.md",permalink:"/docs/reference/cli/op/create",editUrl:"https://github.com/opctl/opctl/edit/main/website/docs/reference/cli/op/create.md",lastUpdatedBy:"Chris Dostert",lastUpdatedAt:1604531186,sidebar_label:"create",sidebar:"docs",previous:{title:"opctl op",permalink:"/docs/reference/cli/op/index"},next:{title:"opctl op install",permalink:"/docs/reference/cli/op/install"}},p=[{value:"Arguments",id:"arguments",children:[{value:"<code>NAME</code>",id:"name",children:[]}]},{value:"Options",id:"options",children:[{value:"<code>-d</code> or <code>--description</code>",id:"-d-or---description",children:[]},{value:"<code>--path</code> <em>default: <code>.opspec</code></em>",id:"--path-default-opspec",children:[]}]},{value:"Global Options",id:"global-options",children:[]},{value:"Examples",id:"examples",children:[]}],l={rightToc:p},d="wrapper";function s(e){var t=e.components,r=Object(o.a)(e,["components"]);return Object(c.b)(d,Object(n.a)({},l,r,{components:t,mdxType:"MDXLayout"}),Object(c.b)("pre",null,Object(c.b)("code",Object(n.a)({parentName:"pre"},{className:"language-sh"}),"opctl op create [OPTIONS] NAME\n")),Object(c.b)("p",null,"Create an op."),Object(c.b)("h2",{id:"arguments"},"Arguments"),Object(c.b)("h3",{id:"name"},Object(c.b)("inlineCode",{parentName:"h3"},"NAME")),Object(c.b)("p",null,"Name of the op"),Object(c.b)("h2",{id:"options"},"Options"),Object(c.b)("h3",{id:"-d-or---description"},Object(c.b)("inlineCode",{parentName:"h3"},"-d")," or ",Object(c.b)("inlineCode",{parentName:"h3"},"--description")),Object(c.b)("p",null,"Description of the op"),Object(c.b)("h3",{id:"--path-default-opspec"},Object(c.b)("inlineCode",{parentName:"h3"},"--path")," ",Object(c.b)("em",{parentName:"h3"},"default: ",Object(c.b)("inlineCode",{parentName:"em"},".opspec"))),Object(c.b)("p",null,"Path to create the op at"),Object(c.b)("h2",{id:"global-options"},"Global Options"),Object(c.b)("p",null,"see ",Object(c.b)("a",Object(n.a)({parentName:"p"},{href:"/docs/reference/cli/global-options"}),"global options")),Object(c.b)("h2",{id:"examples"},"Examples"),Object(c.b)("pre",null,Object(c.b)("code",Object(n.a)({parentName:"pre"},{className:"language-sh"}),'opctl op create -d "my awesome op description" --path some/path my-awesome-op-name\n')))}s.isMDXComponent=!0},216:function(e,t,r){"use strict";r.d(t,"a",(function(){return s})),r.d(t,"b",(function(){return f}));var n=r(0),o=r.n(n);function c(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function a(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function i(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?a(Object(r),!0).forEach((function(t){c(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):a(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function p(e,t){if(null==e)return{};var r,n,o=function(e,t){if(null==e)return{};var r,n,o={},c=Object.keys(e);for(n=0;n<c.length;n++)r=c[n],t.indexOf(r)>=0||(o[r]=e[r]);return o}(e,t);if(Object.getOwnPropertySymbols){var c=Object.getOwnPropertySymbols(e);for(n=0;n<c.length;n++)r=c[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(o[r]=e[r])}return o}var l=o.a.createContext({}),d=function(e){var t=o.a.useContext(l),r=t;return e&&(r="function"==typeof e?e(t):i({},t,{},e)),r},s=function(e){var t=d(e.components);return(o.a.createElement(l.Provider,{value:t},e.children))},b="mdxType",u={inlineCode:"code",wrapper:function(e){var t=e.children;return o.a.createElement(o.a.Fragment,{},t)}},m=Object(n.forwardRef)((function(e,t){var r=e.components,n=e.mdxType,c=e.originalType,a=e.parentName,l=p(e,["components","mdxType","originalType","parentName"]),s=d(r),b=n,m=s["".concat(a,".").concat(b)]||s[b]||u[b]||c;return r?o.a.createElement(m,i({ref:t},l,{components:r})):o.a.createElement(m,i({ref:t},l))}));function f(e,t){var r=arguments,n=t&&t.mdxType;if("string"==typeof e||n){var c=r.length,a=new Array(c);a[0]=m;var i={};for(var p in t)hasOwnProperty.call(t,p)&&(i[p]=t[p]);i.originalType=e,i[b]="string"==typeof e?e:n,a[1]=i;for(var l=2;l<c;l++)a[l]=r[l];return o.a.createElement.apply(null,a)}return o.a.createElement.apply(null,r)}m.displayName="MDXCreateElement"}}]);