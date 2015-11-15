<div style="margin:0px auto; width:380px;">
<pre style="font-size:0.7em; background-color:#fff;color:#99cc00; line-height: 75%;font-weight:bold;">
                   .`.....
                :::;++:;;+`
              .:+@@@@@@:@+`.``
             ::+@@;:.'@@;';'#,;.
            .:+@;;;;;,.#':@
            :+@,;'#',;@++#      ..
            ;+,`#`      @'     `..    ``````
           .'#;#        ...'++:`,: `.......```
           ,+@'         ..'::;+`.,`.,:.,,,,.```
           .+##   :#+#  ..,;@#+;.,.,;;;:,;;,,.``
           ,+#+. ;`,;;#....+@.`.:,:;;;;..,;;:,```
           ,+@,.`+;,@:+'.,;,;::::;;;;;,``.;;;,.```
           ,#+...+:@@;+@.:::::;;;;;;;;,..,;;;;,.``
           :'+#..@+;#+#;,;:';;;+';;;;;::,:;;';,```
          `.++;;::@++;#::::';;;;;';''';;;;;'';:.`.
            :+++;:;@@+;:::;;;;;;;;'''''''''''';,..
            .'+;;'';;;:::#;;;;;;'''''''''''''';,`.
             ;+':+.:;;:;++++++++'''''''''''''';,..
              :'+'::'+''++++++++++'''''''''''';,``
               .+@+++++#++++++++#++''''''''''';,.`
                .+@@@++@@+` `:+++;++'''''''''';,.
                  ;@@@@@@+    .++;;'#''''''''';.
                   ;@@@###     ,#+++##''''''+':`
                    '@@++++`    :###@;+''''+';.
                    :#@;;+++:   `.,,` `+++#++.
                    :++;;::;,;         ###++.
                    :++;,,,,...        ,++:`
                   .+++;,,....,.        ,`
   .. .     ...,..;'+++:,.....``.
  ..:.` `` ....,;++++++;,......`
 .,   .`. ..;:':;++++++;,....... `
 `...`....:'++.`.;++++'+........
  ;,..;;.:'++:`,:::++++++,...... .
 .';,.:.;.++++`,,,:@#++++':.....
 `;++,,:.  +  .,,,;@++++++++,,..
  ,;;:''.. : .,,,,:++++++##@@,..
    .+++,. ``,,,,,;:+++++@@@@,..
     @@; ,...,,,::;;+++#@@@@@,.
     #@+ ,::,,,,,;;,###@@@@@@,`
     +@@.:;,,,,,:;'+@@@@@@@@+.
      +@@';:::;;'':@@@@@@@@@,
       .@@;+'++'';@@@@@@@@@@
          @@@++@@@+@@##@''+'
           ;@@@###@+':   ';'     ,`
              @##       ,'':    +';
              #';`+:+::#;+` `:;;+:+;.
           ..:;:':;'###'#@:'';::':;:+
 ..::::',;::;;;#:;::::::.:';'#++++':,
:::::;:;;;':;+:::,..;:;+:;''+;'+,'.::::;::,.
::::''';:;#'##:;'++'','':`::,''''#::::::::::::..
:::;:;:;:;::;;:+#'':;;;::::;:::;:::::;:::::::;;..
..:::::;::;;,;'':;,,,.;;:;;:::;::::::::::::::::.
    ;:::;;:;:::;;:::,'.;:;;;;;:::::::::::::..
        ..;;:::;;;'#+';:::.


      _                     _           _
     | |                   | |         | |
   __| |_ __  ___ ______ __| | ___   __| | ___
  / _` | '_ \/ __|______/ _` |/ _ \ / _` |/ _ \
 | (_| | | | \__ \     | (_| | (_) | (_| | (_) |
  \__,_|_| |_|___/      \__,_|\___/ \__,_|\___/
</pre>
</div>
###DNS-Digital Ocean DO - a DNS sub-domain IP updater for Digital Ocean

dns-dodo's main purpose is to update a single dns 'A' record to the public IP address of the system dns-dojo is run on.
This is very similar to the dynamic dns clients that you can download for no-ip, dyndns, etc. but
dns-dodo allows you to use your existing Digital Ocean account for this service.

In addition dns-dodo allows you to show the public ip, and show the current dns records (all types) associated with a domain on your digital ocean account.


### Who / What is Digital Ocean
Digital Ocean, https://www.digitalocean.com/ provide a Simple Cloud Infrastructure for Developers. You can setup a virtual server in seconds (about 20) and all are SSD based so are responsive.
They are affordable for the casual developer too...

### Setting up on the Digital Ocean side
1) Add a A record to an existing domain with the name set to the sub-domain you would like to use and initially set it to your droplets IP address so you can validate that dns-dodo works.
2) Get your Personal Access Token (PAT) that provides you with the authentication to talk to your Digital Ocean account. **NOTE** If someone gets hold of your PAT they have full api access to create/delete droplets and make your life a misery so please be careful.
You can get this from the applications page https://cloud.digitalocean.com/settings/applications



#### Usage

Get Help

    dns-dodo help [command]

----

Show the external IP address using the default External IP Service

    dns-dodo show-ip

Show the external IP address using an alternate IP service

    dns-dodo --ext-ip=https://api.ipify.org show-ip

----


Show the current DNS entries for a specific domain on your Digital Ocean Account

    dns-dodo show-dns --pat=[your-long-personal-access-token-here] --domain=[domain-name-to-update-dns-record-for]

Filter the current DNS entries for a specific sub-domain and domain on your Digital Ocean Account

    dns-dodo show-dns --name=home --pat=[your-pat] --domain=[domain-name]

Filter the current DNS entries for a specific record type and domain on your Digital Ocean Account

    dns-dodo show-dns --type=A --pat=[your-pat] --domain=[domain-name]


----

Update the DNS A record for a specific domain and sub-domain to the External IP Adddress

    dns-dodo update-dns --pat=[your-pat] --domain=[domain-name] --sub-domain=home


----