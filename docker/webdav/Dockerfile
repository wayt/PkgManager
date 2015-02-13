FROM dockerfile/nginx

RUN mkdir -p /var/www/data && \
    mkdir -p /var/www/tmp && \
    chown -R www-data:www-data /var/www

ADD webdav.conf /etc/nginx/sites-available/default

