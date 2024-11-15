"use strict";(self.webpackChunkopctl=self.webpackChunkopctl||[]).push([[1841],{2713:(e,t,o)=>{o.r(t),o.d(t,{assets:()=>c,contentTitle:()=>i,default:()=>p,frontMatter:()=>s,metadata:()=>l,toc:()=>a});const l=JSON.parse('{"id":"reference/cli/op/validate","title":"opctl op validate","description":"Validate an op according to:","source":"@site/docs/reference/cli/op/validate.md","sourceDirName":"reference/cli/op","slug":"/reference/cli/op/validate","permalink":"/docs/reference/cli/op/validate","draft":false,"unlisted":false,"editUrl":"https://github.com/opctl/opctl/edit/main/website/docs/reference/cli/op/validate.md","tags":[],"version":"current","lastUpdatedBy":"Chris Dostert","lastUpdatedAt":1604531186000,"frontMatter":{"sidebar_label":"validate","title":"opctl op validate"},"sidebar":"docs","previous":{"title":"kill","permalink":"/docs/reference/cli/op/kill"},"next":{"title":"run","permalink":"/docs/reference/cli/run"}}');var n=o(4848),r=o(8453);const s={sidebar_label:"validate",title:"opctl op validate"},i=void 0,c={},a=[{value:"Arguments",id:"arguments",level:2},{value:"<code>OP_REF</code>",id:"op_ref",level:3},{value:"Examples",id:"examples",level:2},{value:"Global Options",id:"global-options",level:2},{value:"Notes",id:"notes",level:2},{value:"op source username/password prompt",id:"op-source-usernamepassword-prompt",level:4}];function d(e){const t={a:"a",blockquote:"blockquote",code:"code",h2:"h2",h3:"h3",h4:"h4",li:"li",p:"p",pre:"pre",ul:"ul",...(0,r.R)(),...e.components};return(0,n.jsxs)(n.Fragment,{children:[(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-sh",children:"opctl op validate [OPTIONS] OP_REF\n"})}),"\n",(0,n.jsx)(t.p,{children:"Validate an op according to:"}),"\n",(0,n.jsxs)(t.ul,{children:["\n",(0,n.jsxs)(t.li,{children:["existence of ",(0,n.jsx)(t.code,{children:"op.yml"})]}),"\n",(0,n.jsxs)(t.li,{children:["validity of ",(0,n.jsx)(t.code,{children:"op.yml"})," (per\n",(0,n.jsx)(t.a,{href:"https://opctl.io/0.1.6/op.yml.schema.json",children:"schema"}),")"]}),"\n"]}),"\n",(0,n.jsx)(t.h2,{id:"arguments",children:"Arguments"}),"\n",(0,n.jsx)(t.h3,{id:"op_ref",children:(0,n.jsx)(t.code,{children:"OP_REF"})}),"\n",(0,n.jsxs)(t.p,{children:["Op reference (either ",(0,n.jsx)(t.code,{children:"relative/path"}),", ",(0,n.jsx)(t.code,{children:"/absolute/path"}),", ",(0,n.jsx)(t.code,{children:"host/path/repo#tag"}),", or ",(0,n.jsx)(t.code,{children:"host/path/repo#tag/path"}),")."]}),"\n",(0,n.jsx)(t.h2,{id:"examples",children:"Examples"}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-sh",children:"opctl op validate myop\n"})}),"\n",(0,n.jsx)(t.h2,{id:"global-options",children:"Global Options"}),"\n",(0,n.jsxs)(t.p,{children:["see ",(0,n.jsx)(t.a,{href:"/docs/reference/cli/global-options",children:"global options"})]}),"\n",(0,n.jsx)(t.h2,{id:"notes",children:"Notes"}),"\n",(0,n.jsx)(t.h4,{id:"op-source-usernamepassword-prompt",children:"op source username/password prompt"}),"\n",(0,n.jsx)(t.p,{children:"If auth w/ the op source fails the cli will (re)prompt for username & password."}),"\n",(0,n.jsxs)(t.blockquote,{children:["\n",(0,n.jsx)(t.p,{children:"in non-interactive terminals, the cli will note that it can't prompt and exit with a non zero exit code."}),"\n"]})]})}function p(e={}){const{wrapper:t}={...(0,r.R)(),...e.components};return t?(0,n.jsx)(t,{...e,children:(0,n.jsx)(d,{...e})}):d(e)}},8453:(e,t,o)=>{o.d(t,{R:()=>s,x:()=>i});var l=o(6540);const n={},r=l.createContext(n);function s(e){const t=l.useContext(r);return l.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function i(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(n):e.components||n:s(e.components),l.createElement(r.Provider,{value:t},e.children)}}}]);