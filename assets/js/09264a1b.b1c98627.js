"use strict";(self.webpackChunkopctl=self.webpackChunkopctl||[]).push([[2615],{3905:function(e,t,r){r.d(t,{Zo:function(){return s},kt:function(){return m}});var n=r(7294);function o(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function l(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function a(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?l(Object(r),!0).forEach((function(t){o(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):l(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function i(e,t){if(null==e)return{};var r,n,o=function(e,t){if(null==e)return{};var r,n,o={},l=Object.keys(e);for(n=0;n<l.length;n++)r=l[n],t.indexOf(r)>=0||(o[r]=e[r]);return o}(e,t);if(Object.getOwnPropertySymbols){var l=Object.getOwnPropertySymbols(e);for(n=0;n<l.length;n++)r=l[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(o[r]=e[r])}return o}var p=n.createContext({}),c=function(e){var t=n.useContext(p),r=t;return e&&(r="function"==typeof e?e(t):a(a({},t),e)),r},s=function(e){var t=c(e.components);return n.createElement(p.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},d=n.forwardRef((function(e,t){var r=e.components,o=e.mdxType,l=e.originalType,p=e.parentName,s=i(e,["components","mdxType","originalType","parentName"]),d=c(r),m=o,f=d["".concat(p,".").concat(m)]||d[m]||u[m]||l;return r?n.createElement(f,a(a({ref:t},s),{},{components:r})):n.createElement(f,a({ref:t},s))}));function m(e,t){var r=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var l=r.length,a=new Array(l);a[0]=d;var i={};for(var p in t)hasOwnProperty.call(t,p)&&(i[p]=t[p]);i.originalType=e,i.mdxType="string"==typeof e?e:o,a[1]=i;for(var c=2;c<l;c++)a[c]=r[c];return n.createElement.apply(null,a)}return n.createElement.apply(null,r)}d.displayName="MDXCreateElement"},3597:function(e,t,r){r.r(t),r.d(t,{assets:function(){return s},contentTitle:function(){return p},default:function(){return m},frontMatter:function(){return i},metadata:function(){return c},toc:function(){return u}});var n=r(3117),o=r(102),l=(r(7294),r(3905)),a=["components"],i={sidebar_label:"ls",title:"opctl ls"},p=void 0,c={unversionedId:"reference/cli/ls",id:"reference/cli/ls",title:"opctl ls",description:"List ops in a local or remote directory.",source:"@site/docs/reference/cli/ls.md",sourceDirName:"reference/cli",slug:"/reference/cli/ls",permalink:"/docs/reference/cli/ls",draft:!1,editUrl:"https://github.com/opctl/opctl/edit/main/website/docs/reference/cli/ls.md",tags:[],version:"current",lastUpdatedBy:"Chris Dostert",lastUpdatedAt:1604531186,formattedLastUpdatedAt:"Nov 4, 2020",frontMatter:{sidebar_label:"ls",title:"opctl ls"},sidebar:"docs",previous:{title:"global-options",permalink:"/docs/reference/cli/global-options"},next:{title:"node",permalink:"/docs/reference/cli/node/"}},s={},u=[{value:"Arguments",id:"arguments",level:3},{value:"<code>DIR_REF</code> <em>default: <code>.opspec</code></em>",id:"dir_ref-default-opspec",level:4},{value:"Global Options",id:"global-options",level:2},{value:"Examples",id:"examples",level:3},{value:"<code>.opspec</code> dir",id:"opspec-dir",level:4},{value:"remote dir",id:"remote-dir",level:4}],d={toc:u};function m(e){var t=e.components,r=(0,o.Z)(e,a);return(0,l.kt)("wrapper",(0,n.Z)({},d,r,{components:t,mdxType:"MDXLayout"}),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-sh"},"opctl ls [DIR_REF=.opspec]\n")),(0,l.kt)("p",null,"List ops in a local or remote directory."),(0,l.kt)("h3",{id:"arguments"},"Arguments"),(0,l.kt)("h4",{id:"dir_ref-default-opspec"},(0,l.kt)("inlineCode",{parentName:"h4"},"DIR_REF")," ",(0,l.kt)("em",{parentName:"h4"},"default: ",(0,l.kt)("inlineCode",{parentName:"em"},".opspec"))),(0,l.kt)("p",null,"Reference to dir ops will be listed from (either ",(0,l.kt)("inlineCode",{parentName:"p"},"relative/path"),", ",(0,l.kt)("inlineCode",{parentName:"p"},"/absolute/path"),", ",(0,l.kt)("inlineCode",{parentName:"p"},"host/path/repo#tag"),", or ",(0,l.kt)("inlineCode",{parentName:"p"},"host/path/repo#tag/path"),")"),(0,l.kt)("h2",{id:"global-options"},"Global Options"),(0,l.kt)("p",null,"see ",(0,l.kt)("a",{parentName:"p",href:"/docs/reference/cli/global-options"},"global options")),(0,l.kt)("h3",{id:"examples"},"Examples"),(0,l.kt)("h4",{id:"opspec-dir"},(0,l.kt)("inlineCode",{parentName:"h4"},".opspec")," dir"),(0,l.kt)("p",null,"lists ops from ",(0,l.kt)("inlineCode",{parentName:"p"},"./.opspec")),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-sh"},"opctl ls\n")),(0,l.kt)("h4",{id:"remote-dir"},"remote dir"),(0,l.kt)("p",null,"lists ops from ",(0,l.kt)("a",{parentName:"p",href:"https://github.com/opctl/opctl/tree/0.1.24"},"github.com/opctl/opctl#0.1.24")),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-sh"},"opctl ls github.com/opctl/opctl#0.1.24/\n")))}m.isMDXComponent=!0}}]);