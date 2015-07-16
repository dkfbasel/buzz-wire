# BUILD FOR RASPBERRY PI
> GOARM=6 GOARCH=arm GOOS=linux go build -o=pi-uninacht .

# COPY TO RASPBERRY PI
> scp pi-uninacht pi@192.168.2.22:/home/pi

# RUN ON RASPBERRY VIA SSH
ssh pi@192.168.2.22
./pi-uninacht
