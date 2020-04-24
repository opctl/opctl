(window.webpackJsonp=window.webpackJsonp||[]).push([[65],{166:function(e,t,r){"use strict";r.r(t),r.d(t,"frontMatter",(function(){return o})),r.d(t,"metadata",(function(){return i})),r.d(t,"rightToc",(function(){return p})),r.d(t,"default",(function(){return u}));var n=r(1),a=r(6),c=(r(0),r(193)),o={sidebar_label:"Index",title:"Op [object]"},i={id:"reference/opspec/op-directory/op/index",title:"Op [object]",description:"An object which defines an operations inputs, outputs, call graph... etc.",source:"@site/docs/reference/opspec/op-directory/op/index.md",permalink:"/docs/reference/opspec/op-directory/op/index",editUrl:"https://github.com/opctl/opctl/edit/master/docs/docs/reference/opspec/op-directory/op/index.md",lastUpdatedBy:"Chris Dostert",lastUpdatedAt:1587672399,sidebar_label:"Index",sidebar:"docs",previous:{title:"Op [directory]",permalink:"/docs/reference/opspec/op-directory/index"},next:{title:"Call [object]",permalink:"/docs/reference/opspec/op-directory/op/call/index"}},p=[{value:"Properties",id:"properties",children:[{value:"name",id:"name",children:[]},{value:"description",id:"description",children:[]},{value:"inputs",id:"inputs",children:[]},{value:"outputs",id:"outputs",children:[]},{value:"opspec",id:"opspec",children:[]},{value:"run",id:"run",children:[]},{value:"version",id:"version",children:[]}]}],b={rightToc:p},l="wrapper";function u(e){var t=e.components,r=Object(a.a)(e,["components"]);return Object(c.b)(l,Object(n.a)({},b,r,{components:t,mdxType:"MDXLayout"}),Object(c.b)("p",null,"An object which defines an operations inputs, outputs, call graph... etc."),Object(c.b)("h2",{id:"properties"},"Properties"),Object(c.b)("ul",null,Object(c.b)("li",{parentName:"ul"},"must have",Object(c.b)("ul",{parentName:"li"},Object(c.b)("li",{parentName:"ul"},Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"#name"}),"name")))),Object(c.b)("li",{parentName:"ul"},"may have",Object(c.b)("ul",{parentName:"li"},Object(c.b)("li",{parentName:"ul"},Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"#description"}),"description")),Object(c.b)("li",{parentName:"ul"},Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"#inputs"}),"inputs")),Object(c.b)("li",{parentName:"ul"},Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"#opspec"}),"opspec")),Object(c.b)("li",{parentName:"ul"},Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"#outputs"}),"outputs")),Object(c.b)("li",{parentName:"ul"},Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"#run"}),"run")),Object(c.b)("li",{parentName:"ul"},Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"#version"}),"version"))))),Object(c.b)("h3",{id:"name"},"name"),Object(c.b)("p",null,"A string defining a human friendly identifier for the operation."),Object(c.b)("blockquote",null,Object(c.b)("p",{parentName:"blockquote"},"It's considered good practice to make ",Object(c.b)("inlineCode",{parentName:"p"},"name")," unique by using domain\n&/or path based namespacing.")),Object(c.b)("p",null,"Ops MAY be network resolvable; therefore ",Object(c.b)("inlineCode",{parentName:"p"},"name")," MUST be a valid\n",Object(c.b)("a",Object(n.a)({parentName:"p"},{href:"https://tools.ietf.org/html/rfc3986#section-4.1"}),"uri-reference")),Object(c.b)("p",null,"example:"),Object(c.b)("pre",null,Object(c.b)("code",Object(n.a)({parentName:"pre"},{className:"language-yaml"}),"name: `github.com/opspec-pkgs/jwt.encode`\n")),Object(c.b)("h3",{id:"description"},"description"),Object(c.b)("p",null,"A ",Object(c.b)("a",Object(n.a)({parentName:"p"},{href:"/docs/reference/opspec/op-directory/op/markdown"}),"markdown [string]")," defining a human friendly description of the op (since v0.1.6)."),Object(c.b)("h3",{id:"inputs"},"inputs"),Object(c.b)("p",null,"An object defining input parameters of the operation."),Object(c.b)("p",null,"For each property:"),Object(c.b)("ul",null,Object(c.b)("li",{parentName:"ul"},"key is an ",Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"/docs/reference/opspec/op-directory/op/identifier"}),"identifier [string]")," defining the name of the input."),Object(c.b)("li",{parentName:"ul"},"value is a ",Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"/docs/reference/opspec/op-directory/op/parameter/index"}),"parameter [object]")," defining the output. ")),Object(c.b)("h3",{id:"outputs"},"outputs"),Object(c.b)("p",null,"An object defining output parameters of the operation."),Object(c.b)("p",null,"For each property:"),Object(c.b)("ul",null,Object(c.b)("li",{parentName:"ul"},"key is an ",Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"/docs/reference/opspec/op-directory/op/identifier"}),"identifier [string]")," defining the name of the output."),Object(c.b)("li",{parentName:"ul"},"value is a ",Object(c.b)("a",Object(n.a)({parentName:"li"},{href:"/docs/reference/opspec/op-directory/op/parameter/index"}),"parameter [object]")," defining the output.")),Object(c.b)("h3",{id:"opspec"},"opspec"),Object(c.b)("p",null,"A ",Object(c.b)("a",Object(n.a)({parentName:"p"},{href:"https://semver.org/spec/v2.0.0.html"}),"semver v2.0.0 [string]")," which defines the version of opspec used to define the operation."),Object(c.b)("h3",{id:"run"},"run"),Object(c.b)("p",null,"A ",Object(c.b)("a",Object(n.a)({parentName:"p"},{href:"/docs/reference/opspec/op-directory/op/call/index"}),"call [object]")," defining the ops call graph; i.e. what gets run by the operation. "),Object(c.b)("h3",{id:"version"},"version"),Object(c.b)("p",null,"A ",Object(c.b)("a",Object(n.a)({parentName:"p"},{href:"https://semver.org/spec/v2.0.0.html"}),"semver v2.0.0 [string]")," which defines the version of the operation. "),Object(c.b)("blockquote",null,Object(c.b)("p",{parentName:"blockquote"},"If the op is published remotely, this MUST correspond to a ","[git]"," tag on the containing repo.")))}u.isMDXComponent=!0},193:function(e,t,r){"use strict";r.d(t,"a",(function(){return u})),r.d(t,"b",(function(){return f}));var n=r(0),a=r.n(n);function c(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function o(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function i(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?o(Object(r),!0).forEach((function(t){c(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):o(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function p(e,t){if(null==e)return{};var r,n,a=function(e,t){if(null==e)return{};var r,n,a={},c=Object.keys(e);for(n=0;n<c.length;n++)r=c[n],t.indexOf(r)>=0||(a[r]=e[r]);return a}(e,t);if(Object.getOwnPropertySymbols){var c=Object.getOwnPropertySymbols(e);for(n=0;n<c.length;n++)r=c[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(a[r]=e[r])}return a}var b=a.a.createContext({}),l=function(e){var t=a.a.useContext(b),r=t;return e&&(r="function"==typeof e?e(t):i({},t,{},e)),r},u=function(e){var t=l(e.components);return(a.a.createElement(b.Provider,{value:t},e.children))},s="mdxType",d={inlineCode:"code",wrapper:function(e){var t=e.children;return a.a.createElement(a.a.Fragment,{},t)}},O=Object(n.forwardRef)((function(e,t){var r=e.components,n=e.mdxType,c=e.originalType,o=e.parentName,b=p(e,["components","mdxType","originalType","parentName"]),u=l(r),s=n,O=u["".concat(o,".").concat(s)]||u[s]||d[s]||c;return r?a.a.createElement(O,i({ref:t},b,{components:r})):a.a.createElement(O,i({ref:t},b))}));function f(e,t){var r=arguments,n=t&&t.mdxType;if("string"==typeof e||n){var c=r.length,o=new Array(c);o[0]=O;var i={};for(var p in t)hasOwnProperty.call(t,p)&&(i[p]=t[p]);i.originalType=e,i[s]="string"==typeof e?e:n,o[1]=i;for(var b=2;b<c;b++)o[b]=r[b];return a.a.createElement.apply(null,o)}return a.a.createElement.apply(null,r)}O.displayName="MDXCreateElement"}}]);