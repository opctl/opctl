(window.webpackJsonp=window.webpackJsonp||[]).push([[64],{204:function(e,t,n){"use strict";n.r(t),n.d(t,"frontMatter",(function(){return l})),n.d(t,"metadata",(function(){return i})),n.d(t,"rightToc",(function(){return c})),n.d(t,"default",(function(){return u}));var a=n(1),r=n(9),o=(n(0),n(216)),l={sidebar_label:"run",title:"opctl run"},i={id:"reference/cli/run",title:"opctl run",description:"```sh",source:"@site/docs/reference/cli/run.md",permalink:"/docs/reference/cli/run",editUrl:"https://github.com/opctl/opctl/edit/main/website/docs/reference/cli/run.md",lastUpdatedBy:"=",lastUpdatedAt:1628883780,sidebar_label:"run",sidebar:"docs",previous:{title:"opctl op validate",permalink:"/docs/reference/cli/op/validate"},next:{title:"opctl self-update",permalink:"/docs/reference/cli/self-update"}},c=[{value:"Arguments",id:"arguments",children:[{value:"<code>OP_REF</code>",id:"op_ref",children:[]}]},{value:"Options",id:"options",children:[{value:"<code>-a</code>",id:"-a",children:[]},{value:"<code>--arg-file</code> <em>default: <code>.opspec/args.yml</code></em>",id:"--arg-file-default-opspecargsyml",children:[]},{value:"<code>--no-progress</code> <em>default: <code>false</code></em>",id:"--no-progress-default-false",children:[]}]},{value:"Global Options",id:"global-options",children:[]},{value:"Examples",id:"examples",children:[{value:"local op ref w/out args",id:"local-op-ref-wout-args",children:[]},{value:"remote op ref w/ args",id:"remote-op-ref-w-args",children:[]}]},{value:"Notes",id:"notes",children:[{value:"op source username/password prompt",id:"op-source-usernamepassword-prompt",children:[]},{value:"input sources",id:"input-sources",children:[]},{value:"input prompts",id:"input-prompts",children:[]},{value:"validation",id:"validation",children:[]},{value:"caching",id:"caching",children:[]},{value:"image updates",id:"image-updates",children:[]},{value:"container networking",id:"container-networking",children:[]},{value:"container cleanup",id:"container-cleanup",children:[]}]}],p={rightToc:c},b="wrapper";function u(e){var t=e.components,n=Object(r.a)(e,["components"]);return Object(o.b)(b,Object(a.a)({},p,n,{components:t,mdxType:"MDXLayout"}),Object(o.b)("pre",null,Object(o.b)("code",Object(a.a)({parentName:"pre"},{className:"language-sh"}),"opctl run [OPTIONS] OP_REF\n")),Object(o.b)("p",null,"Start and wait on an op."),Object(o.b)("blockquote",null,Object(o.b)("p",{parentName:"blockquote"},"if a node isn't running, one will be automatically created")),Object(o.b)("h2",{id:"arguments"},"Arguments"),Object(o.b)("h3",{id:"op_ref"},Object(o.b)("inlineCode",{parentName:"h3"},"OP_REF")),Object(o.b)("p",null,"Op reference (either ",Object(o.b)("inlineCode",{parentName:"p"},"relative/path"),", ",Object(o.b)("inlineCode",{parentName:"p"},"/absolute/path"),", ",Object(o.b)("inlineCode",{parentName:"p"},"host/path/repo#tag"),", or ",Object(o.b)("inlineCode",{parentName:"p"},"host/path/repo#tag/path"),")"),Object(o.b)("h2",{id:"options"},"Options"),Object(o.b)("h3",{id:"-a"},Object(o.b)("inlineCode",{parentName:"h3"},"-a")),Object(o.b)("p",null,"Explicitly pass args to op in format ",Object(o.b)("inlineCode",{parentName:"p"},"-a NAME1=VALUE1 -a NAME2=VALUE2")),Object(o.b)("h3",{id:"--arg-file-default-opspecargsyml"},Object(o.b)("inlineCode",{parentName:"h3"},"--arg-file")," ",Object(o.b)("em",{parentName:"h3"},"default: ",Object(o.b)("inlineCode",{parentName:"em"},".opspec/args.yml"))),Object(o.b)("p",null,"Read in a file of args in yml format"),Object(o.b)("h3",{id:"--no-progress-default-false"},Object(o.b)("inlineCode",{parentName:"h3"},"--no-progress")," ",Object(o.b)("em",{parentName:"h3"},"default: ",Object(o.b)("inlineCode",{parentName:"em"},"false"))),Object(o.b)("p",null,"Disable live call graph for the op"),Object(o.b)("h2",{id:"global-options"},"Global Options"),Object(o.b)("p",null,"see ",Object(o.b)("a",Object(a.a)({parentName:"p"},{href:"/docs/reference/cli/global-options"}),"global options")),Object(o.b)("h2",{id:"examples"},"Examples"),Object(o.b)("h3",{id:"local-op-ref-wout-args"},"local op ref w/out args"),Object(o.b)("pre",null,Object(o.b)("code",Object(a.a)({parentName:"pre"},{className:"language-sh"}),"opctl run myop\n")),Object(o.b)("h3",{id:"remote-op-ref-w-args"},"remote op ref w/ args"),Object(o.b)("pre",null,Object(o.b)("code",Object(a.a)({parentName:"pre"},{className:"language-sh"}),'opctl run -a apiToken="my-token" -a channelName="my-channel" -a msg="hello!" github.com/opspec-pkgs/slack.chat.post-message#0.1.1\n')),Object(o.b)("h2",{id:"notes"},"Notes"),Object(o.b)("h3",{id:"op-source-usernamepassword-prompt"},"op source username/password prompt"),Object(o.b)("p",null,"If auth w/ the op source fails the cli will (re)prompt for username &\npassword."),Object(o.b)("blockquote",null,Object(o.b)("p",{parentName:"blockquote"},"in non-interactive terminals, the cli will note that it can't prompt\ndue to being in a non-interactive terminal and exit with a non zero\nexit code.")),Object(o.b)("h3",{id:"input-sources"},"input sources"),Object(o.b)("p",null,"Input sources are checked according to the following precedence:"),Object(o.b)("ul",null,Object(o.b)("li",{parentName:"ul"},"arg provided via ",Object(o.b)("inlineCode",{parentName:"li"},"-a")," option"),Object(o.b)("li",{parentName:"ul"},"arg file"),Object(o.b)("li",{parentName:"ul"},"env var"),Object(o.b)("li",{parentName:"ul"},"default"),Object(o.b)("li",{parentName:"ul"},"prompt")),Object(o.b)("h3",{id:"input-prompts"},"input prompts"),Object(o.b)("p",null,"Inputs which are invalid or missing will result in the cli prompting for\nthem."),Object(o.b)("blockquote",null,Object(o.b)("p",{parentName:"blockquote"},"in non-interactive terminals, the cli will provide details about the\ninvalid or missing input, note that it's giving up due to being in a\nnon-interactive terminal and exit with a non zero exit code.")),Object(o.b)("p",null,"example:"),Object(o.b)("pre",null,Object(o.b)("code",Object(a.a)({parentName:"pre"},{className:"language-sh"}),"\n-\n  Please provide value for parameter.\n  Name: version\n  Description: version of app being compiled\n-\n")),Object(o.b)("h3",{id:"validation"},"validation"),Object(o.b)("p",null,"When inputs don't meet constraints, the cli will (re)prompt for the\ninput until a satisfactory value is obtained."),Object(o.b)("h3",{id:"caching"},"caching"),Object(o.b)("p",null,"All pulled ops/image layers will be cached"),Object(o.b)("h3",{id:"image-updates"},"image updates"),Object(o.b)("p",null,"Prior to container creation, updates to the referenced image will be\npulled and applied."),Object(o.b)("p",null,"If checking for or applying updated image layers fails, graceful\nfallback to cached image layers will occur"),Object(o.b)("h3",{id:"container-networking"},"container networking"),Object(o.b)("p",null,"All containers created by opctl will be attached to a single managed\nnetwork."),Object(o.b)("blockquote",null,Object(o.b)("p",{parentName:"blockquote"},"the network is visible from ",Object(o.b)("inlineCode",{parentName:"p"},"docker network ls")," as ",Object(o.b)("inlineCode",{parentName:"p"},"opctl"),".")),Object(o.b)("h3",{id:"container-cleanup"},"container cleanup"),Object(o.b)("p",null,"Containers will be removed as they exit."))}u.isMDXComponent=!0},216:function(e,t,n){"use strict";n.d(t,"a",(function(){return u})),n.d(t,"b",(function(){return O}));var a=n(0),r=n.n(a);function o(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function l(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?l(Object(n),!0).forEach((function(t){o(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):l(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function c(e,t){if(null==e)return{};var n,a,r=function(e,t){if(null==e)return{};var n,a,r={},o=Object.keys(e);for(a=0;a<o.length;a++)n=o[a],t.indexOf(n)>=0||(r[n]=e[n]);return r}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(a=0;a<o.length;a++)n=o[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(r[n]=e[n])}return r}var p=r.a.createContext({}),b=function(e){var t=r.a.useContext(p),n=t;return e&&(n="function"==typeof e?e(t):i({},t,{},e)),n},u=function(e){var t=b(e.components);return(r.a.createElement(p.Provider,{value:t},e.children))},s="mdxType",d={inlineCode:"code",wrapper:function(e){var t=e.children;return r.a.createElement(r.a.Fragment,{},t)}},m=Object(a.forwardRef)((function(e,t){var n=e.components,a=e.mdxType,o=e.originalType,l=e.parentName,p=c(e,["components","mdxType","originalType","parentName"]),u=b(n),s=a,m=u["".concat(l,".").concat(s)]||u[s]||d[s]||o;return n?r.a.createElement(m,i({ref:t},p,{components:n})):r.a.createElement(m,i({ref:t},p))}));function O(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var o=n.length,l=new Array(o);l[0]=m;var i={};for(var c in t)hasOwnProperty.call(t,c)&&(i[c]=t[c]);i.originalType=e,i[s]="string"==typeof e?e:a,l[1]=i;for(var p=2;p<o;p++)l[p]=n[p];return r.a.createElement.apply(null,l)}return r.a.createElement.apply(null,n)}m.displayName="MDXCreateElement"}}]);