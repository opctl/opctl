(window.webpackJsonp=window.webpackJsonp||[]).push([[21],{123:function(e,t,n){"use strict";n.r(t),n.d(t,"frontMatter",(function(){return l})),n.d(t,"metadata",(function(){return i})),n.d(t,"rightToc",(function(){return c})),n.d(t,"default",(function(){return u}));var r=n(1),a=n(6),o=(n(0),n(176)),l={title:"Bare Metal",sidebar_label:"Bare Metal"},i={id:"setup/bare-metal",title:"Bare Metal",description:"## Installation",source:"@site/docs/setup/bare-metal.md",permalink:"/docs/setup/bare-metal",editUrl:"https://github.com/opctl/opctl/edit/master/website/docs/setup/bare-metal.md",lastUpdatedBy:"Chris Dostert",lastUpdatedAt:1605583704,sidebar_label:"Bare Metal",sidebar:"docs",previous:{title:"Introduction",permalink:"/docs/introduction"},next:{title:"Azure Pipelines",permalink:"/docs/setup/azure-pipelines"}},c=[{value:"Installation",id:"installation",children:[{value:"Prerequisites",id:"prerequisites",children:[]},{value:"OSX",id:"osx",children:[]},{value:"Linux",id:"linux",children:[]},{value:"Windows",id:"windows",children:[]}]},{value:"Updating",id:"updating",children:[]},{value:"IDE Plugins",id:"ide-plugins",children:[{value:"VSCode",id:"vscode",children:[]}]}],s={rightToc:c},p="wrapper";function u(e){var t=e.components,n=Object(a.a)(e,["components"]);return Object(o.b)(p,Object(r.a)({},s,n,{components:t,mdxType:"MDXLayout"}),Object(o.b)("h2",{id:"installation"},"Installation"),Object(o.b)("p",null,"opctl is distributed as a self-contained executable, so installation generally consists of:"),Object(o.b)("ol",null,Object(o.b)("li",{parentName:"ol"},"Downloading the OS specific binary"),Object(o.b)("li",{parentName:"ol"},"Adding it to your path")),Object(o.b)("h3",{id:"prerequisites"},"Prerequisites"),Object(o.b)("p",null,"The default container runtime interface implementation relies on API access to a docker daemon to run containers.\n",Object(o.b)("a",Object(r.a)({parentName:"p"},{href:"https://docs.docker.com/install/"}),"Install Docker for your platform")),Object(o.b)("h3",{id:"osx"},"OSX"),Object(o.b)("pre",null,Object(o.b)("code",Object(r.a)({parentName:"pre"},{className:"language-bash"}),"curl -L https://github.com/opctl/opctl/releases/download/0.1.44/opctl0.1.44.darwin.tgz | tar -xzv -C /usr/local/bin\n")),Object(o.b)("h3",{id:"linux"},"Linux"),Object(o.b)("pre",null,Object(o.b)("code",Object(r.a)({parentName:"pre"},{className:"language-bash"}),"curl -L https://github.com/opctl/opctl/releases/download/0.1.44/opctl0.1.44.linux.tgz | sudo tar -xzv -C /usr/local/bin\n")),Object(o.b)("h3",{id:"windows"},"Windows"),Object(o.b)("p",null,"download and run the ",Object(o.b)("a",Object(r.a)({parentName:"p"},{href:"https://github.com/opctl/opctl/releases/download/0.1.44/opctl0.1.44.windows.msi"}),"windows installer")),Object(o.b)("h2",{id:"updating"},"Updating"),Object(o.b)("p",null,"to get the newest release of opctl"),Object(o.b)("pre",null,Object(o.b)("code",Object(r.a)({parentName:"pre"},{className:"language-bash"}),"opctl self-update\n")),Object(o.b)("h2",{id:"ide-plugins"},"IDE Plugins"),Object(o.b)("h3",{id:"vscode"},"VSCode"),Object(o.b)("ol",null,Object(o.b)("li",{parentName:"ol"},"install ",Object(o.b)("a",Object(r.a)({parentName:"li"},{href:"https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml"}),"vscode-yaml plugin")),Object(o.b)("li",{parentName:"ol"},"add to your user or workspace settings",Object(o.b)("pre",{parentName:"li"},Object(o.b)("code",Object(r.a)({parentName:"pre"},{className:"language-json"}),'"yaml.schemas": {\n "https://raw.githubusercontent.com/opctl/opctl/main/opspec/opfile/jsonschema.json": "/op.yml"\n }\n'))),Object(o.b)("li",{parentName:"ol"},"edit or create an op.yml w/ your fancy intellisense.")))}u.isMDXComponent=!0},176:function(e,t,n){"use strict";n.d(t,"a",(function(){return u})),n.d(t,"b",(function(){return O}));var r=n(0),a=n.n(r);function o(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function l(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?l(Object(n),!0).forEach((function(t){o(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):l(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function c(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},o=Object.keys(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var s=a.a.createContext({}),p=function(e){var t=a.a.useContext(s),n=t;return e&&(n="function"==typeof e?e(t):i({},t,{},e)),n},u=function(e){var t=p(e.components);return(a.a.createElement(s.Provider,{value:t},e.children))},b="mdxType",d={inlineCode:"code",wrapper:function(e){var t=e.children;return a.a.createElement(a.a.Fragment,{},t)}},m=Object(r.forwardRef)((function(e,t){var n=e.components,r=e.mdxType,o=e.originalType,l=e.parentName,s=c(e,["components","mdxType","originalType","parentName"]),u=p(n),b=r,m=u["".concat(l,".").concat(b)]||u[b]||d[b]||o;return n?a.a.createElement(m,i({ref:t},s,{components:n})):a.a.createElement(m,i({ref:t},s))}));function O(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var o=n.length,l=new Array(o);l[0]=m;var i={};for(var c in t)hasOwnProperty.call(t,c)&&(i[c]=t[c]);i.originalType=e,i[b]="string"==typeof e?e:r,l[1]=i;for(var s=2;s<o;s++)l[s]=n[s];return a.a.createElement.apply(null,l)}return a.a.createElement.apply(null,n)}m.displayName="MDXCreateElement"}}]);