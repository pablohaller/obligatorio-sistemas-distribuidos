FROM node:lts

# create & set working directory
RUN mkdir -p /usr/src
WORKDIR /usr/src

# copy source files
COPY . /usr/src

COPY package*.json ./
COPY prisma ./prisma/

# install dependencies
RUN npm install
RUN npm install bcrypt
RUN npm install @prisma/client

COPY . .
RUN npx prisma generate --schema ./prisma/schema.prisma
# start app
EXPOSE 3000 5555
CMD [ "npm", "run", "dev:prisma" ] 