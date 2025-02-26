# Use an official Node.js image for building
FROM node:18 AS build

# Set working directory
WORKDIR /app

# Copy package.json and package-lock.json (or yarn.lock)
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# Build the Vite project
RUN npm run build

# Use a lightweight web server for serving static files
FROM nginx:alpine AS production

# Copy the built React Vite files to the Nginx directory
COPY --from=build /app/dist /usr/share/nginx/html

# Copy a custom Nginx config file (optional)
COPY nginx.conf /etc/nginx/nginx.conf

# Expose port 80
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
