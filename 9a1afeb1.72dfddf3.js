(window.webpackJsonp=window.webpackJsonp||[]).push([[43],{183:function(e,t,r){"use strict";r.r(t),r.d(t,"frontMatter",(function(){return c})),r.d(t,"metadata",(function(){return i})),r.d(t,"rightToc",(function(){return p})),r.d(t,"default",(function(){return u}));var n=r(1),o=r(9),a=(r(0),r(216)),c={title:"Markdown [string]"},i={id:"reference/opspec/op-directory/op/markdown",title:"Markdown [string]",description:"A [Commonmark](http://commonmark.org/) (plus table extensions) string used to create markup such as descriptions.",source:"@site/docs/reference/opspec/op-directory/op/markdown.md",permalink:"/docs/reference/opspec/op-directory/op/markdown",editUrl:"https://github.com/opctl/opctl/edit/main/website/docs/reference/opspec/op-directory/op/markdown.md",lastUpdatedBy:"Chris Dostert",lastUpdatedAt:1587672399,sidebar:"docs",previous:{title:"Initializer [array|boolean|number|string|object]",permalink:"/docs/reference/opspec/op-directory/op/initializer"},next:{title:"Variable Reference [string]",permalink:"/docs/reference/opspec/op-directory/op/variable-reference"}},p=[],l={rightToc:p},s="wrapper";function u(e){var t=e.components,r=Object(o.a)(e,["components"]);return Object(a.b)(s,Object(n.a)({},l,r,{components:t,mdxType:"MDXLayout"}),Object(a.b)("p",null,"A ",Object(a.b)("a",Object(n.a)({parentName:"p"},{href:"http://commonmark.org/"}),"Commonmark")," (plus table extensions) string used to create markup such as descriptions."),Object(a.b)("blockquote",null,Object(a.b)("p",{parentName:"blockquote"},"relative &/or absolute paths will be resolved from the root of the op")),Object(a.b)("h4",{id:"examples"},"Examples"),Object(a.b)("pre",null,Object(a.b)("code",Object(n.a)({parentName:"pre"},{className:"language-markdown"}),"checkout [this op's op.yml](op.yml)\ncheckout this image ![my image](/my-image.png)\n# h1\n## h2\n### h3\n#### h4\n##### h5\n###### h6\n**bolded**\n*italicized*\n~~striken~~\n- unordered item\n  - unordered subitem\n1. ordered item\n  1. ordered subitem\n| title 1 | title 2 | title 3 |\n|:-------:|:-------:|:-------:|\n| entry 1 | entry 2 | entry 3 |\n")))}u.isMDXComponent=!0},216:function(e,t,r){"use strict";r.d(t,"a",(function(){return u})),r.d(t,"b",(function(){return f}));var n=r(0),o=r.n(n);function a(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function c(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function i(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?c(Object(r),!0).forEach((function(t){a(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):c(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function p(e,t){if(null==e)return{};var r,n,o=function(e,t){if(null==e)return{};var r,n,o={},a=Object.keys(e);for(n=0;n<a.length;n++)r=a[n],t.indexOf(r)>=0||(o[r]=e[r]);return o}(e,t);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(n=0;n<a.length;n++)r=a[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(o[r]=e[r])}return o}var l=o.a.createContext({}),s=function(e){var t=o.a.useContext(l),r=t;return e&&(r="function"==typeof e?e(t):i({},t,{},e)),r},u=function(e){var t=s(e.components);return(o.a.createElement(l.Provider,{value:t},e.children))},m="mdxType",d={inlineCode:"code",wrapper:function(e){var t=e.children;return o.a.createElement(o.a.Fragment,{},t)}},b=Object(n.forwardRef)((function(e,t){var r=e.components,n=e.mdxType,a=e.originalType,c=e.parentName,l=p(e,["components","mdxType","originalType","parentName"]),u=s(r),m=n,b=u["".concat(c,".").concat(m)]||u[m]||d[m]||a;return r?o.a.createElement(b,i({ref:t},l,{components:r})):o.a.createElement(b,i({ref:t},l))}));function f(e,t){var r=arguments,n=t&&t.mdxType;if("string"==typeof e||n){var a=r.length,c=new Array(a);c[0]=b;var i={};for(var p in t)hasOwnProperty.call(t,p)&&(i[p]=t[p]);i.originalType=e,i[m]="string"==typeof e?e:n,c[1]=i;for(var l=2;l<a;l++)c[l]=r[l];return o.a.createElement.apply(null,c)}return o.a.createElement.apply(null,r)}b.displayName="MDXCreateElement"}}]);