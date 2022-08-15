"use strict";(self.webpackChunkopctl=self.webpackChunkopctl||[]).push([[7179],{3905:function(e,n,t){t.d(n,{Zo:function(){return u},kt:function(){return m}});var r=t(7294);function o(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function a(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);n&&(r=r.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,r)}return t}function i(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?a(Object(t),!0).forEach((function(n){o(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):a(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function c(e,n){if(null==e)return{};var t,r,o=function(e,n){if(null==e)return{};var t,r,o={},a=Object.keys(e);for(r=0;r<a.length;r++)t=a[r],n.indexOf(t)>=0||(o[t]=e[t]);return o}(e,n);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(r=0;r<a.length;r++)t=a[r],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(o[t]=e[t])}return o}var l=r.createContext({}),p=function(e){var n=r.useContext(l),t=n;return e&&(t="function"==typeof e?e(n):i(i({},n),e)),t},u=function(e){var n=p(e.components);return r.createElement(l.Provider,{value:n},e.children)},s={inlineCode:"code",wrapper:function(e){var n=e.children;return r.createElement(r.Fragment,{},n)}},d=r.forwardRef((function(e,n){var t=e.components,o=e.mdxType,a=e.originalType,l=e.parentName,u=c(e,["components","mdxType","originalType","parentName"]),d=p(t),m=o,f=d["".concat(l,".").concat(m)]||d[m]||s[m]||a;return t?r.createElement(f,i(i({ref:n},u),{},{components:t})):r.createElement(f,i({ref:n},u))}));function m(e,n){var t=arguments,o=n&&n.mdxType;if("string"==typeof e||o){var a=t.length,i=new Array(a);i[0]=d;var c={};for(var l in n)hasOwnProperty.call(n,l)&&(c[l]=n[l]);c.originalType=e,c.mdxType="string"==typeof e?e:o,i[1]=c;for(var p=2;p<a;p++)i[p]=t[p];return r.createElement.apply(null,i)}return r.createElement.apply(null,t)}d.displayName="MDXCreateElement"},1354:function(e,n,t){t.r(n),t.d(n,{assets:function(){return u},contentTitle:function(){return l},default:function(){return m},frontMatter:function(){return c},metadata:function(){return p},toc:function(){return s}});var r=t(3117),o=t(102),a=(t(7294),t(3905)),i=["components"],c={title:"How do I run a container?"},l=void 0,p={unversionedId:"training/containers/how-do-i-run-a-container",id:"training/containers/how-do-i-run-a-container",title:"How do I run a container?",description:"TLDR;",source:"@site/docs/training/containers/how-do-i-run-a-container.md",sourceDirName:"training/containers",slug:"/training/containers/how-do-i-run-a-container",permalink:"/docs/training/containers/how-do-i-run-a-container",draft:!1,editUrl:"https://github.com/opctl/opctl/edit/main/website/docs/training/containers/how-do-i-run-a-container.md",tags:[],version:"current",lastUpdatedBy:"=",lastUpdatedAt:1660606934,formattedLastUpdatedAt:"Aug 15, 2022",frontMatter:{title:"How do I run a container?"},sidebar:"docs",previous:{title:"How do I get opctl containers to communicate?",permalink:"/docs/training/containers/how-do-i-get-opctl-containers-to-communicate"},next:{title:"Hello World",permalink:"/docs/training/hello-world"}},u={},s=[{value:"TLDR;",id:"tldr",level:2},{value:"Example",id:"example",level:2}],d={toc:s};function m(e){var n=e.components,t=(0,o.Z)(e,i);return(0,a.kt)("wrapper",(0,r.Z)({},d,t,{components:n,mdxType:"MDXLayout"}),(0,a.kt)("h2",{id:"tldr"},"TLDR;"),(0,a.kt)("p",null,"Opctl supports running an ",(0,a.kt)("a",{parentName:"p",href:"https://opencontainers.org/"},"OCI")," image based container by defining a ",(0,a.kt)("a",{parentName:"p",href:"/docs/reference/opspec/op-directory/op/call/container/"},"container call"),"."),(0,a.kt)("blockquote",null,(0,a.kt)("p",{parentName:"blockquote"},"Note: a common place to obtain ",(0,a.kt)("a",{parentName:"p",href:"https://opencontainers.org/"},"OCI")," images is ",(0,a.kt)("a",{parentName:"p",href:"https://hub.docker.com/"},"Docker Hub"),".")),(0,a.kt)("h2",{id:"example"},"Example"),(0,a.kt)("ol",null,(0,a.kt)("li",{parentName:"ol"},"Start this op: ",(0,a.kt)("pre",{parentName:"li"},(0,a.kt)("code",{parentName:"pre",className:"language-yaml"},"name: runAContainer\nrun:\n  container:\n    cmd: [sh, -ce, 'echo hello!']\n    image: { ref: alpine }\n"))),(0,a.kt)("li",{parentName:"ol"},"Observe the container is run and ",(0,a.kt)("inlineCode",{parentName:"li"},"hello!")," logged.")))}m.isMDXComponent=!0}}]);