#!/bin/bash
# Check if the server IP argument is provided
if [ $# -ne 1 ]; then
    echo "Usage: $0 <server_ip>"
    exit 1
fi

#------------------------------- Configure initial set up if not done yet --------------------
# Set the remote address
REMOTE_SERVER="root@$1"
REMOTE_PATH="/root/flight-prices"
NGINX_CONFIG_LOCAL="configs/nginx/sites-available/flight-prices.conf"
NGINX_CONFIG_REMOTE="/etc/nginx/sites-available/flight-prices"

# Check if Nginx is installed, if not install it
ssh ${REMOTE_SERVER} "which nginx || sudo apt-get install nginx"

# Check if directories are set up, if not set them up
ssh ${REMOTE_SERVER} "[ -d ${REMOTE_PATH} ] || mkdir -p ${REMOTE_PATH}"

# Check if SSL certificates exist
SSL_CERTIFICATE_PATH="/etc/letsencrypt/live/twowayflights.com/fullchain.pem"
SSL_SETUP=$(ssh ${REMOTE_SERVER} "[ -f ${SSL_CERTIFICATE_PATH} ] && echo '1' || echo '0'")

if [ "$SSL_SETUP" -eq "0" ]; then
    # Certificates are not set up, run Certbot
    read -p "Enter your email address for Let's Encrypt: " EMAIL_ADDRESS
    echo "Running Certbot..."
    ssh ${REMOTE_SERVER} "sudo certbot --nginx -d twowayflights.com -d www.twowayflights.com --email $EMAIL_ADDRESS"
    echo "Certbot setup completed."
else
    echo "SSL certificates found. Skipping Certbot setup."
fi

# source the local environment variables
source .env

#------------------------------- Build and transfer application files -------------------------
# Step 1: Build the Linux binary
echo "Building the Linux binary..."
GOOS=linux GOARCH=amd64 go build -o flight-prices-binary main.go
if [ $? -ne 0 ]; then
    echo "Error building the binary. Exiting."
    exit 1
fi

# Step 2: SSH into the remote server and stop the service
echo "Stopping the remote service..."
ssh ${REMOTE_SERVER} "sudo systemctl stop flight-prices"
if [ $? -ne 0 ]; then
    echo "Error stopping the service. Exiting."
    exit 1
fi

# Function to update and check the environment variable in the service file
update_env_var() {
    VAR_NAME="$1"
    VAR_VALUE="$2"
    
    # Check and potentially update the environment variable on the remote server
    ssh ${REMOTE_SERVER} "
      # Ensure the variable is in the [Service] section
      if grep -q 'Environment=\"${VAR_NAME}=' /etc/systemd/system/flight-prices.service; then
        # If it exists but is different, then replace
        if ! grep -q 'Environment=\"${VAR_NAME}=${VAR_VALUE}\"' /etc/systemd/system/flight-prices.service; then
          sed -i '/Environment=\"${VAR_NAME}=/c\Environment=\"${VAR_NAME}=${VAR_VALUE}\"' /etc/systemd/system/flight-prices.service
          echo "${VAR_NAME} updated"
        fi
      else
        # If it doesn't exist, insert it after [Service]
        sed -i '/\[Service\]/a Environment=\"${VAR_NAME}=${VAR_VALUE}\"' /etc/systemd/system/flight-prices.service
        echo "${VAR_NAME} added"
      fi
    "
}

# Update the environment variables
update_env_var "AMADEUS_API_TEST_KEY" "$AMADEUS_API_TEST_KEY"
update_env_var "AMADEUS_API_TEST_SECRET" "$AMADEUS_API_TEST_SECRET"
update_env_var "AMADEUS_API_PROD_KEY" "$AMADEUS_API_PROD_KEY"
update_env_var "AMADEUS_API_PROD_SECRET" "$AMADEUS_API_PROD_SECRET"

# Update feature flags
update_env_var "USE_PROD_API" "$USE_PROD_API"

# Step 3: Push the binary to the remote server
echo "Transferring the binary..."
scp flight-prices-binary ${REMOTE_SERVER}:${REMOTE_PATH}
if [ $? -ne 0 ]; then
    echo "Error transferring the binary. Exiting."
    exit 1
fi

# Step 4: Push the UI to the remote server
echo "Transferring the UI..."
scp -r ui ${REMOTE_SERVER}:${REMOTE_PATH}
if [ $? -ne 0 ]; then
    echo "Error transferring the UI. Exiting."
    exit 1
fi

# Step 5: Push the development UI to the remote server
echo "Transferring the development UI..."
scp -r ui-dev ${REMOTE_SERVER}:${REMOTE_PATH}
if [ $? -ne 0 ]; then
    echo "Error transferring the development UI. Exiting."
    exit 1
fi

# Step 6: Push the Nginx config to the remote server
echo "Transferring the Nginx configuration..."
scp ${NGINX_CONFIG_LOCAL} ${REMOTE_SERVER}:${NGINX_CONFIG_REMOTE}
if [ $? -ne 0 ]; then
    echo "Error transferring the Nginx configuration. Exiting."
    exit 1
fi

#------------------------------- Reload services -------------------------

# Step 7: SSH into the remote server and run Certbot with option 1 (reinstall)
echo "Reinstalling the existing Let's Encrypt certificate on the remote server..."
ssh ${REMOTE_SERVER} "echo 1 | sudo certbot --nginx -d twowayflights.com -d www.twowayflights.com > /dev/null 2>&1"
CERTBOT_EXIT_CODE=$?

# Check the Certbot exit code
if [ $CERTBOT_EXIT_CODE -ne 0 ]; then
    echo "Error running Certbot. Exiting."
    exit 1
fi

# Step 8: SSH into the remote server and reload Nginx
echo "Reloading Nginx on the remote server..."
ssh ${REMOTE_SERVER} "sudo nginx -t && sudo systemctl reload nginx"
if [ $? -ne 0 ]; then
    echo "Error reloading Nginx. Exiting."
    exit 1
fi

# Step 9: Ensure the Nginx configuration is linked to sites-enabled
echo "Ensuring the Nginx configuration is linked to sites-enabled..."
ssh ${REMOTE_SERVER} "sudo ln -sf ${NGINX_CONFIG_REMOTE} /etc/nginx/sites-enabled/"
if [ $? -ne 0 ]; then
    echo "Error linking the Nginx configuration. Exiting."
    exit 1
fi

# Step 10: SSH into the remote server and start the service
echo "Starting the remote service..."
ssh ${REMOTE_SERVER} "systemctl daemon-reload && sudo systemctl start flight-prices"
if [ $? -ne 0 ]; then
    echo "Error starting the service. Exiting."
    exit 1
fi

echo "Deployment complete!"
