; ### This is just an example file, please change accordingly ###
$TTL 10s
@    IN    SOA    localhost. hostmaster.portal-spartan.com ( 2018032501 1d 2h 4w 1h )
@    IN    A      192.168.33.10
     IN    A      192.168.33.11