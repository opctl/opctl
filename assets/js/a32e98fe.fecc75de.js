"use strict";(self.webpackChunkopctl=self.webpackChunkopctl||[]).push([[1738],{3905:function(e,t,n){n.d(t,{Zo:function(){return u},kt:function(){return f}});var r=n(7294);function o(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function c(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function l(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?c(Object(n),!0).forEach((function(t){o(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):c(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function i(e,t){if(null==e)return{};var n,r,o=function(e,t){if(null==e)return{};var n,r,o={},c=Object.keys(e);for(r=0;r<c.length;r++)n=c[r],t.indexOf(n)>=0||(o[n]=e[n]);return o}(e,t);if(Object.getOwnPropertySymbols){var c=Object.getOwnPropertySymbols(e);for(r=0;r<c.length;r++)n=c[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(o[n]=e[n])}return o}var a=r.createContext({}),p=function(e){var t=r.useContext(a),n=t;return e&&(n="function"==typeof e?e(t):l(l({},t),e)),n},u=function(e){var t=p(e.components);return r.createElement(a.Provider,{value:t},e.children)},d={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},s=r.forwardRef((function(e,t){var n=e.components,o=e.mdxType,c=e.originalType,a=e.parentName,u=i(e,["components","mdxType","originalType","parentName"]),s=p(n),f=o,m=s["".concat(a,".").concat(f)]||s[f]||d[f]||c;return n?r.createElement(m,l(l({ref:t},u),{},{components:n})):r.createElement(m,l({ref:t},u))}));function f(e,t){var n=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var c=n.length,l=new Array(c);l[0]=s;var i={};for(var a in t)hasOwnProperty.call(t,a)&&(i[a]=t[a]);i.originalType=e,i.mdxType="string"==typeof e?e:o,l[1]=i;for(var p=2;p<c;p++)l[p]=n[p];return r.createElement.apply(null,l)}return r.createElement.apply(null,n)}s.displayName="MDXCreateElement"},3031:function(e,t,n){n.r(t),n.d(t,{assets:function(){return u},contentTitle:function(){return a},default:function(){return f},frontMatter:function(){return i},metadata:function(){return p},toc:function(){return d}});var r=n(3117),o=n(102),c=(n(7294),n(3905)),l=["components"],i={sidebar_label:"create",title:"opctl node create"},a=void 0,p={unversionedId:"reference/cli/node/create",id:"reference/cli/node/create",title:"opctl node create",description:"Start an in-process node which inherits current",source:"@site/docs/reference/cli/node/create.md",sourceDirName:"reference/cli/node",slug:"/reference/cli/node/create",permalink:"/docs/reference/cli/node/create",draft:!1,editUrl:"https://github.com/opctl/opctl/edit/main/website/docs/reference/cli/node/create.md",tags:[],version:"current",lastUpdatedBy:"=",lastUpdatedAt:1639604483,formattedLastUpdatedAt:"Dec 15, 2021",frontMatter:{sidebar_label:"create",title:"opctl node create"},sidebar:"docs",previous:{title:"node",permalink:"/docs/reference/cli/node/"},next:{title:"delete",permalink:"/docs/reference/cli/node/delete"}},u={},d=[{value:"Global Options",id:"global-options",level:2},{value:"Notes",id:"notes",level:2},{value:"lockfile",id:"lockfile",level:3},{value:"concurrency",id:"concurrency",level:3}],s={toc:d};function f(e){var t=e.components,n=(0,o.Z)(e,l);return(0,c.kt)("wrapper",(0,r.Z)({},s,n,{components:t,mdxType:"MDXLayout"}),(0,c.kt)("pre",null,(0,c.kt)("code",{parentName:"pre",className:"language-sh"},"opctl node create\n")),(0,c.kt)("p",null,"Start an in-process node which inherits current\nstderr/stdout/stdin/PGId (process group id) and blocks until interrupted or deleted."),(0,c.kt)("blockquote",null,(0,c.kt)("p",{parentName:"blockquote"},"There can be only one node running at a time on a given machine.")),(0,c.kt)("h2",{id:"global-options"},"Global Options"),(0,c.kt)("p",null,"see ",(0,c.kt)("a",{parentName:"p",href:"/docs/reference/cli/global-options"},"global options")),(0,c.kt)("h2",{id:"notes"},"Notes"),(0,c.kt)("h3",{id:"lockfile"},"lockfile"),(0,c.kt)("p",null,"Upon creation, nodes populate a lockfile at ",(0,c.kt)("inlineCode",{parentName:"p"},"DATA_DIR/pid.lock"),"\ncontaining their PId (process id)."),(0,c.kt)("h3",{id:"concurrency"},"concurrency"),(0,c.kt)("p",null,"Prior to node creation, if a lockfile exists, the existing lock holder\nwill be liveness tested."),(0,c.kt)("p",null,"Only in the event the existing lock holder is dead will creation of a\nnew node occur."))}f.isMDXComponent=!0}}]);