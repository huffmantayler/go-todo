# Use the official Nginx image as the base
FROM nginx:latest

# Copy custom Nginx configuration file (optional)
COPY nginx.conf /etc/nginx/nginx.conf

# Expose port 80
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
