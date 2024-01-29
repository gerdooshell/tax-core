FROM nginx:latest

COPY ./nginx.conf /etc/nginx/nginx.conf
RUN rm -rf /usr/share/nginx/html/*
COPY /etc/letsencrypt/live/gerdooshell.com /usr/share/gerdooshell.com

EXPOSE 80 443

ENTRYPOINT ["nginx", "-g", "daemon off;"]