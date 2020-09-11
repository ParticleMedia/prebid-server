package endpoints

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// NewStatusEndpoint returns a handler which writes the given response when the app is ready to serve requests.
func NewMockEndpoint(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var payload = []byte(`{
		"id": "6b2ea487-9678-425f-88a4-eb774c384795",
		"seatbid": [
		  {
			"bid": [
			  {
				"id": "1B75CCF7-2CF0-4EAE-B24F-BB42EB4E24C7",
				"impid": "PrebidMobile",
				"price": 1,
				"adm": "<span class=\"PubAPIAd\"  id=\"1B75CCF7-2CF0-4EAE-B24F-BB42EB4E24C7\"><DIV STYLE=\"position: absolute; left: 0px; top: 0px; visibility: hidden;\"><IMG SRC=\"https://pagead2.googlesyndication.com/pagead/gen_204?id=xbid&dbm_b=AKAmf-BfqXTXV24vN64p66aDFhh_PZ4MLxyateM2QujZpEh_kAYMIdKfQDa9yw4Febhg1ar8lcKQK-w7tspv8QfDH_uhK23UjixCABWuiRCriEfdmWK9ook\" BORDER=0 WIDTH=1 HEIGHT=1 ALT=\"\" STYLE=\"display:none\"></DIV><iframe title=\"Blank\" src=\"https://googleads.g.doubleclick.net/xbbe/pixel?d=CPvjKBDbyD0Yt6XZIzAB&v=APEucNWTUGtGFdIhLrSmqz8NIaFUXU-wNPhvCBWaFa724EIubbeF3GpQ_0OMzUoDbsx89y_EA3PHQCUlaQZQUnI1zcSwKkPovg\" style=\"display:none\" aria-hidden=\"true\"></iframe><div><div style=\"position:relative; display:inline-block;\"><script data-jc=\"75\" data-jc-version=\"r20200818\">(function(){/*  Copyright The Closure Library Authors. SPDX-License-Identifier: Apache-2.0 */ var k=this||self;var l=Array.prototype.forEach?function(a,b){Array.prototype.forEach.call(a,b,void 0)}:function(a,b){for(var d=a.length,c=\"string\"===typeof a?a.split(\"\"):a,e=0;e<d;e++)e in c&&b.call(void 0,c[e],e,a)},m=Array.prototype.map?function(a,b){return Array.prototype.map.call(a,b,void 0)}:function(a,b){for(var d=a.length,c=Array(d),e=\"string\"===typeof a?a.split(\"\"):a,f=0;f<d;f++)f in e&&(c[f]=b.call(void 0,e[f],f,a));return c},n=Array.prototype.reduce?function(a,b,d){return Array.prototype.reduce.call(a,b, d)}:function(a,b,d){var c=d;l(a,function(e,f){c=b.call(void 0,c,e,f,a)});return c};function p(a){for(var b=[],d=0;d<a;d++)b[d]=\"\";return b};function q(a){q[\" \"](a);return a}q[\" \"]=function(){};function r(a,b){if(a)for(var d in a)Object.prototype.hasOwnProperty.call(a,d)&&b.call(void 0,a[d],d,a)}var t=/https?:\\/\\/[^\\/]+/;function v(a){return(a=t.exec(a))&&a[0]||\"\"};var w=/^https?:\\/\\/(\\w|-)+\\.cdn\\.ampproject\\.(net|org)(\\?|\\/|$)/;function x(a){a=a||y();for(var b=new z(k.location.href,!1),d=null,c=a.length-1,e=c;0<=e;--e){var f=a[e];!d&&w.test(f.url)&&(d=f);if(f.url&&!f.b){b=f;break}}e=null;f=a.length&&a[c].url;0!=b.depth&&f&&(e=a[c]);return new A(b,e,d)} function y(){var a=k,b=[],d=null;do{var c=a;try{var e;if(e=!!c&&null!=c.location.href)b:{try{q(c.foo);e=!0;break b}catch(h){}e=!1}var f=e}catch(h){f=!1}if(f){var g=c.location.href;d=c.document&&c.document.referrer||null}else g=d,d=null;b.push(new z(g||\"\"));try{a=c.parent}catch(h){a=null}}while(a&&c!=a);c=0;for(a=b.length-1;c<=a;++c)b[c].depth=a-c;c=k;if(c.location&&c.location.ancestorOrigins&&c.location.ancestorOrigins.length==b.length-1)for(a=1;a<b.length;++a)g=b[a],g.url||(g.url=c.location.ancestorOrigins[a- 1]||\"\",g.b=!0);return b}function A(a,b,d){this.c=a;this.f=b;this.a=d}function z(a,b){this.url=a;this.b=!!b;this.depth=null};function B(a,b,d,c,e){var f=[];r(a,function(g,h){(g=C(g,b,d,c,e))&&f.push(h+\"=\"+g)});return f.join(b)}function C(a,b,d,c,e){if(null==a)return\"\";b=b||\"&\";d=d||\",$\";\"string\"==typeof d&&(d=d.split(\"\"));if(a instanceof Array){if(c=c||0,c<d.length){for(var f=[],g=0;g<a.length;g++)f.push(C(a[g],b,d,c+1,e));return f.join(d[c])}}else if(\"object\"==typeof a)return e=e||0,2>e?encodeURIComponent(B(a,b,d,c,e+1)):\"...\";return encodeURIComponent(String(a))};function D(a,b){this.a=a;this.depth=b}function E(){function a(h,u){return null==h?u:h}var b=y(),d=Math.max(b.length-1,0),c=x(b);b=c.c;var e=c.f,f=c.a,g=[];f&&g.push(new D([f.url,f.b?2:0],a(f.depth,1)));e&&e!=f&&g.push(new D([e.url,2],0));b.url&&b!=f&&g.push(new D([b.url,0],a(b.depth,d)));c=m(g,function(h,u){return g.slice(0,g.length-u)});!b.url||(f||e)&&b!=f||(e=v(b.url))&&c.push([new D([e,1],a(b.depth,d))]);c.push([]);return m(c,function(h){return F(d,h)})} function F(a,b){var d=n(b,function(e,f){return Math.max(e,f.depth)},-1),c=p(d+2);c[0]=a;l(b,function(e){return c[e.depth+1]=e.a});return c}function G(){var a=E();return m(a,function(b){return C(b)})};function H(a){try{var b=G();b.pop();for(var d=2083-a.length-5,c=0;c<b.length;c++){var e=encodeURIComponent(b[c]);if(e.length<=d)return setTimeout(function(){var f=void 0===f?.01:f;if(!(Math.random()>f)){var g=document.currentScript;g=(g=void 0===g?null:g)&&75==g.getAttribute(\"data-jc\")?g:document.querySelector('[data-jc=\"75\"]');f=\"https://pagead2.googlesyndication.com/pagead/gen_204?id=jca&jc=75&version=\"+(g&&g.getAttribute(\"data-jc-version\")||\"unknown\")+\"&sample=\"+f;g=window;var h;if(h=g.navigator)h= g.navigator.userAgent,h=/Chrome/.test(h)&&!/Edge/.test(h)?!0:!1;h&&g.navigator.sendBeacon?g.navigator.sendBeacon(f):(g.google_image_requests||(g.google_image_requests=[]),h=g.document.createElement(\"img\"),h.src=f,g.google_image_requests.push(h))}},0),a+\"&rfl=\"+e}return a}catch(f){}return a}var I=[\"rfl\"],J=k;I[0]in J||\"undefined\"==typeof J.execScript||J.execScript(\"var \"+I[0]);for(var K;I.length&&(K=I.shift());)I.length||void 0===H?J[K]&&J[K]!==Object.prototype[K]?J=J[K]:J=J[K]={}:J[K]=H;}).call(this);</script><script>var url = 'https://googleads.g.doubleclick.net/dbm/ad?dbm_c=AKAmf-C5U68OCLVzkMy9SMQYWNYDXR_BAYgiJuUJYDF2ymSVIazkCqvdmBhx7kOmf2PPnkZGUjLr1URCQbW_E1LyQUE7wuHYiXkDZ-WWREv_UAPSBqGVeNS8jnrfZg5lxRgbOc8sTaAA_UV-IklNRXGJFup8eAef1g&cry=1&dbm_d=AKAmf-CWV_aMdbM-4-uvm84Lbr5UhZuD0Jf53Cv6TUvFxWQ6-LRpNZ2Vd8AVzEBgvSBOrOGNpHnRmt52k26PIcp3h48WgZnHNwk8PDabVJY--8_fF6vdG6UxE6j9K8GdnysCkyPQelrpPGbPRu3VG06v6MWc1ybHe9GAPVd_XmlUPGXdVHilEb7LdtaEYoaEk7M-xRpForEmlRutv4GuZTRsTvrDR613ZXG6gijgzJtzeOfcCapNp4lAcF9igpwUvJ6rW-vsKGtGnSnWsw-0ibYP8YZaXXgZ3AOB_1zXhpZ9N6wBIdiKl3KMO42as9TcCn8LBOcGHnbbtAop4psxDWbQrKObIibjrdjsTPfubUKULD1iwrBd5SZGys4lMv_JB9DlssjnZAQkC11jyh59fx2Sqncr21-eaDMWDYp_YkLclEl37hjFs6y6u2jLbxZixJwXED0zlyD0VjZLXLio1k7JDryXk8F5lNAJqmiqTQovu9bD0UVXUZQc2F_3ZzNPc99duVLrHGu7D8RbieBI32-tmY8jb8Y0czUV-D0fdnTGPn3qjjgib50I30F5sRygwRe39ZtOQEbWUSEOfK_ZxtdgqLuIe0dnZhkAyYeHM7NNUtb_ZHgruyQvUMI8uO4vCnbJ23XokN5KBRN1Q2usne8om8zR9X6Hh1iC63CPp0FdUvgCVzV0FSen8_sAgB69NS49Q0RtN7ok3XYRD0TWZUSMWfOENrtOzo9mu89cPGofZCpm7x07R8Hv-9OOGVyiUODAVrS9K76jiXwETtJL-Q3xBI_xSRGks-hHFHGkfFusI7TeekswyfcEa6Nkx0A8d-lYIh-3pTFQ97qLePBjJ5jh5V7-uoCmo_0ToQF1plGzE6n-HE6VbxMnJ_PgI1Sk--QdH8gbf5fx1WYIq_DyOPRECiQuOOdPwG2BG_gVkyhbvbeepgiXV9v-iRZIcP8iYlj4-NLSlAUE8xh7mrwxOoh7rPqdVclPaF01GFNp5kleS6GY4xdnMJtZJzTbam-FSw8koibrhLqzaTczedZPTzqNo96rs0rggmhe-ojyMqfwrjBh4Rd5kCst0ssy9zYbpomGcD6XQ_D4TmIiibpRYWmZLqXJ_O-D5Y1VN2q7u7QEsnun3pYPeksKbk0JyXriZGOYtk1tQdbtklmB-yDv3iTQ2O82gh6i562XUAx-U9YZT2ZAnBvGxp5OtWnLBCnNBUjU2UMdVcS8i1NHHqdvWSbyJ_5Tt8wdYFbRq2kpuAcDREQ0QQzcNNxxsJBNcTHuVgwmpZ-dGzY-20H4s81rpnekRg84dLJeY0l8v2V0LPhclAXaGEcDHqIbx7tQMbnmHTcP3rnEx447o5rgAPmvUuYAdgRq5Ozu1SmfS8wKJGyHvK1e-0icmkpFwp3fjdwXXbxhnVvwhY-Ja4Muq5fBXdLt-zsr6Wo17UOSZLjQyzDZ3kP6zgH0uGPVsfdATRWe1bs1n_lFL3lXL4S4STITFzNw7-7NjMh1SqwjiIosI1mVXCU9pW-IFHLizwVxadm0eKCOUfybN7WBF948tHL4-r8P37hAltv6FQlHLK_TiBoNL7rlhaIOq0sgH3Ji33tRnDO5lRVptuTM8KYkN-epohrjFQPNYHVNgGUTG_tgPxu1fL-PJ1vC21AKGoBRurU0Kit24NbetbQaUclRQfyByL7XWhO6iDXuNe8G4sgft9r6efYwq8UAc99b7Xg96z-z-afBzm2s1Gao69gWFCrL4jYGeG8A6GsPEZ_mseYttKPxHIbFScDu-3q17Kap5btcwwBAcF9rXMgkqFtxOShwsG3wnJbVZSjUZVYHs5dzWaedYvC0IsI72bkfJ1lPnDsD1klYqXVLb1ltnWreUyhrpSTFcYQDIs8WScIymHwWjtF9TklIMk9sxI02YjtXhtZrmdYADeYin-Qa3R7ktE_aCCHxYur7-m5AYDcCNa5FZclBvyqR5qg-t6vDNlsdOA4QRpPQfCbAUEGjdCWQmmF-L_9JTKeEWjuTJm8okigYR9fyD25_wkIfcqxNm69W-cDgoOwdjwu7oKzV-0qUlYAMPqrZDxtAEIgpBiaZRa1k3j-wYU0JnfgNkxLuB4LsT34I88VFa6qGTodldgk7fcVdvHxdwNUDVWaw0tanSXmQ6QQIXhTj9hR43s-YrIz1QlWpgjxnNr7ZqtXfx0N6YcT83W_K09Qj1HAfqD9Yz6x2xRx1BuVbIuUGCab1vVdLcPxZARMujNmOM2pgD4e1Rka1c5C6YfqiGZtIJ5WBIuHuJB5piqAwa94f3pH37km1XRzDZZ1VDFlZ9lSIYqrwEHjQX3PNQS5tfFVC-PrJhsBK4A3ckFukSWP2EsRkBjjJDHoGL0CDjuXqa3hReUaXVrhOwbpt1k2KWi9MXgIR1Yw7gXTfvzVWfshWR7hbrUQfkFxC2o5xRN5V3bKL6kLNNTAGyn2YS15mWkTMJX-ydeVqOibLPJ07jvFNMq5PrrrXY-3Qvs9My73rfpgjuVUiefLM4WV5xlGjdHfMCJiZXlQxC8-I6cMVt1y0YYSPn2dYTD02QucNH6rQpYlJBikocUNHG4SdtdQgMBdH3pYvj6zojfzqWvihSSl8PL0OWlNDGDJyOTxjQB7mgLhn-IMipadfvv2viBErD9mwWwOJ4Kp2RXdVrS91UjFZghAkedlxFA2B0Bk00N_0XEGsQMMy1UevlLKmTTXZfvmnvCnyMeM-bkL8cgpnMWiPJfxiShxaMazSZJJB2rXMo4-oqwidOXTqbcWyV5EwLoVvOKtxgfFwSjgNycfBgUGvGmogpZGQR2vq6JXkz3Tsevr-jFCJsAKF-dWIgFIiqHiAtdTIe_0sym3XZ6mqnw0IxTHwKwFqnCeopYE_OqPA_QP1-DaoUuOhSmRsvSoD98ZFa4l7DpcAxTqSThHo3fwNiQu6L0Yz7B4x_zfL-aMzT9z1h_TzNkiKe67nRdty1GQiO_Xwh4u00vNRR_gaClHt9Lk5s7akgXlvSTUfIA0DCXxoIe7SpFqmC4b9wyBTfpSi2bNt-KopAVpYZusR_lIi9R_LNWDSDCISRTrmTHSkp-B8SAW_svuNV9QNXo2-wtunusZeVJvw7AXQ-IZtpcKYy8OZNilDO6mU4nxjc9TOF1OVSitMMoXhn2Pt3Lj84yQh_dcaFRlpGO-7oZL8xA65Rt_dekngLeoyPnYb6r0Go0LUESOymeUFRg6L6c_6iFBx95U6llyHAwA_Gu2r7cgC2QO7ufNcDODisFNWOT-IfsqXnPAG54ViyelhGCMNuxNo-2HSRXP1DqSy53gFWlRVeuUEEpdTv-zIR6j0UvsrWB_ZEiQJ3JTlGwdO4skR4cS1_PZ9UsOpPO1EYcCuIHdWZ-1drKJ3ANejkr63QhLEZC-YRboncs8o3JMnp2cFf6riZMRICxvYiLxGp5I44PEW4pbJMbroxnO31fhbm74U9tiNEYEx2SqreTygfNGYyppErp_55Zkhyxt1aLiMR3hsIaicpkWue0kJiIYq9u4jh61kdqHV_2QomfUOpJL4flkEbMaSuIU7a-Feny6_4VP7xyRrhMwy5UfJaJlIoBOvgOqCX-jtkkDoBikXVyYtJi_T2hE6O4n_ehkpfHMMpnjCdec&pr=6:1.939933&cid=CAASUeRoI1beu7HEg0sEjnDULsqGnKzEU_eD38I45p5HjOFz_FPpyq0nrRLuJglrlB2eJupISJicIzYyZbNXIrJbGRR01PRZ8LMolZWHqhQSgQmjFg&xfc=https%3A%2F%2Fclicktrack.pubmatic.com%2FAdServer%2FAdDisplayTrackerServlet%3FclickData%3DJnB1YklkPTE1OTQxMyZzaXRlSWQ9NzA1MzM4JmFkSWQ9Mjg5NDM4OSZrYWRzaXplaWQ9MzEmdGxkSWQ9MCZjYW1wYWlnbklkPTIyOTg3JmNyZWF0aXZlSWQ9MCZ1Y3JpZD0xNTYyMDM2NjI0MzIzNjg0MzI3MiZhZFNlcnZlcklkPTI0MyZpbXBpZD0xQjc1Q0NGNy0yQ0YwLTRFQUUtQjI0Ri1CQjQyRUI0RTI0QzcmcGFzc2JhY2s9MA%3D%3D_url%3D';document.write('<script src=\"' + (window.rfl ? window.rfl(url) : url) + '\"></s' + 'cript>');</script></div></div><iframe width=\"0\" scrolling=\"no\" height=\"0\" frameborder=\"0\" src=\"https://aktrack.pubmatic.com/AdServer/AdDisplayTrackerServlet?operId=1&pubId=159413&siteId=705338&adId=2894389&adType=10&adServerId=243&kefact=1.939933&kaxefact=1.939933&kadNetFrequecy=0&kadwidth=320&kadheight=50&kadsizeid=31&kltstamp=1598222593&indirectAdId=0&adServerOptimizerId=2&ranreq=0.1&kpbmtpfact=1.939933&dcId=1&tldId=0&passback=0&svr=BID77419U&adsver=_2762913499&adsabzcid=0&ekefact=AfFCX14CAAAzOv4qjAoo8AcLRPBmm2lsiPWZtfmVNFYz3sn1&ekaxefact=AfFCX3UCAADCftfsq6olStgb70yCpQqh-YoX1iFdokWWLpxn&ekpbmtpfact=AfFCX4sCAAByLHJT0YW0lqyxDVes_sk1us5Ou3pQu4T4uLTD&enpp=AfFCX6ACAAD-C85RKn-5BK_BUHXNaeOYWh3pD4CRvqNlRqfJ&pubBuyId=10411&crID=74863287&lpu=farmers.com&ucrid=15620366243236843272&campaignId=22987&creativeId=0&pctr=0.000000&wDSPByrId=668155&wDspId=80&wbId=6&wrId=2941032&wAdvID=8445&wDspCampId=41351549&isRTB=1&rtbId=D9CC6F3A-D75D-421B-AD4B-5D0A243E149E&imprId=1B75CCF7-2CF0-4EAE-B24F-BB42EB4E24C7&oid=1B75CCF7-2CF0-4EAE-B24F-BB42EB4E24C7&mobflag=1&ismobileapp=1&modelid=16883&osid=206&udidtype=1&cntryId=232&sec=1&pAuSt=0\" style=\"position:absolute;top:-15000px;left:-15000px\" vspace=\"0\" hspace=\"0\" marginwidth=\"0\" marginheight=\"0\" allowtransparency=\"true\" name=\"pbeacon\"></iframe></span> <!-- PubMatic Ad Ends -->",
				"adomain": [
				  "farmers.com"
				],
				"iurl": "http://googleads.g.doubleclick.net/dbm/ad?dbm_p=AKAmf-DjAq8yBeR5axwcq3rMhCJ50mvk83rWY_nRo1LsdZOqF-17h3M8xDp_aIRs4qBmfVRmM8ykBqY3yDP9MKeI8TFdE_rSX6LUfqJuyH9dkzwGOLhWvaPhIb5PzNoAfn5XLYEirrlQ",
				"cid": "22987",
				"crid": "74863287",
				"w": 320,
				"h": 50,
				"ext": {
				  "prebid": {
					"cache": {
					  "key": "",
					  "url": "",
					  "bids": {
						"url": "prebid-cache-stage.newsbreak.com/cache?uuid=6b380143-1aa5-4ab8-bca0-7dd236c9c35d",
						"cacheId": "6b380143-1aa5-4ab8-bca0-7dd236c9c35d"
					  }
					},
					"targeting": {
					  "hb_bidder": "pubmatic",
					  "hb_bidder_pubmatic": "pubmatic",
					  "hb_cache_host": "prebid-cache-stage.newsbreak.com",
					  "hb_cache_host_pubmat": "prebid-cache-stage.newsbreak.com",
					  "hb_cache_id": "6b380143-1aa5-4ab8-bca0-7dd236c9c35d",
					  "hb_cache_id_pubmatic": "6b380143-1aa5-4ab8-bca0-7dd236c9c35d",
					  "hb_cache_path": "/cache",
					  "hb_cache_path_pubmat": "/cache",
					  "hb_env": "mobile-app",
					  "hb_env_pubmatic": "mobile-app",
					  "hb_pb": "0.75",
					  "hb_pb_pubmatic": "0.75",
					  "hb_size": "320x50",
					  "hb_size_pubmatic": "320x50"
					},
					"type": "banner"
				  },
				  "bidder": {
					"dspid": 466,
					"advid": 93051
				  }
				}
			  }
			],
			"seat": "pubmatic"
		  }
		],
		"cur": "USD",
		"ext": {
		  "responsetimemillis": {
			"openx": 192,
			"pubmatic": 151,
			"triplelift": 195,
			"vrtcal": 422
		  },
		  "tmaxrequest": 1000
		}
	  }`)
	w.Write(payload)

}
