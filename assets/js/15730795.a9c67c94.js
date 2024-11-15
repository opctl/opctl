"use strict";(self.webpackChunkopctl=self.webpackChunkopctl||[]).push([[2214],{2184:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>a,contentTitle:()=>s,default:()=>p,frontMatter:()=>c,metadata:()=>o,toc:()=>l});const o=JSON.parse('{"id":"training/containers/container-execution","title":"Container execution","description":"TLDR;","source":"@site/docs/training/containers/container-execution.md","sourceDirName":"training/containers","slug":"/training/containers/container-execution","permalink":"/docs/training/containers/container-execution","draft":false,"unlisted":false,"editUrl":"https://github.com/opctl/opctl/edit/main/website/docs/training/containers/container-execution.md","tags":[],"version":"current","lastUpdatedBy":"=","lastUpdatedAt":1682978842000,"frontMatter":{"title":"Container execution"},"sidebar":"docs","previous":{"title":"Setup","permalink":"/docs/setup"},"next":{"title":"Container networking","permalink":"/docs/training/containers/container-networking"}}');var i=t(4848),r=t(8453);const c={title:"Container execution"},s=void 0,a={},l=[{value:"TLDR;",id:"tldr",level:2},{value:"Example",id:"example",level:2}];function d(e){const n={a:"a",blockquote:"blockquote",code:"code",h2:"h2",li:"li",ol:"ol",p:"p",pre:"pre",...(0,r.R)(),...e.components};return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(n.h2,{id:"tldr",children:"TLDR;"}),"\n",(0,i.jsxs)(n.p,{children:["Opctl supports using ",(0,i.jsx)(n.a,{href:"/docs/reference/opspec/op-directory/op/call/container/",children:"container"})," statements to make your op run ",(0,i.jsx)(n.a,{href:"https://opencontainers.org/",children:"OCI"})," image based containers."]}),"\n",(0,i.jsxs)(n.blockquote,{children:["\n",(0,i.jsxs)(n.p,{children:["Note: a common place to obtain ",(0,i.jsx)(n.a,{href:"https://opencontainers.org/",children:"OCI"})," images is ",(0,i.jsx)(n.a,{href:"https://hub.docker.com/",children:"Docker Hub"}),"."]}),"\n"]}),"\n",(0,i.jsx)(n.h2,{id:"example",children:"Example"}),"\n",(0,i.jsxs)(n.ol,{children:["\n",(0,i.jsxs)(n.li,{children:["Start this op:","\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-yaml",children:"name: containerExecution\nrun:\n  container:\n    cmd: [echo, 'hello!']\n    image: { ref: alpine }\n"})}),"\n"]}),"\n",(0,i.jsxs)(n.li,{children:["Observe the container is started, ",(0,i.jsx)(n.code,{children:"hello!"})," is logged, and the container exits."]}),"\n"]})]})}function p(e={}){const{wrapper:n}={...(0,r.R)(),...e.components};return n?(0,i.jsx)(n,{...e,children:(0,i.jsx)(d,{...e})}):d(e)}},8453:(e,n,t)=>{t.d(n,{R:()=>c,x:()=>s});var o=t(6540);const i={},r=o.createContext(i);function c(e){const n=o.useContext(r);return o.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function s(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(i):e.components||i:c(e.components),o.createElement(r.Provider,{value:n},e.children)}}}]);