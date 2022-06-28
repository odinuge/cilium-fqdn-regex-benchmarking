for _i in {1..600}
do
    bar="$(openssl rand -base64 20 | sed 's/[0-9A-Z\/=+]*//g')"
    domain="${bar}.this-very-long-super-long-domain.no"
    [ -n "$domain" ] && echo "$domain"
done

for _i in {1..200}
do
    bar="$(openssl rand -base64 20 | sed 's/[0-9A-Z\/=+]*//g')"
    bar2="$(openssl rand -base64 20 | sed 's/[0-9A-Z\/=+]*//g')"
    bar3="$(openssl rand -base64 2 | sed 's/[0-9A-Z\/=+]*//g')"
    domain="${bar}.${bar2}.${bar3}"
    [ -n "$domain" ] && echo "$domain"
done

for _i in {1..300}
do
    bar="$(openssl rand -base64 20 | sed 's/[0-9A-Z\/=+]*//g')"
    bar2="$(openssl rand -base64 20 | sed 's/[0-9A-Z\/=+]*//g')"
    bar3="$(openssl rand -base64 2 | sed 's/[0-9A-Z\/=+]*//g')"
    domain="${bar}.${bar2}.com"
    [ -n "$domain" ] && echo "$domain"
done

for _i in {1..10}
do
    bar="$(openssl rand -base64 20 | sed 's/[0-9A-Z\/=+]*//g')"
    domain="*.${bar}.domain.no"
    [ -n "$domain" ] && echo "$domain"
done

for _i in {1..10}
do
    bar="$(openssl rand -base64 20 | sed 's/[0-9A-Z\/=+]*//g')"
    bar2="$(openssl rand -base64 20 | sed 's/[0-9A-Z\/=+]*//g')"
    domain="${bar}*.${bar2}.no"
    [ -n "$domain" ] && echo "$domain"
done

for _i in {1..10}
do
    bar="$(openssl rand -base64 20 | sed 's/[0-9A-Z\/=+]*//g')"
    bar2="$(openssl rand -base64 20 | sed 's/[0-9A-Z\/=+]*//g')"
    domain="*.*.${bar}.doksmain.no"
    [ -n "$domain" ] && echo "$domain"
done
