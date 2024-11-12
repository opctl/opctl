"use strict";(self.webpackChunkopctl=self.webpackChunkopctl||[]).push([[1682],{3905:function(e,n,t){t.d(n,{Zo:function(){return u},kt:function(){return f}});var r=t(7294);function o(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function i(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);n&&(r=r.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,r)}return t}function a(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?i(Object(t),!0).forEach((function(n){o(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):i(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function l(e,n){if(null==e)return{};var t,r,o=function(e,n){if(null==e)return{};var t,r,o={},i=Object.keys(e);for(r=0;r<i.length;r++)t=i[r],n.indexOf(t)>=0||(o[t]=e[t]);return o}(e,n);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)t=i[r],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(o[t]=e[t])}return o}var c=r.createContext({}),p=function(e){var n=r.useContext(c),t=n;return e&&(t="function"==typeof e?e(n):a(a({},n),e)),t},u=function(e){var n=p(e.components);return r.createElement(c.Provider,{value:n},e.children)},d={inlineCode:"code",wrapper:function(e){var n=e.children;return r.createElement(r.Fragment,{},n)}},s=r.forwardRef((function(e,n){var t=e.components,o=e.mdxType,i=e.originalType,c=e.parentName,u=l(e,["components","mdxType","originalType","parentName"]),s=p(t),f=o,m=s["".concat(c,".").concat(f)]||s[f]||d[f]||i;return t?r.createElement(m,a(a({ref:n},u),{},{components:t})):r.createElement(m,a({ref:n},u))}));function f(e,n){var t=arguments,o=n&&n.mdxType;if("string"==typeof e||o){var i=t.length,a=new Array(i);a[0]=s;var l={};for(var c in n)hasOwnProperty.call(n,c)&&(l[c]=n[c]);l.originalType=e,l.mdxType="string"==typeof e?e:o,a[1]=l;for(var p=2;p<i;p++)a[p]=t[p];return r.createElement.apply(null,a)}return r.createElement.apply(null,t)}s.displayName="MDXCreateElement"},3854:function(e,n,t){t.r(n),t.d(n,{assets:function(){return u},contentTitle:function(){return c},default:function(){return f},frontMatter:function(){return l},metadata:function(){return p},toc:function(){return d}});var r=t(3117),o=t(102),i=(t(7294),t(3905)),a=["components"],l={title:"Conditional execution"},c=void 0,p={unversionedId:"training/flow/conditional-execution",id:"training/flow/conditional-execution",title:"Conditional execution",description:"TLDR;",source:"@site/docs/training/flow/conditional-execution.md",sourceDirName:"training/flow",slug:"/training/flow/conditional-execution",permalink:"/docs/training/flow/conditional-execution",draft:!1,editUrl:"https://github.com/opctl/opctl/edit/main/website/docs/training/flow/conditional-execution.md",tags:[],version:"current",lastUpdatedBy:"=",lastUpdatedAt:1682978842,formattedLastUpdatedAt:"May 1, 2023",frontMatter:{title:"Conditional execution"},sidebar:"docs",previous:{title:"Container networking",permalink:"/docs/training/containers/container-networking"},next:{title:"Serial and parallel execution",permalink:"/docs/training/flow/serial-and-parallel-execution"}},u={},d=[{value:"TLDR;",id:"tldr",level:2},{value:"Example",id:"example",level:2}],s={toc:d};function f(e){var n=e.components,t=(0,o.Z)(e,a);return(0,i.kt)("wrapper",(0,r.Z)({},s,t,{components:n,mdxType:"MDXLayout"}),(0,i.kt)("h2",{id:"tldr"},"TLDR;"),(0,i.kt)("p",null,"Opctl supports using ",(0,i.kt)("a",{parentName:"p",href:"/docs/reference/opspec/op-directory/op/call/#if"},"if")," statements and ",(0,i.kt)("a",{parentName:"p",href:"/docs/reference/opspec/op-directory/op/call/predicate"},"predicates")," to make parts of your op run conditionally."),(0,i.kt)("h2",{id:"example"},"Example"),(0,i.kt)("ol",null,(0,i.kt)("li",{parentName:"ol"},"Start this op: ",(0,i.kt)("pre",{parentName:"li"},(0,i.kt)("code",{parentName:"pre",className:"language-yaml"},"name: conditionalExecution\ninputs:\n  shouldRunContainer:\n    description: whether to run the container or not\n    boolean: {}\nrun:\n  if:\n    - eq: [true, $(shouldRunContainer)]\n  container:\n    cmd: [echo, 'hello!']\n    image: { ref: alpine }\n"))),(0,i.kt)("li",{parentName:"ol"},"When prompted, enter ",(0,i.kt)("inlineCode",{parentName:"li"},"true")," or ",(0,i.kt)("inlineCode",{parentName:"li"},"false")),(0,i.kt)("li",{parentName:"ol"},"Observe you only see the container run and ",(0,i.kt)("inlineCode",{parentName:"li"},"hello!")," logged when you enter ",(0,i.kt)("inlineCode",{parentName:"li"},"true"),".")))}f.isMDXComponent=!0}}]);