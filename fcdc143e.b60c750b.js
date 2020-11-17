(window.webpackJsonp=window.webpackJsonp||[]).push([[70],{172:function(e,t,n){"use strict";n.r(t),n.d(t,"frontMatter",(function(){return c})),n.d(t,"metadata",(function(){return i})),n.d(t,"rightToc",(function(){return l})),n.d(t,"default",(function(){return u}));var r=n(1),o=n(6),a=(n(0),n(176)),c={title:"How do I get opctl containers to communicate?"},i={id:"training/containers/how-do-i-get-opctl-containers-to-communicate",title:"How do I get opctl containers to communicate?",description:"## TLDR;",source:"@site/docs/training/containers/how-do-i-get-opctl-containers-to-communicate.md",permalink:"/docs/training/containers/how-do-i-get-opctl-containers-to-communicate",editUrl:"https://github.com/opctl/opctl/edit/master/website/docs/training/containers/how-do-i-get-opctl-containers-to-communicate.md",lastUpdatedBy:"Chris Dostert",lastUpdatedAt:1605583704,sidebar:"docs",previous:{title:"How do I communicate with an opctl container?",permalink:"/docs/training/containers/how-do-i-communicate-with-an-opctl-container"},next:{title:"How do I visualize an op?",permalink:"/docs/training/ui/how-do-i-visualize-an-op"}},l=[{value:"TLDR;",id:"tldr",children:[]},{value:"Example",id:"example",children:[]}],p={rightToc:l},s="wrapper";function u(e){var t=e.components,n=Object(o.a)(e,["components"]);return Object(a.b)(s,Object(r.a)({},p,n,{components:t,mdxType:"MDXLayout"}),Object(a.b)("h2",{id:"tldr"},"TLDR;"),Object(a.b)("p",null,"Opctl attaches all containers to a virtual overlay network.  "),Object(a.b)("p",null,"Adding a ",Object(a.b)("a",Object(r.a)({parentName:"p"},{href:"../../reference/opspec/op-directory/op/call/container/index#name"}),"name")," attribute to container(s) adds a corresponding network wide DNS A record which resolves to the assigned ip(s) of the container(s)."),Object(a.b)("p",null,"Whether containers are defined in the same op or not makes no difference, they can still reach each other."),Object(a.b)("p",null,"If multiple containers have the same ",Object(a.b)("a",Object(r.a)({parentName:"p"},{href:"../../reference/opspec/op-directory/op/call/container/index#name"}),"name")," requests will be load balanced across them."),Object(a.b)("h2",{id:"example"},"Example"),Object(a.b)("ol",null,Object(a.b)("li",{parentName:"ol"},Object(a.b)("p",{parentName:"li"},"Run this op:"),Object(a.b)("pre",{parentName:"li"},Object(a.b)("code",Object(r.a)({parentName:"pre"},{className:"language-yaml"}),"name: ping\nrun:\n  parallel:\n    - container:\n        image: { ref: alpine }\n        name: container1\n        cmd: [sleep, 1000000]\n    - container:\n        image: { ref: alpine }\n        # ping container1 by its name\n        cmd: [ping, container1]\n"))),Object(a.b)("li",{parentName:"ol"},Object(a.b)("p",{parentName:"li"},"Observe the second container succeeds in ",Object(a.b)("inlineCode",{parentName:"p"},"ping"),"ing ",Object(a.b)("inlineCode",{parentName:"p"},"container1"),". "))))}u.isMDXComponent=!0},176:function(e,t,n){"use strict";n.d(t,"a",(function(){return u})),n.d(t,"b",(function(){return f}));var r=n(0),o=n.n(r);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function c(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?c(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):c(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function l(e,t){if(null==e)return{};var n,r,o=function(e,t){if(null==e)return{};var n,r,o={},a=Object.keys(e);for(r=0;r<a.length;r++)n=a[r],t.indexOf(n)>=0||(o[n]=e[n]);return o}(e,t);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(r=0;r<a.length;r++)n=a[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(o[n]=e[n])}return o}var p=o.a.createContext({}),s=function(e){var t=o.a.useContext(p),n=t;return e&&(n="function"==typeof e?e(t):i({},t,{},e)),n},u=function(e){var t=s(e.components);return(o.a.createElement(p.Provider,{value:t},e.children))},m="mdxType",d={inlineCode:"code",wrapper:function(e){var t=e.children;return o.a.createElement(o.a.Fragment,{},t)}},b=Object(r.forwardRef)((function(e,t){var n=e.components,r=e.mdxType,a=e.originalType,c=e.parentName,p=l(e,["components","mdxType","originalType","parentName"]),u=s(n),m=r,b=u["".concat(c,".").concat(m)]||u[m]||d[m]||a;return n?o.a.createElement(b,i({ref:t},p,{components:n})):o.a.createElement(b,i({ref:t},p))}));function f(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var a=n.length,c=new Array(a);c[0]=b;var i={};for(var l in t)hasOwnProperty.call(t,l)&&(i[l]=t[l]);i.originalType=e,i[m]="string"==typeof e?e:r,c[1]=i;for(var p=2;p<a;p++)c[p]=n[p];return o.a.createElement.apply(null,c)}return o.a.createElement.apply(null,n)}b.displayName="MDXCreateElement"}}]);