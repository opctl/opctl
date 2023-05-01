"use strict";(self.webpackChunkopctl=self.webpackChunkopctl||[]).push([[7119],{3905:function(e,n,t){t.d(n,{Zo:function(){return u},kt:function(){return f}});var r=t(7294);function o(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function a(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);n&&(r=r.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,r)}return t}function i(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?a(Object(t),!0).forEach((function(n){o(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):a(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function l(e,n){if(null==e)return{};var t,r,o=function(e,n){if(null==e)return{};var t,r,o={},a=Object.keys(e);for(r=0;r<a.length;r++)t=a[r],n.indexOf(t)>=0||(o[t]=e[t]);return o}(e,n);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(r=0;r<a.length;r++)t=a[r],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(o[t]=e[t])}return o}var p=r.createContext({}),c=function(e){var n=r.useContext(p),t=n;return e&&(t="function"==typeof e?e(n):i(i({},n),e)),t},u=function(e){var n=c(e.components);return r.createElement(p.Provider,{value:n},e.children)},d={inlineCode:"code",wrapper:function(e){var n=e.children;return r.createElement(r.Fragment,{},n)}},s=r.forwardRef((function(e,n){var t=e.components,o=e.mdxType,a=e.originalType,p=e.parentName,u=l(e,["components","mdxType","originalType","parentName"]),s=c(t),f=o,m=s["".concat(p,".").concat(f)]||s[f]||d[f]||a;return t?r.createElement(m,i(i({ref:n},u),{},{components:t})):r.createElement(m,i({ref:n},u))}));function f(e,n){var t=arguments,o=n&&n.mdxType;if("string"==typeof e||o){var a=t.length,i=new Array(a);i[0]=s;var l={};for(var p in n)hasOwnProperty.call(n,p)&&(l[p]=n[p]);l.originalType=e,l.mdxType="string"==typeof e?e:o,i[1]=l;for(var c=2;c<a;c++)i[c]=t[c];return r.createElement.apply(null,i)}return r.createElement.apply(null,t)}s.displayName="MDXCreateElement"},7285:function(e,n,t){t.r(n),t.d(n,{assets:function(){return u},contentTitle:function(){return p},default:function(){return f},frontMatter:function(){return l},metadata:function(){return c},toc:function(){return d}});var r=t(3117),o=t(102),a=(t(7294),t(3905)),i=["components"],l={title:"How do I make parts of my op run conditionally?"},p=void 0,c={unversionedId:"training/flow/how-do-i-make-parts-of-my-op-run-conditionally",id:"training/flow/how-do-i-make-parts-of-my-op-run-conditionally",title:"How do I make parts of my op run conditionally?",description:"TLDR;",source:"@site/docs/training/flow/how-do-i-make-parts-of-my-op-run-conditionally.md",sourceDirName:"training/flow",slug:"/training/flow/how-do-i-make-parts-of-my-op-run-conditionally",permalink:"/docs/training/flow/how-do-i-make-parts-of-my-op-run-conditionally",draft:!1,editUrl:"https://github.com/opctl/opctl/edit/main/website/docs/training/flow/how-do-i-make-parts-of-my-op-run-conditionally.md",tags:[],version:"current",lastUpdatedBy:"=",lastUpdatedAt:1682904716,formattedLastUpdatedAt:"May 1, 2023",frontMatter:{title:"How do I make parts of my op run conditionally?"},sidebar:"docs",previous:{title:"How do I run a container?",permalink:"/docs/training/containers/how-do-i-run-a-container"},next:{title:"How do I make parts of my op run in a loop?",permalink:"/docs/training/flow/how-do-i-make-parts-of-my-op-run-in-a-loop"}},u={},d=[{value:"TLDR;",id:"tldr",level:2},{value:"Example",id:"example",level:2}],s={toc:d};function f(e){var n=e.components,t=(0,o.Z)(e,i);return(0,a.kt)("wrapper",(0,r.Z)({},s,t,{components:n,mdxType:"MDXLayout"}),(0,a.kt)("h2",{id:"tldr"},"TLDR;"),(0,a.kt)("p",null,"Opctl supports using ",(0,a.kt)("a",{parentName:"p",href:"/docs/reference/opspec/op-directory/op/call/#if"},"if")," statements and ",(0,a.kt)("a",{parentName:"p",href:"/docs/reference/opspec/op-directory/op/call/predicate"},"predicates")," to make parts of your op run conditionally."),(0,a.kt)("h2",{id:"example"},"Example"),(0,a.kt)("ol",null,(0,a.kt)("li",{parentName:"ol"},"Start this op: ",(0,a.kt)("pre",{parentName:"li"},(0,a.kt)("code",{parentName:"pre",className:"language-yaml"},"name: conditionalContainer\ninputs:\n  shouldRunContainer:\n    description: whether to run the container or not\n    boolean: {}\nrun:\n  if:\n    - eq: [true, $(shouldRunContainer)]\n  container:\n    cmd: [echo, 'hello!']\n    image: { ref: alpine }\n"))),(0,a.kt)("li",{parentName:"ol"},"When prompted, enter ",(0,a.kt)("inlineCode",{parentName:"li"},"true")," or ",(0,a.kt)("inlineCode",{parentName:"li"},"false")),(0,a.kt)("li",{parentName:"ol"},"Observe you only see the container run and ",(0,a.kt)("inlineCode",{parentName:"li"},"hello!")," logged when you enter ",(0,a.kt)("inlineCode",{parentName:"li"},"true"),".")))}f.isMDXComponent=!0}}]);