# oci-sd Docker Image

This Docker image provides a streamlined environment for running the `oci-sd` application. It utilizes multi-stage builds to optimize the image size, ensuring it contains only the necessary components.

## Building the Image

To build the `oci-sd` Docker image, follow these steps:

1. Ensure you have Docker installed on your system.

2. Clone the `oci-sd` repository:

   ```bash
   git clone https://github.com/sw-samuraj/oci-sd.git
   ```

3. Navigate to the `oci-sd/dockerfile` directory:

   ```bash
   cd oci-sd/dockerfile
   ```

4. Create the Docker image using the provided Dockerfile:

   ```bash
   docker build -t oci-sd .
   ```

## Running the Container

After building the image, you can run a container based on it. For example:

```bash
docker run -d \
  -v /home/user/:/oci-sd/.oci \                # Mounting the host directory to the container for OCI API key file
  -e USER="my_user_id" \                       # Oracle Cloud user ID
  -e FINGERPRINT="42:42:42" \                  # Unique identifier for the API key
  -e KEYFILE="/oci-sd/.oci/my_file.pem" \      # Path to the API key file within the container
  -e TENANCY="my_tenancy_id" \                 # Oracle Cloud tenancy ID
  -e REGION="us-phoenix-1" \                   # Oracle Cloud region
  -e COMPARTMENT="my_compartment_id" \         # Oracle Cloud compartment ID
  -e REFRESH_INTERVAL="600s" \                 # Time interval for resource discovery (default is 60s)
  -e OCI_DISCOVERY_OUTPUT="my_output.json" \   # Output file for discovered resources
  -e PASSPHRASE="my_passphrase" \              # Passphrase for the API key (if applicable)
  -e PORT="9200" \                             # Port for the application (default is 9100)
  --name oci-discovery oci-sd
```

In this example, the container is started with the environment variable `OCI_DISCOVERY_OUTPUT` set to `"my_output.json"`. The container also mounts a directory for the API key and specifies the key file location.

Feel free to customize the command as per your specific use case.

## Additional Notes

It's recommended to replace `"my_output.json"` with the actual desired output file name.

For more information about the `oci-sd` application, refer to the [official repository](https://github.com/sw-samuraj/oci-sd).

---
