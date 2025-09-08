FROM node:20-alpine
WORKDIR /app
COPY package.json package-lock.json* .npmrc* ./
RUN npm ci || npm install
COPY . .
RUN npm run build || true # { SPECULATION } allow dev
EXPOSE 3000
CMD ["npm","run","dev"]
