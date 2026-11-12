FROM node:20-alpine

WORKDIR /app

# Copiamos package.json primero para aprovechar el cache
COPY package*.json ./

RUN npm install --omit=dev

# Copiamos el resto del código
COPY . .

EXPOSE 3000

CMD ["npm", "start"]