# Build frontend dist.
FROM node:18-alpine AS frontend

WORKDIR /frontend-build

COPY . .

WORKDIR /frontend-build/frontend

RUN npm install

EXPOSE 5173

CMD ["npm", "run", "dev"]