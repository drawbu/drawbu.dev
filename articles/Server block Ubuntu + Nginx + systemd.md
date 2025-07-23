---
title: How to setup a server block with Nginx + systemd on Ubuntu with HTTPS
date: 2023-05-02
uri: server-block-nginx-systemd-ubuntu-https
---

In this guide, we are going to use the keyword for you to replace. There will be `REPOSITORY ` and `DOMAIN_NAME` to replace.

> For example:
> - REPOSITORY: https://github.com/Sigmanificient/1l.is
> - DOMAIN_NAME: 1l.is

We'll use Nginx to host the webservers, and CertBot to setup the SSL certificate (free) to access the website using HTTPS.

## Setup
To install Nginx on Ubuntu, just follow [this cool guide](https://www.digitalocean.com/community/tutorials/how-to-install-nginx-on-ubuntu-20-04) from Digital Ocean,
Ubuntu they explain everything.

Then install [CertBot](https://certbot.eff.org/instructions?ws=nginx&os=ubuntufocal) to enable HTTPS certification.

Next don't forget to enable those ports on your firewall:
```bash
sudo ufw enable
sudo ufw allow 'Nginx HTTP'
sudo ufw allow 'Nginx HTTPS'
sudo ufw allow 'OpenSSH'
```

## Clone the repo
First create the folder where all the webpages will live (this is just a convention, but you can place them wherever you want, just remember where).
```bash
sudo mkdir -p /var/www
sudo chown -R $USER /var/www
sudo git clone REPOSITORY /var/www/DOMAIN_NAME
```


## Create the systemd service
> **Note**
> If you are serving a static website, this step is not needed.

The systemd service is a file that is going to generate the .sock file. This file is going to be the middle-man between your web server and Nginx.

Also, this is a Python web server, so maybe you want to change `ExecStart`.
```bash
sudo vi /etc/systemd/system/DOMAIN_NAME.service
```

```toml
[Unit]
Description=Gunicorn instance to serve DOMAIN_NAME

[Service]
User=nginx
Group=www-data
WorkingDirectory=/var/www/DOMAIN_NAME
RuntimeDirectory=DOMAIN_NAME;
Environment="PATH=/var/www/DOMAIN_NAME/venv/bin/"
ExecStart=gunicorn --workers 3 --bind unix:/run/DOMAIN_NAME/DOMAIN_NAME.sock -m 007 wsgi:app
```

```bash
sudo systemctl daemon-reload
```

And then you can start the service and see if it worked fine:
```bash
sudo systemctl start DOMAIN_NAME
sudo systemctl status DOMAIN_NAME
```

If like me, you are generating a socket file, you can check if this step already worked:
```
curl +X GET --unix-socket "/run/DOMAIN_NAME/DOMAIN_NAME.sock" http:/foo
```

## Setup Nginx
Add to Nginx:
```bash
sudo nano /etc/nginx/sites-available/DOMAIN_NAME
```

If it's a static site:
```
server {
    listen 80;
    listen [::]:80;
	server_name DOMAIN_NAME;

	root /var/www/DOMAIN_NAME/html;
	index index.html index.htm index.nginx-debian.html;

	location / {
		try_files $uri $uri/ =404;
	}
}
```
If it's a dynamic (here using a unix socket)
```
server {
    listen 80;
    listen [::]:80;
	server_name DOMAIN_NAME;

    location / {
        include proxy_params;
        proxy_pass http://unix:/run/DOMAIN_NAME/DOMAIN_NAME.sock;
    }
}
```

Add to Nginx, verify, and restart
```bash
sudo ln -s /etc/nginx/sites-available/DOMAIN_NAME /etc/nginx/sites-enabled/  # Enable the Nginx block
sudo nginx -t  # Check Nginx config
sudo systemctl restart nginx  # Restart Nginx
```

Add HTTPS :
```bash
sudo certbot --nginx
```
