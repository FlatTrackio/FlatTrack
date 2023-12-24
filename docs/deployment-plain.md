# Deployment on a plain Ubuntu server

Set up FlatTrack on an Ubuntu 20.04 server using systemd, certbot, and nginx.
This has been tested on Ubuntu 20.04 and will very likely work on later versions.

Commands are run as root.

## Install packages

    apt update && apt upgrade -y
    apt install -y nginx postgresql certbot python3-certbot-nginx

## DNS

Assign a host to your VPS.

e.g:

    flattrack.mydomain.com. 443 IN A	159.89.157.114

## Set up LetsEncrypt with certbot

    certbot --nginx

Notes:

-   recommended redirection of traffic to HTTPS

    Saving debug log to /var/log/letsencrypt/letsencrypt.log
    Plugins selected: Authenticator nginx, Installer nginx
    Enter email address (used for urgent renewal and security notices) (Enter 'c' to
    cancel): EMAIL@ADDRESS.COM
    
    - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
    Please read the Terms of Service at
    https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf. You must
    agree in order to register with the ACME server at
    https://acme-v02.api.letsencrypt.org/directory
    - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
    (A)gree/(C)ancel: a
    
    - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
    Would you be willing to share your email address with the Electronic Frontier
    Foundation, a founding partner of the Let's Encrypt project and the non-profit
    organization that develops Certbot? We'd like to send you email about our work
    encrypting the web, EFF news, campaigns, and ways to support digital freedom.
    - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
    (Y)es/(N)o: n
    No names were found in your configuration files. Please enter in your domain
    name(s) (comma and/or space separated)  (Enter 'c' to cancel): flattrack.mydomain.com
    Obtaining a new certificate
    Performing the following challenges:
    http-01 challenge for flattrack.mydomain.com
    Waiting for verification...
    Cleaning up challenges
    Deploying Certificate to VirtualHost /etc/nginx/sites-enabled/default
    
    Please choose whether or not to redirect HTTP traffic to HTTPS, removing HTTP access.
    - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
    1: No redirect - Make no further changes to the webserver configuration.
    2: Redirect - Make all requests redirect to secure HTTPS access. Choose this for
    new sites, or if you're confident your site works on HTTPS. You can undo this
    change by editing your web server's configuration.
    - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
    Select the appropriate number [1-2] then [enter] (press 'c' to cancel): 2
    Redirecting all traffic on port 80 to ssl in /etc/nginx/sites-enabled/default
    
    - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
    Congratulations! You have successfully enabled
    https://flattrack.mydomain.com
    
    You should test your configuration at:
    https://www.ssllabs.com/ssltest/analyze.html?d=flattrack.mydomain.com
    - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
    
    IMPORTANT NOTES:
     - Congratulations! Your certificate and chain have been saved at:
       /etc/letsencrypt/live/flattrack.mydomain.com/fullchain.pem
       Your key file has been saved at:
       /etc/letsencrypt/live/flattrack.mydomain.com/privkey.pem
       Your cert will expire on 2020-09-04. To obtain a new or tweaked
       version of this certificate in the future, simply run certbot again
       with the "certonly" option. To non-interactively renew *all* of
       your certificates, run "certbot renew"
     - Your account credentials have been saved in your Certbot
       configuration directory at /etc/letsencrypt. You should make a
       secure backup of this folder now. This configuration directory will
       also contain certificates and private keys obtained by Certbot so
       making regular backups of this folder is ideal.
     - If you like Certbot, please consider supporting our work by:
    
       Donating to ISRG / Let's Encrypt:   https://letsencrypt.org/donate
       Donating to EFF:                    https://eff.org/donate-le

## Create FlatTrack user

    useradd -m flattrack

Add a password:

    passwd flattrack

## Set up Postgres

Create the `FlatTrack` role:

    su postgres -c 'createuser flattrack -s'   

Create the `FlatTrack` database:

    su postgres -c 'createdb flattrack'

