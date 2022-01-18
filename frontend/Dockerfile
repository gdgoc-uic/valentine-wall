FROM node:16.13.2
ARG env
ENV NODE_ENV=$env
WORKDIR /app
COPY . .
CMD [ -d "node_modules" ] && npm run start && npm run build

FROM caddy:2.4.6
EXPOSE 3000
ENTRYPOINT ["caddy", "file-server", "--root", "/app/dist", "--listen", ":3000"]