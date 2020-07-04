#!/usr/bin/env bash
rm caddy.log
touch caddy.log
caddy -conf Caddyfile 
