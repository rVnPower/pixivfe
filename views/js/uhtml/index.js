const{isArray:e}=Array,{getPrototypeOf:t,getOwnPropertyDescriptor:n}=Object,r=[],s=()=>document.createRange(),l=(e,t,n)=>(e.set(t,n),n),{setPrototypeOf:i}=Object;let o;var c=(e,t,n)=>(o||(o=s()),n?o.setStartAfter(e):o.setStartBefore(e),o.setEndAfter(t),o.deleteContents(),e);const a=({firstChild:e,lastChild:t},n)=>c(e,t,n);let u=!1;const h=(e,t)=>u&&11===e.nodeType?1/t<0?t?a(e,!0):e.lastChild:t?e.valueOf():e.firstChild:e,d=e=>document.createComment(e);class f extends((e=>{function t(e){return i(e,new.target.prototype)}return t.prototype=e.prototype,t})(DocumentFragment)){#e=d("<>");#t=d("</>");#n=r;constructor(e){super(e),this.replaceChildren(this.#e,...e.childNodes,this.#t),u=!0}get firstChild(){return this.#e}get lastChild(){return this.#t}get parentNode(){return this.#e.parentNode}remove(){a(this,!1)}replaceWith(e){a(this,!0).replaceWith(e)}valueOf(){let{firstChild:e,lastChild:t,parentNode:n}=this;if(n===this)this.#n===r&&(this.#n=[...this.childNodes]);else{if(n)for(this.#n=[e];e!==t;)this.#n.push(e=e.nextSibling);this.replaceChildren(...this.#n)}return this}}const p=(e,t,n)=>e.setAttribute(t,n),g=(e,t)=>e.removeAttribute(t);let m;const v=(t,n,r)=>{r=r.slice(1),m||(m=new WeakMap);const s=m.get(t)||l(m,t,{});let i=s[r];return i&&i[0]&&t.removeEventListener(r,...i),i=e(n)?n:[n,!1],s[r]=i,i[0]&&t.addEventListener(r,...i),n},b=(e,t)=>{const{t:n,n:r}=e;let s=!1;switch(typeof t){case"object":if(null!==t){(r||n).replaceWith(e.n=t.valueOf());break}case"undefined":s=!0;default:n.data=s?"":t,r&&(e.n=null,r.replaceWith(n))}return t},C=(e,t,n)=>e[n]=t,x=(e,t,n)=>C(e,t,n.slice(1)),w=(e,t,n)=>null==t?(g(e,n),t):C(e,t,n),$=(e,t)=>("function"==typeof t?t(e):t.current=e,t),y=(e,t,n)=>(null==t?g(e,n):p(e,n,t),t),N=(e,t,n)=>(e.toggleAttribute(n.slice(1),t),t),O=(e,t,n)=>{const{length:s}=t;if(e.data=`[${s}]`,s)return((e,t,n,r,s)=>{const l=n.length;let i=t.length,o=l,c=0,a=0,u=null;for(;c<i||a<o;)if(i===c){const t=o<l?a?r(n[a-1],-0).nextSibling:r(n[o-a],0):s;for(;a<o;)e.insertBefore(r(n[a++],1),t)}else if(o===a)for(;c<i;)u&&u.has(t[c])||e.removeChild(r(t[c],-1)),c++;else if(t[c]===n[a])c++,a++;else if(t[i-1]===n[o-1])i--,o--;else if(t[c]===n[o-1]&&n[a]===t[i-1]){const s=r(t[--i],-1).nextSibling;e.insertBefore(r(n[a++],1),r(t[c++],-1).nextSibling),e.insertBefore(r(n[--o],1),s),t[i]=n[o]}else{if(!u){u=new Map;let e=a;for(;e<o;)u.set(n[e],e++)}if(u.has(t[c])){const s=u.get(t[c]);if(a<s&&s<o){let l=c,h=1;for(;++l<i&&l<o&&u.get(t[l])===s+h;)h++;if(h>s-a){const l=r(t[c],0);for(;a<s;)e.insertBefore(r(n[a++],1),l)}else e.replaceChild(r(n[a++],1),r(t[c++],-1))}else c++}else e.removeChild(r(t[c++],-1))}return n})(e.parentNode,n,t,h,e);switch(n.length){case 1:n[0].remove();case 0:break;default:c(h(n[0],0),h(n.at(-1),-0),!1)}return r},k=new Map([["aria",(e,t)=>{for(const n in t){const r=t[n],s="role"===n?n:`aria-${n}`;null==r?g(e,s):p(e,s,r)}return t}],["class",(e,t)=>w(e,t,null==t?"class":"className")],["data",(e,t)=>{const{dataset:n}=e;for(const e in t)null==t[e]?delete n[e]:n[e]=t[e];return t}],["ref",$],["style",(e,t)=>null==t?w(e,t,"style"):C(e.style,t,"cssText")]]),A=(e,r,s)=>{switch(r[0]){case".":return x;case"?":return N;case"@":return v;default:return s||"ownerSVGElement"in e?"ref"===r?$:y:k.get(r)||(r in e?r.startsWith("on")?C:((e,r)=>{let s;do{s=n(e,r)}while(!s&&(e=t(e)));return s})(e,r)?.set?w:y:y)}},W=(e,t)=>(e.textContent=null==t?"":t,t),M=(e,t,n)=>({a:e,b:t,c:n}),S=()=>M(null,null,r),E=(e,t)=>t.reduceRight(T,e),T=(e,t)=>e.childNodes[t];var B=e=>(t,n)=>{const{a:s,b:l,c:i}=e(t,n),o=document.importNode(s,!0);let c=r;if(l!==r){c=[];for(let e,t,n=0;n<l.length;n++){const{a:s,b:i,c:f}=l[n],p=s===t?e:e=E(o,t=s);c[n]=(a=i,u=p,h=f,d=i===O?[]:i===b?S():null,{v:r,u:a,t:u,n:h,c:d})}}var a,u,h,d;return((e,t)=>({b:e,c:t}))(i?o.firstChild:new f(o),c)};const j=/^(?:plaintext|script|style|textarea|title|xmp)$/i,D=/^(?:area|base|br|col|embed|hr|img|input|keygen|link|menuitem|meta|param|source|track|wbr)$/i,L=/<([a-zA-Z0-9]+[a-zA-Z0-9:._-]*)([^>]*?)(\/?)>/g,P=/([^\s\\>"'=]+)\s*=\s*(['"]?)\x01/g,z=/[\x01\x02]/g;let F,R,Z=document.createElement("template");var G=(e,t)=>{if(t)return F||(F=document.createElementNS("http://www.w3.org/2000/svg","svg"),R=s(),R.selectNodeContents(F)),R.createContextualFragment(e);Z.innerHTML=e;const{content:n}=Z;return Z=Z.cloneNode(!1),n};const H=e=>{const t=[];let n;for(;n=e.parentNode;)t.push(t.indexOf.call(n.childNodes,e)),e=n;return t},V=()=>document.createTextNode(""),_=(t,n,s)=>{const i=G(((e,t,n)=>{let r=0;return e.join("").trim().replace(L,((e,t,r,s)=>`<${t}${r.replace(P,"=$2$1").trimEnd()}${s?n||D.test(t)?" /":`></${t}`:""}>`)).replace(z,(e=>""===e?`\x3c!--${t+r++}--\x3e`:t+r++))})(t,I,s),s),{length:o}=t;let c=r;if(o>1){const t=[],r=document.createTreeWalker(i,129);let l=0,a=`${I}${l++}`;for(c=[];l<o;){const i=r.nextNode();if(8===i.nodeType){if(i.data===a){const r=e(n[l-1])?O:b;r===b&&t.push(i),c.push(M(H(i),r,null)),a=`${I}${l++}`}}else{let e;for(;i.hasAttribute(a);){e||(e=H(i));const t=i.getAttribute(a);c.push(M(e,A(i,t,s),t)),g(i,a),a=`${I}${l++}`}!s&&j.test(i.localName)&&i.textContent.trim()===`\x3c!--${a}--\x3e`&&(c.push(M(e||H(i),W,null)),a=`${I}${l++}`)}}for(l=0;l<t.length;l++)t[l].replaceWith(V())}const{childNodes:a}=i;let{length:u}=a;return u<1?(u=1,i.appendChild(V())):1===u&&1!==o&&1!==a[0].nodeType&&(u=0),l(q,t,M(i,c,1===u))},q=new WeakMap,I="isµ";var J=e=>(t,n)=>q.get(t)||_(t,n,e);const K=B(J(!1)),Q=B(J(!0)),U=(e,{s:t,t:n,v:r})=>{if(e.a!==n){const{b:s,c:l}=(t?Q:K)(n,r);e.a=n,e.b=s,e.c=l}for(let{c:t}=e,n=0;n<t.length;n++){const e=r[n],s=t[n];switch(s.u){case O:s.v=O(s.t,X(s.c,e),s.v);break;case b:const t=e instanceof Y?U(s.c||(s.c=S()),e):(s.c=null,e);t!==s.v&&(s.v=b(s,t));break;default:e!==s.v&&(s.v=s.u(s.t,e,s.n,s.v))}}return e.b},X=(e,t)=>{let n=0,{length:r}=t;for(r<e.length&&e.splice(r);n<r;n++){const r=t[n];r instanceof Y?t[n]=U(e[n]||(e[n]=S()),r):e[n]=null}return t};class Y{constructor(e,t,n){this.s=e,this.t=t,this.v=n}toDOM(e=S()){return U(e,this)}}const ee=new WeakMap;var te=(e,t)=>{const n=ee.get(e)||l(ee,e,S()),{b:r}=n;return r!==("function"==typeof t?t():t).toDOM(n)&&e.replaceChildren(n.b.valueOf()),e};
/*! (c) Andrea Giammarchi - MIT */const ne=e=>(t,...n)=>new Y(e,t,n),re=ne(!1),se=ne(!0);export{Y as Hole,k as attr,re as html,te as render,se as svg};