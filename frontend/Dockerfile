FROM node:16.13.2-alpine AS base-node
WORKDIR /app

ARG env
ARG frontend_url
ARG backend_url
ARG firebase_api_key
ARG firebase_auth_domain
ARG firebase_storage_bucket
ARG firebase_messaging_sender_id
ARG firebase_app_id
ARG firebase_measurement_id
ARG firebase_project_id
ARG report_api_url
ARG report_api_key
ARG report_api_category_id_key
ARG read_only

ENV VITE_FRONTEND_URL=${frontend_url}
ENV VITE_BACKEND_URL=${backend_url}
ENV VITE_FIREBASE_API_KEY=${firebase_api_key}
ENV VITE_FIREBASE_AUTH_DOMAIN=${firebase_auth_domain}
ENV VITE_FIREBASE_STORAGE_BUCKET=${firebase_storage_bucket}
ENV VITE_FIREBASE_MESSAGING_SENDER_ID=${firebase_messaging_sender_id}
ENV VITE_FIREBASE_APP_ID=${firebase_app_id}
ENV VITE_FIREBASE_MEASUREMENT_ID=${firebase_measurement_id}
ENV VITE_FIREBASE_PROJECT_ID=${firebase_project_id}
ENV VITE_REPORT_API_URL=${report_api_url}
ENV VITE_REPORT_API_KEY=${report_api_key}
ENV VITE_REPORT_API_CATEGORY_ID_KEY=${report_api_category_id_key}
ENV VITE_READ_ONLY=${read_only}

COPY . .
RUN yarn install

ENV NODE_ENV=${env}
RUN yarn build

EXPOSE 3000
ENTRYPOINT ["yarn", "serve"]