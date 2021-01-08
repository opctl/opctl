(window.webpackJsonp=window.webpackJsonp||[]).push([[47],{149:function(e,t,n){"use strict";n.r(t),n.d(t,"frontMatter",(function(){return c})),n.d(t,"metadata",(function(){return p})),n.d(t,"rightToc",(function(){return i})),n.d(t,"default",(function(){return u}));var r=n(1),a=n(6),o=(n(0),n(178)),c={title:"Run a react app",sidebar_label:"run a react app"},p={id:"run-a-react-app",title:"Run a react app",description:"We'll now look at an op to run a sample React application. Please see the code at [https://github.com/opctl/opctl/tree/main/examples/run-a-react-app](https://github.com/opctl/opctl/tree/main/examples/run-a-react-app)",source:"@site/docs/run-a-react-app.md",permalink:"/docs/run-a-react-app",editUrl:"https://github.com/opctl/opctl/edit/master/website/docs/run-a-react-app.md",lastUpdatedBy:"Chris Dostert",lastUpdatedAt:1605586799,sidebar_label:"run a react app",sidebar:"docs",previous:{title:"Run a go service",permalink:"/docs/run-a-go-service"},next:{title:"How do I communicate with an opctl container?",permalink:"/docs/training/containers/how-do-i-communicate-with-an-opctl-container"}},i=[],l={rightToc:i},s="wrapper";function u(e){var t=e.components,n=Object(a.a)(e,["components"]);return Object(o.b)(s,Object(r.a)({},l,n,{components:t,mdxType:"MDXLayout"}),Object(o.b)("p",null,"We'll now look at an op to run a sample React application. Please see the code at ",Object(o.b)("a",Object(r.a)({parentName:"p"},{href:"https://github.com/opctl/opctl/tree/main/examples/run-a-react-app"}),"https://github.com/opctl/opctl/tree/main/examples/run-a-react-app")),Object(o.b)("p",null,"We'll stick to simple node.js conventions of including the run command in an ",Object(o.b)("inlineCode",{parentName:"p"},"npm start")," script in ",Object(o.b)("inlineCode",{parentName:"p"},"package.json"),". Because we used ",Object(o.b)("inlineCode",{parentName:"p"},"create-react-app")," to bootstrap our project, the start script is ",Object(o.b)("inlineCode",{parentName:"p"},"react-scripts start"),", which will launch the webpack dev server to serve our app for development."),Object(o.b)("p",null,"Let's say our frontend React app needs to call the a go api from our previous example. When running that stack we want our services running locally (React app running via webpack, our go service, and the mysql database). We however prefer to have the ops running each service to live with the source code it runs, rather than in a separate place."),Object(o.b)("p",null,"What we need our ",Object(o.b)("inlineCode",{parentName:"p"},"dev")," op to do then is to:\n1. call ",Object(o.b)("inlineCode",{parentName:"p"},"go-svc"),"'s ",Object(o.b)("inlineCode",{parentName:"p"},"dev")," op by remote reference\n3. run our React app in a container"),Object(o.b)("p",null,"So our ops in ",Object(o.b)("inlineCode",{parentName:"p"},"run-a-react-app")," would look like this"),Object(o.b)("h4",{id:"dev"},"dev"),Object(o.b)("pre",null,Object(o.b)("code",Object(r.a)({parentName:"pre"},{className:"language-yaml"}),"name: dev\ndescription: runs react-app for development\nrun:\n  parallel:\n    - op:\n        # reference run-a-go-service dev op\n        ref: ../../../run-a-go-service/.opspec/dev\n    - container:\n        image:\n          ref: 'node:15-alpine'\n        cmd:\n          - sh\n          - -ce\n          - yarn && yarn run start\n        dirs:\n          /src: $(../..)\n        workDir: /src\n        ports:\n          3000: 3000\n")),Object(o.b)("p",null,"going to http://localhost:8080 should show us the ",Object(o.b)("inlineCode",{parentName:"p"},"go-svc")," api being served, and http://localhost:3000 should show us the react app, which in turns is making a call to ",Object(o.b)("inlineCode",{parentName:"p"},"go-svc")," and fetching data to show."),Object(o.b)("p",null,"Notice the following:\n1. we're referencing the ",Object(o.b)("inlineCode",{parentName:"p"},"dev")," op from the go-svc.\n2. we can run our ops in any combination of ",Object(o.b)("inlineCode",{parentName:"p"},"parallel")," and ",Object(o.b)("inlineCode",{parentName:"p"},"serial")," blocks, composing them as needed. for our case the ",Object(o.b)("inlineCode",{parentName:"p"},"dev"),' op can run in the background while we init then run our react app\n3. networking between our services "just works" by referencing containers by name, thanks to how they all are in the same Docker container network, so the webpack dev server proxy configuration in ',Object(o.b)("inlineCode",{parentName:"p"},"package.json")," targets ",Object(o.b)("inlineCode",{parentName:"p"},"go-svc"),":"),Object(o.b)("pre",null,Object(o.b)("code",Object(r.a)({parentName:"pre"},{className:"language-json"}),'"proxy": {\n    "/api": {\n      "target": "http://go-svc:8080"\n    }\n  }\n')))}u.isMDXComponent=!0},178:function(e,t,n){"use strict";n.d(t,"a",(function(){return u})),n.d(t,"b",(function(){return h}));var r=n(0),a=n.n(r);function o(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function c(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function p(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?c(Object(n),!0).forEach((function(t){o(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):c(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function i(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},o=Object.keys(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var l=a.a.createContext({}),s=function(e){var t=a.a.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):p({},t,{},e)),n},u=function(e){var t=s(e.components);return(a.a.createElement(l.Provider,{value:t},e.children))},b="mdxType",d={inlineCode:"code",wrapper:function(e){var t=e.children;return a.a.createElement(a.a.Fragment,{},t)}},m=Object(r.forwardRef)((function(e,t){var n=e.components,r=e.mdxType,o=e.originalType,c=e.parentName,l=i(e,["components","mdxType","originalType","parentName"]),u=s(n),b=r,m=u["".concat(c,".").concat(b)]||u[b]||d[b]||o;return n?a.a.createElement(m,p({ref:t},l,{components:n})):a.a.createElement(m,p({ref:t},l))}));function h(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var o=n.length,c=new Array(o);c[0]=m;var p={};for(var i in t)hasOwnProperty.call(t,i)&&(p[i]=t[i]);p.originalType=e,p[b]="string"==typeof e?e:r,c[1]=p;for(var l=2;l<o;l++)c[l]=n[l];return a.a.createElement.apply(null,c)}return a.a.createElement.apply(null,n)}m.displayName="MDXCreateElement"}}]);