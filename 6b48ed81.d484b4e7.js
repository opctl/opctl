(window.webpackJsonp=window.webpackJsonp||[]).push([[42],{144:function(e,r,t){"use strict";t.r(r),t.d(r,"frontMatter",(function(){return c})),t.d(r,"metadata",(function(){return i})),t.d(r,"rightToc",(function(){return p})),t.d(r,"default",(function(){return f}));var n=t(1),a=t(6),o=(t(0),t(193)),c={title:"Variable Reference [string]"},i={id:"reference/opspec/op-directory/op/variable-reference",title:"Variable Reference [string]",description:"A string referencing the location of/for a value in the form of `$(REFERENCE)` where `REFERENCE` MUST start with an [identifier [string]](identifier.md) and MAY end with one or more:",source:"@site/docs/reference/opspec/op-directory/op/variable-reference.md",permalink:"/docs/reference/opspec/op-directory/op/variable-reference",editUrl:"https://github.com/opctl/opctl/edit/master/docs/docs/reference/opspec/op-directory/op/variable-reference.md",lastUpdatedBy:"Chris Dostert",lastUpdatedAt:1587672399,sidebar:"docs",previous:{title:"Markdown [string]",permalink:"/docs/reference/opspec/op-directory/op/markdown"},next:{title:"Array",permalink:"/docs/reference/opspec/types/array"}},p=[],l={rightToc:p},b="wrapper";function f(e){var r=e.components,t=Object(a.a)(e,["components"]);return Object(o.b)(b,Object(n.a)({},l,t,{components:r,mdxType:"MDXLayout"}),Object(o.b)("p",null,"A string referencing the location of/for a value in the form of ",Object(o.b)("inlineCode",{parentName:"p"},"$(REFERENCE)")," where ",Object(o.b)("inlineCode",{parentName:"p"},"REFERENCE")," MUST start with an ",Object(o.b)("a",Object(n.a)({parentName:"p"},{href:"/docs/reference/opspec/op-directory/op/identifier"}),"identifier [string]")," and MAY end with one or more:"),Object(o.b)("ul",null,Object(o.b)("li",{parentName:"ul"},Object(o.b)("a",Object(n.a)({parentName:"li"},{href:"../../../types/array.md#item-referencing"}),"array item references")),Object(o.b)("li",{parentName:"ul"},Object(o.b)("a",Object(n.a)({parentName:"li"},{href:"../../../types/object.md#property-referencing"}),"object property references")),Object(o.b)("li",{parentName:"ul"},Object(o.b)("a",Object(n.a)({parentName:"li"},{href:"../../../types/dir.md#entry-referencing"}),"dir entry references"))),Object(o.b)("p",null,"References can be used to either define or access values in the current scope. "),Object(o.b)("p",null,"When an op starts, it's initial scope includes:"),Object(o.b)("ul",null,Object(o.b)("li",{parentName:"ul"},Object(o.b)("inlineCode",{parentName:"li"},"/")," with a value of the current op directory i.e. the current op's ",Object(o.b)("inlineCode",{parentName:"li"},"op.yml")," can be accessed via ",Object(o.b)("inlineCode",{parentName:"li"},"$(/op.yml)"),"."),Object(o.b)("li",{parentName:"ul"},"any defined inputs")),Object(o.b)("blockquote",null,Object(o.b)("p",{parentName:"blockquote"},"note: variable references can be escaped by prefixing the ","[would be]"," variable reference with ",Object(o.b)("inlineCode",{parentName:"p"},"\\")," i.e. ",Object(o.b)("inlineCode",{parentName:"p"},"\\\\$(wouldBeVariableReference)")," would not be treated as a variable reference. ")))}f.isMDXComponent=!0},193:function(e,r,t){"use strict";t.d(r,"a",(function(){return f})),t.d(r,"b",(function(){return m}));var n=t(0),a=t.n(n);function o(e,r,t){return r in e?Object.defineProperty(e,r,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[r]=t,e}function c(e,r){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);r&&(n=n.filter((function(r){return Object.getOwnPropertyDescriptor(e,r).enumerable}))),t.push.apply(t,n)}return t}function i(e){for(var r=1;r<arguments.length;r++){var t=null!=arguments[r]?arguments[r]:{};r%2?c(Object(t),!0).forEach((function(r){o(e,r,t[r])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):c(Object(t)).forEach((function(r){Object.defineProperty(e,r,Object.getOwnPropertyDescriptor(t,r))}))}return e}function p(e,r){if(null==e)return{};var t,n,a=function(e,r){if(null==e)return{};var t,n,a={},o=Object.keys(e);for(n=0;n<o.length;n++)t=o[n],r.indexOf(t)>=0||(a[t]=e[t]);return a}(e,r);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(n=0;n<o.length;n++)t=o[n],r.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(a[t]=e[t])}return a}var l=a.a.createContext({}),b=function(e){var r=a.a.useContext(l),t=r;return e&&(t="function"==typeof e?e(r):i({},r,{},e)),t},f=function(e){var r=b(e.components);return(a.a.createElement(l.Provider,{value:r},e.children))},s="mdxType",u={inlineCode:"code",wrapper:function(e){var r=e.children;return a.a.createElement(a.a.Fragment,{},r)}},d=Object(n.forwardRef)((function(e,r){var t=e.components,n=e.mdxType,o=e.originalType,c=e.parentName,l=p(e,["components","mdxType","originalType","parentName"]),f=b(t),s=n,d=f["".concat(c,".").concat(s)]||f[s]||u[s]||o;return t?a.a.createElement(d,i({ref:r},l,{components:t})):a.a.createElement(d,i({ref:r},l))}));function m(e,r){var t=arguments,n=r&&r.mdxType;if("string"==typeof e||n){var o=t.length,c=new Array(o);c[0]=d;var i={};for(var p in r)hasOwnProperty.call(r,p)&&(i[p]=r[p]);i.originalType=e,i[s]="string"==typeof e?e:n,c[1]=i;for(var l=2;l<o;l++)c[l]=t[l];return a.a.createElement.apply(null,c)}return a.a.createElement.apply(null,t)}d.displayName="MDXCreateElement"}}]);