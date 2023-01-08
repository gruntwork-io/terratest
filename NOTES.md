NOTES..




# WSL2 issues with DNS and Provider Registration
See these:
https://github.com/hashicorp/terraform-provider-azurerm/issues/15345
https://github.com/golang/go/issues/51127

https://github.com/microsoft/WSL/issues/5420#issuecomment-1248791740
https://github.com/microsoft/WSL/issues/5420

## To workaround
### check the current configurations
cat /etc/resolv.conf

### create the /etc/wsl.conf file to disable /etc/resolv.conf generation
printf "[network]\ngenerateResolvConf = false" >> /etc/wsl.conf

### remove the symlink WSL2 created
rm /etc/resolv.conf

### finally, add your public DNS server to the /etc/resolv.conf file. If your VPN server has a nameserver, add it too
echo "nameserver 8.8.8.8" >> /etc/resolv.conf
echo "nameserver 1.1.1.1" >> /etc/resolv.conf

>Note: the /etc/resolv.conf file gets removed at each reboot