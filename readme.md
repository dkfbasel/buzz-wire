# BUILD FOR RASPBERRY PI
> GOARM=6 GOARCH=arm GOOS=linux go build -o=pi-uninacht .

# COPY TO RASPBERRY PI
> scp pi-uninacht pi@192.168.2.22:/home/pi/uninacht

# COMBINED INTO ONE COMMAND
> GOARM=6 GOARCH=arm GOOS=linux go build -o=pi-uninacht . | scp pi-uninacht pi@192.168.2.22:/home/pi/uninacht

# INSTALL FONTS ON RASPBERRY
> scp fonts/Open_Sans/* pi@192.168.2.22:~/.fonts
> scp fonts/Source_Code_Pro/* pi@192.168.2.22:~/.fonts

# RUN ON RASPBERRY VIA SSH
> ssh pi@192.168.2.22
> ./pi-uninacht

# RUN FIREFOX (ICEWEASEL) ON THE DISPLAY
> iceweasel --display=:0 http://localhost:8484

# PHYSICAL CONNECTIONS
see pin assignment on http://pi.gadgetoid.com/pinout
use the physical numbers (in black) to match gobot gpio pins
