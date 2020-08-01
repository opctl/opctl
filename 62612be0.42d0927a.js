(window.webpackJsonp=window.webpackJsonp||[]).push([[22],{124:function(e,t,n){"use strict";n.r(t),n.d(t,"frontMatter",(function(){return c})),n.d(t,"metadata",(function(){return i})),n.d(t,"rightToc",(function(){return p})),n.d(t,"default",(function(){return s}));var r=n(1),o=n(6),a=(n(0),n(161)),c={title:"How do I communicate with an opctl container?"},i={id:"training/containers/how-do-i-communicate-with-an-opctl-container",title:"How do I communicate with an opctl container?",description:"## TLDR;",source:"@site/docs/training/containers/how-do-i-communicate-with-an-opctl-container.md",permalink:"/docs/training/containers/how-do-i-communicate-with-an-opctl-container",editUrl:"https://github.com/opctl/opctl/edit/master/website/docs/training/containers/how-do-i-communicate-with-an-opctl-container.md",lastUpdatedBy:"Chris Dostert",lastUpdatedAt:1596248938,sidebar:"docs",previous:{title:"Run a react app",permalink:"/docs/run-a-react-app"},next:{title:"How do I get opctl containers to communicate?",permalink:"/docs/training/containers/how-do-i-get-opctl-containers-to-communicate"}},p=[{value:"TLDR;",id:"tldr",children:[]},{value:"Example",id:"example",children:[]}],l={rightToc:p},u="wrapper";function s(e){var t=e.components,n=Object(o.a)(e,["components"]);return Object(a.b)(u,Object(r.a)({},l,n,{components:t,mdxType:"MDXLayout"}),Object(a.b)("h2",{id:"tldr"},"TLDR;"),Object(a.b)("p",null,"Adding a ",Object(a.b)("a",Object(r.a)({parentName:"p"},{href:"../../reference/opspec/op-directory/op/call/container/index#ports"}),"ports")," attribute to a container binds container ports to the opctl host."),Object(a.b)("h2",{id:"example"},"Example"),Object(a.b)("ol",null,Object(a.b)("li",{parentName:"ol"},"Start this op: ",Object(a.b)("pre",{parentName:"li"},Object(a.b)("code",Object(r.a)({parentName:"pre"},{className:"language-yaml"}),"name: curl\nrun:\n  container:\n    image: { ref: nginx:alpine }\n    ports:\n      # bind container port 80 to host port 8080\n      80: 8080\n"))),Object(a.b)("li",{parentName:"ol"},"On the opctl host, open a web browser to ",Object(a.b)("a",Object(r.a)({parentName:"li"},{href:"localhost:8080"}),"localhost:8080"),"."),Object(a.b)("li",{parentName:"ol"},"Observe the nginx containers default page is returned. ")))}s.isMDXComponent=!0},161:function(e,t,n){"use strict";n.d(t,"a",(function(){return s})),n.d(t,"b",(function(){return f}));var r=n(0),o=n.n(r);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function c(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?c(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):c(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function p(e,t){if(null==e)return{};var n,r,o=function(e,t){if(null==e)return{};var n,r,o={},a=Object.keys(e);for(r=0;r<a.length;r++)n=a[r],t.indexOf(n)>=0||(o[n]=e[n]);return o}(e,t);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(r=0;r<a.length;r++)n=a[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(o[n]=e[n])}return o}var l=o.a.createContext({}),u=function(e){var t=o.a.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):i({},t,{},e)),n},s=function(e){var t=u(e.components);return(o.a.createElement(l.Provider,{value:t},e.children))},m="mdxType",d={inlineCode:"code",wrapper:function(e){var t=e.children;return o.a.createElement(o.a.Fragment,{},t)}},b=Object(r.forwardRef)((function(e,t){var n=e.components,r=e.mdxType,a=e.originalType,c=e.parentName,l=p(e,["components","mdxType","originalType","parentName"]),s=u(n),m=r,b=s["".concat(c,".").concat(m)]||s[m]||d[m]||a;return n?o.a.createElement(b,i({ref:t},l,{components:n})):o.a.createElement(b,i({ref:t},l))}));function f(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var a=n.length,c=new Array(a);c[0]=b;var i={};for(var p in t)hasOwnProperty.call(t,p)&&(i[p]=t[p]);i.originalType=e,i[m]="string"==typeof e?e:r,c[1]=i;for(var l=2;l<a;l++)c[l]=n[l];return o.a.createElement.apply(null,c)}return o.a.createElement.apply(null,n)}b.displayName="MDXCreateElement"}}]);