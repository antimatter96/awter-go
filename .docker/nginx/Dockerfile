FROM nginx:1.19.0-alpine
EXPOSE 80
COPY static /var/www/static
RUN ls /var/www/static
COPY ".docker/nginx/nginx.conf" /etc/nginx/nginx.conf
