"use strict";(self.webpackChunkopctl=self.webpackChunkopctl||[]).push([[5666],{3905:function(e,t,n){n.d(t,{Zo:function(){return s},kt:function(){return d}});var r=n(7294);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function l(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},o=Object.keys(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var c=r.createContext({}),p=function(e){var t=r.useContext(c),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},s=function(e){var t=p(e.components);return r.createElement(c.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},f=r.forwardRef((function(e,t){var n=e.components,a=e.mdxType,o=e.originalType,c=e.parentName,s=l(e,["components","mdxType","originalType","parentName"]),f=p(n),d=a,m=f["".concat(c,".").concat(d)]||f[d]||u[d]||o;return n?r.createElement(m,i(i({ref:t},s),{},{components:n})):r.createElement(m,i({ref:t},s))}));function d(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var o=n.length,i=new Array(o);i[0]=f;var l={};for(var c in t)hasOwnProperty.call(t,c)&&(l[c]=t[c]);l.originalType=e,l.mdxType="string"==typeof e?e:a,i[1]=l;for(var p=2;p<o;p++)i[p]=n[p];return r.createElement.apply(null,i)}return r.createElement.apply(null,n)}f.displayName="MDXCreateElement"},933:function(e,t,n){n.r(t),n.d(t,{assets:function(){return s},contentTitle:function(){return c},default:function(){return d},frontMatter:function(){return l},metadata:function(){return p},toc:function(){return u}});var r=n(3117),a=n(102),o=(n(7294),n(3905)),i=["components"],l={title:"Gitlab",sidebar_label:"Gitlab"},c=void 0,p={unversionedId:"setup/gitlab",id:"setup/gitlab",title:"Gitlab",description:'gitlab ci looks for a .gitlab-ci.yml file at the root of each repo to identify ci "stages".',source:"@site/docs/setup/gitlab.md",sourceDirName:"setup",slug:"/setup/gitlab",permalink:"/docs/setup/gitlab",draft:!1,editUrl:"https://github.com/opctl/opctl/edit/main/website/docs/setup/gitlab.md",tags:[],version:"current",lastUpdatedBy:"Chris Dostert",lastUpdatedAt:1614614385,formattedLastUpdatedAt:"Mar 1, 2021",frontMatter:{title:"Gitlab",sidebar_label:"Gitlab"},sidebar:"docs",previous:{title:"Github",permalink:"/docs/setup/github"},next:{title:"Kubernetes",permalink:"/docs/setup/kubernetes"}},s={},u=[{value:"Examples",id:"examples",level:3}],f={toc:u};function d(e){var t=e.components,n=(0,a.Z)(e,i);return(0,o.kt)("wrapper",(0,r.Z)({},f,n,{components:t,mdxType:"MDXLayout"}),(0,o.kt)("p",null,(0,o.kt)("a",{parentName:"p",href:"https://gitlab.io"},"gitlab")," ci looks for a ",(0,o.kt)("inlineCode",{parentName:"p"},".gitlab-ci.yml"),' file at the root of each repo to identify ci "stages".'),(0,o.kt)("p",null,"Their hosted agents support running the ci process within a docker container so running opctl is\njust a matter of defining your ",(0,o.kt)("inlineCode",{parentName:"p"},".gitlab-ci.yml")," as follows:"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},"using the official ",(0,o.kt)("a",{parentName:"li",href:"https://hub.docker.com/r/opctl/opctl/"},"opctl docker image")," as ",(0,o.kt)("inlineCode",{parentName:"li"},"image")),(0,o.kt)("li",{parentName:"ul"},'adding "stages" with your calls to opctl')),(0,o.kt)("h3",{id:"examples"},"Examples"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},'# .gitlab-ci.yml\nimage: opctl/opctl:0.1.48-dind\nstages:\n  - build\n  - deploy\nbuild:\n  stage: build\n  script:\n    # passes args to opctl from gitlab variables\n    - export gitlabUsername="$CI_REGISTRY_USER"\n    - export gitlabSecret="$CI_REGISTRY_PASSWORD"\n    - opctl run build\ndeploy:\n  stage: deploy\n  only:\n    - master\n  script:\n    - opctl run deploy\n')))}d.isMDXComponent=!0}}]);