FROM redis:6.0.5-alpine

RUN echo never > /sys/kernel/mm/transparent_hugepage/enable

CMD redis-server