Check that you can authenticate:

    su flattrack -c 'psql -c \\conninfo'

Change the postgres password FlatTrack user:

    su flattrack -c "psql -c \"ALTER USER flattrack WITH PASSWORD 'flattrack';\""

Note: setting `'flattrack'` to the password that you want the role to have

## Set up nginx

Add a customized version of the following to `/etc/nginx/sites-available/default`:

    server {
      listen 443 ssl http2;
      server_name flattrack.mydomain.com;
    
      ssl_certificate     /etc/letsencrypt/live/flattrack.mydomain.com/fullchain.pem;
      ssl_certificate_key /etc/letsencrypt/live/flattrack.mydomain.com/privkey.pem;
      ssl_protocols       TLSv1 TLSv1.1 TLSv1.2;
      ssl_ciphers         HIGH:!aNULL:!MD5;
      add_header          Strict-Transport-Security "max-age=15552000";
    
      fastcgi_hide_header X-Powered-By;
    
      location / {
        proxy_pass http://localhost:8080;
    
        proxy_set_header X-Forwarded-Host     $host;
        proxy_set_header X-Forwarded-Server   $host;
        proxy_set_header X-Real-IP            $remote_addr;
        proxy_set_header X-Forwarded-For      $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto    $scheme;
        proxy_set_header X-Forwarded-Protocol $scheme;
        proxy_set_header X-Forwarded-Port     $server_port;
        proxy_set_header Host                 $http_host;
    
        proxy_redirect  off;
        proxy_buffering off;
    
        proxy_http_version 1.1;
        proxy_set_header Upgrade    $http_upgrade;
        proxy_set_header Connection "upgrade";
      }
    }

Reload nginx

    systemctl reload nginx

## Download FlatTrack

Download the latest zip:

    curl -L -o /tmp/flattrack.zip https://gitlab.com/flattrack/flattrack/-/jobs/artifacts/0.16.1/download?job=build-zip

Unpack the archives into /opt/flattrack

    unzip -d /tmp /tmp/flattrack.zip
    tar xvf /tmp/flattrack.tar.gz -C /opt

## Write the environment settings

Install a custom environment file into `/home/flattrack/.env`:

    APP_DB_USER=flattrack
    APP_DB_PASSWORD=flattrack
    APP_DB_HOST=localhost
    APP_DB_DATABASE=flattrack

## Install a systemd service

Install a customized version of the following, in `/etc/systemd/system/flattrack.service`:

    [Unit]
    Description=Collaborate with your flatmates
    After=postgresql.service
    After=nginx.service
    
    [Service]
    Type=simple
    ExecStart=/opt/flattrack/flattrack
    Restart=always
    User=flattrack
    Environment="APP_DB_MIGRATIONS_PATH=/opt/flattrack/migrations"
    Environment="APP_PORT=127.0.0.1:8080"
    Environment="APP_PORT_METRICS=127.0.0.1:2112"
    Environment="APP_PORT_HEALTH=127.0.0.1:8081"
    Environment="APP_DIST_FOLDER=/opt/flattrack/dist"
    Environment="APP_ENV_FILE=/home/flattrack/.env"
    
    [Install]
    WantedBy=default.target

The configuration above configures:

-   ports for FlatTrack, metrics, health
-   the database password; update `APP_DB_PASSWORD` it isn&rsquo;t `flattrack`
-   the location of the built frontend
-   the location of the environment variables file, it is recommended to use this file for fields like database credentials instead of placing them inside the systemd unit file

## Start FlatTrack

    systemctl enable --now flattrack

Check if FlatTrack is running and has started successfully.

    systemctl status flattrack

Woohoo! FlatTrack should now be running. Go to the hostname assigned in the DNS stage in a web browser to access.

## Notes

-   Once the frontend and backend is built, golang and nodejs is no longer needed or used (except for manual updates), so feel free to remove them

# Extra notes

To configure FlatTrack, please refer to the [configuration guide](./configuration.md).

