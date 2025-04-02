"use strict";(self.webpackChunkopctl=self.webpackChunkopctl||[]).push([[9623],{6860:(e,r,n)=>{n.r(r),n.d(r,{assets:()=>s,contentTitle:()=>t,default:()=>p,frontMatter:()=>l,metadata:()=>i,toc:()=>a});const i=JSON.parse('{"id":"reference/opspec/op-directory/op/call/container/image","title":"Image [object]","description":"An object which defines the image of a container call.","source":"@site/docs/reference/opspec/op-directory/op/call/container/image.md","sourceDirName":"reference/opspec/op-directory/op/call/container","slug":"/reference/opspec/op-directory/op/call/container/image","permalink":"/docs/reference/opspec/op-directory/op/call/container/image","draft":false,"unlisted":false,"editUrl":"https://github.com/opctl/opctl/edit/main/website/docs/reference/opspec/op-directory/op/call/container/image.md","tags":[],"version":"current","lastUpdatedBy":"Chris Dostert","lastUpdatedAt":1743610274000,"frontMatter":{"title":"Image [object]"},"sidebar":"docs","previous":{"title":"Overview","permalink":"/docs/reference/opspec/op-directory/op/call/container/"},"next":{"title":"Loop Vars [object]","permalink":"/docs/reference/opspec/op-directory/op/call/loop-vars"}}');var c=n(4848),o=n(8453);const l={title:"Image [object]"},t=void 0,s={},a=[{value:"Properties",id:"properties",level:2},{value:"ref",id:"ref",level:3},{value:"Example ref (docker.io/ubuntu:19.10)",id:"example-ref-dockerioubuntu1910",level:3},{value:"Example ref (variable)",id:"example-ref-variable",level:3},{value:"platform",id:"platform",level:3},{value:"pullCreds",id:"pullcreds",level:3}];function d(e){const r={a:"a",code:"code",h2:"h2",h3:"h3",li:"li",p:"p",ul:"ul",...(0,o.R)(),...e.components};return(0,c.jsxs)(c.Fragment,{children:[(0,c.jsx)(r.p,{children:"An object which defines the image of a container call."}),"\n",(0,c.jsx)(r.h2,{id:"properties",children:"Properties"}),"\n",(0,c.jsxs)(r.ul,{children:["\n",(0,c.jsxs)(r.li,{children:["must have","\n",(0,c.jsxs)(r.ul,{children:["\n",(0,c.jsx)(r.li,{children:(0,c.jsx)(r.a,{href:"#ref",children:"ref"})}),"\n"]}),"\n"]}),"\n",(0,c.jsxs)(r.li,{children:["may have","\n",(0,c.jsxs)(r.ul,{children:["\n",(0,c.jsx)(r.li,{children:(0,c.jsx)(r.a,{href:"#platform",children:"platform"})}),"\n",(0,c.jsx)(r.li,{children:(0,c.jsx)(r.a,{href:"#pullcreds",children:"pullCreds"})}),"\n"]}),"\n"]}),"\n"]}),"\n",(0,c.jsx)(r.h3,{id:"ref",children:"ref"}),"\n",(0,c.jsx)(r.p,{children:"A string referencing a local or remote image."}),"\n",(0,c.jsx)(r.p,{children:"Must be one of:"}),"\n",(0,c.jsxs)(r.ul,{children:["\n",(0,c.jsxs)(r.li,{children:["a ",(0,c.jsx)(r.a,{href:"/docs/reference/opspec/op-directory/op/variable-reference",children:"variable-reference [string]"})," evaluating to a ",(0,c.jsxs)(r.a,{href:"https://github.com/opencontainers/image-spec/blob/v1.0.1/image-layout.md",children:["v1.0.1 OCI (Open Container Initiative) ",(0,c.jsx)(r.code,{children:"image-layout"})]}),"."]}),"\n",(0,c.jsxs)(r.li,{children:["a ",(0,c.jsx)(r.a,{href:"/docs/reference/opspec/types/string#initialization",children:"string initializer"})," evaluating to a docker image name i.e. ",(0,c.jsx)(r.code,{children:"[host][repository]image[tag]"})," where by default host is ",(0,c.jsx)(r.code,{children:"docker.io"})," and tag is ",(0,c.jsx)(r.code,{children:"latest"})]}),"\n"]}),"\n",(0,c.jsxs)(r.h3,{id:"example-ref-dockerioubuntu1910",children:["Example ref (",(0,c.jsx)(r.a,{href:"https://hub.docker.com/_/ubuntu",children:"docker.io/ubuntu:19.10"}),")"]}),"\n",(0,c.jsxs)(r.p,{children:[(0,c.jsx)(r.code,{children:"ref: 'ubuntu:19.10'"})," or ",(0,c.jsx)(r.code,{children:"ref: 'docker.io/ubuntu:19.10'"})]}),"\n",(0,c.jsx)(r.h3,{id:"example-ref-variable",children:"Example ref (variable)"}),"\n",(0,c.jsx)(r.p,{children:(0,c.jsx)(r.code,{children:"ref: $(myOCIImageLayoutDir)"})}),"\n",(0,c.jsx)(r.h3,{id:"platform",children:"platform"}),"\n",(0,c.jsxs)(r.p,{children:["An ",(0,c.jsx)(r.a,{href:"/docs/reference/opspec/op-directory/op/call/oci-image-platform",children:"oci-image-platform [object]"})," constraining the image which will be pulled from the source."]}),"\n",(0,c.jsx)(r.h3,{id:"pullcreds",children:"pullCreds"}),"\n",(0,c.jsxs)(r.p,{children:["A ",(0,c.jsx)(r.a,{href:"/docs/reference/opspec/op-directory/op/call/pull-creds",children:"pull-creds [object]"})," defining creds used to pull the image from a private source."]})]})}function p(e={}){const{wrapper:r}={...(0,o.R)(),...e.components};return r?(0,c.jsx)(r,{...e,children:(0,c.jsx)(d,{...e})}):d(e)}},8453:(e,r,n)=>{n.d(r,{R:()=>l,x:()=>t});var i=n(6540);const c={},o=i.createContext(c);function l(e){const r=i.useContext(o);return i.useMemo((function(){return"function"==typeof e?e(r):{...r,...e}}),[r,e])}function t(e){let r;return r=e.disableParentContext?"function"==typeof e.components?e.components(c):e.components||c:l(e.components),i.createElement(o.Provider,{value:r},e.children)}}}]);